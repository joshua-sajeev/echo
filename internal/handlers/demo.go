package handlers

import (
	"net/http"

	"github.com/joshu-sajeev/echo/internal/db"
)

type DemoHandler struct{}

func NewDemoHandler() *DemoHandler { return &DemoHandler{} }

// POST /demo/reset — wipes and re-seeds demo data, then redirects home.
func (h *DemoHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if !db.IsDemo() {
		http.Error(w, "not available", http.StatusForbidden)
		return
	}
	db.SeedDemo(r.Context())
	// HTMX-aware redirect
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
