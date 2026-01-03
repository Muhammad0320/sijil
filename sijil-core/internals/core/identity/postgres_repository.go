package identity

import (
	"context"
	"errors"
	"fmt"
	"sijil-core/internals/core/domain"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetPlanByUserID(ctx context.Context, id int) (*domain.Plan, error) {

	user, err := r.GetByID(ctx, id)
	if err != nil {
		return &domain.Plan{}, fmt.Errorf("faild to get plan %s", err)
	}

	var p domain.Plan
	err = r.db.QueryRow(ctx, `
		SELECT id, name, max_projects, max_members, max_daily_logs, retention_days
 	FROM plans 
	WHERE id = $1
	`, user.PlanID).Scan(&p.ID, &p.Name, &p.MaxProjects, &p.MaxMemebers, &p.MaxDailyLogs, &p.RetentionDays)

	if err != nil {
		return &domain.Plan{}, fmt.Errorf("faild to get plan %s", err)
	}

	return &p, nil
}

func (r *postgresRepository) GetPlanByName(ctx context.Context, name string) (*domain.Plan, error) {
	var p domain.Plan
	err := r.db.QueryRow(ctx, `
		SELECT id, name, max_projects, max_members, max_daily_logs, retention_days
	FROM plans
	WHERE name ILIKE $1
	`, name).Scan(&p.ID, &p.Name, &p.MaxProjects, &p.MaxMemebers, &p.MaxDailyLogs, &p.RetentionDays)

	if err != nil {
		return &domain.Plan{}, fmt.Errorf("failed to get plan by name %s: %w", name, err)
	}

	return &p, nil
}

func (r *postgresRepository) UpdateUserPlan(ctx context.Context, userID, planID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET plan_id = $1
		WHERE id = $2
	`, planID, userID)
	return err
}

func (r *postgresRepository) Create(ctx context.Context, u *User) (int, error) {
	var newUserID int
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (firstname, lastname, email, password_hash ,verification_token, verification_expires, plan_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		u.FirstName, u.LastName, u.Email, u.PasswordHash, u.VerificationToken, u.VerificationTokenExpires, 1,
	).Scan(&newUserID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, errors.New("email already exists")
			}
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return newUserID, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User

	err := r.db.QueryRow(ctx,
		`SELECT id, firstname, lastname, email, password_hash, avatar_url ,plan_id, is_verified 
         FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.PlanID, &u.IsVerified)

	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil

}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*User, error) {

	var u User

	err := r.db.QueryRow(ctx,
		`SELECT id, firstname, lastname, email, password_hash, avatar_url, plan_id, is_verified 
         FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.PlanID, &u.IsVerified)

	if err != nil {
		return &u, errors.New("user not found")
	}

	return &u, nil

}

func (r *postgresRepository) VerifyUserAccount(ctx context.Context, token string) error {
	commonTag, err := r.db.Exec(ctx, `
	UPDATE users 
	SET is_verified = TRUE, 
		verification_token = NULL,
		verification_expires = NULL
	WHERE verification_token = $1
	AND verification_expires > NOW()`, token)

	if err != nil {
		return fmt.Errorf("db verification error: %w", err)
	}

	if commonTag.RowsAffected() == 0 {
		return errors.New("invalid or expired verification token")
	}

	return nil
}

func (r *postgresRepository) SetPasswordResetToken(ctx context.Context, email string, token string, expiry time.Time) error {

	_, err := r.db.Exec(ctx, `
	UPDATE users
	SET password_reset_token = $1, password_reset_expired = $2
	WHERE email = $3
	`, token, expiry, email)

	return err

}

func (r *postgresRepository) ResetPasswordByToken(ctx context.Context, token, passwordHash string) error {

	commandTag, err := r.db.Exec(ctx, `
		UPDATE users
		SET password_hash = $1, password_reset_token = NULL, password_reset_expires = NULL 
		WHERE password_reset_token = $2 AND password_reset_expires > NOW()
	`, passwordHash, token)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("invalid or expired token")
	}

	return nil
}

func (r *postgresRepository) UpdateUserAvatar(ctx context.Context, userID int, avatarUrl string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET avatar_url = $1 WHERE id = $2`, avatarUrl, userID)

	return err
}
