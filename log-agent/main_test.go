package main

import (
	"testing"
	"time"
)

func TestLogParsing(t *testing.T) {
	const timeLayout = "2006-01-02 15:04:05"
	defaultService := "log-agent-v1"

	testCases := []struct {
		name          string
		logLine       string
		expectedLog   Log
		expectMatch   bool
		expectedError bool
	}{
		{
			name:    "Valid log line with service",
			logLine: "2025-11-14 10:00:00 [my-service] [info] This is a test message",
			expectedLog: Log{
				Timestamp: mustParseTime(timeLayout, "2025-11-14 10:00:00"),
				Level:     "info",
				Service:   "my-service",
				Message:   "This is a test message",
			},
			expectMatch: true,
		},
		{
			name:    "Valid log line without service",
			logLine: "2025-11-14 10:00:00 [debug] Another test message",
			expectedLog: Log{
				Timestamp: mustParseTime(timeLayout, "2025-11-14 10:00:00"),
				Level:     "debug",
				Service:   defaultService,
				Message:   "Another test message",
			},
			expectMatch: true,
		},
		{
			name:    "Line that does not match pattern",
			logLine: "this is just a plain log message",
			expectedLog: Log{
				Level:   "info",
				Service: defaultService,
				Message: "this is just a plain log message",
			},
			expectMatch: false,
		},
		{
			name:    "Log line with different log level",
			logLine: "2025-11-14 10:00:00 [my-app] [error] Something went wrong",
			expectedLog: Log{
				Timestamp: mustParseTime(timeLayout, "2025-11-14 10:00:00"),
				Level:     "error",
				Service:   "my-app",
				Message:   "Something went wrong",
			},
			expectMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newLog := Log{}
			matches := logRegex.FindStringSubmatch(tc.logLine)

			if (matches != nil) != tc.expectMatch {
				t.Fatalf("Expected match: %v, but got: %v", tc.expectMatch, matches != nil)
			}

			if matches == nil {
				newLog.Level = "info"
				newLog.Message = tc.logLine
				newLog.Service = defaultService
			} else {
				parsedTime, err := time.Parse(timeLayout, matches[1])
				if err == nil {
					newLog.Timestamp = parsedTime
				}

				if matches[2] == "" {
					newLog.Service = defaultService
				} else {
					newLog.Service = matches[2]
				}

				newLog.Level = matches[3]
				newLog.Message = matches[4]
			}

			if !tc.expectedLog.Timestamp.IsZero() && !newLog.Timestamp.Equal(tc.expectedLog.Timestamp) {
				t.Errorf("Timestamp mismatch: expected %v, got %v", tc.expectedLog.Timestamp, newLog.Timestamp)
			}
			if newLog.Level != tc.expectedLog.Level {
				t.Errorf("Level mismatch: expected %s, got %s", tc.expectedLog.Level, newLog.Level)
			}
			if newLog.Service != tc.expectedLog.Service {
				t.Errorf("Service mismatch: expected %s, got %s", tc.expectedLog.Service, newLog.Service)
			}
			if newLog.Message != tc.expectedLog.Message {
				t.Errorf("Message mismatch: expected %s, got %s", tc.expectedLog.Message, newLog.Message)
			}
		})
	}
}

// mustParseTime is a helper to panic on time parsing errors in test setup.
func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
