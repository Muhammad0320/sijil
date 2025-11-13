package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nxadm/tail"
)

type Log struct {
	Level string
	Message string
	Service string
}

func main() {

	fmt.Println("starting to tail test.log")
	
	t, err := tail.TailFile("test.log", tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Failed to tail file: %v", err)
	}
	
	
	for line := range t.Lines {
		
		newLog := Log{Level: "info", Service: "log-agent-v1", Message: line.Text}
		jsonData, err := json.Marshal(newLog)
		if err != nil {
			fmt.Printf("error marchaling json: %s", err)
			break
		}


	// Create a new HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/logs", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Printf("Response Body: %s\n", string(body))


	}
}