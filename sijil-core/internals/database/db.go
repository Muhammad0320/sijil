package database

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/goccy/go-json"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var EmailExists = errors.New("email already exists")
var NameExists = errors.New("project name already exists")

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	ProjectID int                    `json:"project_id"`
	Data      map[string]interface{} `json:"data,omitempty"`
	SegmentID int                    `json:"-"`
}

func (l *LogEntry) Serialize() ([]byte, error) {
	return json.Marshal(l)
}

func (l *LogEntry) Deserialize(data []byte) error {
	return json.Unmarshal(data, l)
}

// ConnectDB tries to connect to the database and returns the connection.
func ConnectDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	fmt.Println("Attempting to connect with connection string:", connString)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, config)
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

	createPlansTableSQL := `

	CREATE TABLE IF NOT EXISTS plans (
		id 	SERIAL PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		max_projects INT NOT NULL DEFAULT 3,
		max_members INT NOT NULL DEFAULT 1,
		retention_days INT NOT NULL DEFAULT 3,
		max_daily_logs BIGINT NOT NULL,
		price_usd DECIMAL(10, 2) NOT NULL,
		price_ngn DECIMAL(10, 2) NOT NULL
		);
	`

	_, err = db.Exec(ctx, createPlansTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create plans table %w\n", err)
	}

	seedQuery := `
    INSERT INTO plans (id, name, max_projects, max_members, max_daily_logs, retention_days, price_usd, price_ngn)
    VALUES 
        (1, 'Hobby', 1, 1, 10000, 3, 0.00, 0.00),          -- Free, 10k logs/day
        (2, 'Pro', 10, 10, 1000000, 14, 20.00, 12500.00),       -- $20, 1M logs/day (Solid for startups)
        (3, 'Ultra', -1, -1, 25000000, 30, 100.00, 95000.00)  -- $100, High volume
    ON CONFLICT (id) DO NOTHING;
    `

	if _, err := db.Exec(ctx, seedQuery); err != nil {
		return fmt.Errorf("failed to seed plans table %w \n", err)
	}

	// Create user tables
	createUserTableSQL := `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
	is_verified BOOLEAN DEFAULT FALSE,
	verification_token TEXT,
	verification_expires TIMESTAMP,
	token_expires_at TIMESTAMP,
	password_reset_token TEXT,
	password_reset_expires TIMESTAMP,
    avatar_url TEXT,
	plan_id INTEGER  NOT NULL REFERENCES plans(id),
	plan_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

`
	_, err = db.Exec(ctx, createUserTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create 'users' table :%w ", err)
	}

	createProjectTableSQL := `
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) UNIQUE NOT NULL,
			api_key VARCHAR(255) UNIQUE NOT NULL,
			api_secret_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(user_id, name)
		);
		CREATE INDEX IF NOT EXISTS idx_projects_api_key ON projects(api_key);
	`
	_, err = db.Exec(ctx, createProjectTableSQL)
	if err != nil {
		return fmt.Errorf("faild to create project: %w", err)
	}

	createMemebersTableSQL := `
		CREATE TABLE IF NOT EXISTS project_members (

			project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL DEFAULT 'viewer',
			joined_at TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (project_id, user_id)

		);
	`
	_, err = db.Exec(ctx, createMemebersTableSQL)
	if err != nil {
		return fmt.Errorf("failed to  create 'project_member' table: %w", err)
	}

	// Create the main logs table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		timestamp   TIMESTAMPTZ       NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
		level       VARCHAR(50)       NOT NULL,
		message     TEXT,
		service     VARCHAR(100),
		data 		JSONB		DEFAULT '{}'::jsonb,
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
	_, err = db.Exec(ctx, createFunctionSQL)
	if err != nil {
		return fmt.Errorf("failed to create trigger function: %w", err)
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

	createGinIndexSQL := `CREATE INDEX IF NOT EXISTS idx_logs_data ON logs USING GIN (data);`
	_, err = db.Exec(ctx, createGinIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create GIN index: %w", err)
	}

	fmt.Println("Database schema is ready!")
	return nil
}

type Project struct {
	ID            int
	UserID        int
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

	return project, nil
}

func CheckProjectIDOwners(ctx context.Context, db *pgxpool.Pool, userID, projectID int) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
	SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)
	`, projectID, userID).Scan(&exists)

	return exists, err
}

func CreateProject(ctx context.Context, db *pgxpool.Pool, userID int, name, apiKey, apiSecretHash string) (int, error) {
	var projectID int

	err := db.QueryRow(ctx, `
	INSERT INTO projects (user_id, name, api_key, api_secret_hash) VALUES ($1, $2, $3, $4) RETURNING id
	`, userID, name, apiKey, apiSecretHash).Scan(&projectID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, NameExists
			}
		}

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return projectID, nil
}

func GetUserProjects(ctx context.Context, db *pgxpool.Pool, userID int) ([]struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}, error) {

	rows, err := db.Query(ctx, `SELECT id, name FROM projects WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	for rows.Next() {

		var p struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		if err := rows.Scan(&p.ID, p.Name); err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

type LogStat struct {
	Bucket time.Time `json:"time"`
	Count  int       `json:"count"`
}

func GetLogStats(ctx context.Context, db *pgxpool.Pool, projectID int, fromTime, toTime time.Time, bucket string) ([]LogStat, error) {

	validBuckets := map[string]bool{
		"1 minutes": true, "5 minutes": true, "15 minutes": true, "30 minutes": true,
		"1 hour": true, "6 hours": true, "12 hours": true, "1 day": true,
	}

	if !validBuckets[bucket] {
		return nil, fmt.Errorf("Invalid bucket interval: %s\n", bucket)
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

	return stats, nil
}

type LogSummary struct {
	TotalLogs    int     `json:"total_logs"`
	ErrorCount   int     `json:"error_count"`
	ServiceCount int     `json:"service_count"`
	ErrorRate    float64 `json:"error_rate"`
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
		return summary, err
	}

	if summary.TotalLogs > 0 {
		summary.ErrorRate = (float64(summary.ErrorCount) / float64(summary.TotalLogs)) * 100
	} else {
		summary.ErrorRate = 0
	}

	return summary, nil
}

func RunRetentionPolicy(ctx context.Context, db *pgxpool.Pool) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("🧹 Starting Retention Policy Job")

			_, err := db.Exec(ctx, `
				DELETE FROM logs
				WHERE timestamp < NOW() - INTERVAL '3 days' 
				AND project_id IN (
					SELECT p.id FROM projects
					JOIN users u ON p.user_id = u.id
					WHERE u.plan_id = 1
				);
			`)

			if err != nil {
				fmt.Printf("⚠️ Failed to clean Free tier logs. %v\n", err)
			}

			_, err = db.Exec(ctx, `
				DELETE FROM logs
				WHERE timestamp < NOW() - INTERVAL '14 days' 
				AND project_id IN (
					SELECT p.id FROM projects
					JOIN users u ON p.user_id = u.id
					WHERE u.plan_id = 2
				);
			`)
			if err != nil {
				fmt.Printf("⚠️ Failed to clean Pro tier logs. %v\n", err)
			}

			// Ideally we should save into cold storage like S3 before deleting but this is a okay for v1
			_, err = db.Exec(ctx, `
				DELETE FROM logs
				WHERE timestamp < NOW() - INTERVAL '30 days' 
				AND project_id IN (
					SELECT p.id FROM projects
					JOIN users u ON p.user_id = u.id
					WHERE u.plan_id = 3
				);
			`)

			if err != nil {
				fmt.Printf("⚠️ Failed to clean Ultra  tier logs. %v\n", err)
			}

			// 4. Global safety Net. Just in case we missed something
			_, err = db.Exec(ctx, "SELECT drop_chunks('logs', INTERVAL '60 days');")
			if err != nil {
				fmt.Printf("⚠️ Gloabal aafety net failedL %v\n", err)
			}

			fmt.Println("✅ Retention Policy Job completed")

		case <-ctx.Done():
			return
		}
	}
}

func CheckProjectAccess(ctx context.Context, db *pgxpool.Pool, userID, projectID int) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 FROM projects WHERE id = $1 AND user_id = $2
			UNION 
			SELECT 1 FROM project_members WHERE project_id = $1 AND user_id = $2
		)
	`
	err := db.QueryRow(ctx, query, projectID, userID).Scan(&exists)

	return exists, err
}

func GetProjectRole(ctx context.Context, db *pgxpool.Pool, userID, projectID int) (string, error) {
	// 1. Check if Owner
	var isOwner bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)", projectID, userID).Scan(&isOwner)

	if err != nil {
		// Log it or return it. Don't let it fall through silently!
		return "", fmt.Errorf("failed to check owner status: %w", err)
	}

	if isOwner {
		return "owner", nil
	}

	// 2. Check if Member
	var role string
	err = db.QueryRow(ctx, "SELECT role FROM project_members WHERE user_id = $1 AND project_id = $2", userID, projectID).Scan(&role)
	if err != nil {
		return "", err // Not found in either
	}
	return role, nil
}

func GetMemberCountByProjectID(ctx context.Context, db *pgxpool.Pool, projectID int) (int, error) {
	var count int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM project_members WHERE project_id = $1`, projectID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get memebers count: %w", err)
	}

	return count, err
}

type ProjectMember struct {
	UserID    int       `json:"user_id"`
	ProjectID int       `json:"project_id"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	// Avatar might be needed in the future
}
