package step

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
	Execute(out *string)
}

func FromJSON(s map[string]interface{}) (IStep, error) {
	val, ok := s["Type"].(string)
	if !ok {
		return nil, errors.New("could not find step type")
	}

	slog.Info("Read step type", "type", val)
	switch val {
	case "PS":
		return ReadPSType(s)
	default:
		slog.Error("Unknown step type", "type", val)
		return nil, errors.New("unknown step type")
	}
}
