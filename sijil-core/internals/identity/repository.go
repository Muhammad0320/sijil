package identity

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, u *User) (int, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetPlanByUserID(ctx context.Context, id int) (*Plan, error)
	VerifyUserAccount(ctx context.Context, token string) error
	SetPasswordResetToken(ctx context.Context, email string, token string, expiry time.Time) error
	ResetPasswordByToken(ctx context.Context, token, passwordHash string) error
}
