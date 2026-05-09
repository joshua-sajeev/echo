// Package handlers
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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

	w.Header().Set("HX-Trigger", `{"accountCreated": true}`)
	w.Write([]byte(`
		<div class="bg-green-500/10 border border-green-500/30 text-green-400 p-3 rounded-xl text-sm">
			Account created
		</div>
	`))
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
		w.Write([]byte(accountRow(a, false)))
	}
}

func (h *AccountHandler) ListArchived(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.ListArchivedWithBalances(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(accounts) == 0 {
		w.Write([]byte(`<p class="text-zinc-600 text-xs">No archived accounts</p>`))
		return
	}

	for _, a := range accounts {
		w.Write([]byte(accountRow(a, true)))
	}
}

func (h *AccountHandler) Rename(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Rename(r.Context(), id, name); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			http.Error(w, fmt.Sprintf(`"%s" already exists`, name), http.StatusConflict)
			return
		}
		http.Error(w, "rename failed", http.StatusInternalServerError)
		return
	}

	// return the updated row so HTMX can swap it in place
	w.Header().Set("HX-Trigger", `{"accountCreated": true}`)
	w.Write([]byte(`<div></div>`)) // list reloads via event
}

func (h *AccountHandler) Archive(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Archive(r.Context(), id); err != nil {
		http.Error(w, "archive failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"accountCreated": true}`)
	w.Write([]byte(`<div></div>`))
}

func (h *AccountHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Unarchive(r.Context(), id); err != nil {
		http.Error(w, "unarchive failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"accountCreated": true}`)
	w.Write([]byte(`<div></div>`))
}

// accountRow renders a single account row HTML
// archived=true gives it a muted archived style with unarchive button
func accountRow(a models.AccountWithBalance, archived bool) string {
	balanceColor := "text-zinc-300"
	if a.Balance > 0 {
		balanceColor = "text-emerald-400"
	} else if a.Balance < 0 {
		balanceColor = "text-red-400"
	}

	if archived {
		return fmt.Sprintf(`
			<div class="flex justify-between items-center bg-zinc-950/50 border border-zinc-800/50 rounded-lg p-4 opacity-50"
			     id="account-%d">
				<p class="font-medium text-zinc-500 line-through">%s</p>
				<div class="flex items-center gap-3">
					<p class="text-sm %s">₹%.2f</p>
					<button
						hx-patch="/accounts/%d/unarchive"
						hx-swap="none"
						class="text-xs text-zinc-500 hover:text-zinc-300 px-2 py-1 rounded border border-zinc-700"
					>Restore</button>
				</div>
			</div>
		`, a.ID, a.Name, balanceColor, a.Balance, a.ID)
	}

	return fmt.Sprintf(`
		<div class="account-row flex justify-between items-center bg-zinc-950 border border-zinc-800 rounded-lg p-4 select-none"
		     id="account-%d"
		     data-id="%d"
		     data-name="%s">
			<div class="name-display flex-1 font-medium">%s</div>
			<div class="rename-input hidden flex-1">
				<input
					class="w-full bg-zinc-900 border border-zinc-700 rounded px-2 py-1 text-sm outline-none focus:border-zinc-500"
					value="%s"
					onkeydown="if(event.key==='Enter') submitRename(this, %d)"
					onblur="cancelRename(this)"
				/>
			</div>
			<div class="flex items-center gap-3">
				<p class="text-sm font-semibold %s">₹%.2f</p>
			</div>
		</div>
	`, a.ID, a.ID, a.Name, a.Name, a.Name, a.ID, balanceColor, a.Balance)
}

func writeError(w http.ResponseWriter, msg string) {
	fmt.Fprintf(w, `
		<div class="bg-red-500/10 border border-red-500/30 text-red-400 p-3 rounded-xl text-sm">
			%s
		</div>
	`, msg)
}

func parseID(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}
