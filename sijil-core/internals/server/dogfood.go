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
		return func(ctx *gin.Context) { ctx.Next() }
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

		if s.ingestEngine == nil {
			return
		}

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
		message := fmt.Sprintf("%s %s | %d | %v", c.Request.Method, c.Request.URL.Path, status, latency)

		logEntry := database.LogEntry{
			Timestamp: time.Now(),
			Level:     level,
			Service:   "sijil-api",
			Message:   message,
			ProjectID: internalPID,
			Data: map[string]interface{}{
				"latency_ns": latency.Nanoseconds(),
				"ip":         clientIP,
				"status":     status,
				"userAgent":  c.Request.UserAgent(),
			},
		}

		go func(entry database.LogEntry) {

			if s.ingestEngine == nil {
				return
			}

			if err := s.ingestEngine.Wal.WriteBatch([]database.LogEntry{entry}); err != nil {
				fmt.Printf("DogFood wal error: %v\n", err)
				return
			}

			s.ingestEngine.LogQueue <- []database.LogEntry{entry}

		}(logEntry)

	}
}
