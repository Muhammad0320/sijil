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
	apiKeyPtr := flag.String("pk", "", "Public API key (pk_live_...)")	
	secretKeyPtr := flag.String("sk", "", "Secret API key (pk_live_...)")	

	flag.Parse()

	if *apiKeyPtr == "" || *secretKeyPtr == "" {
		log.Fatal("Error: you must provide both pk ans sk flags")
	} 

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("starting to tail test.log")
	
	t, err := tail.TailFile(*filePtr, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Failed to tail file: %v", err)
	}
	
	
	for line := range t.Lines {
		newLog := Log{}

		matches := logRegex.FindStringSubmatch(line.Text)
		if matches == nil {

			fmt.Println("Line didn't match pattern, sending as info")
			newLog.Level = "info"
			newLog.Message = line.Text
			newLog.Service = *servicePtr
		} else {
			fmt.Println("Line matched! parsing...")
			parsedTime, err := time.Parse(timeLayout, matches[1])
			if err == nil {
				newLog.Timestamp = parsedTime
			}

			if matches[2] == ""{
				newLog.Service = *servicePtr
			} else {
				newLog.Service = matches[2]
			}

			newLog.Level = matches[3]
			newLog.Message = matches[4]
		}

		jsonData, err := json.Marshal(newLog)
		if err != nil {
			fmt.Printf("error marchaling json: %s", err)
			continue
		}
	// Create a new HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/logs", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		continue
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-Api-Key", *apiKeyPtr)
	req.Header.Set("Authorization", "Bearer "+*secretKeyPtr)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		continue
	}
	
	io.ReadAll(req.Body)
	resp.Body.Close()
	
	fmt.Printf("Send log, server responded with status: %d\n", resp.StatusCode)
	}
}