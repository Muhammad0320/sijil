package domain

type Plan struct {
	ID            int
	Name          string
	MaxDailyLogs  int
	RetentionDays int
	MaxMemebers   int
	MaxProjects   int
}
