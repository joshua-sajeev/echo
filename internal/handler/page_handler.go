package handler

import (
	"html/template"
	"net/http"
)

type PageHandler struct {
	tmpl *template.Template
}

func NewPageHandler(t *template.Template) *PageHandler {
	return &PageHandler{tmpl: t}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
