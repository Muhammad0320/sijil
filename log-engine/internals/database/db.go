package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var EmailExists = errors.New("email already exists")

// LogEntry struct now lives here, as it defines our database model.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level   string `json:"level"`
	Message string `json:"message"`
	Service string `json:"service"`
	ProjectID int `json:"-"`
}


// ConnectDB tries to connect to the database and returns the connection.
func ConnectDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	fmt.Println("Attempting to connect with connection string:", connString)
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return db, nil
}

// CreateSchema sets up the logs table and hypertable.
func CreateSchema(ctx context.Context, db *pgxpool.Pool) error {
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
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
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
			api_secret_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
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
func InsertLog(ctx context.Context, db *pgxpool.Pool, log LogEntry) error {
	// I wish there's a ternary expression equivalent in Go
	var logTime = log.Timestamp
	if logTime.IsZero() {
		logTime = time.Now()
	}

	fmt.Printf("The log to insert ------------ : %v \n", log)

	insertSQL := `
		INSERT INTO logs (timestamp, level, message, service, project_id) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(ctx, insertSQL,
		logTime,
		log.Level,
		log.Message,
		log.Service,
		log.ProjectID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}
	
	return nil
}

func GetLogs(ctx context.Context, db *pgxpool.Pool,  projectID ,limit, offset int, searchQuery string) ([]LogEntry, error) { 

	args := make([]interface{}, 0)
	argsCounter := 1

	var queryBuilder strings.Builder
	queryBuilder.WriteString(fmt.Sprintf(`SELECT timestamp, level, message, service 
		FROM logs WHERE project_id = $%d`, argsCounter))
	args = append(args, projectID)
	argsCounter++

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

type User struct {
	ID int
	PasswordHash string
}

func CreateUser(ctx context.Context, db *pgxpool.Pool, name, email, hashpassword string) (int, error) {

	var newUserID int 
	err := db.QueryRow(ctx, 
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		name, email, hashpassword,
	).Scan(&newUserID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, EmailExists
			}
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return newUserID, nil 
}

func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (User, error) {

	var user User
	err := db.QueryRow(ctx, 
	`SELECT id, password_hash FROM users WHERE email = $1`,
	email).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return user, fmt.Errorf("failed to get user: %w", err)
	}

	return  user, nil 
}

type Project struct {
	ID int 
	UserID int 
	ApiSecretHash string 
}

func GetProductByApiKey(ctx context.Context, db *pgxpool.Pool, apiKey string) (Project, error) {
	var project Project	
	err := db.QueryRow(ctx, `
	SELECT id, user_id, api_secret_hash FROM projects WHERE api_key = $1
	`, apiKey).Scan(&project.ID, &project.UserID, &project.ApiSecretHash)
	if err != nil {
		return project, fmt.Errorf("failed to gett project: %w", err)
	}

	return  project, nil 
}

func CheckProjectIDOwners(ctx context.Context, db *pgxpool.Pool, userID, projectID int) (bool, error) {
	var exists bool 
	err := db.QueryRow(ctx, `
	SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)
	`, projectID, userID).Scan(&exists) 

	return  exists, err
}

func CreateProject(ctx context.Context, db *pgxpool.Pool, userID int, name, apiKey, apiSecretHash string) (int, error) {
	var projectID int 

	err := db.QueryRow(ctx, `
	INSERT INTO projects (user_id, name, api_key, api_secret_hash) VALUES ($1, $2, $3, $4) RETURNING id
	`, userID, name, apiKey, apiSecretHash).Scan(&projectID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return projectID, nil 	
}

func GetUserProjects(ctx context.Context, db *pgxpool.Pool, userID int) ( []struct{
	ID int 		`json:"id"`
	Name string `json:"name"`
} , error) {

	rows, err := db.Query(ctx, `SELECT id, name FROM projects WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []struct{
	ID int 		`json:"id"`
	Name string `json:"name"`
	}

	for rows.Next() {

		var p struct{
	ID int 		`json:"id"`
	Name string `json:"name"`
	}
	
	if err := rows.Scan(&p.ID, p.Name); err != nil {
		return nil, err
	}

	projects = append(projects, p)
	}

	return  projects, nil 
}

type LogStat struct { 
	Bucket time.Time `json:"time"`
	Count int `json:"count"`
}

func GetLogStats(ctx context.Context, db *pgxpool.Pool, projectID int, fromTime, toTime time.Time, bucket string) ([]LogStat, error) {

	validBuckets := map[string]bool {
		"1 minutes": true, "5 minutes": true, "15 minutes": true, "30 minutes": true,
		"1 hour": true, "6 hours": true, "12 hours": true, "1 day": true,
	}

	if !validBuckets[bucket] {
		return  nil, fmt.Errorf("Invalid bucket interval: %s\n", bucket)
	}

	query := fmt.Sprintf(`
		SELECT time_bucket('%s', timestamp) AS bucket, COUNT(*)
		FROM logs 
		WHERE project_id = $1
		  AND  timestamp >= $2
		  AND  timestamp <= $3
		GROUP BY bucket
		ORDER BY bucket ASC;
	`, bucket)
	rows, err := db.Query(ctx, query, projectID, fromTime, toTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []LogStat
	for rows.Next() {
		var s LogStat

		if err := rows.Scan(&s.Bucket, &s.Count); err != nil {
			return nil, err
		}

		stats = append(stats, s)
	}

	return  stats, nil 
} 

type LogSummary struct {
	TotalLogs int `json:"total_logs"`
	ErrorCount int `json:"error_count"`
	ServiceCount int `json:"service_count"`
	ErrorRate float64 `json:"error_rate"`
}

func GetlogSummary(ctx context.Context, db *pgxpool.Pool, projectID int, fromTime, toTime time.Time) (LogSummary, error) {
	var summary LogSummary
	
	query := `
		WITH subset AS (
		SELECT level, service
		FROM logs
		WHERE project_id = $1
		  AND timestamp >= $2
		  AND timestamp <= $3
	) 
	SELECT 
		COUNT(*) as total,
		COUNT(*) FILTER (WHERE level ILIKE 'error' OR level ILIKE 'critical') as errors,
		COUNT(DISTICT service) as services
	FROM subset;
	`

	err := db.QueryRow(ctx, query, projectID, fromTime, toTime).Scan(
		&summary.TotalLogs, &summary.ErrorCount, &summary.ServiceCount,
	)

	if err != nil {
		return  summary, err
	}

	if summary.TotalLogs > 0 {
		summary.ErrorRate = ( float64(summary.ErrorCount) / float64(summary.TotalLogs)) * 100
	} else {
		summary.ErrorRate = 0
	}

	return summary, nil 
}