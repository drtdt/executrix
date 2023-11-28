package executrix

import (
	"errors"
	"executrix/data"
	"executrix/pipeline"
	"log/slog"
	"os/exec"
)

type Execution struct {
	pipeline   *pipeline.Pipeline
	stepInfo   []data.StepInfo
	outputs    map[string]*string
	currentCmd *exec.Cmd
	finished   bool
}

func NewExecution(p *pipeline.Pipeline, stepInfo []data.StepInfo) (*Execution, error) {
	if p == nil {
		return nil, errors.New("pipeline must not be nil")
	}

	return &Execution{
		pipeline:   p,
		stepInfo:   stepInfo,
		outputs:    make(map[string]*string),
		currentCmd: nil,
		finished:   false,
	}, nil
}

func (e *Execution) PipelineName() string {
	return e.pipeline.Name
}

func (e *Execution) SetFinished() {
	e.finished = true
}

func (e Execution) IsFinished() bool {
	return e.finished
}

func (e Execution) StepOutput(step string) (string, error) {
	output, ok := e.outputs[step]
	if !ok {
		return "", errors.New("step not found")
	}

	return *output, nil
}

func (e *Execution) Execute() {
	slog.Info("Starting pipeline")

	for _, step := range e.stepInfo {
		if !step.Checked {
			slog.Info("Skipping unchecked step", "step", step.StepName)
			continue
		}

		pStep := e.pipeline.FindStep(step.StepName)
		if pStep == nil {
			slog.Error("Could not find Pipeline Step!")
			// todo error handling
		}

		s := ""
		e.outputs[step.StepName] = &s

		pStep.Execute(e.outputs[step.StepName])
	}

	slog.Info("Pipeline finished")
}
