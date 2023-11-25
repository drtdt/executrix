package helper

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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
		result = append(result, filepath.Join(path, e.Name()))
	}

	return result, nil
}

func CreateIfNotExisting(path string) error {
	pathExists, err := Exists(path)
	if err != nil {
		slog.Error("Failed checking for path", "path", path)
		return err
	}

	if !pathExists {
		slog.Info("Path not found - creating...", "path", path)
		if err := os.Mkdir(path, 0644); err != nil {
			slog.Error("Failed to create directory", "path", path)
			return err
		}

		// check again
		pathExists, err = Exists(path)
		if err != nil {
			slog.Error("Failed checking for path", "path", path)
			return err
		}
		if !pathExists {
			slog.Error("Directory still not found", "path", path)
			return err
		}
	}

	return nil
}

func ReadFile(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	slog.Debug("Successfully opened", "file", path)

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func CleanUpString(s string) string {
	strings.ReplaceAll(s, `'`, `"`)
	return strings.ReplaceAll(s, `\`, `\\`)
}
