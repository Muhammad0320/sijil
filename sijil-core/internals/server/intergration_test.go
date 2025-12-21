package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sijil-core/internals/auth"
	"sijil-core/internals/database"
	"sijil-core/internals/hub"
	"sijil-core/internals/identity"
	"sijil-core/internals/ingest"
	"testing"

	"github.com/joho/godotenv"
)

// Helper to setup a fresh server and clean DB for testing
func setupTestServer(t *testing.T) *Server {
	// 1. Load Env (or hardcode for test)
	_ = godotenv.Load("../../.env")

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "logpassword123"
	}
	connString := fmt.Sprintf("postgres://postgres:%s@127.0.0.1:5433/log_db?sslmode=disable", dbPassword)

	ctx := context.Background()
	db, err := database.ConnectDB(ctx, connString)
	if err != nil {
		t.Fatalf("Could not connect to DB: %v", err)
	}

	// 2. CLEAN THE DB (Ruthless Reset)
	// We truncate tables to ensure a clean slate for every test run
	_, err = db.Exec(ctx, "TRUNCATE TABLE users, projects, logs, project_members RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Could not truncate DB: %v", err)
	}

	// 3. Initialize dependencies
	h := hub.NewHub()
	go h.Run()

	jwtSecret := "test-secret-key"

	idRepo := identity.NewRepository(db)
	idService := identity.NewService(idRepo, jwtSecret)
	idHandler := identity.NewHandler(idService)

	// We don't need a real WAL for RBAC testing, but the engine needs one
	// Create a temp file
	tmpWal, _ := os.CreateTemp("", "test_wal")
	wal, _ := ingest.NewWal(tmpWal.Name())
	engine := ingest.NewIngestionEngine(db, wal, h)
	authCache := auth.NewAuthCache(db)

	return NewServer(db, engine, h, authCache, jwtSecret, idHandler)
}

// Helper to make requests
func makeRequest(s *Server, method, url, token string, body interface{}) *httptest.ResponseRecorder {
	var jsonBody []byte
	if body != nil {
		jsonBody, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {

		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	return w
}

func TestRBAC_EndToEnd(t *testing.T) {
	s := setupTestServer(t)
	// Ensure you have access to your DB pool here.
	// If 's' contains the DB, use s.DB. If it's a global variable, use that.
	db := s.db

	// --- 1. SETUP USERS ---
	// Register User A (Admin)
	w := makeRequest(s, "POST", "/api/v1/auth/register", "", map[string]string{
		"firstname": "Admin", "lastname": "User", "email": "admin@test.com", "password": "password123",
	})
	if w.Code != 201 {
		t.Fatalf("Failed to register Admin: %v", w.Body.String())
	}

	var adminAuth struct {
		Token string `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &adminAuth)
	if err != nil {
		t.Fatalf("error unmarshalling admin token: %s", err)
	}

	// ----------------------------------------------------------------------
	// [INSERTED FIX]: FORCE UPGRADE ADMIN TO PRO
	// We do this here so Step 4 doesn't fail on "Quota Exceeded" before hitting "Conflict"
	// ----------------------------------------------------------------------
	_, err = db.Exec(context.Background(), `
       UPDATE users
	   SET plan = $1
	   WHERE email = $2
    `, "pro", "admin@test.com")
	if err != nil {
		t.Fatalf("Failed to force upgrade admin user: %v", err)
	}
	// ----------------------------------------------------------------------

	// Register User B (Viewer)
	w = makeRequest(s, "POST", "/api/v1/auth/register", "", map[string]string{
		"firstname": "Viewer", "lastname": "User", "email": "viewer@test.com", "password": "password123",
	})

	if w.Code != 201 {
		t.Fatalf("Failed to register Viewer")
	}
	var viewerAuth struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &viewerAuth)

	// --- 2. CREATE PROJECT (As Admin) ---
	w = makeRequest(s, "POST", "/api/v1/projects", adminAuth.Token, map[string]string{
		"name": "Death Star Logs",
	})

	if w.Code != 201 {
		t.Fatalf("Failed to create project: %v", w.Body.String())
	}
	var projectRes struct {
		ProjectID int `json:"project_id"`
	}
	json.Unmarshal(w.Body.Bytes(), &projectRes)
	pID := projectRes.ProjectID

	// --- 3. ADD MEMBER (Success Case) ---
	// Admin invites Viewer.
	// Because we upgraded the plan earlier, this uses 1/100 slots. Safe.
	w = makeRequest(s, "POST", fmt.Sprintf("/api/v1/projects/%d/members", pID), adminAuth.Token, map[string]string{
		"email": "viewer@test.com",
		"role":  "viewer",
	})
	if w.Code != 200 {
		t.Errorf("Expected 200 OK for adding member, got %d. Body: %s", w.Code, w.Body.String())
	}

	// --- 4. ADD MEMBER (Conflict Case) ---
	// Admin invites Viewer AGAIN.
	// The quota check passes (1/100 used), so it hits the Duplicate Check logic -> 409.
	w = makeRequest(s, "POST", fmt.Sprintf("/api/v1/projects/%d/members", pID), adminAuth.Token, map[string]string{
		"email": "viewer@test.com",
		"role":  "viewer",
	})

	if w.Code != 409 {
		t.Errorf("Expected 409 Conflict for duplicate member, got %d. Body: %s", w.Code, w.Body.String())
	}

	// --- 5. UNAUTHORIZED INVITE (Forbidden Case) ---
	// Viewer tries to invite Admin (Circular, but testing permissions)
	w = makeRequest(s, "POST", fmt.Sprintf("/api/v1/projects/%d/members", pID), viewerAuth.Token, map[string]string{
		"email": "admin@test.com",
		"role":  "admin",
	})
	if w.Code != 403 {
		t.Errorf("Expected 403 Forbidden for viewer inviting others, got %d", w.Code)
	}

	// --- 6. VIEW LOGS (Access Granted Case) ---
	// Viewer tries to read logs
	w = makeRequest(s, "GET", fmt.Sprintf("/api/v1/logs?project_id=%d", pID), viewerAuth.Token, nil)
	if w.Code != 200 {
		t.Errorf("Expected 200 OK for viewer reading logs, got %d", w.Code)
	}

	// --- 7. VIEW LOGS (Access Denied Case) ---
	// Register Random User C
	w = makeRequest(s, "POST", "/api/v1/auth/register", "", map[string]string{
		"firstname": "Stranger", "lastname": "Danger", "email": "stranger@test.com", "password": "password123",
	})
	var strangerAuth struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &strangerAuth)

	// Stranger tries to read logs
	w = makeRequest(s, "GET", fmt.Sprintf("/api/v1/logs?project_id=%d", pID), strangerAuth.Token, nil)
	if w.Code != 403 {
		t.Errorf("Expected 403 Forbidden for stranger reading logs, got %d", w.Code)
	}
}
