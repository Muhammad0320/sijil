package projects

import (
	"errors"
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

func getPlan(c *gin.Context) *domain.Plan {
	planRaw, exists := c.Get("plan")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "plan not found in context"})
		return nil
	}

	plan := planRaw.(*domain.Plan)

	return plan
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetInt("userID")
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := getPlan(c)

	resp, err := h.service.CreateProject(c.Request.Context(), userID, req, plan)
	if err != nil {
		if errors.Is(err, ErrLimitReached) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrProjectNameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetInt("userID")
	projects, err := h.service.ListProjects(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *Handler) AddMember(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID, _ := strconv.Atoi(c.Param("id"))

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := getPlan(c)

	err := h.service.AddMember(c.Request.Context(), userID, projectID, req, plan)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrLimitReached):
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add member"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member invited"})
}

func (h *Handler) GetMembers(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID, _ := strconv.Atoi(c.Param("id"))

	members, err := h.service.GetMembers(c.Request.Context(), userID, projectID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"members": members})
}
