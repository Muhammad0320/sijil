package identity

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {

	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) VerifyEmail(c *gin.Context) {

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})

}

func (h *Handler) ForgotPassword(c *gin.Context) {

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	_ = h.service.ForgotPassword(c.Request.Context(), req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "If account exists, a reset link has been sent."})
}
