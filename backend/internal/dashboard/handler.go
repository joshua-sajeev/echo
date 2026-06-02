package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/transactions"
)

type Handler struct {
	AccountService     *accounts.AccountService
	JarService         *jars.JarService
	TransactionService *transactions.TransactionService
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

	// TODO: Use Recent Transactions later
	txs, err := h.TransactionService.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"accounts":     accounts,
		"jars":         jars,
		"transactions": txs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
