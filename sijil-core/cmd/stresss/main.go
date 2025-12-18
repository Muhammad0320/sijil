package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	URL         = "http://localhost:8080/api/v1/logs"
	Concurrency = 20    // 20 parallel connections
	TotalLogs   = 10000 // Total logs to send
	BatchSize   = 50    // Logs per HTTP request
	APIKey      = "pk_live_REPLACE_WITH_YOURS" // <--- REPLACE WITH YOURS
	APISecret   = "sk_live_REPLACE_WITH_YOURS"// <--- REPLACE WITH YOURS
)

func main() {
	var success int64
	var fail int64
	
	// Calculate how many HTTP requests we need to send
	// If we want 10,000 logs, and we send 50 at a time, we need 200 Requests.
	totalRequests := TotalLogs / BatchSize
	reqsPerWorker := totalRequests / Concurrency

	fmt.Printf("🔥 Stress Test: Sending %d logs in batches of %d.\n", TotalLogs, BatchSize)
	fmt.Printf("🔥 Total HTTP Requests: %d (%d per worker)\n", totalRequests, reqsPerWorker)

	// 1. PRE-GENERATE THE PAYLOAD
	// We do this once so we measure Network speed, not JSON marshalling speed.
	var batch []map[string]string
	for k := 0; k < BatchSize; k++ {
		batch = append(batch, map[string]string{
			"level":   "info",
			"message": fmt.Sprintf("stress test log %d", k),
			"service": "stress-bot",
		})
	}
	payload, _ := json.Marshal(batch)

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(Concurrency)

	for i := 0; i < Concurrency; i++ {
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 30 * time.Second}

			for j := 0; j < reqsPerWorker; j++ {
				req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Api-Key", APIKey)
				req.Header.Set("Authorization", "Bearer "+APISecret)
				
				resp, err := client.Do(req)
				if err != nil {
					atomic.AddInt64(&fail, 1)
					fmt.Printf("E") // Connection error
				} else {
					if resp.StatusCode > 299 {
						atomic.AddInt64(&fail, 1)
						fmt.Printf("S%d ", resp.StatusCode) // Status error
					} else {
						// Success! We sent 50 logs.
						atomic.AddInt64(&success, int64(BatchSize))
					}
					resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("\n\n--- RESULTS ---")
	fmt.Printf("Time: %v\n", duration)
	fmt.Printf("RPS:  %.2f logs/sec\n", float64(TotalLogs)/duration.Seconds())
	fmt.Printf("✅ Logs Ingested: %d\n", success)
	fmt.Printf("❌ Batches Failed: %d\n", fail)
}
