package server_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sijil-core/internals/auth"
	"sijil-core/internals/core/identity"
	"sijil-core/internals/core/observability"
	"sijil-core/internals/core/projects"
	"sijil-core/internals/database"
	"sijil-core/internals/hub"
	"sijil-core/internals/ingest"
	"sijil-core/internals/server"
	"sijil-core/internals/shared"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	testDB *pgxpool.Pool
	router *gin.Engine
)

func setupTestDB() {
	if os.Getenv("SKIP_DB") == "true" {
		return
	}

	godotenv.Load("../../../.env")

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	connString := fmt.Sprintf("postgres://postgres:%s@127.0.0.1:5434/log_db?sslmode=disable", dbPassword)

	var err error
	testDB, err = database.ConnectDB(context.Background(), connString)
	if err != nil {
		fmt.Printf("Skipping tests requiring DB: %v\n", err)
		return
	}
}

func setupServer() {
	if testDB == nil {
		return
	}

	jwtSecret := "testsecret"

	h := hub.NewHub()
	go h.Run()

	authCache := auth.NewAuthCache(testDB)

	os.MkdirAll("test_wal_data", 0755)
	wal, _ := ingest.NewWal("test_wal_data")
	engine := ingest.NewIngestionEngine(testDB, wal, h)
	ctx, _ := context.WithCancel(context.Background())
	go engine.Start(ctx)

	mailerFunc := func(email, subject, body string) error {
		fmt.Printf("Mock Email to %s: %s\n", email, subject)
		return nil
	}

	identityRepo := identity.NewRepository(testDB)
	identityService := identity.NewService(identityRepo, jwtSecret, mailerFunc)
	identityHandler := identity.NewHandler(identityService)

	projectsRepo := projects.NewRepository(testDB)
	projectService := projects.NewService(projectsRepo, mailerFunc)
	projectHandler := projects.NewHandler(projectService)

	observabilityRepo := observability.NewRepository(testDB)
	observabilityService := observability.NewService(observabilityRepo, projectsRepo, engine)
	observabilityHandler := observability.NewHandler(observabilityService)

	handlers := shared.Handlers{
		IdentityRepo:    identityRepo,
		IdentityService: identityService,
		Identity:        identityHandler,
		Projects:        projectHandler,
		Observability:   observabilityHandler,
	}

	srv := server.NewServer(testDB, engine, h, authCache, jwtSecret, handlers)
	router = srv.Router
}

func TestFullFlow(t *testing.T) {
	setupTestDB()
	if testDB == nil {
		t.Skip("Skipping integration test: DB not available")
	}
	setupServer()

	// 1. Register User
	email := fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano())
	password := "password123"

	registerBody := map[string]string{
		"firstname": "Test",
		"lastname":  "User",
		"email":      email,
		"password":   password,
	}
	body, _ := json.Marshal(registerBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var registerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResp)
	token := registerResp["token"].(string)
	user := registerResp["user"].(map[string]interface{})
	userID := int(user["id"].(float64))

	// 2. Login
	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Create Project
	createProjectBody := map[string]string{
		"name": fmt.Sprintf("Test Project %d", time.Now().UnixNano()),
	}
	body, _ = json.Marshal(createProjectBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var projectResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &projectResp)
	projectID := int(projectResp["id"].(float64))
	apiKey := projectResp["api_key"].(string)
	apiSecret := projectResp["api_secret"].(string)

	// 4. Ingest Logs
	logEntries := []map[string]interface{}{
		{
			"timestamp": time.Now(),
			"level":     "info",
			"message":   "Test log message",
			"service":   "test-service",
			"data":      map[string]string{"foo": "bar"},
		},
	}
	body, _ = json.Marshal(logEntries)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/logs", bytes.NewBuffer(body))
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiSecret)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	time.Sleep(1 * time.Second)

	// 5. Search Logs
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/logs?project_id=%d&q=Test", projectID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 6. Simulate Paystack Payment (Upgrade to Plan 2)
	paystackPayload := map[string]interface{}{
		"event": "charge.success",
		"data": map[string]interface{}{
			"amount": 1250000,
			"reference": "ref_123",
			"status": "success",
			"customer": map[string]string{"email": email},
			"metadata": map[string]interface{}{"user_id": userID},
		},
	}
	payloadBytes, _ := json.Marshal(paystackPayload)

	secret := os.Getenv("PAYSTACK_SECRET_KEY")
	if secret == "" {
		secret = "test_secret"
		os.Setenv("PAYSTACK_SECRET_KEY", secret)
	}

	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write(payloadBytes)
	signature := hex.EncodeToString(hash.Sum(nil))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/webhooks/paystack", bytes.NewBuffer(payloadBytes))
	req.Header.Set("X-Paystack-Signature", signature)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Simulate Lemon Squeezy Webhook (Upgrade to Plan 3)
	lemonPayload := map[string]interface{}{
		"meta": map[string]interface{}{
			"event_name": "order_created",
			"custom_data": map[string]interface{}{
				"user_id": userID,
			},
		},
		"data": map[string]interface{}{
			"attributes": map[string]interface{}{
				"variant_id": 456, // defined in .env as LEMON_VARIANT_ULTRA
				"status": "paid",
				"total": 10000,
			},
		},
	}
	lemonBytes, _ := json.Marshal(lemonPayload)

	lemonSecret := os.Getenv("LEMONSQUEEZY_WEBHOOK_SECRET")
	if lemonSecret == "" {
		lemonSecret = "test_secret"
		os.Setenv("LEMONSQUEEZY_WEBHOOK_SECRET", lemonSecret)
	}

	hash = hmac.New(sha256.New, []byte(lemonSecret))
	hash.Write(lemonBytes)
	lemonSig := hex.EncodeToString(hash.Sum(nil))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/webhooks/lemonsqueezy", bytes.NewBuffer(lemonBytes))
	req.Header.Set("X-Signature", lemonSig)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
