package ingest

import "expvar"

var (

	// Counters (only go up)
	mReceived = expvar.NewInt("ingest_logs_received_total")
	mQueued = expvar.NewInt("ingest_log_queued_total")
	mFlushed = expvar.NewInt("ingest_log_written_total")
	mErrors = expvar.NewInt("ingest_log_errors_total")
	
	// Gauges (Go up and down)
	// We handle queue depth dynamically, we don't need a variable here
)


func RecordReceived(count int) {
	mReceived.Add(int64(count))
}

func RecordError() {
	mErrors.Add(1)
}