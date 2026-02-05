package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/nxadm/tail"
)

type Log struct {
	Timestamp time.Time              `json:"timestamp,omitempty"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

var logRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(?:\[(.*?)\]\s+)?\[(.*?)\]\s+(.*)$`)

const timeLayout = "2006-01-02 15:04:05"

func main() {

	filePtr := flag.String("f", "test.log", "log file to tail")
	servicePtr := flag.String("s", "log-agent-v1", "service name to tag logs with")
	formatPtr := flag.String("format", "regex", "Log format: 'regex' or 'json' ")
	apiKeyPtr := flag.String("pk", "", "Public API key (pk_live_...)")
	secretKeyPtr := flag.String("sk", "", "Secret API key (sk_live_...)")
	urlPtr := flag.String("url", "https://api.sijil.dev/v1/ingest", "Sijil Ingest Endpoint")

	flag.Parse()

	if *apiKeyPtr == "" || *secretKeyPtr == "" {
		log.Fatal("Error: you must provide both pk and sk flags")
	}

	var parser Parser
	switch *formatPtr {
	case "regex":
		parser = NewRegexParser(*servicePtr)
	case "json":
		parser = NewJsonParser(*servicePtr)
	default:
		log.Fatalf("FATAL: Unknown format '%s'", *formatPtr)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("Starting agent on %s...\n", *filePtr)
	t, err := tail.TailFile(*filePtr, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatalf("Failed to tail file: %v", err)
	}

	// --- 	BATCHING CONFIG ----
	var batch []Log
	batchSize := 50
	flushInterval := 1 * time.Second

	// Robust Sender
	sendBatch := func(b []Log) bool {
		jsonData, _ := json.Marshal(b)
		req, err := http.NewRequest("POST", *urlPtr, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return false // Retry? Logic error usually
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Api-Key", *apiKeyPtr)
		req.Header.Set("Authorization", "Bearer "+*secretKeyPtr)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("⚠️ Network/Server Error: %v. Retrying...\n", err)
			return false // Retry
		}
		defer resp.Body.Close()
		io.Copy(io.Discard, resp.Body)

		if resp.StatusCode >= 500 {
			fmt.Printf("⚠️ Server Error %d. Retrying...\n", resp.StatusCode)
			return false
		}

		if resp.StatusCode >= 400 {
			fmt.Printf("❌ Rejected %d. Possible Config Error. Dropping batch.\n", resp.StatusCode)
			return true // Drop, don't retry bad requests forever
		}

		fmt.Printf("✅ Sent %d logs.\n", len(b))
		return true // Success
	}

	// Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	fmt.Printf("Agent started. Watching %s [%s mode] -> %s\n", *filePtr, *formatPtr, *urlPtr)

	for {
		select {
		case line, ok := <-t.Lines:
			if !ok {
				// File closed?
				return
			}

			parsedLogs, err := parser.Parse(line.Text)
			if err != nil {
				// fmt.Printf("Parse Error: %v\n", err)
				continue
			}

			batch = append(batch, parsedLogs)
			if len(batch) >= batchSize {
				if sendBatch(batch) {
					batch = batch[:0]
				} else {
					// Failed to send, keep accumulating? Or pause?
					// Simple backoff: sleep 1s then continue (batch grows)
					time.Sleep(1 * time.Second)
				}
			}

		case <-ticker.C:
			if len(batch) > 0 {
				if sendBatch(batch) {
					batch = batch[:0]
				}
			}

		case <-sigChan:
			fmt.Println("\nStopping... Flushing remaining logs.")
			if len(batch) > 0 {
				sendBatch(batch)
			}
			return
		}
	}

}