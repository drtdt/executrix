package routes

import (
	"log/slog"
	"net/http"
	"text/template"
)

type IndexHandler struct {
	page template.Template
	data any
}

func NewIndexHandler(page template.Template, data any) IndexHandler {
	return IndexHandler{
		page: page,
		data: data,
	}
}

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request to index page")
	slog.Debug("Request to index page", "request", *r)

	// todo check there's nothing after '/'

	// reload pipeline files (if pipelines no are running)

	h.page.Execute(w, h.data)
}
