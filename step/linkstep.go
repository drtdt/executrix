package step

import (
	"errors"
	"executrix/helper"
	"executrix/server/config"
	"log/slog"
)

type LinkStep struct {
	Name  string // todo: this is public so it can be read in html template - should become decoupled
	Link  string
	state State
}

func (s LinkStep) Type() string {
	return "Link"
}

func (s *LinkStep) ShowAs() string {
	return s.Name
}

func (s *LinkStep) GetState() State {
	return s.state
}

func (s *LinkStep) SetState(state State) {
	s.state = state
}

func (s *LinkStep) Kill() error {
	// nothing to do here
	return nil
}

func (step *LinkStep) Execute(out *string) {
	// nothing to do here (so far)
}

func ReadLinkType(s map[string]interface{}, cfg config.GlobalConfig) (*LinkStep, error) {
	step := LinkStep{}

	if val, ok := s["Name"].(string); !ok {
		return nil, errors.New("could not find step name")
	} else {
		step.Name = val
		slog.Info("Read step name", "s", step.Name)
	}

	if val, ok := s["Link"].(string); !ok {
		return nil, errors.New("could not find link")
	} else {
		step.Link = helper.ReplaceAll(val, cfg.GetVars())
		slog.Info("Read link", "s", step.Link)
	}

	step.state = Waiting

	return &step, nil
}
