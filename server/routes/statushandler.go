package routes

import (
	"encoding/json"
	"executrix/data"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type StatusHandler struct {
	state data.IServerState
}

func NewStatusHandler(state data.IServerState) StatusHandler {
	return StatusHandler{
		state: state,
	}
}

func (h StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to status endpoint")
	slog.Debug("Request to status endpoint", "request", *r)

	w.Header().Set("Content-Type", "application/json")

	name := strings.TrimPrefix(r.URL.Path, "/status/")
	pipeline := h.state.PipelineFromName(name)
	if pipeline == nil {
		slog.Error("Could not find pipeline", "name", name)
		fmt.Fprint(w, `{"running": false}`) // todo error handling
		return
	}

	bytes, err := json.Marshal(pipeline.GetStepStates())
	if err != nil {
		slog.Error("Could not create status data")
		fmt.Fprint(w, `{"running": false}`) // todo error handling
		return
	}

	fmt.Fprint(w, "{"+
		`"running": `+strconv.FormatBool(h.state.IsRunning())+", "+
		`"stepStates": `+string(bytes)+
		"}")
}
