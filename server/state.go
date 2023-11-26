package server

import (
	"errors"
	"executrix/data"
	"executrix/executrix"
	"executrix/helper"
	"log/slog"
	"slices"
)

type ServerState struct {
	Pipelines []data.Pipeline
	execution *executrix.Execution
}

func NewServerState(pipelineDir string) (ServerState, error) {
	state := ServerState{}

	if err := state.reloadPipelines(pipelineDir); err != nil {
		slog.Error("Error while reloading pipeline configs", "err", err)
		return ServerState{}, errors.New("error loading pipeline configs")
	}

	return state, nil
}

func (s *ServerState) PipelineFromName(name string) *data.Pipeline {
	if idx := slices.IndexFunc(s.Pipelines, func(p data.Pipeline) bool { return p.Name == name }); idx < 0 {
		return nil
	} else {
		return &s.Pipelines[idx]
	}
}

func (s *ServerState) IsRunning() bool {
	return s.execution != nil && !s.execution.IsFinished()
}

func (s *ServerState) HasExecution() bool {
	return s.execution != nil
}

func (s *ServerState) StepOutput(step string) (string, error) {
	if !s.HasExecution() {
		return "", errors.New("no perforing or performed execution")
	}

	return s.execution.StepOutput(step)
}

func (s *ServerState) NewExecution(p *data.Pipeline, stepInfo []data.StepInfo) error {
	exec, err := executrix.NewExecution(p, stepInfo)
	if err != nil {
		return errors.New("failed to create new execution")
	}

	s.execution = exec

	return nil
}

func (s *ServerState) Execute() {
	if s.execution == nil {
		slog.Warn("Trying to call ServerStore::Execute when ServerStore::Execution is nil")
		return
	}

	s.execution.Execute()

	s.execution.SetFinished()
}

func (s *ServerState) reloadPipelines(pipelineDir string) error {
	s.Pipelines = nil
	slog.Debug("Cleared piplines before reloading")

	result, err := helper.FindAllFiles(pipelineDir)
	if err != nil {
		return err
	}

	for _, file := range result {
		pipeline, err := data.FromJson(file)
		if err != nil {
			slog.Error("Error reading pipline configuration", "file", file, "error", err)
			// todo - put info to html?
			continue
		}

		s.Pipelines = append(s.Pipelines, pipeline)
	}

	return nil
}
