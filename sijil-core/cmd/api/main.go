package main

// Fundamental rule in Go: "Never let one I/O op block another"
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sijil-core/internals/auth"
	"sijil-core/internals/database"
	"sijil-core/internals/hub"
	"sijil-core/internals/identity"
	"sijil-core/internals/ingest"
	"sijil-core/internals/projects"
	"sijil-core/internals/server"
	"sijil-core/internals/shared"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environent variables")
	}

	// 1. Set up the "root" context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Connect to the database
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("FATAL: DB_PASSWORD environment variable is not set ")
	}

	connString := fmt.Sprintf("postgres://postgres:%s@127.0.0.1:5434/log_db?sslmode=disable", dbPassword)

	db, err := database.ConnectDB(ctx, connString)
	if err != nil {
		log.Fatalf("Fatal error: could not connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// 3. Set up the database schema
	if err := database.CreateSchema(ctx, db); err != nil {
		log.Fatalf("Fatal error: could not set up database schema: %v", err)
	}

	// -- Retention policy
	go database.RunRetentionPolicy(ctx, db)

	// -- WebSocket --
	h := hub.NewHub()
	go h.Run()
	fmt.Println("Websocket has started! -------")
	// -- End websocket

	// -- Auth Catch
	authCache := auth.NewAuthCache(db)

	// -- WAL
	wal, err := ingest.NewWal("wal_data")
	if err != nil {
		log.Fatal("Could not open wal %w", err)
	}
	defer wal.Close()

	// --- Recovery Logic
	fmt.Println("Checking was for unsaved logs")
	recoveredLogs, err := wal.Recover()
	if err != nil {
		log.Fatalf("Fatal: WAL recovery failed: %v", err)
	}

	if len(recoveredLogs) > 0 {
		fmt.Printf("⚠️ Found %d unsaved logs in WAL. Replaying... \n", len(recoveredLogs))

		// We use a temporary context for recovery
		recoverCtx := context.Background()

		rows := make([][]any, len(recoveredLogs))
		for i, log := range recoveredLogs {
			rows[i] = []any{
				log.Timestamp,
				log.Level,
				log.Message,
				log.Service,
				log.ProjectID,
				log.Data,
			}
		}

		_, err := db.CopyFrom(
			recoverCtx,
			pgx.Identifier{"logs"},
			[]string{"timestamp", "level", "message", "service", "project_id", "data"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			log.Fatalf("⚠️ Failed to save recovered logs: %v", err)
		}

		fmt.Println("Recover successful. clearing WAL...")

		if err := wal.Reset(); err != nil {
			log.Fatalf("Fatal: Failed to clear WAL: %v", err)
		}

		fmt.Println("✅ Recover complete. Wal reset to segment 1")
	} else {
		wal.Reset()
		fmt.Println("✅ Wal is empty. Clean startup.")
	}
	// -- End Recovery

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is not set")
	}

	mailer := func(email, subject, body string) error {
		fmt.Printf("[Real mock] To %s | Token %s\n", email, body)
		return nil
	}

	identityRepo := identity.NewRepository(db)
	identityService := identity.NewService(identityRepo, jwtSecret, mailer)
	identityHandler := identity.NewHandler(identityService)

	projectsRepo := projects.NewRepository(db)
	projectService := projects.NewService(projectsRepo, mailer)
	projectHandler := projects.NewHandler(projectService)

	handlers := shared.Handlers{
		Identity: identityHandler,
		Projects: projectHandler,
	}

	// -- Ingesting engine
	engine := ingest.NewIngestionEngine(db, wal, h)
	engine.Start(ctx)
	srv := server.NewServer(db, engine, h, authCache, jwtSecret, handlers)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      srv.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Println("🚀 High-Throughtput Log Engine running on port :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crashed: %v", err)
		}
	}()

	// Graceful shutdown Block: listens for Ctlr + C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	cancel()
	fmt.Println("Waiting for ingestion engine to drain...")
	engine.Shutdown()

	fmt.Println("Cleanup complete. Bye!")
}
