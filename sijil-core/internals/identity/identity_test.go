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
	"github.com/stretchr/testify/assert"
)

// setupTestDB creates a temporary connection.
// Ideally, this uses a separate TEST_DB or cleans up tables.
func setupTestDB() *pgxpool.Pool {
	connString := os.Getenv("TEST_DB_URL") // e.g. postgres://user:pass@localhost:5432/sijil_test
	if connString == "" {
		// Fallback for local dev speed
		connString = "postgres://postgres:password@localhost:5433/log_db?sslmode=disable"
	}

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

	mailer := func(email, subject, body string) error {
		fmt.Printf("[Real mock] To %s | Token %s\n", email, body)
		return nil
	}

	svc := identity.NewService(repo, "test-secret", mailer)
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
	t.Log("STEP 1: Registering User...")
	regBody := `{"firstname":"Ali","lastname":"Test","email":"test@sijil.dev","password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(regBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// B. GRAB TOKEN FROM DB (The "Cheat" that Postman can't do)
	var verificationHash string
	err := db.QueryRow(context.Background(),
		"SELECT verification_token FROM users WHERE email='test@sijil.dev'").Scan(&verificationHash)
	assert.NoError(t, err)

	// Wait! We stored the HASH in the DB, but we need the RAW token to verify.
	// Oh snap. We designed it securely so even the DB admin can't verify users.
	// This means our TEST cannot complete the flow unless we:
	// 1. Mock the Random Generator (Hard)
	// 2. Use a "Test Mode" in Service (Messy)
	// 3. Return the token in the API response ONLY for dev/test builds.

	// FIX: For this test, we will Manually Update the user to verify them,
	// verifying the 'Login' check works.
	// OR: We intercept the logic.
}
