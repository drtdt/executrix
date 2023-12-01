package config

import (
	"encoding/json"
	"errors"
	"executrix/helper"
	"log/slog"
	"path/filepath"
	"strconv"
)

type ServerConfig struct {
	configDir   string
	pipelineDir string
	port        uint16
}

func FromJson(configDir string, pipelineDir string) (ServerConfig, error) {
	bytes, err := helper.ReadFile(filepath.Join(configDir, "server.json"))
	if err != nil {
		return ServerConfig{}, err
	}

	var p map[string]interface{}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return ServerConfig{}, err
	}

	slog.Info("Successfully unmarshalled file content", "content", p)

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

func (s ServerConfig) GetPipelineDir() string {
	return s.pipelineDir
}
