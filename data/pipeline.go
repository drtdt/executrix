package data

import (
	"encoding/json"
	"errors"
	"log/slog"
	"slices"

	"executrix/helper"
)

type Pipeline struct {
	Name        string
	Description string
	Steps       []IStep
}

type StateInfo struct {
	Step  string
	State State
}

func (p Pipeline) FindStep(name string) IStep {
	if idx := slices.IndexFunc(p.Steps, func(s IStep) bool { return s.ShowAs() == name }); idx < 0 {
		return nil
	} else {
		return p.Steps[idx]
	}
}

func (p Pipeline) GetStepStates() []StateInfo {
	var list []StateInfo
	for _, s := range p.Steps {
		list = append(list, StateInfo{
			Step:  s.ShowAs(),
			State: s.GetState(),
		})
	}

	return list
}

func FromJson(path string) (Pipeline, error) {
	bytes, err := helper.ReadFile(path)
	if err != nil {
		return Pipeline{}, err
	}

	var p map[string]interface{}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return Pipeline{}, err
	}

	slog.Debug("Successfully unmarshalled file content", "content", p)

	var pipeline Pipeline

	if val, ok := p["Name"].(string); !ok {
		return Pipeline{}, errors.New("could not find pipeline name")
	} else {
		slog.Debug("Read pipeline name", "s", val)
		pipeline.Name = val
	}

	if val, ok := p["Description"].(string); !ok {
		return Pipeline{}, errors.New("could not find pipeline description")
	} else {
		slog.Debug("Read pipeline description", "s", val)
		pipeline.Description = val
	}

	if val, ok := p["Steps"].([]interface{}); !ok {
		return Pipeline{}, errors.New("error reading pipeline steps")
	} else {
		slog.Debug("Read pipeline steps", "steps", val)
		for _, elem := range val {
			slog.Debug("Read pipeline steps", "steps", elem)

			val, ok := elem.(map[string]interface{})
			if !ok {
				return Pipeline{}, errors.New("unexpected type for step")
			}

			step, err := FromJSON(val)
			if err != nil {
				return Pipeline{}, err
			}

			pipeline.Steps = append(pipeline.Steps, step)
		}
	}

	return pipeline, nil
}
