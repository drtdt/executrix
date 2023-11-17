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
			pStep := p.FindStep(step.StepName)
			if pStep == nil {
				slog.Error("Could not find Pipeline Step!")
				// todo error handling
			}

			slog.Info("Excuting step", "step", step.StepName)
			pStep.SetState(data.Running)

			// dummy implementation
			time.Sleep(10000 * time.Millisecond)

			pStep.SetState(data.Success)
		} else {
			slog.Info("Skipping unchecked step", "step", step.StepName)
		}
	}

	p.IsRunning = false
	slog.Info("Pipeline finished")
}
