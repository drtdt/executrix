package executrix

import (
	"executrix/data"
	"log/slog"
	"os/exec"
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
			pStep.SetState(execPS(s))
		default:
			slog.Error("Could not find Pipeline Step!")
			pStep.SetState(data.Failed)
			continue
		}

	}

	p.IsRunning = false
	slog.Info("Pipeline finished")
}

func execPS(step *data.PSStep) data.State {
	slog.Info("Excuting PS step", "step", step.Name, "script", step.ScriptPath)

	cmd := exec.Command("Powershell", "-nologo", "-noprofile", "-noninteractive", step.ScriptPath)

	_, err := cmd.StdoutPipe() // todo
	if err != nil {
		slog.Error("Error getting stdout pipe in PS step", "error", err)
		return data.Failed
	}

	if err := cmd.Start(); err != nil {
		slog.Error("Error starting PS step", "error", err)
		return data.Failed
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("Error waiting for PS step", "error", err)
		return data.Failed
	}

	//slog.Info("PS output", "out", stdout)

	slog.Info("Finished executing PS step", "step", step.Name)

	return data.Success
}
