package step

import (
	"bufio"
	"errors"
	"log/slog"
	"os/exec"
	"strings"
	"sync"

	"executrix/helper"
	"executrix/server/config"
)

type PSStep struct {
	Name       string
	ScriptPath string
	state      State
	Args       []string
	DependsOn  []string
}

func (s *PSStep) ShowAs() string {
	return s.Name
}

func (s *PSStep) GetState() State {
	return s.state
}

func (s *PSStep) SetState(state State) {
	s.state = state
}

func (step *PSStep) Execute(out *string) {
	step.SetState(Running)

	slog.Info("Excuting PS step", "step", step.Name, "script", step.ScriptPath)
	helper.AppendLine(out, "Excuting PS step: "+step.Name)

	args := []string{"-nologo", "-noprofile", "-noninteractive", step.ScriptPath}
	args = append(args, step.Args...)
	helper.AppendLine(out, "Excution: powershell "+strings.Join(args, " "))
	helper.AppendLine(out, "")

	cmd := exec.Command("powershell", args...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error getting stdout pipe in PS step", "error", err)
		helper.AppendLine(out, "Error getting stdout pipe in PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("Error getting stderr pipe in PS step", "error", err)
		helper.AppendLine(out, "Error getting stderr pipe in PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(2)

	if err := cmd.Start(); err != nil {
		slog.Error("Error starting PS step", "error", err)
		helper.AppendLine(out, "Error starting PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	go func() {
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			helper.AppendLine(out, helper.CleanUpString(scanner.Text()))
		}
		waitgroup.Done()
	}()

	go func() {
		scanner := bufio.NewScanner(errPipe)
		for scanner.Scan() {
			helper.AppendLine(out, helper.CleanUpString(scanner.Text()))
		}
		waitgroup.Done()
	}()

	if err := cmd.Wait(); err != nil {
		slog.Error("Error waiting for PS step", "error", err)
		helper.AppendLine(out, "Error waiting for PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	waitgroup.Wait()

	slog.Info("Finished executing PS step", "step", step.Name)
	helper.AppendLine(out, "")
	helper.AppendLine(out, "")
	helper.AppendLine(out, "Successfully finished PS step: "+step.Name)

	step.SetState(Success)
}

func ReadPSType(s map[string]interface{}, cfg config.GlobalConfig) (*PSStep, error) {
	step := PSStep{}

	if val, ok := s["Name"].(string); !ok {
		return nil, errors.New("could not find step name")
	} else {
		step.Name = val
		slog.Info("Read step name", "s", step.Name)
	}

	if val, ok := s["ScriptPath"].(string); !ok {
		return nil, errors.New("could not find script path")
	} else {
		step.ScriptPath = helper.ReplaceAll(val, cfg.GetVars())
		slog.Info("Read script path", "path", step.ScriptPath)
	}

	if val, ok := s["Arguments"].([]interface{}); !ok {
		return nil, errors.New("could not find script args")
	} else {
		for _, v := range val {
			step.Args = append(step.Args, helper.ReplaceAll(v.(string), cfg.GetVars()))
		}
		slog.Info("Read script args", "args", step.Args)
	}

	if val, ok := s["DependsOn"].([]interface{}); !ok {
		return nil, errors.New("could not find script dependencies")
	} else {
		for _, v := range val {
			step.DependsOn = append(step.DependsOn, v.(string))
		}
		slog.Info("Read script dependencies", "dependencies", step.DependsOn)
	}

	step.state = Waiting

	return &step, nil
}
