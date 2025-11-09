package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// LogEntry struct now lives here, as it defines our database model.
type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Service string `json:"service"`
}

// ConnectDB tries to connect to the database and returns the connection.
func ConnectDB(ctx context.Context, connString string) (*pgx.Conn, error) {
	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return db, nil
}

// CreateSchema sets up the logs table and hypertable.
func CreateSchema(ctx context.Context, db *pgx.Conn) error {
	// Enable the TimescaleDB extension
	enableExtensionSQL := `CREATE EXTENSION IF NOT EXISTS timescaledb;`
	_, err := db.Exec(ctx, enableExtensionSQL)
	if err != nil {
		return fmt.Errorf("failed to enable timescaledb extension: %w", err)
	}

	// Create the main logs table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		timestamp   TIMESTAMPTZ       NOT NULL,
		level       VARCHAR(50)       NOT NULL,
		message     TEXT,
		service     VARCHAR(100)
	);`

	_, err = db.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create 'logs' table: %w", err)
	}

	// Turn it into a Hypertable
	createHypertableSQL := `SELECT create_hypertable('logs', 'timestamp', if_not_exists => TRUE);`
	_, err = db.Exec(ctx, createHypertableSQL)
	
		if err != nil {
		return fmt.Errorf("failed to create hypertable: %w", err)
	}
	
	fmt.Println("Database schema is ready!")
	return nil
}

// InsertLog writes a new LogEntry to the database.
func InsertLog(ctx context.Context, db *pgx.Conn, log LogEntry) error {
	insertSQL := `
		INSERT INTO logs (timestamp, level, message, service) 
		VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(ctx, insertSQL,
		time.Now(),
		log.Level,
		log.Message,
		log.Service,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}
	
	return nil
}

func GetLogs(ctx context.Context, db *pgx.Conn) ([]LogEntry, error) { 

	getSQL := `SELECT * FROM logs`

	rows, err := db.Query(ctx, getSQL)
	if err != nil {
	return nil,	fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close();

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		err := rows.Scan(&log.Level, &log.Message, &log.Service)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}