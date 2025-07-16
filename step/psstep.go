package step

import (
	"bufio"
	"errors"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"executrix/helper"
	"executrix/server/config"
)

type PSStep struct {
	Name       string // todo: this is public so it can be read in html template - should become decoupled
	DependsOn  []string
	Default    bool
	scriptPath string
	state      State
	args       []string
	cmd        *exec.Cmd
}

func (s PSStep) Type() string {
	return "PS"
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

func (s *PSStep) Kill() error {
	if s.cmd != nil {
		slog.Info("Trying to kill PS step", "step", s.Name)
		return s.cmd.Process.Kill()
	}

	return nil
}

func (step *PSStep) Execute(out *string) {
	step.SetState(Running)

	start := time.Now()

	slog.Info("Excuting PS step", "step", step.Name, "script", step.scriptPath)
	helper.AppendLine(out, "Excuting PS step: "+step.Name)

	args := []string{"-nologo", "-noprofile", "-noninteractive", step.scriptPath}
	args = append(args, step.args...)
	helper.AppendLine(out, "Excution: powershell "+strings.Join(args, " "))
	helper.AppendLine(out, "")

	g, err := helper.NewProcessExitGroup()
	if err != nil {
		slog.Error("Error getting creating process exit group", "error", err)
		helper.AppendLine(out, "Error getting creating process exit group: "+err.Error())
		step.SetState(Failed)
		return
	}
	defer g.Dispose()

	step.cmd = exec.Command("powershell", args...)
	defer func() { step.cmd = nil }()

	outPipe, err := step.cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error getting stdout pipe in PS step", "error", err)
		helper.AppendLine(out, "Error getting stdout pipe in PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	errPipe, err := step.cmd.StderrPipe()
	if err != nil {
		slog.Error("Error getting stderr pipe in PS step", "error", err)
		helper.AppendLine(out, "Error getting stderr pipe in PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(2)

	if err := step.cmd.Start(); err != nil {
		slog.Error("Error starting PS step", "error", err)
		helper.AppendLine(out, "Error starting PS step: "+err.Error())
		step.SetState(Failed)
		return
	}

	if err := g.AddProcess(step.cmd.Process); err != nil {
		slog.Error("Error adding process to process exit group", "error", err)
		helper.AppendLine(out, "Error adding process to process exit group: "+err.Error())
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

	if err := step.cmd.Wait(); err != nil {
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
	helper.AppendLine(out, "Duration: "+strconv.FormatFloat(time.Since(start).Seconds(), 'f', -1, 64)+"secs")

	step.SetState(Success)
}

func ReadPSType(s map[string]interface{}, cfg config.GlobalConfig) (*PSStep, error) {
	step := PSStep{}
	step.cmd = nil

	if val, ok := s["Name"].(string); !ok {
		return nil, errors.New("could not find step name")
	} else {
		step.Name = val
		slog.Info("Read step name", "s", step.Name)
	}

	if val, ok := s["Default"].(bool); !ok {
		step.Default = false
	} else {
		step.Default = val
		slog.Info("Read step name", "s", step.Default)
	}

	if val, ok := s["ScriptPath"].(string); !ok {
		return nil, errors.New("could not find script path")
	} else {
		step.scriptPath = helper.ReplaceAll(val, cfg.GetVars())
		slog.Info("Read script path", "path", step.scriptPath)
	}

	if val, ok := s["Arguments"].([]interface{}); !ok {
		return nil, errors.New("could not find script args")
	} else {
		for _, v := range val {
			step.args = append(step.args, helper.ReplaceAll(v.(string), cfg.GetVars()))
		}
		slog.Info("Read script args", "args", step.args)
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
