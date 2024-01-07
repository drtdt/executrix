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

	if err := h.state.Reset(name); err != nil {
		slog.Error("Could not reset pipeline", "error", err)
		fmt.Fprint(w, `{"success": false}`) // todo error handling
	} else {
		fmt.Fprint(w, `{"success": true}`)
	}
}
