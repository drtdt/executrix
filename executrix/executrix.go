package executrix

import (
	"executrix/data"
	"log/slog"
	"time"
)

func ExecutePipeline(p *data.Pipeline, stepInfo []data.StepInfo) {
	slog.Info("Starting pipeline")

	for _, step := range stepInfo {
		if step.Checked {
			slog.Info("Excuting step", "step", step.StepName)

			// dummy implementation
			time.Sleep(10000 * time.Millisecond)
		} else {
			slog.Info("Skipping unchecked step", "step", step.StepName)
		}
	}

	p.IsRunning = false
	slog.Info("Pipeline finished")
}
