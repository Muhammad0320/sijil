package ingest

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sijil-core/internals/database"
	"sort"
	"strings"
	"sync"
)

const (
	MaxSegmentSize = 10 * 1024 * 1024
	WalDir         = "wal_data"
)

type WAL struct {
	dir         string
	mu          sync.Mutex
	activeFile  *os.File
	activeSeq   int
	currentSize int64
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
			_, err := fmt.Sscanf(e.Name(), "segment-%d.log", &seq)
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

	filename := filepath.Join(w.dir, fmt.Sprintf("segment-%06d", seq))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open segment %d, %w", seq, err)
	}

	stat, _ := f.Stat()
	w.activeFile = f
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

	for _, entry := range batch {
		data, err := entry.Serialize()
		if err != nil {
			continue
		}

		// Write Length Prefix (for easy reading later)
		length := int32(len(data))
		if err := binary.Write(w.activeFile, binary.LittleEndian, length); err != nil {
			return err
		}

		// Write Data
		n, err := w.activeFile.Write(data)
		if err != nil {
			return err
		}

		w.currentSize += int64(4 + n)
	}

	return w.activeFile.Sync()
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

func (w *WAL) CleanupOldSegments() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		var seq int
		fmt.Sscanf(e.Name(), "segment-%d.log", &seq)

		if seq < w.activeSeq {
			path := filepath.Join(w.dir, e.Name())
			os.Remove(path)
		}
	}
	return nil
}
