package data

import (
	"errors"
	"log/slog"
)

type PSStep struct {
	Name       string
	ScriptPath string
}

type Step interface {
	ShowAs() string
}

func (s PSStep) ShowAs() string {
	return s.Name
}

func readPSType(s map[string]interface{}) (PSStep, error) {
	step := PSStep{}
	if val, ok := s["Name"].(string); !ok {
		return PSStep{}, errors.New("could not find step name")
	} else {
		slog.Info("Read step name", "s", val)
		step.Name = val
	}

	return step, nil
}

func FromJSON(s map[string]interface{}) (Step, error) {
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
