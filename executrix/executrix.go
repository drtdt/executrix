package executrix

import (
	"executrix/data"
	"log/slog"
	"time"
)

func ExecutePipeline(p *data.Pipeline) {
	slog.Info("Starting pipeline")

	// dummy implementation
	time.Sleep(10000 * time.Millisecond)

	p.IsRunning = false
	slog.Info("Pipeline finished")
}
