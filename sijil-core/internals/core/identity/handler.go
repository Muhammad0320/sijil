package identity

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
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

	fmt.Println("Did you reach here? -------------------")

	token, user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	_ = h.service.ForgotPassword(c.Request.Context(), req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "If account exists, a reset link has been sent."})
}

func (h *Handler) ResetPassword(c *gin.Context) {

	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *Handler) handleUploadAvatar(c *gin.Context) {

	userID := c.GetInt("userID")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid image format. Use PNG  or JPEG"})
		return
	}

	newImage := resize.Thumbnail(256, 256, img, resize.Lanczos3)

	// Save to disk
	err = os.MkdirAll("data/avatars", 0755)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create avatar directory"})
		return
	}

	filename := fmt.Sprintf("user-%d.jpg", userID)
	outPath := filepath.Join("data/avatars", filename)

	out, err := os.Create(outPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save image"})
		return
	}
	defer out.Close()

	jpeg.Encode(out, newImage, &jpeg.Options{Quality: 75})

	avatarUrl := fmt.Sprintf("/static/avatars/%s", filename)

	if err = h.service.repo.UpdateUserAvatar(c.Request.Context(), userID, avatarUrl); err != nil {
		c.JSON(500, gin.H{"error": "Database update failed"})
		return
	}

	c.JSON(200, gin.H{"data": avatarUrl})
}
