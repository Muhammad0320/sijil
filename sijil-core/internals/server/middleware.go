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


type IPRateLimiter struct {

	ips map[string]*rate.Limiter
	mu *sync.RWMutex
	r rate.Limit
	b int

}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {

	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu: &sync.RWMutex{},
		r: r,
		b: b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {

	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists :=  i.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func (s *Server) rateLimitMiddleware(limit *IPRateLimiter) gin.HandlerFunc {

	return func (c *gin.Context) {

		limiter := limit.GetLimiter(c.ClientIP())
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

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error)  {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return  nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
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
		c.Set("userID", userID)

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

