package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/goals"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/transactions"
)

type Handler struct {
	AccountService     *accounts.AccountService
	JarService         *jars.JarService
	TransactionService *transactions.TransactionService
	GoalService        *goals.GoalService
}

func (h *Handler) GetDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	accounts, err := h.AccountService.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jars, err := h.JarService.ListJarAllocations(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	txs, err := h.TransactionService.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	goalsWithProgress, err := h.GoalService.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"accounts":     accounts,
		"jars":         jars,
		"transactions": txs,
		"goals":        goalsWithProgress, // ← ADD THIS
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
