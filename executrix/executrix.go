package executrix

import (
	"bufio"
	"errors"
	"executrix/data"
	"executrix/helper"
	"log/slog"
	"os/exec"
	"sync"
)

type Execution struct {
	pipeline   *data.Pipeline
	stepInfo   []data.StepInfo
	outputs    map[string]*string
	currentCmd *exec.Cmd
	finished   bool
}

func NewExecution(p *data.Pipeline, stepInfo []data.StepInfo) (*Execution, error) {
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
		pStep.SetState(data.Running)

		switch s := pStep.(type) {
		case *data.PSStep:
			pStep.SetState(execPS(s, e.outputs[step.StepName]))
		default:
			slog.Error("Could not find Pipeline Step!")
			pStep.SetState(data.Failed)
			continue
		}
	}

	slog.Info("Pipeline finished")
}

// todo make this a member of PSStep?
func execPS(step *data.PSStep, out *string) data.State {
	slog.Info("Excuting PS step", "step", step.Name, "script", step.ScriptPath)

	cmd := exec.Command("Powershell", "-nologo", "-noprofile", "-noninteractive", step.ScriptPath)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error getting stdout pipe in PS step", "error", err)
		return data.Failed
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("Error getting stderr pipe in PS step", "error", err)
		return data.Failed
	}

	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(2)

	if err := cmd.Start(); err != nil {
		slog.Error("Error starting PS step", "error", err)
		return data.Failed
	}

	go func() {
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			//slog.Info("OUT FROM PS", "out", scanner.Text())
			*out += helper.CleanUpString(scanner.Text()) + "\\n"
		}
		waitgroup.Done()
	}()

	go func() {
		scanner := bufio.NewScanner(errPipe)
		for scanner.Scan() {
			//slog.Info("ERR FROM PS", "out", scanner.Text())
			*out += helper.CleanUpString(scanner.Text()) + "\\n"
		}
		waitgroup.Done()
	}()

	if err := cmd.Wait(); err != nil {
		slog.Error("Error waiting for PS step", "error", err)
		return data.Failed
	}

	waitgroup.Wait()

	slog.Info("Finished executing PS step", "step", step.Name)

	return data.Success
}
