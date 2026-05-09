package handlers

import (
	"html/template"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/repository"
)

type PageHandler struct {
	tmpl   *template.Template
	txRepo repository.TransactionRepository
}

func NewPageHandler(tmpl *template.Template, txRepo repository.TransactionRepository) *PageHandler {
	return &PageHandler{tmpl: tmpl, txRepo: txRepo}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	stats, err := h.txRepo.Stats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"TotalBalance":    stats.TotalBalance,
		"MonthlyExpenses": stats.MonthlyExpenses,
		"Savings":         stats.Savings,
	}

	if err := h.tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
