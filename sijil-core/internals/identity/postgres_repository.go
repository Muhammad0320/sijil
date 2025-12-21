package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, u *User) (int, error) {
	var newUserID int
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (firstname, lastname, email, password_hash, verification_token, verification_token_expires_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		u.FirstName, u.LastName, u.Email, u.PasswordHash, u.VerificationToken, u.VerificationTokenExpires,
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
		`SELECT id, firstname, lastname, email, password_hash, plan, is_verified 
         FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.Plan, &u.IsVerified)

	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil

}
