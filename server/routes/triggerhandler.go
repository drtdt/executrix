package routes

import (
	"encoding/json"
	"executrix/data"
	server "executrix/server/state"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type TriggerHandler struct {
	state server.IServerState
}

func NewTriggerHandler(state server.IServerState) TriggerHandler {
	return TriggerHandler{
		state: state,
	}
}

func (h TriggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to trigger endpoint")
	slog.Debug("Request to trigger endpoint", "request", *r)

	if h.state.HasExecution() {
		// todo queueing?
		slog.Error("Already running a pipeline")
		fmt.Fprint(w, `{"started": false}`) // todo give reason
		return
	}

	w.Header().Set("Content-Type", "application/json")

	name := strings.TrimPrefix(r.URL.Path, "/trigger/")
	pipeline := h.state.PipelineFromName(name)
	if pipeline == nil {
		slog.Error("Could not find pipeline", "name", name)
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

	if err := h.state.NewExecution(pipeline, stepInfo); err != nil {
		slog.Error("Could not create new execution", "err", err)
		fmt.Fprint(w, `{"started": false}`) // todo give reason
		return
	}

	go h.state.Execute()

	fmt.Fprint(w, `{"started": true}`)
}
