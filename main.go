package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"executrix/data"
	"executrix/helper"
)

type ServerConfig struct {
	configDir   string
	pipelineDir string
	pages       map[string]*template.Template
}

var serverConfig ServerConfig

type IndexPageData struct {
	Pipelines []data.Pipeline
}

var indexPageData IndexPageData

func reloadPipelines() error {
	indexPageData.Pipelines = nil
	slog.Debug("Cleared piplines before reloading")

	result, err := helper.FindAllFiles(serverConfig.pipelineDir)
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

		indexPageData.Pipelines = append(indexPageData.Pipelines, pipeline)
	}

	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to index page")
	slog.Debug("Request to index page", "request", *r)
	// reload pipeline files
	if err := reloadPipelines(); err != nil {
		slog.Error("Error while reloading pipeline configs", "err", err.Error())
		// todo
	}

	serverConfig.pages["index"].Execute(w, indexPageData)
}

func pipelineHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request pipeline page")

	id := strings.TrimPrefix(r.URL.Path, "/pipeline/")
	slog.Info("Found Pipeline ID", "id", id)

	if idx := slices.IndexFunc(indexPageData.Pipelines, func(p data.Pipeline) bool { return p.Name == id }); idx < 0 {
		slog.Error("Could not find pipeline", "name", id)
		// todo
	} else {
		serverConfig.pages["pipeline"].Execute(w, indexPageData.Pipelines[idx])
	}
}

func main() {
	const PORT uint16 = 8080
	const CONFIG_DIR_NAME = "Executrix"
	const PIPELINE_DIR_NAME = "pipelines"

	configBaseDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Failed to determine user default config location", "error", err.Error())
		os.Exit(-1)
	}

	slog.Info("Found default config location", "path", configBaseDir)

	configDir := filepath.Join(configBaseDir, CONFIG_DIR_NAME)
	if err = helper.CreateIfNotExisting(configDir); err != nil {
		slog.Error("Error while checking for config path", "error", err.Error())
		os.Exit(-1)
	}
	slog.Info("Found config directory", "path", configDir)

	pipelineDir := filepath.Join(configDir, PIPELINE_DIR_NAME)
	if err = helper.CreateIfNotExisting(pipelineDir); err != nil {
		slog.Error("Error while checking for pipeline path", "error", err.Error())
		os.Exit(-1)
	}
	slog.Info("Found pipeline directory", "path", pipelineDir)

	serverConfig = ServerConfig{
		configDir:   configDir,
		pipelineDir: pipelineDir,
		pages:       make(map[string]*template.Template),
	}

	indexTemplate, err := template.ParseFiles("html/index.html")
	if err != nil {
		slog.Error("Failed to parse index.html", "error", err.Error())
		os.Exit(-1)
	}

	pipelineTemplate, err := template.ParseFiles("html/pipeline.html")
	if err != nil {
		slog.Error("Failed to parse pipeline.html", "error", err.Error())
		os.Exit(-1)
	}

	serverConfig.pages["index"] = indexTemplate
	serverConfig.pages["pipeline"] = pipelineTemplate

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pipeline/", pipelineHandler)

	slog.Info("Start listening", "port", PORT)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(-1)
	}
}
