package data

import (
	"encoding/json"
	"errors"
	"executrix/helper"
	"log/slog"
)

type Step interface {
	ShowAs() string
}

type PSStep struct {
	Name       string
	ScriptPath string
}

type Pipeline struct {
	Name        string
	Description string
	Steps       []Step
}

func (s PSStep) ShowAs() string {
	return s.Name
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

	slog.Info("Successfully read", "pipeline", p)

	var pipeline Pipeline

	if val, ok := p["Name"].(string); !ok {
		return Pipeline{}, errors.New("could not find pipeline name")
	} else {
		pipeline.Name = val
	}

	if val, ok := p["Description"].(string); !ok {
		return Pipeline{}, errors.New("could not find pipeline description")
	} else {
		pipeline.Description = val
	}

	return pipeline, nil
}
