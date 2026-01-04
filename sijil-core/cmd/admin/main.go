package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func mian() {

	_ = godotenv.Load()

	action := flag.String("action", "", "Actions to perform: promote | ban | reset-pass | stats")
	email := flag.String("email", "", "Target user email")
	planID := flag.Int("paln", 0, "Plan ID (1=Hobby, 2=Pro, 3=Ultra)")

	flag.Parse()

	if *action == "" {
		printHelp()
		return
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl != "" {

		dbPass := os.Getenv("DATABASE_PASSWORD")
		if dbPass != "" {
			dbUrl = fmt.Sprintf("postgres://postgres:%s@127.0.0.1:5434/log_db?sslmode=disable", dbPass)
		} else {
			log.Fatal("DATABASE_URL or DB_PASSWORD is required")
		}
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch *action {

	case "promote":
		if *email == "" || *planID == 0 {
			log.Fatal("Usage: -action=promote email=... plan_id=2")
		}

		promoteUser(ctx, db, *email, *planID)

	case "ban":
		if *email == "" {
			log.Fatal("Usage -action=ban email=...")
		}
		banUser(ctx, db, *email)
	case "stats":
		showStats(ctx, db)
	default:
		printHelp()

	}
}

func promoteUser(ctx context.Context, db *pgxpool.Pool, email string, planID int) {

	tag, err := db.Exec(ctx, "UPDATE users SET plan_id = $1 WHERE email = $2", planID, email)
	if err != nil {
		log.Fatalf("❌ Error: %v", err)
	}

	if tag.RowsAffected() == 0 {
		fmt.Println("⚠️ User not found.")
	} else {
		fmt.Printf("✅ User %s upgraded to Plan %d\n", email, planID)
	}
}

func banUser(ctx context.Context, db *pgxpool.Pool, email string) {

	_, err := db.Exec(ctx, "UPDATE users SET password_hash = 'BANNED' WHERE email = $1", email)
	if err != nil {
		log.Fatalf("❌ Error %v", err)
	}

	fmt.Printf("User %s has been banned (password scrambled) \n", email)

}

func showStats(ctx context.Context, db *pgxpool.Pool) {

	var userCount, logCount, projectCount int
	var dbSize string

	db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	db.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&projectCount)
	db.QueryRow(ctx, "SELECT COUNT(*) FROM logs").Scan(&logCount)
	db.QueryRow(ctx, "SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)

	fmt.Println("\n ------- 📊 SIJIL ADMIN STATS ------------")
	fmt.Printf("⚔️ Total users: %d\n", projectCount)
	fmt.Printf("👥 Total users: %d\n", userCount)
	fmt.Printf("📃 Total log: %d\n", logCount)
	fmt.Printf("💾 DB Size: %s\n", dbSize)
	fmt.Println("----------------------------------")

}

func printHelp() {
	fmt.Println("Sijil Admin CLI tool")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}
