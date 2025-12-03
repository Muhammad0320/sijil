package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/nxadm/tail"
)

type Log struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Level string `json:"level"`
	Service string `json:"service"`
	Message string `json:"message"`
	
}

var logRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(?:\[(.*?)\]\s+)?\[(.*?)\]\s+(.*)$`)

const timeLayout = "2006-01-02 15:04:05"

func main() {

	filePtr := flag.String("f", "test.log", "log file to tail")
	servicePtr := flag.String("s", "log-agent-v1", "service name to tag logs with")
	formarPtr := flag.String("format", "regex", "Log format: 'regex' or 'json' ")
	apiKeyPtr := flag.String("pk", "", "Public API key (pk_live_...)")	
	secretKeyPtr := flag.String("sk", "", "Secret API key (pk_live_...)")	

	flag.Parse()

	if *apiKeyPtr == "" || *secretKeyPtr == "" {
		log.Fatal("Error: you must provide both pk ans sk flags")
	} 

	var parser Parser
	switch *formarPtr {
	case "regex":
		parser = NewRegexParser(*servicePtr)
	case "json":
		parser = NewJsonParser(*servicePtr)
	default:
		log.Fatalf("FATAL: Unknow format '%s'", *formarPtr)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("starting to tail test.log")
	t, err := tail.TailFile(*filePtr, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatalf("Failed to tail file: %v", err)
	}
	
	// --- 	BATCHING CONFIG ----
	var batch []Log
	batchSize := 50 
	flushInterval := 1 * time.Second

	flush := func () {
		if len(batch) == 0 {return} 

		// Serialize the whole batch
		jsonData, _ := json.Marshal(batch)

		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/logs", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-Api-Key", *apiKeyPtr)
	req.Header.Set("Authorization", "Bearer "+*secretKeyPtr)
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send batch: %v\n", err)
		
	} else {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 400 {
			fmt.Printf("⚠️ Server error %d\n", resp.StatusCode)
		} else {	
			fmt.Printf("✅ Sent batch of %d logs.", len(batch))
		}
	}

	batch = batch[:0]
}
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	fmt.Printf("Agent started.  Watching %s [%s mode]...\n", *filePtr, *formarPtr)

	for  {
	   select {
	   case line, ok := <- t.Lines:
			if !ok {return}
			
			parsedLogs, err := parser.Parse(line.Text)
			if err != nil {
				fmt.Printf("Parse Error: %v\n", err)
				continue
			}
			
			batch = append(batch, parsedLogs)
			if len(batch) >= batchSize {
				flush()
			}

		case <- ticker.C: 
			flush()
	   }
	}

}