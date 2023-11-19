package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"executrix/data"
	"executrix/executrix"
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

func (d IndexPageData) IndexFromName(name string) int {
	return slices.IndexFunc(d.Pipelines, func(p data.Pipeline) bool { return p.Name == name })
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

	// todo check there's nothing after '/'

	// todo no reload if pipelines are running!

	// reload pipeline files
	if err := reloadPipelines(); err != nil {
		slog.Error("Error while reloading pipeline configs", "err", err)
		// todo
	}

	serverConfig.pages["index"].Execute(w, indexPageData)
}

func pipelineHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request pipeline page")
	slog.Debug("Request pipeline page", "request", *r)

	id := strings.TrimPrefix(r.URL.Path, "/pipeline/")
	if idx := indexPageData.IndexFromName(id); idx < 0 {
		slog.Error("Could not find pipeline", "name", id)
		// todo
	} else {
		serverConfig.pages["pipeline"].Execute(w, indexPageData.Pipelines[idx])
	}
}

func triggerHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to trigger endpoint")
	slog.Debug("Request to trigger endpoint", "request", *r)

	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/trigger/")
	idx := indexPageData.IndexFromName(id)
	if idx < 0 {
		slog.Error("Could not find pipeline", "name", id)
		fmt.Fprint(w, `{"started": false}`) // todo give reason
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Could not read body from request", "err", err)
		fmt.Fprint(w, `{"started": false}`) // todo give reason
		return
	}

	slog.Debug("recieved body", "body", body)

	var stepInfo []data.StepInfo
	if err = json.Unmarshal(body, &stepInfo); err != nil {
		slog.Error("Could not unmarshall body from request", "err", err)
		fmt.Fprint(w, `{"started": false}`) // todo give reason
		return
	}

	slog.Debug("parsed body", "body", stepInfo)

	indexPageData.Pipelines[idx].IsRunning = true

	go executrix.ExecutePipeline(&indexPageData.Pipelines[idx], stepInfo)

	fmt.Fprint(w, `{"started": true}`)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to status endpoint")
	slog.Debug("Request to status endpoint", "request", *r)

	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/status/")
	idx := indexPageData.IndexFromName(id)
	if idx < 0 {
		slog.Error("Could not find pipeline", "name", id)
		fmt.Fprint(w, `{"running": false}`) // todo error handling
		return
	}

	bytes, err := json.Marshal(indexPageData.Pipelines[idx].GetStepStates())
	if err != nil {
		slog.Error("Could not create status data")
		fmt.Fprint(w, `{"running": false}`) // todo error handling
		return
	}

	fmt.Fprint(w, "{"+
		`"running": `+strconv.FormatBool(indexPageData.Pipelines[idx].IsRunning)+", "+
		`"stepStates": `+string(bytes)+
		"}")
}

func main() {
	const PORT uint16 = 8080
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

	serverConfig = ServerConfig{
		configDir:   configDir,
		pipelineDir: pipelineDir,
		pages:       make(map[string]*template.Template),
	}

	indexTemplate, err := template.ParseFiles("html/index.html")
	if err != nil {
		slog.Error("Failed to parse index.html", "error", err)
		os.Exit(-1)
	}

	pipelineTemplate, err := template.ParseFiles("html/pipeline.html")
	if err != nil {
		slog.Error("Failed to parse pipeline.html", "error", err)
		os.Exit(-1)
	}

	serverConfig.pages["index"] = indexTemplate
	serverConfig.pages["pipeline"] = pipelineTemplate

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pipeline/", pipelineHandler)
	http.HandleFunc("/trigger/", triggerHandler)
	http.HandleFunc("/status/", statusHandler)

	slog.Info("Start listening", "port", PORT)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(-1)
	}
}
