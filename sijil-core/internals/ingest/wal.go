package ingest

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sijil-core/internals/database"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	MaxSegmentSize = 10 * 1024 * 1024
	WalDir         = "wal_data"
)

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	ProjectID int                    `json:"-"`
	Data      map[string]interface{} `json:"data,omitempty"`
	SegmentID int                    `json:"-"`
}

type WAL struct {
	dir         string
	mu          sync.Mutex
	activeFile  *os.File
	activeSeq   int
	currentSize int64
	bufWriter   *bufio.Writer
}

func NewWal(dir string) (*WAL, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create wal dir: %w", err)
	}

	w := &WAL{dir: dir}

	// Find the latest segment to continue writing, or start fresh
	lastSeq, err := w.findLastSegment()
	if err != nil {
		return nil, err
	}

	if err := w.openSegment(lastSeq); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WAL) findLastSegment() (int, error) {

	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return 0, err
	}

	maxSeq := 0
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "segment-") && strings.HasSuffix(e.Name(), ".log") {
			var seq int
			_, err := fmt.Sscanf(e.Name(), "segment-%d.wal", &seq)
			if err == nil && seq > maxSeq {
				maxSeq = seq
			}
		}
	}
	// If no files exist, start at 1
	if maxSeq == 0 {
		return 1, nil
	}

	return maxSeq, nil
}

func (w *WAL) openSegment(seq int) error {

	filename := filepath.Join(w.dir, fmt.Sprintf("segment-%06d.wal", seq))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open segment %d, %w", seq, err)
	}

	stat, _ := f.Stat()
	w.activeFile = f
	w.bufWriter = bufio.NewWriterSize(f, 64*1024)
	w.activeSeq = seq
	w.currentSize = stat.Size()

	return nil
}

func (w *WAL) WriteBatch(batch []database.LogEntry) error {

	w.mu.Lock()
	defer w.mu.Unlock()

	// 1. Check Rotation
	if w.currentSize > MaxSegmentSize {
		if err := w.rotate(); err != nil {
			return err
		}
	}

	for i := range batch {

		batch[i].SegmentID = w.activeSeq

		data, err := batch[i].Serialize()
		if err != nil {
			continue
		}

		// Write Length Prefix (for easy reading later)
		length := int32(len(data))
		if err := binary.Write(w.bufWriter, binary.LittleEndian, length); err != nil {
			return err
		}

		// Write Data
		_, err = w.bufWriter.Write(data)
		if err != nil {
			return err
		}

		w.currentSize += int64(4 + len(data))
	}

	return w.bufWriter.Flush()
}

func (w *WAL) Sync() error {

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.activeFile != nil {
		return w.activeFile.Sync()
	}
	return nil

}

// rotate closes the current file and opens the next sequence
func (w *WAL) rotate() error {
	if w.activeFile != nil {
		w.activeFile.Close()
	}

	newSeq := w.activeSeq + 1
	return w.openSegment(newSeq)
}

func (w *WAL) Recover() ([]database.LogEntry, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var allLogs []database.LogEntry

	// Get all segment files sorted
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "segment-") {
			filenames = append(filenames, filepath.Join(w.dir, e.Name()))
		}
	}
	sort.Strings(filenames)

	for _, fname := range filenames {
		logs, err := w.readSegment(fname)
		if err != nil {
			fmt.Printf("⚠️ Corrupt segment %s: %v\n", fname, err)
			continue // Try to read the rest
		}
		allLogs = append(allLogs, logs...)
	}

	return allLogs, nil
}

func (w *WAL) readSegment(path string) ([]database.LogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var logs []database.LogEntry
	for {
		// Read Length
		var length int32
		if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return logs, err
		}

		// Read Data
		data := make([]byte, length)
		if _, err := io.ReadFull(f, data); err != nil {
			return logs, err
		}

		// Deserialize (Assuming JSON for now)
		// You need to ensure LogEntry has a generic deserializer or use json.Unmarshal
		var entry database.LogEntry
		if err := entry.Deserialize(data); err != nil {
			continue
		}
		logs = append(logs, entry)
	}
	return logs, nil
}

// Close gracefully shuts down the WAL
func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.activeFile != nil {
		return w.activeFile.Close()
	}
	return nil
}

func (w *WAL) Reset() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 1. Close active file if open
	if w.activeFile != nil {
		w.activeFile.Close()
		w.activeFile = nil
	}

	// 2. Delete all segment files
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "segment-") {
			os.Remove(filepath.Join(w.dir, e.Name()))
		}
	}

	// 3. Start fresh at segment 1
	return w.openSegment(1)
}

func (w *WAL) CleanupUntil(maxSeqToDelete int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	fmt.Println("Were you even called")

	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return err
	}

	fmt.Printf("🧹 Janitor Waking Up: Safe to delete up to Segment %d\n", maxSeqToDelete)

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "segment-") && strings.HasSuffix(e.Name(), ".wal") {
			var seq int
			fmt.Sscanf(e.Name(), "segment-%d.wal", &seq)

			if seq <= maxSeqToDelete {
				path := filepath.Join(w.dir, e.Name())
				os.Remove(path)
				fmt.Printf("🗑️ Deterministic Cleanup: Deleted committed segment %d\n", seq)
			}
		}
	}
	return nil
}
