package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type LogEntry struct {
	
	Message string `json:"message"`
	Level string `json:"level"`
	Service string `json:"service"`

}

var db *pgx.Conn

func createSchema(ctx context.Context, db *pgx.Conn) error {

	// Enable the timescaledb extension
	_, err := db.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS timescaledb;")
	if err != nil {
		return fmt.Errorf("failed to create timescaledb extension: %w", err)
	}

	createTableSchema := `
		CREATE TABLE IF NOT EXISTS logs (
			timestamp TIMESTAMPTZ  NOT NULL,
			level VARCHAR(50)  NOT NULL,
			message TEXT,
			service VARCHAR(100)
		);
	`

	_, err = db.Exec(ctx, createTableSchema)

	if err != nil {
		return fmt.Errorf("failed to create 'logs' table: %w", err)

	}

	createHypertableSQL := `
	SELECT create_hypertable('logs', 'timestamp', if_not_exists => TRUE);`

	_, err = db.Exec(ctx, createHypertableSQL)
	if err != nil {
		return fmt.Errorf("failed to create hypertable: %w", err)
	}

	fmt.Println("Database schema is ready.")
	return nil
}


func main() {
	ctx := context.Background()
	// ----- DATABASE CONNECTION ------ 
	// "postgres://[USER]:[PASSWORD]@[HOST]:[PORT]/[DB_NAME]"
	connString := "postgres://postgres:logpassword123@localhost:5433/log_db?sslmode=disable"
	
	var err error
	db, err = pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err) 
	}
	defer db.Close(ctx)
	fmt.Println("Successfully connected to db!")

	if err := createSchema(ctx, db);err != nil {
		log.Fatalf("failed to set up db schema: %v\n", err)
	}

	router := gin.Default()

	router.POST("/api/v1/log", func(c *gin.Context) {

		var log LogEntry

		if err := c.BindJSON(&log); err != nil {
			c.JSON(400, gin.H{
				"message": "Bad Request",
			})
			return
		}

		fmt.Printf("LOG RECEIVED: %+v\n", log)

		ctx := c.Request.Context()

		insertSql := `
		INSERT INTO logs (timestamp, level, message, service) VALUES ($1, $2, $3, $4);
		`
		_, err := db.Exec(ctx, insertSql, time.Now(), log.Level, log.Message, log.Service)

		if err != nil {
			fmt.Printf("failed to insert to db: %v\n", err)

			c.JSON(500, gin.H{"error": "internal server error"})
		}

		c.JSON(200, gin.H{
			 "message": "Log received!",
		})
	})
	router.Run(":8080")
}

