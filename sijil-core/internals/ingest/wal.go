package ingest

import (
	"fmt"
	"os"
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
