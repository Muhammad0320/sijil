package server

import (
	"errors"
	"expvar"
	"fmt"
	"log"
	"log-engine/internals/auth"
	"log-engine/internals/database"
	"log-engine/internals/hub"
	"log-engine/internals/ingest"
	"log-engine/internals/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db *pgxpool.Pool
	ingestEngine *ingest.IngestionEngine
	hub *hub.Hub
	authCache *auth.AuthCache
	jwtSecret string
	// -----
	Router *gin.Engine
}

func NewServer (db *pgxpool.Pool, ingestEngine *ingest.IngestionEngine, hub *hub.Hub, authCache *auth.AuthCache , jwtSecret string) *Server {
	s := &Server{
		db: db,
		ingestEngine: ingestEngine,
		hub: hub,
		authCache: authCache,
		jwtSecret:  jwtSecret,
	}
	
	router := gin.Default() 
	s.registerRoutes(router)
	s.Router = router

	return s
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

var validate = validator.New()

func (s *Server) registerRoutes(router *gin.Engine) {
	router.Use(s.SecurityHeadersMiddleware()) 
	router.Use(s.DogFoodMiddleware())

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
		protected.POST("/projects/:id/members", s.handleAddMember)
		protected.GET("/projects/:id/members", s.handleGetMembers)

		// Logs
		protected.GET("logs", s.handleGetLogs)
		protected.GET("logs/stats", s.handleGetStats)
		protected.GET("logs/summary", s.handleGetSummary)
		protected.GET("logs/ws", s.handleWsLogic)
	}
	
	// For the Loading Dock (Agent route)
	apiv1.POST("/logs", s.apiKeyAuthMiddleware(), s.handleLogIngest)

	adminGroup := router.Group("/debug", gin.BasicAuth(gin.Accounts{
		os.Getenv("ADMIN_USER"): os.Getenv("ADMIN_PASS"),
	}))
	adminGroup.GET("/vars", gin.WrapH(expvar.Handler()))
}

func (s *Server) handleLogIngest(c *gin.Context) {
	ProjectID := c.GetInt("projectID")


	var logEntries []database.LogEntry
	if err:= c.BindJSON(&logEntries); err != nil {
		c.JSON(400, gin.H{"error": "bad request; expected an array"})
		return 
	};

	// Enrichments
	for i := range logEntries {
	 logEntries[i].ProjectID = ProjectID
	 if logEntries[i].Timestamp.IsZero() {
		logEntries[i].Timestamp = time.Now()
	 }
	}

	//  WAL (DURABILITY)
	 if err := s.ingestEngine.Wal.WriteBatch(logEntries); err != nil {
		ingest.RecordError()
		ingest.RecordDropped(1) 
		c.JSON(500, gin.H{"error": "durability failure"})
		return 
	 }

	 for _, entry := range logEntries {
		 s.ingestEngine.LogQueue <- entry 
		 ingest.RecordQueued(1)
	}
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
   plan, _ := database.GetUserPlan(ctx, s.db, userID.(int))

  limits := database.GetPlanLimits(plan) 
  rentention := limits.RetentionDays

	logs, err := database.GetLogs(ctx, s.db, projectID ,limit, offset, serachQuery, rentention)
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

func (s *Server) handleUserRegister(c *gin.Context) {

	var req struct {
		FirstName string `json:"firstname" validate:"required,min=2"`
		LastName string `json:"lastname" validate:"required,min=2"`
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

	newUserId, err := database.CreateUser(c.Request.Context(), s.db, req.FirstName, req.LastName, req.Email, hash)
	if err != nil {
		if errors.Is(err, database.EmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "409 Confilict: this email is already registered"})
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

	plan, _ := database.GetUserPlan(c.Request.Context(), s.db, userID.(int))

	// Check plan limit
	count, _ := database.GetProjectCountByUserID(c.Request.Context(), s.db, userID.(int))
	
	if count >= database.GetPlanLimits(plan).MaxProject {
		// I don't know the right status code to send
		c.JSON(400, gin.H{"error": fmt.Sprintf("You've reached you limit of %d members! kindly unpgrade your plan.", count)})
		return
	}


	projectID, err := database.CreateProject(c.Request.Context(), s.db, userID.(int), req.Name, apiKey, secretHash)
	if err != nil {
		if errors.Is(err, database.NameExists){
			c.JSON(409, gin.H{"error": "409 Conflict: project name already exists"})
			return
		}

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

	c.JSON(200, gin.H{"summary": summary})
}

func (s *Server) handleAddMember(c *gin.Context) {

	userID := c.GetInt("userID")
	projectIDStr := c.Param("id")
	projectID, _ := strconv.Atoi(projectIDStr)

	role, err:= database.GetProjectRole(c.Request.Context(), s.db, userID, projectID)
	if err != nil {
		log.Fatalf("Cannot get role: %s", err)
		return
	}

	if role != "owner" && role != "admin" {
		c.JSON(403, gin.H{"error": "only admins can invite members"})
		return
	}

	plan, _ := database.GetUserPlan(c.Request.Context(), s.db, userID)

	// Check plan limit
	count, _ := database.GetMemberCountByProjectID(c.Request.Context(), s.db, projectID) 
	
	if count >= database.GetPlanLimits(plan).MaxMemebers  {
		// I don't know the right status code to send
		c.JSON(402, gin.H{"error": fmt.Sprintf("You've reached you limit of %d members! kindly unpgrade your plan.", count)})
		return
	}

	// Add member
	var req struct { Email string `json:"email"`; Role string `json:"role"` }
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