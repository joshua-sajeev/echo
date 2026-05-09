package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

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

		// jar required if not master income
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

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	txs, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	for _, tx := range txs {
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

		// use created_at if date is zero/invalid
		d := tx.Date
		if d.IsZero() || d.Year() < 2000 {
			d = tx.CreatedAt
		}

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

		// label: jar name if set, otherwise type
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

		// account shown bottom-right
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

		if _, err := fmt.Fprintf(w, `
<div class="flex items-center justify-between py-3">
  <div class="min-w-0">
    <p class="text-sm font-medium text-zinc-100 truncate">%s</p>
    <p class="text-xs text-zinc-500 mt-0.5">%s • %s</p>
  </div>
  <div class="text-right shrink-0 ml-4">
    <p class="text-sm font-semibold %s">%s</p>
    <p class="text-xs text-zinc-600">%s</p>
  </div>
</div>
`, tx.Name, label, dateStr, amountColor, amountStr, account); err != nil {
			log.Println("write error:", err)
		}
	}
}

func (h *TransactionHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	// serves the full new transaction page
	if err := h.tmpl.ExecuteTemplate(w, "new_transaction", nil); err != nil {
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
