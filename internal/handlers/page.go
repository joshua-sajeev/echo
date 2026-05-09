package handlers

import (
	"html/template"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/repository"
)

type PageHandler struct {
	tmpl    *template.Template
	txRepo  repository.TransactionRepository
	jarRepo repository.JarRepository
}

func NewPageHandler(tmpl *template.Template, txRepo repository.TransactionRepository, jarRepo repository.JarRepository) *PageHandler {
	return &PageHandler{tmpl: tmpl, txRepo: txRepo, jarRepo: jarRepo}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	stats, err := h.txRepo.Stats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jars, err := h.jarRepo.ListWithMonthlySpend(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build per-jar budget: each jar's allocation = (target_amount / sum_of_targets) * master_income
	// If jar has an explicit target_amount, use it directly as the budget cap.
	// Master income is the monthly budget ceiling.
	budget := stats.MonthlyMasterIncome

	type JarData struct {
		Name    string
		Spent   float64
		Budget  float64
		Percent int
	}

	var jarData []JarData
	for _, j := range jars {
		b := j.TargetAmount // jar's own fixed target
		if b <= 0 {
			b = budget // fallback: whole budget if no target set
		}
		pct := 0
		if b > 0 {
			pct = int((j.MonthlySpend / b) * 100)
			if pct > 100 {
				pct = 100
			}
		}
		jarData = append(jarData, JarData{
			Name:    j.Name,
			Spent:   j.MonthlySpend,
			Budget:  b,
			Percent: pct,
		})
	}

	totalSpent := stats.MonthlyExpenses
	budgetPct := 0
	if budget > 0 {
		budgetPct = int((totalSpent / budget) * 100)
		if budgetPct > 100 {
			budgetPct = 100
		}
	}

	data := map[string]any{
		"TotalBalance":    stats.TotalBalance,
		"MonthlyExpenses": stats.MonthlyExpenses,
		"Savings":         stats.Savings,
		"Jars":            jarData,
		"MonthlyBudget":   budget,
		"SpentThisMonth":  totalSpent,
		"BudgetPercent":   budgetPct,
	}

	if err := h.tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
