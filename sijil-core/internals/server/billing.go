package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
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

type LemonSqueezyEvent struct {
	Meta struct {
		EventName string `json:"event_name"`
		CustomData struct {
			UserID int `json:"user_id"`
		} `json:"custom_data"`
	} `json:"meta"`
	Data struct {
		Attributes struct {
			VariantID int `json:"variant_id"`
			Status    string `json:"status"`
			Total     int    `json:"total"`
		} `json:"attributes"`
	} `json:"data"`
}

func (s *Server) handleLemonSqueezyWebhook(c *gin.Context) {
	// 1. Verify Signature
	secret := os.Getenv("LEMONSQUEEZY_WEBHOOK_SECRET")
	if secret == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(payload)
	expectedSig := hex.EncodeToString(hash.Sum(nil))
	incomingSig := c.GetHeader("X-Signature")

	if !hmac.Equal([]byte(expectedSig), []byte(incomingSig)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 2. Parse Event
	var event LemonSqueezyEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 3. Process Event
	if event.Meta.EventName == "order_created" || event.Meta.EventName == "subscription_created" || event.Meta.EventName == "subscription_payment_success" {
		// Only process paid status
		if event.Data.Attributes.Status != "paid" && event.Data.Attributes.Status != "active" {
			c.Status(200)
			return
		}

		variantID := event.Data.Attributes.VariantID

		var planID int
		var duration time.Duration

		// Map Variant IDs (from Env) to Plan IDs
		proVariantID, _ := strconv.Atoi(os.Getenv("LEMON_VARIANT_PRO"))
		ultraVariantID, _ := strconv.Atoi(os.Getenv("LEMON_VARIANT_ULTRA"))

		// For simplicity, assuming monthly. Can be improved by checking variant/product logic.
		duration = 30 * 24 * time.Hour

		switch variantID {
		case proVariantID:
			planID = 2
		case ultraVariantID:
			planID = 3
		}

		if planID > 1 && event.Meta.CustomData.UserID > 0 {
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
