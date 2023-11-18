package executrix

import (
	"executrix/data"
	"log/slog"
	"time"
)

func ExecutePipeline(p *data.Pipeline, stepInfo []data.StepInfo) {
	slog.Info("Starting pipeline")

	for _, step := range stepInfo {
		if !step.Checked {
			slog.Info("Skipping unchecked step", "step", step.StepName)
			continue
		}

		pStep := p.FindStep(step.StepName)
		if pStep == nil {
			slog.Error("Could not find Pipeline Step!")
			// todo error handling
		}

		pStep.SetState(data.Running)

		switch s := pStep.(type) {
		case *data.PSStep:
			execPS(s)
		default:
			slog.Error("Could not find Pipeline Step!")
			pStep.SetState(data.Failed)
			continue
		}
	}

	p.IsRunning = false
	slog.Info("Pipeline finished")
}

func execPS(step *data.PSStep) {
	slog.Info("Excuting step", "step", step.Name)
	// dummy implementation
	time.Sleep(10000 * time.Millisecond)

	step.SetState(data.Success)
}
