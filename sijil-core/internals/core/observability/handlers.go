package observability

import (
	"encoding/json"
	"net/http"
	"sijil-core/internals/core/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// POST /api/v1/logs (API Key Auth)
func (h *Handler) Ingest(c *gin.Context) {
	projectID := c.GetInt("projectID")

	decoder := json.NewDecoder(c.Request.Body)

	token, err := decoder.Token()
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	// If it's an array '[', we loop until we hit ']'
	if delim, ok := token.(json.Delim); ok && delim == '[' {
		for decoder.More() {

			var log LogEntry

			if err := decoder.Decode(&log); err != nil {
				c.JSON(400, gin.H{"error": "bad log entry"})
				return
			}

			h.service.ProcessAndQueue(c.Request.Context(), projectID, &log)
		}

		// Consuming the closing ']'
		decoder.Token()
	} else {
		// Handle case when user sends a single JSON object '{...}'
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "accepted"})
}

// GET /api/v1/logs?project_id=1&q=error
func (h *Handler) Search(c *gin.Context) {
	userID := c.GetInt("userID") // Set by JWT Middleware
	projectID, _ := strconv.Atoi(c.Query("project_id"))

	plan := c.MustGet("plan").(*domain.Plan)

	if projectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id required"})
		return
	}

	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.service.Search(c.Request.Context(), userID, projectID, query, limit, offset, plan.RetentionDays)
	if err != nil {
		if err == ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
