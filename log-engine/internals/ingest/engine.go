package ingest

import (
	"context"
	"fmt"
	"log"
	"log-engine/internals/database"
	"log-engine/internals/hub"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	BatchSize     = 500
	FlushInterval = 1 * time.Second
	WorkerCount = 5
	QueueSize = 10_000
)

type IngestionEngine struct {
	db *pgxpool.Pool 
	Wal *WAL
	hub *hub.Hub
	LogQueue chan database.LogEntry
	wg sync.WaitGroup
}

func NewIngestionEngine(db *pgxpool.Pool, wal *WAL, h *hub.Hub) *IngestionEngine {
	return &IngestionEngine{
		db: db,
		Wal: wal, 
		hub: h,
		LogQueue: make(chan database.LogEntry, QueueSize),
	}
}

func (e *IngestionEngine) Start(ctx context.Context) {
	fmt.Printf("Staring ingesting engine with %d workers", WorkerCount)
	
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
		case entry := <- e.LogQueue: 
			batch = append(batch, entry) 

			if len(batch) >= BatchSize {
				e.flush(ctx, batch)
				batch = batch[:0]
			}
		case <- ticker.C:
			if len(batch) > 0 {
				e.flush(ctx, batch)
				batch = batch[:0]
			}
		case <- ctx.Done(): 
			if len(batch) > 0 {
				e.flush(ctx, batch)
			}
			return
		}

	}

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
		[]string{"timestamp", "level", "message", "service", "prorject_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Printf("⚠️ BATCH INSERT FAILED %s", err)
		return
	}

	for _, row := range batch {
		e.hub.BroadcastLog(row)
	}
}