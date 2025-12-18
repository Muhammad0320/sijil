package server

import (
	"fmt"
	"os"
	"sijil-core/internals/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) DogFoodMiddleware() gin.HandlerFunc {

	pidStr := os.Getenv("INTERNAL_LOG_PROJECT_ID")
	if pidStr == "" {
		fmt.Println("⚠️ INTERNAL_LOG_PROJECT_ID not set. Dog Fooding disabled")
		return func(ctx *gin.Context) {ctx.Next()}
	}
	internalPID, _ := strconv.Atoi(pidStr)



	return func(c *gin.Context) {
		start := time.Now() 
		path := c.Request.URL.Path
		method := c.Request.Method

		if method == "POST" && path == "/api/v1/logs" {
			c.Next()
			return 
		}

		c.Next()

		if internalPID == 0 {
			return 
		}

		latency := time.Since(start)
		clientIP := c.ClientIP()
		status := c.Writer.Status()

		// Smart level detection
		level := "info"
		if status >= 500 {
			level = "error"
		} else if status >= 400 {
			level = "warn"
		} 

		// Rich message 
		message := fmt.Sprintf("[%d] %s %s | %v | IP: %s | UA: %s", 
			status, method, path, latency, clientIP, c.Request.UserAgent(),
		)	

		logEntry := database.LogEntry{
			Timestamp: time.Now(),
			Level: level,
			Service: "sijil-core-internal",
			Message: message,
			ProjectID: internalPID,
		}
		
		if err := s.ingestEngine.Wal.WriteLog(logEntry); err == nil {
			s.ingestEngine.LogQueue <- logEntry
		}
	}
}