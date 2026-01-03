package identity

import (
	"context"
	"sijil-core/internals/core/domain"
	"time"
)

type Repository interface {
	Create(ctx context.Context, u *User) (int, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetPlanByUserID(ctx context.Context, id int) (*domain.Plan, error)
	GetPlanByName(ctx context.Context, name string) (*domain.Plan, error)
	UpdateUserPlan(ctx context.Context, userID, planID int) error
	VerifyUserAccount(ctx context.Context, token string) error
	SetPasswordResetToken(ctx context.Context, email string, token string, expiry time.Time) error
	ResetPasswordByToken(ctx context.Context, token, passwordHash string) error
}
