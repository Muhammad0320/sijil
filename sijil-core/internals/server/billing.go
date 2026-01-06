package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type PayStackEvent struct {
	Event string `json:"event"`
	Data  struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
		Amount    int    `json:"amount"`
		Customer  struct {
			Email string `json:"email"`
		} `json:"customer"`

		Metadata struct {
			UserID int `json:"user_id"`
		} `json:"metadata"`
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
		amount := event.Data.Amount

		var planID int
		var duration time.Duration

		switch {
		case amount == 12_500_00:
			planID = 2
			duration = 30 * 24 * time.Hour
		case amount == 125_000_00: // 2 months free
			planID = 2
			duration = 365 * 24 * time.Hour
		case amount == 95_000_00:
			planID = 3
			duration = 30 * 24 * time.Hour
		case amount == 950_000_00: // 2 months free
			planID = 3
			duration = 365 * 24 * time.Hour
		}

		if planID > 1 {
			expiry := time.Now().Add(duration)

			err := s.identityRepo.UpdateUserPlan(c.Request.Context(), event.Data.Metadata.UserID, planID, expiry)
			if err != nil {
				c.Status(500)
				return
			}
		}

	}

	c.Status(200)
}
