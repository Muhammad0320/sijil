package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type PayStackEvent struct {
	Event string `json:"event"`
	Data  struct {
		Reference string `json:"references"`
		Status    string `json:"status"`
		Amount    int    `json:"amount"`
		Customer  struct {
			Email string `json:"email"`
		} `json:"customer"`

		Metadata struct {
			UserID int `json:"user_id"`
		} `json:"metadate"`
	} `json:"data"`
}

func (s *Server) handlePayStackWebhook(c *gin.Context) {

	// Verify signature
	secret := os.Getenv("PAYSTACK_SECRET_KEY")
	payload, _ := io.ReadAll(c.Request.Body)

	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write(payload)
	expectedSig := hex.EncodeToString(hash.Sum(nil))
	incomingSig := c.GetHeader("X-Paystack-Signature")

	if expectedSig != incomingSig {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 2. Parse Events
	var event PayStackEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 3. Update the database
	if event.Event == "charge.success" {
		var planName string
		if event.Data.Amount == 12_500_00 {
			planName = "Pro"
		} else if event.Data.Amount == 95_000_00 {
			planName = "Ultra"
		}

		if planName != "" {
			err := s.identityService.UpgradePlan(c.Request.Context(), event.Data.Metadata.UserID, planName)
			if err != nil {
				c.Status(500)
				return
			}
		}

	}

	c.Status(200)
}
