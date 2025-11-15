package database

import (
	"context"
	"fmt"
	"strings"
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

	// Create user tables
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMPZ DEFAULT NOW()
	)
	`
	_, err = db.Exec(ctx, createUserTableSQL)
	if err != nil {
		return  fmt.Errorf("failed to create 'users' table :%w ", err)
	}

	createProjectTableSQL := `
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) UNIQUE NOT NULL,
			api_key VARCHAR(255) UNIQUE NOT NULL,
			api_secret_hashed VARCHAR(255) NOT NULL,
			created_at TIMESTAMPZ DEFAULT NOW()
		)
	`
	_, err = db.Exec(ctx, createProjectTableSQL)
	if err != nil {
		return fmt.Errorf("faild to create project: %w", err)
	}

	// Create the main logs table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		timestamp   TIMESTAMPTZ       NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
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
	createHypertableSQL := `SELECT create_hypertable('logs', 'timestamp',
	partitioning_column => 'project_id',
	number_partitions => 10, 
	chunk_time_interval => INTERVAL '1 day', 
	if_not_exists => TRUE);`
	_, err = db.Exec(ctx, createHypertableSQL)
		if err != nil {
		return fmt.Errorf("failed to create hypertable: %w", err)
	  }

	  createFunctionSQL := `
	  			CREATE OR REPLACE FUNCTION update_log_search_vector()
				RETURNS TRIGGER AS $$
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
	  _, err = db.Exec(ctx, createFunctionSQL);
	  if err != nil {
		return  fmt.Errorf("failed to create trigger function: %w", err)
	  }

	  createTriggerSQL := `
	  		DROP TRIGGER IF EXISTS ts_vector_update ON logs; 
			CREATE TRIGGER ts_vector_update
			BEFORE INSERT ON logs
			FOR EACH ROW
	  		EXECUTE FUNCTION update_log_search_vector();
	  `

	  _, err = db.Exec(ctx, createTriggerSQL)
	  if err != nil {
		return fmt.Errorf("failed to create trigger : %w", err)
	  }
	
	fmt.Println("Database FTS is ready!")  
	
	fmt.Println("Database schema is ready!")
	return nil
}

// InsertLog writes a new LogEntry to the database.
func InsertLog(ctx context.Context, db *pgx.Conn, log LogEntry) error {
	// I wish there's a ternary expression equivalent in Go
	var logTime = log.Timestamp
	if logTime.IsZero() {
		logTime = time.Now()
	}
	
	insertSQL := `
		INSERT INTO logs (timestamp, level, message, service) 
		VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(ctx, insertSQL,
		logTime,
		log.Level,
		log.Message,
		log.Service,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}
	
	return nil
}

func GetLogs(ctx context.Context, db *pgx.Conn, limit int, offset int, searchQuery string) ([]LogEntry, error) { 

	args := make([]interface{}, 0)
	argsCounter := 1

	
	
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT timestamp, level, message, service 
		FROM logs
	`)

	// Conditionally add the WHERE clause for search
	if searchQuery != "" {
		// This is the FTS part
		queryBuilder.WriteString(fmt.Sprintf(" WHERE search_vector @@ plainto_tsquery('simple', $%d)", argsCounter))
		args = append(args, searchQuery)
		argsCounter++
	}

	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", argsCounter, argsCounter+1))
	args = append(args, limit, offset)

	getSQL := queryBuilder.String()
	fmt.Println("Final Query:", getSQL)

	rows, err := db.Query(ctx, getSQL, args...)
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