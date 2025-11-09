package server

import (
	"fmt"
	"log-engine/internals/database"

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
		apiv1.POST("/log", s.handleLogIngest)
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