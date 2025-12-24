package shared

import (
	"context"
	"fmt"
)

type Plan struct {
	Name          string
	MaxDailyLogs  int
	RetentionDays int
	MaxMemebers   int
	MaxProjects   int
}

type UserContextKey string

const PlanKey UserContextKey = "user_plan"

func GetPlanFromContext(ctx context.Context) (Plan, error) {

	plan, ok := ctx.Value(PlanKey).(Plan)
	if !ok {
		return Plan{}, fmt.Errorf("no plan found in this context")
	}

	return plan, nil
}
