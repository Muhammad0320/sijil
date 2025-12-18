package ingest

import (
	"expvar"
	"sync/atomic"
)

var (

	// Counters (only go up)
	mReceived = expvar.NewInt("ingest_logs_received_total")
	mQueued = expvar.NewInt("ingest_logs_queued_total")
	mFlushed = expvar.NewInt("ingest_logs_written_total")
	mErrors = expvar.NewInt("ingest_errors_total")
	mDropped = expvar.NewInt("ingest_logs_dropped_total")
	
	// Gauges (Go up and down)
	// we use atomic int64 becaue expvar.Int doesn't have it. 
	activeWorkers int64 
)

func init() {
	expvar.Publish("ingest_active_workers", expvar.Func( func() any {
		return atomic.LoadInt64(&activeWorkers)
	}))
}

func RecordReceived(n int) { mReceived.Add(int64(n)) }
func RecordQueued(n int) { mQueued.Add(int64(n)) }
func RecordFlushed(n int) {mFlushed.Add(int64(n)) }
func RecordDropped(n int) {mDropped.Add(int64(n))}
func RecordError() { mErrors.Add(1) }

func WorkerWake() {atomic.AddInt64(&activeWorkers, 1)}
func WorkerSleep() {atomic.AddInt64(&activeWorkers, -1)}
