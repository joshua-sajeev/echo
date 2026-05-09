package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joshu-sajeev/echo/internal/models"
	"github.com/joshu-sajeev/echo/internal/repository"
)

type TransactionHandler struct {
	repo        repository.TransactionRepository
	accountRepo repository.AccountRepository
	jarRepo     repository.JarRepository
	tmpl        *template.Template
}

func NewTransactionHandler(
	repo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	jarRepo repository.JarRepository,
	tmpl *template.Template,
) *TransactionHandler {
	return &TransactionHandler{repo: repo, accountRepo: accountRepo, jarRepo: jarRepo, tmpl: tmpl}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	txType := r.FormValue("type")
	name := r.FormValue("name")
	amount := r.FormValue("amount")
	dateStr := r.FormValue("date")

	if txType == "" || amount == "" || name == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	var amt float64
	if _, err := fmt.Sscanf(amount, "%f", &amt); err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	date := time.Now()
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	tx := models.Transaction{
		Type:   txType,
		Amount: amt,
		Name:   name,
		Date:   date,
	}

	switch txType {
	case "income":
		toID, err := parseID(r.FormValue("account"))
		if err != nil {
			http.Error(w, "missing to account", http.StatusBadRequest)
			return
		}
		tx.ToAccountID = &toID
		tx.IsMasterIncome = r.FormValue("master_income") == "true"
		if !tx.IsMasterIncome {
			jarID, err := parseID(r.FormValue("jar_id"))
			if err != nil {
				http.Error(w, "jar required for non-master income", http.StatusBadRequest)
				return
			}
			tx.JarID = &jarID
		}
	case "expense":
		fromID, err := parseID(r.FormValue("account"))
		if err != nil {
			http.Error(w, "missing from account", http.StatusBadRequest)
			return
		}
		tx.FromAccountID = &fromID
		jarID, err := parseID(r.FormValue("jar_id"))
		if err != nil {
			http.Error(w, "jar required for expense", http.StatusBadRequest)
			return
		}
		tx.JarID = &jarID
	case "transfer":
		fromID, err := parseID(r.FormValue("from"))
		if err != nil {
			http.Error(w, "missing from account", http.StatusBadRequest)
			return
		}
		toID, err := parseID(r.FormValue("to"))
		if err != nil {
			http.Error(w, "missing to account", http.StatusBadRequest)
			return
		}
		tx.FromAccountID = &fromID
		tx.ToAccountID = &toID
	default:
		http.Error(w, "invalid transaction type", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), tx); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(`
		<div class="bg-zinc-800 p-4 rounded-xl text-green-400">
		Transaction added successfully
		</div>
		`)); err != nil {
		log.Println("write error:", err)
	}
}

// renderTxRow renders a single transaction as an HTML row with action menu support.
// Used by both recent list and view-all list.
func renderTxRow(tx repository.TransactionRow) string {
	var amountStr, amountColor string
	switch tx.Type {
	case "income":
		amountStr = fmt.Sprintf("+ ₹%.2f", tx.Amount)
		amountColor = "text-emerald-400"
	case "expense":
		amountStr = fmt.Sprintf("- ₹%.2f", tx.Amount)
		amountColor = "text-red-400"
	case "transfer":
		amountStr = fmt.Sprintf("₹%.2f", tx.Amount)
		amountColor = "text-sky-400"
	}

	d := tx.Date
	if d.IsZero() || d.Year() < 2000 {
		d = tx.CreatedAt
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	txDay := d.Truncate(24 * time.Hour)

	var dateStr string
	switch {
	case txDay.Equal(today):
		dateStr = "Today"
	case txDay.Equal(yesterday):
		dateStr = "Yesterday"
	default:
		dateStr = d.Format("02 Jan 2006")
	}

	label := ""
	switch {
	case tx.JarName != "":
		label = tx.JarName
	case tx.Type == "income":
		label = "Income"
	case tx.Type == "expense":
		label = "Expense"
	case tx.Type == "transfer":
		label = "Transfer"
	}

	account := ""
	switch tx.Type {
	case "income":
		account = tx.ToAccountName
	case "expense":
		account = tx.FromAccountName
	case "transfer":
		if tx.FromAccountName != "" && tx.ToAccountName != "" {
			account = tx.FromAccountName + " → " + tx.ToAccountName
		}
	}

	subtitle := label + " · " + dateStr

	return fmt.Sprintf(`
<div class="tx-row flex items-center justify-between py-3 cursor-pointer select-none group"
     data-id="%d"
     data-name="%s"
     data-amount="%.2f"
     data-type="%s"
     data-date="%s"
     onclick="openTxMenu(event, this)">
  <div class="min-w-0">
    <p class="text-sm font-medium text-zinc-100 truncate">%s</p>
    <p class="text-xs text-zinc-500 mt-0.5 truncate">%s</p>
  </div>
  <div class="text-right shrink-0 ml-4">
    <p class="text-sm font-semibold %s">%s</p>
    <p class="text-xs text-zinc-600 mt-0.5 truncate">%s</p>
  </div>
</div>
`, tx.ID, tx.Name, tx.Amount, tx.Type, d.Format("2006-01-02"),
		tx.Name, subtitle, amountColor, amountStr, account)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	txs, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if len(txs) == 0 {
		w.Write([]byte(`<p class="text-zinc-600 text-xs py-4">No transactions yet</p>`))
		return
	}

	var sb strings.Builder
	for _, tx := range txs {
		sb.WriteString(renderTxRow(tx))
	}
	w.Write([]byte(sb.String()))
}

func (h *TransactionHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var amtMin, amtMax float64
	fmt.Sscanf(q.Get("amt_min"), "%f", &amtMin)
	fmt.Sscanf(q.Get("amt_max"), "%f", &amtMax)

	var accountID, jarID int64
	if v := q.Get("account_id"); v != "" {
		fmt.Sscanf(v, "%d", &accountID)
	}
	if v := q.Get("jar_id"); v != "" {
		fmt.Sscanf(v, "%d", &jarID)
	}

	f := repository.TxFilters{
		Type:      q.Get("type"),
		Search:    q.Get("search"),
		AccountID: accountID,
		JarID:     jarID,
		AmountMin: amtMin,
		AmountMax: amtMax,
		DateFrom:  q.Get("date_from"),
		DateTo:    q.Get("date_to"),
	}

	txs, err := h.repo.ListAll(r.Context(), f)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if len(txs) == 0 {
		w.Write([]byte(`<p class="text-zinc-600 text-xs py-4">No transactions found</p>`))
		return
	}

	var sb strings.Builder
	for _, tx := range txs {
		sb.WriteString(renderTxRow(tx))
	}
	w.Write([]byte(sb.String()))
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	// Signal both lists to refresh
	w.Header().Set("HX-Trigger", `{"txChanged": true}`)
	w.WriteHeader(http.StatusOK)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	existing, err := h.repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	amountStr := r.FormValue("amount")
	dateStr := r.FormValue("date")

	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	var amt float64
	if _, err := fmt.Sscanf(amountStr, "%f", &amt); err != nil || amt <= 0 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	date := existing.Date
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	tx := models.Transaction{
		ID:            id,
		Name:          name,
		Amount:        amt,
		Date:          date,
		FromAccountID: existing.FromAccountID,
		ToAccountID:   existing.ToAccountID,
		JarID:         existing.JarID,
	}

	if err := h.repo.Update(r.Context(), tx); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"txChanged": true}`)
	w.WriteHeader(http.StatusOK)
}

func (h *TransactionHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "new_transaction", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// FilterOptions returns JSON-like HTML data for populating filter dropdowns.
// Returns accounts and jars as <option> lists for use by the filter panel.
func (h *TransactionHandler) FilterOptions(w http.ResponseWriter, r *http.Request) {
	accounts, _ := h.accountRepo.List(r.Context())
	jars, _ := h.jarRepo.List(r.Context())

	var accountOpts, jarOpts string
	for _, a := range accounts {
		accountOpts += fmt.Sprintf(`<option value="%d">%s</option>`, a.ID, a.Name)
	}
	for _, j := range jars {
		jarOpts += fmt.Sprintf(`<option value="%d">%s</option>`, j.ID, j.Name)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<select id="filter-account" onchange="loadList()"
		class="bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-1.5 text-xs outline-none focus:border-zinc-600 text-zinc-300">
		<option value="">All accounts</option>%s
	</select>
	<select id="filter-jar" onchange="loadList()"
		class="bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-1.5 text-xs outline-none focus:border-zinc-600 text-zinc-300">
		<option value="">All jars</option>%s
	</select>`, accountOpts, jarOpts)
}

func (h *TransactionHandler) AllPage(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "all_transactions", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TransactionHandler) Fields(w http.ResponseWriter, r *http.Request) {
	txType := r.URL.Query().Get("type")
	isMaster := r.URL.Query().Get("master") == "true"

	accounts, _ := h.accountRepo.List(r.Context())
	jars, _ := h.jarRepo.List(r.Context())

	var opts, jarOpts string
	for _, a := range accounts {
		opts += fmt.Sprintf(`<option value="%d">%s</option>`, a.ID, a.Name)
	}
	for _, j := range jars {
		jarOpts += fmt.Sprintf(`<option value="%d">%s</option>`, j.ID, j.Name)
	}

	sel := func(name, label, o string) string {
		return fmt.Sprintf(`
			<div>
				<label class="block text-xs text-zinc-500 mb-2">%s</label>
				<select name="%s" required class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
					%s
				</select>
			</div>`, label, name, o)
	}

	jarSel := fmt.Sprintf(`
		<div>
			<label class="block text-xs text-zinc-500 mb-2">Category (Jar) *</label>
			<select name="jar_id" required class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-4 py-3 outline-none focus:border-zinc-600">
				<option value="">Select jar</option>
				%s
			</select>
		</div>`, jarOpts)

	var html string
	switch txType {
	case "expense":
		html = sel("account", "From Account *", opts) + jarSel
	case "income":
		masterCheck := `
			<div>
				<label class="flex items-center gap-3 cursor-pointer">
					<input type="checkbox" name="master_income" value="true"
						onchange="onMasterIncomeChange(this)"
						class="w-4 h-4 accent-emerald-400"` +
			func() string {
				if isMaster {
					return " checked"
				}
				return ""
			}() + `/>
					<span class="text-sm text-zinc-300">Master Income
						<span class="text-xs text-zinc-500">(auto-split to jars)</span>
					</span>
				</label>
			</div>`
		html = sel("account", "To Account *", opts) + masterCheck
		if !isMaster {
			html += jarSel
		}
	case "transfer":
		html = sel("from", "From Account *", opts) + sel("to", "To Account *", opts)
	}

	w.Write([]byte(html))
}
