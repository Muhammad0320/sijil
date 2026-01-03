package server

import (
	"errors"
	"expvar"
	"fmt"
	"log"
	"os"
	"sijil-core/internals/auth"
	"sijil-core/internals/core/identity"
	"sijil-core/internals/core/observability"
	"sijil-core/internals/core/projects"
	"sijil-core/internals/database"
	"sijil-core/internals/hub"
	"sijil-core/internals/ingest"
	"sijil-core/internals/shared"
	"sijil-core/internals/utils"
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
	identityService      *identity.Service
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
		identityService:      handler.IdentityService,
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

	if s.ingestEngine != nil {
		router.Use(s.DogFoodMiddleware())
	}

	router.Static("/static/avatars", "./data/avatars")

	apiv1 := router.Group("/api/v1")

	authGroup := apiv1.Group("/auth")
	{
		authGroup.POST("/register", s.identityHandler.Register)
		authGroup.POST("/login", s.identityHandler.Login)

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

	// For the Loading Dock (Agent route)
	apiv1.POST("/logs", s.apiKeyAuthMiddleware(), s.observabilityHandler.Ingest)

	adminGroup := router.Group("/debug", gin.BasicAuth(gin.Accounts{
		os.Getenv("ADMIN_USER"): os.Getenv("ADMIN_PASS"),
	}))
	adminGroup.GET("/vars", gin.WrapH(expvar.Handler()))
}

func (s *Server) handleVerifyEmail(c *gin.Context) {

	token := c.Query("token")
	if token == "" {

		c.JSON(400, gin.H{"error": "token required"})
		return
	}

	tokenHash := utils.Hashtoken(token)

	success, err := database.VerifyUserAccount(c.Request.Context(), s.db, tokenHash)
	if err != nil || !success {
		c.JSON(400, gin.H{"error": "Invalid token or token expired"})
		return
	}

	c.JSON(200, gin.H{"message": "email verified successfully"})
}

func (s *Server) handleForgotPassword(c *gin.Context) {

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.BindJSON(&req); err != nil {
		return
	}

	// Token management
	token, _ := utils.GenerateRandomString(32)
	hashToken := utils.Hashtoken(token)

	// Set in DB
	expires := time.Now().Add(15 * time.Minute)
	_ = database.SetPasswordResetToken(c.Request.Context(), s.db, hashToken, expires, req.Email)

	fmt.Printf("DEBUG: Send reset link to %s :/reset-password?:token=%s\n", req.Email, token)

	c.JSON(200, gin.H{"message": "If an account exists, a reset link has been sent"})
}

func (s *Server) handleResetPassword(c *gin.Context) {

	var req struct {
		Token      string `json:"token" binding:"required"`
		NewPassord string `json:"password" binding:"required,min=8"`
	}

	hashToken := utils.Hashtoken(req.Token)
	hashPassword, _ := auth.HashPasswod(req.NewPassord)

	success, err := database.ResetPasswordByToken(c.Request.Context(), s.db, hashToken, hashPassword)
	if err != nil || !success {
		c.JSON(400, gin.H{"error": "invalid token or token expired"})
		return
	}

	c.JSON(200, gin.H{"message": "password reset successfully"})
}

func (s *Server) handleGetLogs(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(400, gin.H{"error": "project_id query param is required"})
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid project_id"})
		return
	}

	// Security check: Does this user own this project
	isOwner, err := database.CheckProjectAccess(c.Request.Context(), s.db, userID.(int), projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "server error checking ownership"})
		return
	}

	if !isOwner {
		c.JSON(403, gin.H{"error": "you do not have access to this project"})
		return
	}

	limitString := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitString)
	if err != nil || limit < 1 {
		limit = 100
	}

	offsetString := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetString)
	if err != nil || offset < 0 {
		offset = 0
	}

	serachQuery := c.DefaultQuery("q", "")

	ctx := c.Request.Context()
	plan, _ := s.identityRepo.GetPlanByUserID(ctx, userID.(int))

	logs, err := database.GetLogs(ctx, s.db, projectID, limit, offset, serachQuery, plan.RetentionDays)
	if err != nil {
		fmt.Printf("Failed to get logs: %v\n", err)

		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"logs": logs})
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

func (s *Server) handleCreateProject(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name string `json:"name" validate:"required,min=3"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	randkey, _ := utils.GenerateRandomString(16)
	apiKey := "pk_live_" + randkey

	randSecret, _ := utils.GenerateRandomString(16)
	apiSecret := "sk_live_" + randSecret

	secretHash, err := auth.HashPasswod(apiSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to secure peoject keys"})
		return
	}

	plan, _ := s.identityRepo.GetPlanByUserID(c.Request.Context(), userID.(int))

	// Check plan limit
	count, _ := database.GetProjectCountByUserID(c.Request.Context(), s.db, userID.(int))

	if count >= plan.MaxProjects {
		// I don't know the right status code to send
		c.JSON(400, gin.H{"error": fmt.Sprintf("You've reached you limit of %d members! kindly unpgrade your plan.", count)})
		return
	}

	projectID, err := database.CreateProject(c.Request.Context(), s.db, userID.(int), req.Name, apiKey, secretHash)
	if err != nil {
		if errors.Is(err, database.NameExists) {
			c.JSON(409, gin.H{"error": "409 Conflict: project name already exists"})
			return
		}

		c.JSON(500, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(201, gin.H{
		"message":    "project created",
		"project_id": projectID,
		"api_key":    apiKey,
		"api_secret": apiSecret, // This is the first and last time to send the api key
	})
}

func (s *Server) handleListProjects(c *gin.Context) {

	userID, _ := c.Get("userID")
	projects, err := database.GetUserProjects(c.Request.Context(), s.db, userID.(int))
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"projects": projects})
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

func (s *Server) handleAddMember(c *gin.Context) {

	userID := c.GetInt("userID")
	projectIDStr := c.Param("id")
	projectID, _ := strconv.Atoi(projectIDStr)

	role, err := database.GetProjectRole(c.Request.Context(), s.db, userID, projectID)
	if err != nil {
		log.Fatalf("Cannot get role: %s", err)
		return
	}

	if role != "owner" && role != "admin" {
		c.JSON(403, gin.H{"error": "only admins can invite members"})
		return
	}

	plan, _ := s.identityRepo.GetPlanByUserID(c.Request.Context(), userID)

	// Check plan limit
	count, _ := database.GetMemberCountByProjectID(c.Request.Context(), s.db, projectID)

	if count >= plan.MaxMemebers {
		// I don't know the right status code to send
		c.JSON(402, gin.H{"error": fmt.Sprintf("You've reached you limit of %d members! kindly unpgrade your plan.", count)})
		return
	}

	// Add member
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	c.BindJSON(&req)

	if err := database.AddProjectMember(c.Request.Context(), s.db, projectID, req.Email, req.Role); err != nil {
		if err.Error() == "user is already a member" {
			c.JSON(409, gin.H{"error": "user is already a member"})
			return
		}

		c.JSON(500, gin.H{"error": "failed to add member"})
		return
	}
	c.JSON(200, gin.H{"message": "member invited"})
}

func (s *Server) handleGetMembers(c *gin.Context) {
	userID := c.GetInt("userID")
	projectIDStr := c.Param("id")
	projectID, _ := strconv.Atoi(projectIDStr)

	hasAccess, _ := database.CheckProjectAccess(c.Request.Context(), s.db, userID, projectID)
	if !hasAccess {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	members, err := database.GetProjectMembers(c.Request.Context(), s.db, projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch members"})
		return
	}

	c.JSON(200, gin.H{"members": members})
}
