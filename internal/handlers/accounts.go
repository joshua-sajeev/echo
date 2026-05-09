// Package handlers
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/joshu-sajeev/echo/internal/models"
	"github.com/joshu-sajeev/echo/internal/repository"
)

type AccountHandler struct {
	repo   repository.AccountRepository
	txRepo repository.TransactionRepository
}

func NewAccountHandler(repo repository.AccountRepository, txRepo repository.TransactionRepository) *AccountHandler {
	return &AccountHandler{repo: repo, txRepo: txRepo}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	opening := r.FormValue("opening_balance")

	if name == "" {
		writeError(w, "Account name is required")
		return
	}

	accountID, err := h.repo.Create(r.Context(), name)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			writeError(w, fmt.Sprintf(`Account "%s" already exists`, name))
			return
		}
		writeError(w, "Failed to create account")
		log.Println("create account error:", err)
		return
	}

	if opening != "" && opening != "0" {
		var amt float64
		if _, err := fmt.Sscanf(opening, "%f", &amt); err != nil {
			writeError(w, "Invalid opening balance")
			return
		}
		tx := models.Transaction{
			Type:        "income",
			Amount:      amt,
			Name:        "Opening Balance",
			ToAccountID: &accountID,
		}
		if err := h.txRepo.Create(r.Context(), tx); err != nil {
			writeError(w, "Account created but failed to record opening balance")
			log.Println("opening balance error:", err)
			return
		}
	}

	// tell HTMX to reload the account list, then show inline success + collapse form
	w.Header().Set("HX-Trigger", `{"accountCreated": true}`)
	if _, err := w.Write([]byte(`
		<div class="bg-green-500/10 border border-green-500/30 text-green-400 p-3 rounded-xl text-sm">
			Account created
		</div>
	`)); err != nil {
		log.Println("write error:", err)
	}
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.ListWithBalances(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(accounts) == 0 {
		w.Write([]byte(`<p class="text-zinc-600 text-xs">No accounts yet</p>`))
		return
	}

	for _, a := range accounts {
		balanceColor := "text-zinc-300"
		if a.Balance > 0 {
			balanceColor = "text-emerald-400"
		} else if a.Balance < 0 {
			balanceColor = "text-red-400"
		}

		if _, err := fmt.Fprintf(w, `
			<div class="flex justify-between items-center bg-zinc-950 border border-zinc-800 rounded-lg p-4">
				<p class="font-medium">%s</p>
				<p class="text-sm font-semibold %s">₹%.2f</p>
			</div>
		`, a.Name, balanceColor, a.Balance); err != nil {
			log.Println("write error:", err)
		}
	}
}

func writeError(w http.ResponseWriter, msg string) {
	fmt.Fprintf(w, `
		<div class="bg-red-500/10 border border-red-500/30 text-red-400 p-3 rounded-xl text-sm">
			%s
		</div>
	`, msg)
}
