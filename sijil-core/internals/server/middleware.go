package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

func (s *Server) rateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	// We store the limiters in a map, keyed by IP
	var limiters = make(map[string]*rate.Limiter)
	var mu sync.RWMutex

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}

// Control Room Guard
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth header is required!"})
			return
		}

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 || parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(s.jwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token" + err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "incalid user ID in token"})
			return
		}

		userID := int(userIDFloat)

		plan, err := s.identityRepo.GetPlanByUserID(c.Request.Context(), userID)
		if err != nil {
			fmt.Println("This is the stupid error ----------------", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch plan"})
			return
		}

		c.Set("userID", userID)
		c.Set("plan", plan)

		c.Next()
	}
}

// The Loading Dock guard
func (s *Server) apiKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		apiKey := c.GetHeader("X-Api-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-Api-Key header is required"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization error format"})
			return
		}
		apiSecret := parts[1]

		projectID, valid := s.authCache.ValidateAPIKey(c.Request.Context(), apiKey, apiSecret)
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api credentials"})
			return
		}

		c.Set("projectID", projectID)
		c.Next()
	}
}

func (s *Server) SecurityHeadersMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Prevent clicking jacking
		c.Header("X-Frame-Options", "Deny")
		// Stop browserd from guessing content type
		c.Header("X-Content-Type-Options", "nosniff")
		// Enable XSS protection filter in older browsers
		c.Header("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (Strict transport Securoty) - Important for Production
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}
