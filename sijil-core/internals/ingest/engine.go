package ingest

import (
	"context"
	"expvar"
	"fmt"
	"log"
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
	QueueSize     = 500_000
)

type IngestionEngine struct {
	db            *pgxpool.Pool
	Wal           *WAL
	hub           *hub.Hub
	LogQueue      chan database.LogEntry
	wg            sync.WaitGroup
	activeWorkers int32
}

func NewIngestionEngine(db *pgxpool.Pool, wal *WAL, h *hub.Hub) *IngestionEngine {
	return &IngestionEngine{
		db:       db,
		Wal:      wal,
		hub:      h,
		LogQueue: make(chan database.LogEntry, QueueSize),
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
		case entry := <-e.LogQueue:
			batch = append(batch, entry)

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

	atomic.AddInt32(&e.activeWorkers, 1)  // "I'm active"
	atomic.AddInt32(&e.activeWorkers, -1) // "I'm Done"

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

	for _, row := range batch {
		e.hub.BroadcastLog(row)
	}
}

func (e *IngestionEngine) walJanitor(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if len(e.LogQueue) == 0 {
				if err := e.Wal.CleanupSafeSegments(4); err != nil {
					log.Printf("⚠️ Wal cleanup failed: %v", err)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
