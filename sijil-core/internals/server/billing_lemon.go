package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type LemonEvent struct {
	Meta struct {
		EventName  string `json:"event_name"`
		CustomData struct {
			UserID int `json:"user_id,string"`
		} `json:"custom_data"`
	} `json:"meta"`

	Data struct {
		Attributes struct {
			Status    string `json:"status"`
			Total     int    `json:"total"`
			UserEmail string `json:"user_email"`
		} `json:"attributes"`
	} `json:"data"`
}

func (s *Server) handleLemonWebhook(c *gin.Context) {

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	secret := os.Getenv("LEMONSQUEEZY_WEBHOOK_SECRET")
	signature := c.GetHeader("X-Signature")

	if !checkLemonSignature(payload, secret, signature) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Parse event
	var event LemonEvent
	if err = json.Unmarshal(payload, &event); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if event.Meta.EventName == "subscription_created" || event.Meta.EventName == "order_created" {

		var planID int
		var duration time.Duration
		amount := event.Data.Attributes.Total

		switch {
		// PRO
		case amount == 20_00:
			planID = 2
			duration = 30 * 24 * time.Hour
		case amount == 200_00: // 2 months free
			planID = 2
			duration = 365 * 24 * time.Hour
		// Ultra
		case amount == 100_00:
			planID = 3
			duration = 30 * 24 * time.Hour
		case amount == 1_000_00: // 2 months free
			planID = 3
			duration = 365 * 24 * time.Hour

		}

		if planID > 1 {
			expiry := time.Now().Add(duration)
			err := s.identityRepo.UpdateUserPlan(c.Request.Context(), event.Meta.CustomData.UserID, planID, expiry)
			if err != nil {
				c.Status(500)
				return
			}
		}

	}
	c.Status(200)

}

func checkLemonSignature(payload []byte, secret string, signature string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
