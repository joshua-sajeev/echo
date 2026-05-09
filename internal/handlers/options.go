package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/repository"
)

// OptionsHandler serves HTMX partials for dynamic form fields
type OptionsHandler struct {
	accountRepo repository.AccountRepository
	jarRepo     repository.JarRepository
}

func NewOptionsHandler(accountRepo repository.AccountRepository, jarRepo repository.JarRepository) *OptionsHandler {
	return &OptionsHandler{accountRepo: accountRepo, jarRepo: jarRepo}
}

// GET /accounts/options?type=income|expense|transfer
// Returns the right account field(s) for the transaction type
func (h *OptionsHandler) AccountFields(w http.ResponseWriter, r *http.Request) {
	txType := r.URL.Query().Get("type")

	accounts, err := h.accountRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// build <option> tags
	var opts string
	for _, a := range accounts {
		opts += fmt.Sprintf(`<option value="%d">%s</option>`, a.ID, a.Name)
	}

	var html string

	switch txType {
	case "transfer":
		html = fmt.Sprintf(`
			<div>
				<label class="block text-xs text-zinc-500 mb-2">From Account</label>
				<select name="from" class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
					%s
				</select>
			</div>
			<div>
				<label class="block text-xs text-zinc-500 mb-2">To Account</label>
				<select name="to" class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
					%s
				</select>
			</div>
		`, opts, opts)

	case "expense":
		html = fmt.Sprintf(`
			<div>
				<label class="block text-xs text-zinc-500 mb-2">From Account</label>
				<select name="account" class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
					%s
				</select>
			</div>
		`, opts)

	default: // income
		html = fmt.Sprintf(`
			<div>
				<label class="block text-xs text-zinc-500 mb-2">To Account</label>
				<select name="account" class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
					%s
				</select>
			</div>
		`, opts)
	}

	if _, err := w.Write([]byte(html)); err != nil {
		log.Println("write error:", err)
	}
}

// GET /jars/options?type=income|expense&master=true|false
// Returns jar selector (or nothing for transfer / master income)
func (h *OptionsHandler) JarFields(w http.ResponseWriter, r *http.Request) {
	txType := r.URL.Query().Get("type")
	master := r.URL.Query().Get("master")

	// transfer → no jar
	// master income → no jar
	if txType == "transfer" || master == "true" {
		w.Write([]byte(""))
		return
	}

	jars, err := h.jarRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var opts string
	for _, j := range jars {
		opts += fmt.Sprintf(`<option value="%d">%s</option>`, j.ID, j.Name)
	}

	required := ""
	if txType == "expense" || txType == "income" {
		required = "required"
	}

	html := fmt.Sprintf(`
		<div>
			<label class="block text-xs text-zinc-500 mb-2">Jar %s</label>
			<select name="jar_id" %s class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
				<option value="">Select jar</option>
				%s
			</select>
		</div>
	`, func() string {
		if required == "required" {
			return "*"
		}
		return ""
	}(), required, opts)

	if _, err := w.Write([]byte(html)); err != nil {
		log.Println("write error:", err)
	}
}
