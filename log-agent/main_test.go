package main

import (
	"reflect"
	"testing"
	"time"
)

func TestLogParsing(t *testing.T) {

	type expectedResult struct {
		Timestamp time.Time 
		Service string
		Level string
		Message string
	}

	testCases := []struct{
		name string
		logLine string
		shouldMatch bool 
		expexted expectedResult
	}{
		{
			name: "Full log with service",
			logLine: "2025-01-10 09:15:00 [user-api] [INFO] User 'admin' logged in.",
			shouldMatch: true,
			expexted: expectedResult{
				Timestamp: time.Date(2025, 1, 10, 9, 15, 0, 0, time.UTC),
				Service: "user-api",
				Level: "INFO",
				Message: "User 'admin' logged in.",
			},
		},

		{

			name: "Log without a service",
			logLine: "2025-01-10 09:17:20 [ERROR] Failed to process payment: 12345",
			shouldMatch: true,
			expexted: expectedResult{
				Timestamp: time.Date(2025, 1, 10, 9, 17, 20, 0 ,time.UTC),
				Service: "",
				Level: "ERROR",
				Message: "Failed to process payment: 12345",
			},
		}, 
		{

			name: "Non-matching line",
			logLine: "This log line should match out regex pattern",
			shouldMatch: false,
			expexted: expectedResult{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			matches := logRegex.FindStringSubmatch(tc.logLine)

			// -- Did it match when it should? --
			if matches == nil && tc.shouldMatch {
				t.Fatalf("Regex failed to match a line when it should have: %q", tc.logLine)
			}

			if matches != nil && !tc.shouldMatch {
				t.Fatalf("Regex matched a line when it should not have: %q", tc.logLine)
			}

			if !tc.shouldMatch {
				return
			}

			// -- Check the extracted values --
			parsedTime, err := time.ParseInLocation(timeLayout, matches[1], time.UTC)
			if err != nil {
				t.Fatalf("Failed to parse timestamp: %v", err)
			}

			got := expectedResult{
				Timestamp: parsedTime,
				Service: matches[2],
				Level: matches[3],
				Message: matches[4],
			}

			if !reflect.DeepEqual(got, tc.expexted) {
				t.Errorf("Parsed log does not match expected." + 
				"\n Got:  %+v" + 
				"\n Wanted: %+v", got, tc.expexted)
			}
		})	
	}	
}