package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// LogEntry struct now lives here, as it defines our database model.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
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
		service     VARCHAR(100),
		search_vector TSVECTOR
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

	  createFunctionSQL := `
	  			CREATE OR REPLACE FUNCTIONS update_log_search_vector()
				RETURN TRIGGER AS $$
				BEGIN
	  				-- Combine level, service, and message into one text block
					-- and convert it into a tsvector
					NEW.search_vector = to_tsvector('simple', 
	  					COALESCE(NEW.level, '') || ' ' ||
						COALESCE(NEW.service, '') || ' ' ||
						COALESCE(NEW.message, '') 
					
					);
					RETURN NEW;
	  			END;
				$$ LANGUAGE plpgsql
	  `
	  _, err = db.Exec(ctx, createFunctionSQL)
	  if err != nil {
		return  fmt.Errorf("failed to create trigger function: %w", err)
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

func GetLogs(ctx context.Context, db *pgx.Conn, limit int, offset int) ([]LogEntry, error) { 

	getSQL := `SELECT timestamp, level, message, service
	 FROM logs
	 ORDER BY timestamp DESC 
	LIMIT $1
	OFFSET $2;
	`;

	rows, err := db.Query(ctx, getSQL, limit, offset)
	if err != nil {
	return nil,	fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close();
 
	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		err := rows.Scan(&log.Timestamp, &log.Level, &log.Message, &log.Service)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}