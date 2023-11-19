package data

import (
	"errors"
	"log/slog"
)

type State int

const (
	Waiting State = iota
	Running
	Failed
	Success
	Semi
)

type IStep interface {
	ShowAs() string
	GetState() State
	SetState(b State)
}

type PSStep struct {
	Name       string
	ScriptPath string
	state      State
	//Args       string
	//DependsOn  []string
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

func readPSType(s map[string]interface{}) (*PSStep, error) {
	step := PSStep{}

	if val, ok := s["Name"].(string); !ok {
		return nil, errors.New("could not find step name")
	} else {
		slog.Info("Read step name", "s", val)
		step.Name = val
	}

	if val, ok := s["ScriptPath"].(string); !ok {
		return nil, errors.New("could not find script path")
	} else {
		slog.Info("Read script path", "path", val)
		step.ScriptPath = val
	}

	step.state = Waiting

	return &step, nil
}

func FromJSON(s map[string]interface{}) (IStep, error) {
	val, ok := s["Type"].(string)
	if !ok {
		return nil, errors.New("could not find step type")
	}

	slog.Info("Read step type", "type", val)
	switch val {
	case "PS":
		return readPSType(s)
	default:
		slog.Error("Unknown step type", "type", val)
		return nil, errors.New("unknown step type")
	}
}
