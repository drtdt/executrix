package routes

import (
	server "executrix/server/state"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type NewRunHandler struct {
	state server.IServerState
}

func NewNewRunHandler(state server.IServerState) NewRunHandler {
	return NewRunHandler{
		state: state,
	}
}

func (h NewRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to new run endpoint")
	slog.Debug("Request to new run endpoint", "request", *r)

	w.Header().Set("Content-Type", "application/json")

	name := strings.TrimPrefix(r.URL.Path, "/new/")
	h.state.Reset(name)
	pipeline := h.state.PipelineFromName(name)
	if pipeline == nil {
		slog.Error("Could not find pipeline", "name", name)
		fmt.Fprint(w, `{"success": false}`) // todo error handling
		return
	}

	fmt.Fprint(w, `{"success": true}`)
}
