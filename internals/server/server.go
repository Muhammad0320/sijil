package server

import (
	"fmt"
	"log-engine/internals/auth"
	"log-engine/internals/database"
	"log-engine/internals/hub"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	db *pgx.Conn
	logQueue chan <- database.LogEntry
	hub *hub.Hub
	jwtSecret string
}

func NewServer (db *pgx.Conn, logQueue chan <- database.LogEntry, hub *hub.Hub, jwtSecret string) *Server {
	return &Server{
		db: db,
		logQueue: logQueue,
		hub: hub,
		jwtSecret:  jwtSecret,
	}
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

	// For the Loading Dock
	apiv1.POST("/logs", s.apiKeyAuthMiddleware() ,s.handleLogIngest)

	// For our Control Room
	logGroups := apiv1.Group("/logs")	
	logGroups.Use(s.authMiddleware())
	{
		logGroups.GET("", s.handleGetLogs)
		logGroups.GET("/ws", s.handleWsLogic)
	}
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
	isOwner, err := database.CheckProjectIDOwners(c, s.db, userID.(int), projectID)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
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