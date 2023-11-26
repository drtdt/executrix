package routes

import (
	"executrix/data"
	"log/slog"
	"net/http"
	"strings"
	"text/template"
)

type PipelineHandler struct {
	page      template.Template
	pipelines data.IPipelineContainer
}

func NewPipelineHandler(page template.Template, pipelines data.IPipelineContainer) PipelineHandler {
	return PipelineHandler{
		page:      page,
		pipelines: pipelines,
	}
}

func (h PipelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request pipeline page")
	slog.Debug("Request pipeline page", "request", *r)

	name := strings.TrimPrefix(r.URL.Path, "/pipeline/")

	if p := h.pipelines.PipelineFromName(name); p == nil {
		slog.Error("Could not find pipeline", "name", name)
		// todo
	} else {
		h.page.Execute(w, *p)
	}
}
