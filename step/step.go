package step

import (
	"errors"
	"log/slog"

	"executrix/server/config"
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
	Type() string
	GetState() State
	SetState(b State)
	Execute(out *string)
	Kill() error
}

func StepFromJSON(s map[string]interface{}, cfg config.GlobalConfig) (IStep, error) {
	val, ok := s["Type"].(string)
	if !ok {
		return nil, errors.New("could not find step type")
	}

	slog.Info("Read step type", "type", val)
	switch val {
	case "PS":
		return ReadPSType(s, cfg)
	case "Link":
		return ReadLinkType(s, cfg)
	default:
		slog.Error("Unknown step type", "type", val)
		return nil, errors.New("unknown step type")
	}
}
