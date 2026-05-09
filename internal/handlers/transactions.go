package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joshu-sajeev/echo/internal/models"
	"github.com/joshu-sajeev/echo/internal/repository"
)

type TransactionHandler struct {
	repo repository.TransactionRepository
}

func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
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

	for _, tx := range txs {
		color := "text-red-400"
		if tx.Type == "income" {
			color = "text-green-400"
		} else if tx.Type == "transfer" {
			color = "text-blue-400"
		}

		if _, err := fmt.Fprintf(w, `
			<div class="bg-zinc-800 p-4 rounded-xl flex justify-between">
				<div>%s</div>
				<div class="%s">₹%.2f</div>
			</div>
		`, tx.Name, color, tx.Amount); err != nil {
			log.Println("write error:", err)
		}
	}
}
