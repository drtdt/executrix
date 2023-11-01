package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type Pipeline struct {
	Name string
}

type IndexPageData struct {
	Pipelines []Pipeline
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createIfNotExisting(path string) {
	pathExists, err := exists(path)
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
		pathExists, err = exists(path)
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
	// TODO this should probably go somewhere else
	indexTemplate, err := template.ParseFiles("html/index.html")
	if err != nil {
		slog.Error("Failed to parse index.html", "error", err.Error())
		os.Exit(-1)
	}

	data := IndexPageData{
		[]Pipeline{
			//{"Test1"},
			//{"Test2"},
		},
	}

	indexTemplate.Execute(w, data)
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

	http.HandleFunc("/", handler)

	slog.Info("Start listening", "port", PORT)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(-1)
	}
}
