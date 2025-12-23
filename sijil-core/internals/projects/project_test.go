package projects_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"sijil-core/internals/identity"
	"sijil-core/internals/projects"
	"sijil-core/internals/server"
	"sijil-core/internals/shared"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func parseToken(body *bytes.Buffer) string {
	var resp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body.Bytes(), &resp)
	return resp.Token
}

// Setup Helper (Same as Identity)
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

func TestProjectFlow(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	ctx := context.Background()

	db.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	// 3. BOOTSTRAP THE APP (Wiring everything exactly like main.go)
	// Identity Domain
	idRepo := identity.NewRepository(db)
	mailer := func(to, sub, body string) error { return nil } // Silent mailer
	idService := identity.NewService(idRepo, "test-secret", mailer)
	idHandler := identity.NewHandler(idService)

	// Projects Domain
	projRepo := projects.NewRepository(db)
	projService := projects.NewService(projRepo, mailer) // Inject mailer here too
	projHandler := projects.NewHandler(projService)

	// Server
	handlers := shared.Handlers{
		Identity: idHandler,
		Projects: projHandler,
	}

	srv := server.NewServer(db, nil, nil, nil, "test-secret", handlers)

	// --- THE TEST SCENARIO ---

	// A. REGISTER OWNER
	regBody := `{"firstname":"Owner","lastname":"Boss","email":"owner@sijil.dev","password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(regBody))
	srv.Router.ServeHTTP(w, req)
	fmt.Println(w.Body.String(), "------------------------")
	assert.Equal(t, 201, w.Code)

	ownerToken := parseToken(w.Body)
	assert.NotEmpty(t, ownerToken, "Owner should receive a JWT")

	// B. CREATE PROJECT (Authenticated)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBufferString(`{"name":"Sijil Production"}`))
	req2.Header.Set("Authorization", "Bearer "+ownerToken)
	srv.Router.ServeHTTP(w2, req2)
	assert.Equal(t, 201, w2.Code)

	// Extract Project ID
	var projResp struct {
		ID int `json:"id"`
	}
	json.Unmarshal(w2.Body.Bytes(), &projResp)
	projectID := projResp.ID
	assert.NotZero(t, projectID)

	// C. REGISTER VIEWER
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(`{"firstname":"Spy","lastname":"Guy","email":"spy@sijil.dev","password":"password123"}`))
	srv.Router.ServeHTTP(w3, req3)
	viewerToken := parseToken(w3.Body)

	// D. VIEWER TRIES TO SEE MEMBERS (Should Fail)
	// This proves our privacy logic works in the real handler chain
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/projects/%d/members", projectID), nil)
	req4.Header.Set("Authorization", "Bearer "+viewerToken)
	srv.Router.ServeHTTP(w4, req4)

	// Assert 403 Forbidden
	assert.Equal(t, 403, w4.Code, "Viewer should be forbidden from listing members")
}
