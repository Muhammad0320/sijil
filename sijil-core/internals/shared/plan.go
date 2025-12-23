package shared

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserPlan(ctx context.Context, db *pgxpool.Pool, userID int) (string, error) {
	var plan string
	err := db.QueryRow(ctx, "SELECT plan FROM users WHERE id = $1", userID).Scan(&plan)
	return plan, err
}
