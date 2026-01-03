package projects

import (
	"context"
	"errors"
	"fmt"
	"sijil-core/internals/core/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProjectNameExists = errors.New("project name already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrAlreadyMember     = errors.New("user is already a member")
)

type Repository interface {
	Create(ctx context.Context, project *Project) (int, error)
	GetByID(ctx context.Context, projectID int) (*Project, error)
	ListByUserID(ctx context.Context, userID int) ([]Project, error)

	// Membership
	GetRole(ctx context.Context, projectID, userID int) (string, error)
	AddMember(ctx context.Context, projectID int, email, role string) error
	ListMembers(ctx context.Context, projectID int) ([]ProjectMember, error)

	// Limits
	CountProjects(ctx context.Context, userID int) (int, error)
	CountMembers(ctx context.Context, projectID int) (int, error)
	GetUserPlan(ctx context.Context, userID int) (*domain.Plan, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, p *Project) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO projects (user_id, name, api_key, api_secret_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		p.UserID, p.Name, p.APIKey, p.APISecretHash,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrProjectNameExists
		}
		return 0, err
	}
	return id, nil
}

func (r *postgresRepository) ListByUserID(ctx context.Context, userID int) ([]Project, error) {
	// We only select fields safe for listing (usually hide keys)
	rows, err := r.db.Query(ctx, `SELECT id, name, created_at FROM projects WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, projectID int) (*Project, error) {
	var p Project
	err := r.db.QueryRow(ctx, `SELECT id, user_id, name, api_key, api_secret_hash FROM projects WHERE id=$1`, projectID).
		Scan(&p.ID, &p.UserID, &p.Name, &p.APIKey, &p.APISecretHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresRepository) GetRole(ctx context.Context, projectID, userID int) (string, error) {
	// 1. Check Owner
	var isOwner bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM projects WHERE id=$1 AND user_id=$2)", projectID, userID).Scan(&isOwner)
	if err != nil {
		return "", err
	}
	if isOwner {
		return "owner", nil
	}

	// 2. Check Member
	var role string
	err = r.db.QueryRow(ctx, "SELECT role FROM project_members WHERE project_id=$1 AND user_id=$2", projectID, userID).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil // No access
	}
	return role, err
}

func (r *postgresRepository) AddMember(ctx context.Context, projectID int, email, role string) error {
	// Subquery to find user ID by email
	tag, err := r.db.Exec(ctx, `
		INSERT INTO project_members (project_id, user_id, role)
		SELECT $1, id, $2 FROM users WHERE email = $3
		ON CONFLICT DO NOTHING`,
		projectID, role, email,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		// Either user doesn't exist OR they are already a member.
		// Let's distinguish for better error messages
		var exists bool
		_ = r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
		if !exists {
			return ErrUserNotFound
		}
		return ErrAlreadyMember
	}
	return nil
}

func (r *postgresRepository) ListMembers(ctx context.Context, projectID int) ([]ProjectMember, error) {
	// Union: Owner (implicitly admin) + Members
	query := `
		SELECT u.id, u.email, 'owner' as role, p.created_at as joined_at
		FROM projects p JOIN users u ON p.user_id = u.id WHERE p.id = $1
		UNION ALL
		SELECT u.id, u.email, pm.role, pm.joined_at
		FROM project_members pm JOIN users u ON pm.user_id = u.id WHERE pm.project_id = $1
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []ProjectMember
	for rows.Next() {
		var m ProjectMember
		m.ProjectID = projectID
		if err := rows.Scan(&m.UserID, &m.Email, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *postgresRepository) CountProjects(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM projects WHERE user_id=$1", userID).Scan(&count)
	return count, err
}

func (r *postgresRepository) CountMembers(ctx context.Context, projectID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM project_members WHERE project_id=$1", projectID).Scan(&count)
	return count, err
}

func (r *postgresRepository) GetUserPlan(ctx context.Context, userID int) (*domain.Plan, error) {
	var p domain.Plan
	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.name, p.max_projects, p.max_members, p.max_daily_logs, p.retention_days
		FROM users u
		JOIN plans p ON u.plan_id = p.id
		WHERE u.id = $1
	`, userID).Scan(&p.ID, &p.Name, &p.MaxProjects, &p.MaxMemebers, &p.MaxDailyLogs, &p.RetentionDays)

	if err != nil {
		return nil, fmt.Errorf("failed to get plan for user %d: %w", userID, err)
	}
	return &p, nil
}
