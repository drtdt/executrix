package data

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type Pipeline struct {
	Name string
}

func FromJson(path string) (Pipeline, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return Pipeline{}, err
	}

	slog.Debug("Successfully opened", "file", path)

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Pipeline{}, err
	}

	var pipeline Pipeline
	err = json.Unmarshal(byteValue, &pipeline)
	if err != nil {
		return Pipeline{}, err
	}

	return pipeline, nil
}
