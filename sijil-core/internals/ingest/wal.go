package ingest

import (
	"fmt"
	"os"
	"path/filepath"
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

// rotate closes the current file and opens the next sequence
func (w *WAL) rotate() error {
	if w.activeFile != nil {
		w.activeFile.Close()
	}

	newSeq := w.activeSeq + 1
	return w.openSegment(newSeq)
}
