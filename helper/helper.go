package helper

import (
	"log/slog"
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FindAllFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, e := range entries {
		slog.Debug("Found file", "file", e.Name())
		result = append(result, e.Name())
	}

	return result, nil
}
