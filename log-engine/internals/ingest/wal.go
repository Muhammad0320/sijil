package ingest

import (
	"encoding/json"
	"fmt"
	"log-engine/internals/database"
	"os"
	"sync"
)

type WAL struct {
	mu sync.Mutex
	file *os.File 
}

func NewWal(path string) (*WAL, error) {

	file, err := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL: %w", err)
	}

	return &WAL{file: file}, nil
}

func (w *WAL) WriteLog(entry database.LogEntry) error {

	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("wal marshal error: %w", err)
	}

	if _, err := w.file.Write(data); err != nil {
		return fmt.Errorf("wal write error: %w", err)
	}

	// Sync (The paranoid step)
	if err := w.file.Sync(); err != nil {
		return fmt.Errorf("wal sync error: %w", err)
	}

	return  nil 
}

func (w *WAL) Close() error {
	return w.file.Close()
}