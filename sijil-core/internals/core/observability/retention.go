package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunRetentionPolicy(ctx context.Context, db *pgxpool.Pool) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("🧹 Starting Retention Policy Job")

			// Free tier - Plan ID 1
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

			// Pro tier - Plan ID 2
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

			// Ultra tier - Plan ID 3
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
				fmt.Printf("⚠️ Global safety net failed: %v\n", err)
			}

			fmt.Println("✅ Retention Policy Job completed")

		case <-ctx.Done():
			return
		}
	}
}
