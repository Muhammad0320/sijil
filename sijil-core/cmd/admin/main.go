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

}

func printHelp() {
	fmt.Println("Sijil Admin CLI tool")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}
