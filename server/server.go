package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"text/template"

	"executrix/constants"
	"executrix/server/config"
	"executrix/server/routes"
	"executrix/server/state"
)

type Server struct {
	serverConfig config.ServerConfig
	globalConfig config.GlobalConfig
	state        state.ServerState
	indexPage    template.Template
	pipelinePage template.Template
}

func NewServer(serverConfig config.ServerConfig) (Server, error) {
	// loading global configuration
	globalConfig, err := config.GlobalConfigFromJson(filepath.Join(serverConfig.GetConfigDir(), constants.GLOBAL_CONFIG_FILE))
	if err != nil {
		slog.Error("Failed to load gobal config", "error", err)
		return Server{}, err
	}

	slog.Info("Sucessfully read global config")

	// loading html templates
	indexTemplate, err := template.ParseFiles("html/index.html")
	if err != nil {
		slog.Error("Failed to parse index.html", "error", err)
		return Server{}, err
	}

	pipelineTemplate, err := template.ParseFiles("html/pipeline.html")
	if err != nil {
		slog.Error("Failed to parse pipeline.html", "error", err)
		return Server{}, err
	}

	// creating struct for tracking the state of the server
	state, err := state.NewServerState(serverConfig.GetPipelineDir(), globalConfig)
	if err != nil {
		slog.Error("Failed to read pipeline configs", "error", err)
		return Server{}, err
	}

	return Server{
		serverConfig: serverConfig,
		globalConfig: globalConfig,
		state:        state,
		indexPage:    *indexTemplate,
		pipelinePage: *pipelineTemplate,
	}, nil
}

func (s *Server) Serve() error {
	mux := http.NewServeMux()

	indexHandler := routes.NewIndexHandler(s.indexPage, s.state)
	pipelineHandler := routes.NewPipelineHandler(s.pipelinePage, &s.state)
	triggerHandler := routes.NewTriggerHandler(&s.state)
	statusHandler := routes.NewStatusHandler(&s.state)
	outputHandler := routes.NewOutputHandler(&s.state)
	newRunHandler := routes.NewNewRunHandler(&s.state)
	newKillHandler := routes.NewKillHandler(&s.state)

	mux.Handle("/", indexHandler)
	mux.Handle("/pipeline/", pipelineHandler)
	mux.Handle("/trigger/", triggerHandler)
	mux.Handle("/status/", statusHandler)
	mux.Handle("/output/", outputHandler)
	mux.Handle("/new/", newRunHandler)
	mux.Handle("/kill/", newKillHandler)

	slog.Info("Start listening", "port", s.serverConfig.GetPort())
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", s.serverConfig.GetPort()), mux); err != nil {
		slog.Error("Failed to start server", "error", err)
		return err
	}

	return nil
}
