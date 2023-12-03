package state

import (
	"errors"
	"log/slog"
	"slices"

	"executrix/data"
	"executrix/executrix"
	"executrix/helper"
	"executrix/pipeline"
	"executrix/server/config"
)

type IPipelineContainer interface {
	PipelineFromName(name string) *pipeline.Pipeline
}

type IServerState interface {
	IPipelineContainer
	HasExecution() bool
	IsRunning() bool
	StepOutput(name string) (string, error)
	NewExecution(p *pipeline.Pipeline, stepInfo []data.StepInfo) error
	Execute()
}

type ServerState struct {
	Pipelines []pipeline.Pipeline
	execution *executrix.Execution
}

func NewServerState(pipelineDir string, cfg config.GlobalConfig) (ServerState, error) {
	state := ServerState{}

	if err := state.reloadPipelines(pipelineDir, cfg); err != nil {
		slog.Error("Error while reloading pipeline configs", "err", err)
		return ServerState{}, errors.New("error loading pipeline configs")
	}

	return state, nil
}

func (s *ServerState) PipelineFromName(name string) *pipeline.Pipeline {
	if idx := slices.IndexFunc(s.Pipelines, func(p pipeline.Pipeline) bool { return p.Name == name }); idx < 0 {
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

func (s *ServerState) NewExecution(p *pipeline.Pipeline, stepInfo []data.StepInfo) error {
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

func (s *ServerState) reloadPipelines(pipelineDir string, cfg config.GlobalConfig) error {
	s.Pipelines = nil
	slog.Debug("Cleared piplines before reloading")

	result, err := helper.FindAllFiles(pipelineDir)
	if err != nil {
		return err
	}

	for _, file := range result {
		pipeline, err := pipeline.PipelineFromJson(file, cfg)
		if err != nil {
			slog.Error("Error reading pipline configuration", "file", file, "error", err)
			// todo - put info to html?
			continue
		}

		s.Pipelines = append(s.Pipelines, pipeline)
	}

	return nil
}
