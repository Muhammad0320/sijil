package identity_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sijil-core/internals/identity"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// setupTestDB creates a temporary connection.
// Ideally, this uses a separate TEST_DB or cleans up tables.
func setupTestDB() *pgxpool.Pool {
	_ = godotenv.Load("../../.env")

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "logpassword123"
	}
	connString := fmt.Sprintf("postgres://postgres:%s@127.0.0.1:5434/log_db?sslmode=disable", dbPassword)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	return db
}

// THE TEST
func TestIdentityFlow(t *testing.T) {
	// 1. Setup
	db := setupTestDB()
	defer db.Close()

	// Clean DB before test (Ruthless!)
	_, _ = db.Exec(context.Background(), "DELETE FROM users WHERE email = 'test@sijil.dev'")

	// 2. Wiring (Manually, like main.go)
	repo := identity.NewRepository(db)

	var capturedToken string
	mockMailer := func(email, subject, body string) error {
		capturedToken = body
		return nil
	}

	svc := identity.NewService(repo, "test-secret", mockMailer)
	handler := identity.NewHandler(svc)

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	authGroup := r.Group("/auth")
	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
	authGroup.GET("/verify", handler.VerifyEmail)

	// --- SCENARIO START ---

	// A. REGISTER
	t.Log("STEP 1: Registering...")
	regBody := `{"firstname":"Ali","lastname":"Test","email":"test@sijil.dev","password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(regBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	// B. VERIFY (Using captured token!)
	t.Log("STEP 2: Verifying with token:", capturedToken)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/auth/verify?token="+capturedToken, nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)

	// C. LOGIN
	t.Log("STEP 3: Logging in...")
	loginBody := `{"email":"test@sijil.dev","password":"password123"}`
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(loginBody))
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
}
