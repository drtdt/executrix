package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"executrix/helper"
	"executrix/server"
)

func main() {
	const CONFIG_DIR_NAME = "Executrix"
	const PIPELINE_DIR_NAME = "pipelines"

	configBaseDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Failed to determine user default config location", "error", err)
		os.Exit(-1)
	}

	slog.Info("Found default config location", "path", configBaseDir)

	configDir := filepath.Join(configBaseDir, CONFIG_DIR_NAME)
	if err = helper.CreateIfNotExisting(configDir); err != nil {
		slog.Error("Error while checking for config path", "error", err)
		os.Exit(-1)
	}
	slog.Info("Found config directory", "path", configDir)

	pipelineDir := filepath.Join(configDir, PIPELINE_DIR_NAME)
	if err = helper.CreateIfNotExisting(pipelineDir); err != nil {
		slog.Error("Error while checking for pipeline path", "error", err)
		os.Exit(-1)
	}
	slog.Info("Found pipeline directory", "path", pipelineDir)

	server, err := server.NewServer(configDir, pipelineDir)
	if err != nil {
		slog.Error("Error while configuring server", "error", err)
		os.Exit(-1)
	}

	server.Serve()
}
