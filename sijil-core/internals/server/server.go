package server

import (
	"expvar"
	"fmt"
	"os"
	"sijil-core/internals/auth"
	"sijil-core/internals/core/identity"
	"sijil-core/internals/core/observability"
	"sijil-core/internals/core/projects"
	"sijil-core/internals/database"
	"sijil-core/internals/hub"
	"sijil-core/internals/ingest"
	"sijil-core/internals/shared"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

type Server struct {
	db           *pgxpool.Pool
	ingestEngine *ingest.IngestionEngine
	hub          *hub.Hub
	authCache    *auth.AuthCache
	jwtSecret    string
	// -----
	Router               *gin.Engine
	identityRepo         identity.Repository
	identityHandler      *identity.Handler
	projectHandler       *projects.Handler
	observabilityHandler *observability.Handler
}

func NewServer(db *pgxpool.Pool, ingestEngine *ingest.IngestionEngine, hub *hub.Hub, authCache *auth.AuthCache, jwtSecret string, handler shared.Handlers) *Server {
	s := &Server{
		db:           db,
		ingestEngine: ingestEngine,
		hub:          hub,
		authCache:    authCache,
		jwtSecret:    jwtSecret,

		identityRepo:         handler.IdentityRepo,
		identityHandler:      handler.Identity,
		projectHandler:       handler.Projects,
		observabilityHandler: handler.Observability,
	}

	router := gin.Default()
	s.registerRoutes(router)
	s.Router = router

	return s
}

func ParseTime(c *gin.Context) (time.Time, time.Time) {

	fromStr := c.Query("from")
	fromTime := time.Now().Add(-24 * time.Hour)
	if fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			fromTime = parsed
		}
	}

	toStr := c.Query("to")
	toTime := time.Now()
	if toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			toTime = parsed
		}
	}

	return fromTime, toTime
}

func (s *Server) registerRoutes(router *gin.Engine) {
	router.Use(s.SecurityHeadersMiddleware())

	// if s.ingestEngine != nil {
	// 	router.Use(s.DogFoodMiddleware())
	// }

	router.Static("/static/avatars", "./data/avatars")

	apiv1 := router.Group("/api/v1")

	authGroup := apiv1.Group("/auth")
	{
		authGroup.POST("/register", s.rateLimitMiddleware(rate.Limit(5.0/60.0), 5), s.identityHandler.Register)
		authGroup.POST("/login", s.rateLimitMiddleware(rate.Limit(5.0/60.0), 5), s.identityHandler.Login)

		// Sensitive routes
		authGroup.GET("/verify", s.rateLimitMiddleware(rate.Limit(1.0/60.0), 5), s.identityHandler.VerifyEmail)
		authGroup.POST("/forgot-password", s.rateLimitMiddleware(rate.Limit(1.0/60.0), 5), s.identityHandler.ForgotPassword)
		authGroup.POST("/reset-password", s.rateLimitMiddleware(rate.Limit(1.0/60.0), 5), s.identityHandler.ResetPassword)
	}

	// For our Control Room
	protected := apiv1.Group("")
	protected.Use(s.authMiddleware())
	{
		// Projects
		protected.POST("/projects", s.projectHandler.Create)
		protected.GET("/projects", s.projectHandler.List)
		protected.POST("/projects/:id/members", s.projectHandler.AddMember)
		protected.GET("/projects/:id/members", s.projectHandler.GetMembers)

		// Logs
		protected.GET("logs", s.observabilityHandler.Search)
		protected.GET("logs/stats", s.handleGetStats)
		protected.GET("logs/summary", s.handleGetSummary)
		protected.GET("logs/ws", s.handleWsLogic)
	}

	// Webhooks
	apiv1.POST("/webhooks/paystack", s.handlePayStackWebhook)
	apiv1.POST("/webhooks/lemonsqueezy", s.handleLemonWebhook)

	// For the Loading Dock (Agent route)
	apiv1.POST("/logs", s.apiKeyAuthMiddleware(), s.observabilityHandler.Ingest)

	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")
	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		adminPass = "password"
	}

	adminGroup := router.Group("/debug", gin.BasicAuth(gin.Accounts{
		adminUser: adminPass,
	}))
	adminGroup.GET("/vars", gin.WrapH(expvar.Handler()))
}

func (s *Server) handleWsLogic(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	projectIDStr := c.Query("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invaid project_id "})
		return
	}

	isOwner, err := database.CheckProjectAccess(c.Request.Context(), s.db, userID.(int), projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	if !isOwner {
		c.JSON(403, gin.H{"error": "You are not allowed to perform this action"})
		return
	}

	hub.ServeWs(s.hub, c.Writer, c.Request, projectID, userID.(int))
}

func (s *Server) handleGetStats(c *gin.Context) {

	userID, _ := c.Get("userID")
	projectIDStr := c.Query("project_id")
	projectID, _ := strconv.Atoi(projectIDStr)

	isOwner, err := database.CheckProjectAccess(c.Request.Context(), s.db, userID.(int), projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "error checking user's project ownership"})
		return
	}

	if !isOwner {
		c.JSON(403, gin.H{"error": "You are not allowed to view this project"})
		return
	}

	fromTime, toTime := ParseTime(c)

	bucketParams := c.DefaultQuery("bucket", "1hr")
	bucketMap := map[string]string{
		"1m": "1 minutes", "5m": "5 minutes", "15m": "15 minutes", "30m": "30 minutes",
		"1h": "1 hour", "6h": "6 hours", "12h": "12 hours", "1d": "1 day",
	}

	pgBucket, ok := bucketMap[bucketParams]
	if !ok {
		pgBucket = "1 hour"
	}

	stats, err := database.GetLogStats(c.Request.Context(), s.db, projectID, fromTime, toTime, pgBucket)
	if err != nil {
		fmt.Printf("Stats Error: %v \n", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": stats})
}

func (s *Server) handleGetSummary(c *gin.Context) {

	userID, _ := c.Get("userID")
	projectIDStr := c.Query("project_id")
	projectID, _ := strconv.Atoi(projectIDStr)

	isOwner, _ := database.CheckProjectAccess(c.Request.Context(), s.db, userID.(int), projectID)
	if !isOwner {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	fromTime, toTime := ParseTime(c)

	summary, err := database.GetlogSummary(c.Request.Context(), s.db, projectID, fromTime, toTime)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": summary})
}
