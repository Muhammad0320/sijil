package sijil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	batchSize = 250

	maxQueue = 4096

	workerCount = 3

	retryCount = 3
)

type Config struct {
	APIKey    string
	APISecret string
	Endpoint  string // For self hosted users
	Service   string

	// The Tuning knob
	FlushTime time.Duration
	Silent    bool
}

// For lazy devs (myself included)
func DefaultConfig(key, secret string) Config {
	return Config{
		APIKey:    key,
		APISecret: secret,
		Endpoint:  "https://api.sijil.dev/v1/ingest",
		FlushTime: 1 * time.Second,
		Service:   "default",
	}
}

type LogEntry struct {
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type Client struct {
	config   Config
	queue    chan LogEntry
	client   *http.Client
	wg       sync.WaitGroup
	shutdown chan struct{}
	isClosed bool
	mu       sync.Mutex
	service  string

	flushSig chan chan struct{}
}

func NewClient(cfg Config) *Client {

	if cfg.APIKey == "" || cfg.APISecret == "" {
		log.Fatal("Sijil: Credentials missing!")
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://api.sijil.dev/v1/ingest"
	}

	if cfg.FlushTime == 0 || cfg.FlushTime < 250*time.Millisecond {
		cfg.FlushTime = 500 * time.Millisecond
	}

	svc := cfg.Service
	if svc == "" {
		svc = "default"
	}

	c := &Client{
		config:   cfg,
		queue:    make(chan LogEntry, maxQueue),
		client:   &http.Client{Timeout: 5 * time.Second},
		shutdown: make(chan struct{}),
		service:  svc, // can be overridden perlog or global
		flushSig: make(chan chan struct{}),
	}

	for i := range workerCount {
		c.wg.Add(1)
		go c.worker(i)
	}

	return c
}

func (c *Client) SetService(name string) {
	// To prevent multiple updates
	c.mu.Lock()
	defer c.mu.Unlock()
	c.service = name
}

func (c *Client) Info(msg string, data map[string]interface{}) { c.push("info", msg, data) }

func (c *Client) Error(msg string, data map[string]interface{}) { c.push("error", msg, data) }

func (c *Client) Debug(msg string, data map[string]interface{}) { c.push("debug", msg, data) }

func (c *Client) Warn(msg string, data map[string]interface{}) { c.push("warn", msg, data) }

func (c *Client) Critical(msg string, data map[string]interface{}) { c.push("critical", msg, data) }

func (c *Client) push(level, msg string, data map[string]interface{}) {

	c.mu.Lock()
	if c.isClosed {
		c.mu.Unlock()
		return
	}
	svc := c.service
	c.mu.Unlock()

	entry := LogEntry{
		Level:     level,
		Message:   msg,
		Service:   svc,
		Timestamp: time.Now(),
		Data:      data,
	}

	select {
	case c.queue <- entry:
	default:
		fmt.Fprintf(os.Stderr, "Sijil Queue full: Dropping los: %s\n", msg)
	}
}

func (c *Client) Flush() {

	c.mu.Lock()
	defer c.mu.Unlock()

	// Drain the channel buffer into the workers. Look uptil the channel is empty
	timeout := time.After(2 * time.Second)
	for len(c.queue) > 0 {

		select {

		case <-timeout:
			if !c.config.Silent {
				fmt.Println("Sijil: Flush timeout waiting for queue drain")
			}
			return
		default:
			time.Sleep(10 * time.Millisecond) // Tiny sleep to allow workers to pick up
		}

	}

	// 2. Forcr workers to flush their internal bufferss.
	ack := make(chan struct{})
	for range workerCount {
		c.flushSig <- ack
	}

	// 3. Wait for all workers to acknowledge.
	for range workerCount {
		<-ack
	}

}

func (c *Client) Close() {
	c.mu.Lock()
	if c.isClosed {
		c.mu.Unlock()
		return
	}
	c.isClosed = true
	c.mu.Unlock()

	close(c.queue)
	c.wg.Wait()
}

func (c *Client) sendWithRetry(logs []LogEntry) {
	payload, err := json.Marshal(logs)
	if err != nil {
		fmt.Printf("Sijil SDK Error: Failed to marshal batch %v\n", err)
		return
	}

	for attempts := 0; attempts < retryCount; attempts++ {
		// Backoff: 100ms, 200ms, 400ms...
		if attempts > 0 {
			time.Sleep(time.Duration(100*1<<attempts) * time.Millisecond)
		}

		req, _ := http.NewRequest("POST", c.config.Endpoint, bytes.NewReader(payload))

		req.Header.Set("Content-Type", "application/json") //
		req.Header.Set("X-Api-Key", c.config.APIKey)
		req.Header.Set("Authorization", "Bearer "+c.config.APISecret)

		res, err := c.client.Do(req)

		// Network Error (DNS, timeout) -> Retry
		if err != nil {

			if !c.config.Silent {
				fmt.Printf("Sijil SDK Error: Failed to send batch: %v\n", err)
			}

			continue
		}
		defer res.Body.Close()

		// Success
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			return
		}

		// Server error -> Retry
		if res.StatusCode >= 500 {
			if !c.config.Silent {
				fmt.Printf("Sijil SDK Error: Server error %d (attempt %d)\n", res.StatusCode, attempts+1)
			}
			continue
		}

		// Client error -> DO NOT retry (it will never succeed)
		if res.StatusCode >= 400 {
			if !c.config.Silent {
				fmt.Printf("Sijil: Rejected %d (Bad Config/Auth). Dropping batch.\n", res.StatusCode)
			}
			return
		}

	}

	fmt.Fprintf(os.Stderr, "Sijil critical: Failed to send %d logs after multiple retries\n", len(logs))
}

func (c *Client) worker(_ int) {
	defer c.wg.Done()

	buffer := make([]LogEntry, 0, batchSize)
	ticker := time.NewTicker(c.config.FlushTime)
	defer ticker.Stop()

	flush := func() {
		if len(buffer) == 0 {
			return
		}

		// Copy buffer to free for new logs immediately
		batch := make([]LogEntry, len(buffer))
		copy(batch, buffer)
		buffer = buffer[:0]

		c.sendWithRetry(batch)
	}

	for {
		select {
		case entry, ok := <-c.queue:
			if !ok {
				flush()
				return
			}
			buffer = append(buffer, entry)
			if len(buffer) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()

		case ack := <-c.flushSig:
			flush()
			ack <- struct{}{} // Signal back I'm empty

		}
	}

}
