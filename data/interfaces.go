package data

type IPipelineContainer interface {
	PipelineFromName(name string) *Pipeline
}

type IServerState interface {
	IPipelineContainer
	HasExecution() bool
	IsRunning() bool
	StepOutput(name string) (string, error)
	NewExecution(p *Pipeline, stepInfo []StepInfo) error
	Execute()
}
