package server

import (
	"fmt"
	"log-engine/internals/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	db *pgx.Conn
}

func NewServer (db *pgx.Conn) *Server {
	return &Server{db: db}
}

func (s *Server) Run(addr string)  error {

	router := gin.Default() 

	s.registerRoutes(router)

	return router.Run(addr)
}

func (s *Server) registerRoutes(router *gin.Engine) {

	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST("/logs", s.handleLogIngest)
		apiv1.GET("/logs", s.handleGetLogs)
	}
}

func (s *Server) handleLogIngest(c *gin.Context) {
	var logEntry database.LogEntry

	if err:= c.BindJSON(&logEntry); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return 
	}

	ctx := c.Request.Context()

	err := database.InsertLog(ctx, s.db, logEntry)
	if err != nil {
		fmt.Printf("Failed to insert log: %v\n",  err)

		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "log received!"})
}

func (s *Server) handleGetLogs(c *gin.Context)  {

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

	ctx := c.Request.Context()
	logs, err := database.GetLogs(ctx, s.db, limit, offset)
	if err != nil {
		fmt.Printf("Failed to get logs: %v\n", err)

		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"logs": logs})
}
