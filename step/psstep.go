package step

import (
	"bufio"
	"errors"
	"executrix/helper"
	"log/slog"
	"os/exec"
	"sync"
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

	args := []string{"-nologo", "-noprofile", "-noninteractive", step.ScriptPath}
	args = append(args, step.Args...)

	cmd := exec.Command("powershell", args...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error getting stdout pipe in PS step", "error", err)
		step.SetState(Failed)
		return
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("Error getting stderr pipe in PS step", "error", err)
		step.SetState(Failed)
		return
	}

	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(2)

	if err := cmd.Start(); err != nil {
		slog.Error("Error starting PS step", "error", err)
		step.SetState(Failed)
		return
	}

	go func() {
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			*out += helper.CleanUpString(scanner.Text()) + "\\n"
		}
		waitgroup.Done()
	}()

	go func() {
		scanner := bufio.NewScanner(errPipe)
		for scanner.Scan() {
			*out += helper.CleanUpString(scanner.Text()) + "\\n"
		}
		waitgroup.Done()
	}()

	if err := cmd.Wait(); err != nil {
		slog.Error("Error waiting for PS step", "error", err)
		step.SetState(Failed)
		return
	}

	waitgroup.Wait()

	slog.Info("Finished executing PS step", "step", step.Name)

	step.SetState(Success)
}

func ReadPSType(s map[string]interface{}) (*PSStep, error) {
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
		step.ScriptPath = val
		slog.Info("Read script path", "path", step.ScriptPath)
	}

	if val, ok := s["Arguments"].([]interface{}); !ok {
		return nil, errors.New("could not find script args")
	} else {
		for _, v := range val {
			step.Args = append(step.Args, v.(string))
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
