package domain

type Plan struct {
	Name          string
	MaxDailyLogs  int
	RetentionDays int
	MaxMemebers   int
	MaxProjects   int
}
