package routes

import (
	server "executrix/server/state"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type KillHandler struct {
	state server.IServerState
}

func NewKillHandler(state server.IServerState) KillHandler {
	return KillHandler{
		state: state,
	}
}

func (h KillHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to kill handler")
	slog.Debug("Request to hill handler", "request", *r)

	w.Header().Set("Content-Type", "application/json")

	name := strings.TrimPrefix(r.URL.Path, "/kill/")
	if err := h.state.Kill(name); err != nil {
		slog.Error("Could not cancel pipeline", "error", err)
		fmt.Fprint(w, `{"success": false}`) // todo error handling
	} else {
		fmt.Fprint(w, `{"success": true}`)
	}

}
