package observability

import (
	"net/http"
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

	var logs []LogEntry
	if err := c.ShouldBindJSON(&logs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log format"})
		return
	}

	if err := h.service.Ingest(c.Request.Context(), projectID, logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ingestion failed"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "accepted"})
}

// GET /api/v1/logs?project_id=1&q=error
func (h *Handler) Search(c *gin.Context) {
	userID := c.GetInt("userID") // Set by JWT Middleware
	projectID, _ := strconv.Atoi(c.Query("project_id"))

	if projectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id required"})
		return
	}

	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.service.Search(c.Request.Context(), userID, projectID, query, limit, offset)
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
