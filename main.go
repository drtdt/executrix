package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"executrix/helper"
)

type ServerConfig struct {
	configDir   string
	pipelineDir string
	pages       map[string]*template.Template
}

var serverConfig ServerConfig

type Pipeline struct {
	Name string
}

type IndexPageData struct {
	Pipelines []Pipeline
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
		// todo read json file

		indexPageData.Pipelines = append(indexPageData.Pipelines, Pipeline{
			Name: file,
		})
	}

	return nil
}

func createIfNotExisting(path string) {
	pathExists, err := helper.Exists(path)
	if err != nil {
		slog.Error("Failed checking for path", "path", path)
		os.Exit(-1)
	}

	if !pathExists {
		slog.Info("Path not found - creating...", "path", path)
		if err := os.Mkdir(path, 0644); err != nil {
			slog.Error("Failed to create directory", "path", path)
			os.Exit(-1)
		}

		// check again
		pathExists, err = helper.Exists(path)
		if err != nil {
			slog.Error("Failed checking for path", "path", path)
			os.Exit(-1)
		}
		if !pathExists {
			slog.Error("Directory still not found", "path", path)
			os.Exit(-1)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// reload pipeline files
	if err := reloadPipelines(); err != nil {
		slog.Error("Error while reloading pipeline configs", "err", err.Error())
		// todo
	}

	serverConfig.pages["index"].Execute(w, indexPageData)
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
	createIfNotExisting(configDir)
	slog.Info("Found config directory", "path", configDir)

	pipelineDir := filepath.Join(configDir, PIPELINE_DIR_NAME)
	createIfNotExisting(pipelineDir)
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

	serverConfig.pages["index"] = indexTemplate

	http.HandleFunc("/", handler)

	slog.Info("Start listening", "port", PORT)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(-1)
	}
}
