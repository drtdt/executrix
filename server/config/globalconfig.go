package config

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"executrix/helper"
)

type GlobalConfig struct {
	vars map[string]string
}

func (cfg GlobalConfig) ResolveVar(name string) (string, error) {
	if val, ok := cfg.vars[name]; !ok {
		return "", errors.New("var is not defined")
	} else {
		return val, nil
	}
}

func (cfg GlobalConfig) GetVars() map[string]string {
	return cfg.vars
}

func GlobalConfigFromJson(path string) (GlobalConfig, error) {
	pathExists, err := helper.Exists(path)
	if err != nil {
		slog.Error("Failed checking for global config path", "path", path)
		return GlobalConfig{}, err
	}

	if !pathExists {
		slog.Info("Global config not found - creating file with default settings", "path", path)
		if err := createDefaultGlobalConfig(path); err != nil {
			slog.Error("Failed creating default server config", "path", path)
			return GlobalConfig{}, err
		}
	}

	bytes, err := helper.ReadFile(path)
	if err != nil {
		return GlobalConfig{}, err
	}

	var p map[string]interface{}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		return GlobalConfig{}, err
	}

	slog.Debug("Successfully unmarshalled file content", "content", p)

	val, ok := p["vars"].([]interface{})
	if !ok {
		return GlobalConfig{}, errors.New("error reading global vars")
	}

	slog.Info("Read global vars", "vars", val)
	vars := map[string]string{}
	for _, elem := range val {
		slog.Info("Read global var", "var", elem)

		pair, ok := elem.(map[string]interface{})
		if !ok {
			return GlobalConfig{}, errors.New("unexpected type for pair")
		}

		name, ok := pair["name"].(string)
		if !ok {
			return GlobalConfig{}, errors.New("unexpected type for name")
		}

		// check if name has already been used
		if _, ok := vars[name]; ok {
			return GlobalConfig{}, errors.New("found non-unique name in vars")
		}

		value, ok := pair["value"].(string)
		if !ok {
			return GlobalConfig{}, errors.New("unexpected type for value")
		}

		vars[name] = value
	}

	return GlobalConfig{
		vars: vars,
	}, nil
}

func createDefaultGlobalConfig(path string) error {
	// todo: use marshalling
	data := []byte("{\n" +
		"\t\"vars\": []\n" +
		"}\n")
	return os.WriteFile(path, data, 0644)
}
