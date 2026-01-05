package domain

type Plan struct {
	ID            int
	Name          string
	MaxDailyLogs  int
	RetentionDays int
	MaxMembers    int
	MaxProjects   int
}
