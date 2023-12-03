package config

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"executrix/constants"
	"executrix/helper"
)

type ServerConfig struct {
	configDir   string
	pipelineDir string
	port        uint16
}

func ServerConfigFromJson(configDir string) (ServerConfig, error) {
	serverConfigPath := filepath.Join(configDir, constants.SERVER_CONFIG_FILE)
	pipelineDir := filepath.Join(configDir, constants.PIPELINE_DIR_NAME)

	pathExists, err := helper.Exists(serverConfigPath)
	if err != nil {
		slog.Error("Failed checking for server config path", "path", serverConfigPath)
		return ServerConfig{}, err
	}

	if !pathExists {
		slog.Info("Server config not found - creating file with default settings", "path", serverConfigPath)
		if err := createDefaultServerConfig(serverConfigPath); err != nil {
			slog.Error("Failed creating default server config", "path", serverConfigPath)
			return ServerConfig{}, err
		}
	}

	bytes, err := helper.ReadFile(serverConfigPath)
	if err != nil {
		return ServerConfig{}, err
	}

	var p map[string]interface{}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return ServerConfig{}, err
	}

	slog.Debug("Successfully unmarshalled file content", "content", p)

	var config ServerConfig
	config.configDir = configDir
	config.pipelineDir = pipelineDir

	// reading server port from config
	val, ok := p["port"].(string)
	if !ok {
		return ServerConfig{}, errors.New("could not read server port from config")
	} else {
		slog.Debug("Read server port", "port", val)
	}

	if port, err := strconv.ParseUint(val, 10, 16); err != nil {
		return ServerConfig{}, errors.New("port has wrong format")
	} else {
		config.port = uint16(port)
	}

	return config, nil
}

func (s ServerConfig) GetPort() uint16 {
	return s.port
}

func (s ServerConfig) GetConfigDir() string {
	return s.configDir
}

func (s ServerConfig) GetPipelineDir() string {
	return s.pipelineDir
}

func createDefaultServerConfig(path string) error {
	// todo: use marshalling
	data := []byte("{\n" +
		"\t\"port\": \"8111\"\n" +
		"}\n")
	return os.WriteFile(path, data, 0644)
}
