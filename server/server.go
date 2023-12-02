package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"executrix/server/config"
	"executrix/server/routes"
	"executrix/server/state"
)

type Server struct {
	config       config.ServerConfig
	state        state.ServerState
	indexPage    template.Template
	pipelinePage template.Template
}

func NewServer(config config.ServerConfig) (Server, error) {

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

	state, err := state.NewServerState(config.GetPipelineDir())
	if err != nil {
		slog.Error("Failed to read pipeline configs", "error", err)
		return Server{}, err
	}

	return Server{
		config:       config,
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

	mux.Handle("/", indexHandler)
	mux.Handle("/pipeline/", pipelineHandler)
	mux.Handle("/trigger/", triggerHandler)
	mux.Handle("/status/", statusHandler)
	mux.Handle("/output/", outputHandler)

	slog.Info("Start listening", "port", s.config.GetPort())
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", s.config.GetPort()), mux); err != nil {
		slog.Error("Failed to start server", "error", err)
		return err
	}

	return nil
}
