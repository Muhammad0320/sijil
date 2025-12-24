package observability

import "time"

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	ProjectID int                    `json:"project_id"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type LogStat struct {
	Bucket time.Time `json:"time"`
	Count  int       `json:"count"`
}

type LogSummary struct {
	TotalLogs    int     `json:"total_logs"`
	ErrorCount   int     `json:"error_count"`
	ServiceCount int     `json:"service_count"`
	ErrorRate    float64 `json:"error_rate"`
}

// IngestRequest allows single or batch ingestion
type IngestRequest []LogEntry
