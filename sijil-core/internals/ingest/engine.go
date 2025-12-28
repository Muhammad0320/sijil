package ingest

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sijil-core/internals/database"
	"sijil-core/internals/hub"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	BatchSize     = 5_000
	FlushInterval = 1 * time.Second
	WorkerCount   = 25
	QueueSize     = 10_000
)

type IngestionEngine struct {
	db            *pgxpool.Pool
	Wal           *WAL
	hub           *hub.Hub
	LogQueue      chan []database.LogEntry
	wg            sync.WaitGroup
	activeWorkers int32
	CommitedSeq   chan int
}

func NewIngestionEngine(db *pgxpool.Pool, wal *WAL, h *hub.Hub) *IngestionEngine {
	return &IngestionEngine{
		db:          db,
		Wal:         wal,
		hub:         h,
		LogQueue:    make(chan []database.LogEntry, QueueSize),
		CommitedSeq: make(chan int, 100),
	}
}

func (e *IngestionEngine) Start(ctx context.Context) {
	fmt.Printf("Staring ingesting engine with %d workers", WorkerCount)

	e.wg.Add(1)
	go e.walJanitor(ctx)

	expvar.Publish("ingest_queue_depth", expvar.Func(func() interface{} {
		return len(e.LogQueue)
	}))

	// Tell the Tracker how many workers we're hiring
	e.wg.Add(WorkerCount)
	for i := range WorkerCount {
		go e.worker(ctx, i)
	}

	// Background WAL syncer
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				e.Wal.Sync()
			case <-ctx.Done():
				return

			}
		}
	}()

}

func (e *IngestionEngine) Shutdown() {
	e.wg.Wait()

	fmt.Println("All workers finished flushing ✅")
}

func (e *IngestionEngine) worker(ctx context.Context, id int) {

	defer e.wg.Done()

	batch := make([]database.LogEntry, 0, BatchSize)

	ticker := time.NewTicker(FlushInterval)
	defer ticker.Stop()

	for {

		select {
		case microBatch := <-e.LogQueue:

			if err := e.Wal.WriteBatch(microBatch); err != nil {
				fmt.Printf("failed to sync to buffered wal %s \n", err)
			}

			batch = append(batch, microBatch...)

			if len(batch) >= BatchSize {
				WorkerWake()
				e.safeFlush(ctx, batch)
				WorkerSleep()
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				WorkerWake()
				e.safeFlush(ctx, batch)
				WorkerSleep()
				batch = batch[:0]
			}
		case <-ctx.Done():
			if len(batch) > 0 {
				WorkerWake()
				e.safeFlush(ctx, batch)
				WorkerSleep()
			}
			return
		}

	}

}

func (e *IngestionEngine) safeFlush(ctx context.Context, batch []database.LogEntry) {

	atomic.AddInt32(&e.activeWorkers, 1)        // "I'm active"
	defer atomic.AddInt32(&e.activeWorkers, -1) // "I'm Done"

	e.flush(ctx, batch)

}

func (e *IngestionEngine) flush(ctx context.Context, batch []database.LogEntry) {

	rows := make([][]interface{}, len(batch))
	for i, log := range batch {
		rows[i] = []interface{}{
			log.Timestamp,
			log.Level,
			log.Message,
			log.Service,
			log.ProjectID,
		}
	}

	_, err := e.db.CopyFrom(
		ctx,
		pgx.Identifier{"logs"},
		[]string{"timestamp", "level", "message", "service", "project_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Printf("⚠️ BATCH INSERT FAILED %s", err)
		RecordError()
		return
	}

	RecordFlushed(len(batch))

	// The Janitor Logic
	maxSeq := 0
	for _, log := range batch {
		if log.SegmentID > maxSeq {
			maxSeq = log.SegmentID
		}
	}

	if maxSeq > 0 {

		select {
		case e.CommitedSeq <- maxSeq:
		default:

		}
	}
	// Ends here

	for _, row := range batch {
		e.hub.BroadcastLog(row)
	}
}

func (e *IngestionEngine) walJanitor(ctx context.Context) {
	defer e.wg.Done()

	var maxSafeSeq int
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// EVENT: When flush succeeeds, update out knowlege
		case seq := <-e.CommitedSeq:
			if seq > maxSafeSeq {
				maxSafeSeq = seq
			}
		// TIME: Act on that knowlwdge
		case <-ticker.C:

			threshold := maxSafeSeq
			if maxSafeSeq >= e.Wal.activeSeq {
				threshold = e.Wal.activeSeq - 1
			}

			if threshold > 0 {
				e.Wal.CleanupUntil(threshold)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (e *IngestionEngine) CheckDiskPressure() error {

	var size int64
	err := filepath.Walk(e.Wal.dir, func(_ string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			size += info.Size()
		}

		return err
	})

	if err != nil {
		return err
	}

	if size > 2*1024*1024*1024 {
		return fmt.Errorf("disk pressure: WAL size %.2f MB exceeds limit", float64(size)/1024/1024)
	}

	return nil
}
