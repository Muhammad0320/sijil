package server

import (
	"errors"
	"fmt"
	"log-engine/internals/auth"
	"log-engine/internals/database"
	"log-engine/internals/hub"
	"log-engine/internals/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db *pgxpool.Pool
	logQueue chan <- database.LogEntry
	hub *hub.Hub
	jwtSecret string
}

func NewServer (db *pgxpool.Pool, logQueue chan <- database.LogEntry, hub *hub.Hub, jwtSecret string) *Server {
	return &Server{
		db: db,
		logQueue: logQueue,
		hub: hub,
		jwtSecret:  jwtSecret,
	}
}

func ParseTime(c *gin.Context) (time.Time, time.Time) {

	fromStr := c.Query("from")
	fromTime := time.Now().Add(-24 * time.Hour )
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

	return  fromTime, toTime
}

func (s *Server) Run(addr string)  error {

	router := gin.Default() 

	s.registerRoutes(router)

	return router.Run(addr)
}

var validate = validator.New()

func (s *Server) registerRoutes(router *gin.Engine) {

	apiv1 := router.Group("/api/v1")
	authGroup := apiv1.Group("/auth")
	{
		authGroup.POST("/register", s.handleUserRegister)
		authGroup.POST("/login", s.handleUserLogin)
	}

	
	// For our Control Room
	protected := apiv1.Group("")
	protected.Use(s.authMiddleware())
	{
		// Projects
		protected.POST("/projects", s.handleCreateProject)
		protected.GET("/projects", s.handleListProjects)

		// Logs
		protected.GET("logs", s.handleGetLogs)
		protected.GET("logs/stats", s.handleGetStats)
		protected.GET("logs/summary", s.handleGetSummary)
		protected.GET("logs/ws", s.handleWsLogic)
	}
	
	// For the Loading Dock (Agent route)
	apiv1.POST("/logs", s.apiKeyAuthMiddleware() ,s.handleLogIngest)
}

func (s *Server) handleLogIngest(c *gin.Context) {
	var logEntry database.LogEntry

	if err:= c.BindJSON(&logEntry); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return 
	};

	projectID, exists := c.Get("projectID")
	if !exists {
		c.JSON(500, gin.H{"error": "unauthorized context"})
		return
	}

	logEntry.ProjectID = projectID.(int)

	fmt.Println("Logs reach here? ")

	s.logQueue <- logEntry

	c.JSON(202, gin.H{"message": "log received!"})
}


func (s *Server) handleGetLogs(c *gin.Context)  {

	userID, exists := c.Get("userID")
	if !exists{
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
	isOwner, err := database.CheckProjectIDOwners(c.Request.Context(), s.db, userID.(int), projectID)
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
	logs, err := database.GetLogs(ctx, s.db, projectID ,limit, offset, serachQuery)
	if err != nil {
		fmt.Printf("Failed to get logs: %v\n", err)

		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"logs": logs})
}

func (s *Server) handleWsLogic(c *gin.Context) {
	hub.ServeWs(s.hub, c.Writer, c.Request)
}

func (s *Server) handleUserRegister(c *gin.Context) {

	var req struct {
		Name string `json:"name" validate:"required,min=2"`
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := auth.HashPasswod(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to has password"})
		return 
	}

	newUserId, err := database.CreateUser(c.Request.Context(), s.db, req.Name, req.Email, hash)
	if err != nil {
		if errors.Is(err, database.EmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "this email is already registered"})
			return
		}

		fmt.Printf("Internal error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	
	tokenString, err := auth.CreateJWT(s.jwtSecret, newUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return 
	}
	
	c.JSON(http.StatusCreated, gin.H{"token": tokenString})
}

func (s *Server) handleUserLogin(c *gin.Context) {

	var req struct {
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "invalid request body"})
		return
	}

		if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByEmail(c.Request.Context(), s.db, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := auth.ComparePasswordHash(req.Password, user.PasswordHash); !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return 
	} 

	tokenString, err := auth.CreateJWT(s.jwtSecret, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}
	
c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (s *Server) handleCreateProject(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Printf("The user id ------------: %v\n", userID.(int))

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
	}

	projectID, err := database.CreateProject(c.Request.Context(), s.db, userID.(int), req.Name, apiKey, secretHash)
	if err != nil {
		fmt.Printf("Internal Error ------------ : %v\n", err)
		c.JSON(500, gin.H{"error":  "failed to create project"})
		return
	}

	c.JSON(201, gin.H{
		"message": "project created",
		"project_id": projectID,
		"api_key": apiKey,
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

	isOwner, err := database.CheckProjectIDOwners(c.Request.Context(), s.db, userID.(int), projectID) 
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
	bucketMap  := map[string]string {
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

	c.JSON(200, gin.H{"stats": stats})
}


func (s *Server) handleGetSummary(c *gin.Context) {

	userID, _ := c.Get("userID")
	projectIDStr := c.Query("project_id")
	projectID, _ := strconv.Atoi(projectIDStr)

	isOwner, _ := database.CheckProjectIDOwners(c.Request.Context(), s.db, userID.(int), projectID)
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

	c.JSON(200, gin.H{"summary": summary})
}