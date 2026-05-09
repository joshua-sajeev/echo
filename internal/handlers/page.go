package handlers

import (
	"html/template"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/allocation"
	"github.com/joshu-sajeev/echo/internal/repository"
)

type PageHandler struct {
	tmpl    *template.Template
	txRepo  repository.TransactionRepository
	jarRepo repository.JarRepository
}

func NewPageHandler(
	tmpl *template.Template,
	txRepo repository.TransactionRepository,
	jarRepo repository.JarRepository,
) *PageHandler {
	return &PageHandler{tmpl: tmpl, txRepo: txRepo, jarRepo: jarRepo}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	stats, err := h.txRepo.Stats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jars, spent, err := h.jarRepo.ListWithSpend(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allocs := allocation.Calculate(stats.MonthlyMasterIncome, jars, spent)
	summary := allocation.Summarise(allocs)

	// build per-jar data for the template
	type JarData struct {
		Name      string
		Spent     float64
		Allocated float64
		Percent   int
		IsCap     bool
	}

	var jarData []JarData
	for _, a := range allocs {
		pct := 0
		if a.Allocated > 0 {
			pct = int(a.Spent / a.Allocated * 100)
			if pct > 100 {
				pct = 100
			}
		}
		jarData = append(jarData, JarData{
			Name:      a.Name,
			Spent:     a.Spent,
			Allocated: a.Allocated,
			Percent:   pct,
			IsCap:     a.IsCap,
		})
	}

	data := map[string]any{
		"TotalBalance":    stats.TotalBalance,
		"MonthlyExpenses": stats.MonthlyExpenses,
		"Savings":         stats.Savings,
		"Jars":            jarData,
		"MonthlyBudget":   stats.MonthlyMasterIncome,
		"SpentThisMonth":  summary.TotalSpent,
		"BudgetPercent": func() int {
			if stats.MonthlyMasterIncome <= 0 {
				return 0
			}
			p := int(summary.TotalSpent / stats.MonthlyMasterIncome * 100)
			if p > 100 {
				return 100
			}
			return p
		}(),
	}

	if err := h.tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
