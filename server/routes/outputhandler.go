package routes

import (
	server "executrix/server/state"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type OutputHandler struct {
	state server.IServerState
}

func NewOutputHandler(state server.IServerState) OutputHandler {
	return OutputHandler{
		state: state,
	}
}

func (h OutputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to output endpoint")
	slog.Debug("Request to output endpoint", "request", *r)

	w.Header().Set("Content-Type", "application/json")

	name := strings.TrimPrefix(r.URL.Path, "/output/")
	output, err := h.state.StepOutput(name)
	if err != nil {
		slog.Error("Error retrieving step output", "step", name)
		fmt.Fprint(w, `{"text": "Error retrieving step output!"}`)
		return
	}

	slog.Debug("Sending output", "step", name, "text", output)
	fmt.Fprint(w, `{"text": "`+output+`"}`)
}
