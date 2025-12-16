package logengine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	APIKey    string
	APISecret string
	Endpoint  string
	BatchSize int
	Interval  time.Duration
}


type LogEntry struct {

	Level string `json:"level"`
	Service string `json:"service"`
	Message string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data map[string]interface{} `json:"data,omitempty"`

}

type Client struct {

	config Config
	queue chan LogEntry
	client *http.Client
	wg sync.WaitGroup
	shutdown chan struct{}
	service string 

}

func NewClient(cfg Config) *Client {

	if cfg.APIKey == "" || cfg.APISecret == "" {
		log.Fatal("LogEngine: APIKey and APISecret are required")
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = "http://localhost:8080/api/v1/logs"
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	if cfg.Interval == 0 {
		cfg.Interval = 5 * time.Second
	}

	c := &Client{
		config: cfg,
		queue: make(chan LogEntry, 1000),
		client: &http.Client{Timeout: 10 * time.Second},
		shutdown: make(chan struct{}),
		service: "default", // can be overridden perlog or global
	}

	c.wg.Add(1)
	

	return  c
}

func (c *Client) SetService(name string)  {
	c.service = name
}

func (c *Client) Info(msg string, data map[string]interface{}) {

	c.push("info", msg, data)

}

func (c *Client) Error(msg string, data map[string]interface{}) {
	c.push("error", msg, data)
}

func (c *Client) Debug(msg string, data map[string]interface{}) {

	c.push("debug", msg, data)
}

func (c *Client) Warn(msg string, data map[string]interface{}) {

	c.push("warn", msg, data)
}

func (c *Client) push(level, msg string, data map[string]interface{}) {

	entry := LogEntry{
		Level: level,
		Message: msg,
		Service: c.service,
		Timestamp: time.Now(),
		Data: data,
	}

	select {
	case c.queue <- entry:
	default:
		fmt.Fprintf(os.Stderr, "LogEngine Queue full: Dropping los: %s\n", msg)
	}
} 

func (c *Client) Close() {
	close(c.shutdown)
	c.wg.Wait()
}

func (c *Client) Worker()  {
	
}

func (c *Client) sendBatch(logs []LogEntry)  {
	payload, err := json.Marshal(logs) 
	if err != nil {
		fmt.Printf("LogEngine SDK Error: Failed to marshal batch %v\n", err)
		return
	}
	
}