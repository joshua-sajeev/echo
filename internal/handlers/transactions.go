package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/db"
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	txType := r.FormValue("type")
	name := r.FormValue("name")
	amount := r.FormValue("amount")

	if txType == "" || amount == "" || name == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// convert amount
	var amt float64
	if _, err := fmt.Sscanf(amount, "%f", &amt); err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	switch txType {

	case "expense":
		fromID := r.FormValue("account")

		if fromID == "" {
			http.Error(w, "missing account", http.StatusBadRequest)
			return
		}

		if _, err := db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, from_account_id)
			VALUES ('expense', $1, $2, $3)`,
			amt, name, fromID,
		); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

	case "income":
		toID := r.FormValue("account")

		if toID == "" {
			http.Error(w, "missing account", http.StatusBadRequest)
			return
		}

		if _, err := db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, to_account_id)
			VALUES ('income', $1, $2, $3)`,
			amt, name, toID,
		); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

	case "transfer":
		fromID := r.FormValue("from")
		toID := r.FormValue("to")

		if fromID == "" || toID == "" {
			http.Error(w, "missing accounts", http.StatusBadRequest)
			return
		}

		if _, err := db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, from_account_id, to_account_id)
			VALUES ('transfer', $1, $2, $3, $4)`,
			amt, name, fromID, toID,
		); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "invalid transaction type", http.StatusBadRequest)
		return
	}

	// HTMX response
	if _, err := w.Write([]byte(`
		<div class="bg-zinc-800 p-4 rounded-xl text-green-400">
			Transaction added successfully
		</div>
	`)); err != nil {
		log.Println("write error:", err)
	}
}

func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Conn.Query(
		context.Background(),
		`SELECT name, amount, type FROM transactions ORDER BY id DESC`,
	)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name, txType string
		var amount float64

		if err := rows.Scan(&name, &amount, &txType); err != nil {
			log.Println("scan error:", err)
			continue
		}

		color := "text-red-400"
		if txType == "income" {
			color = "text-green-400"
		}

		if _, err := fmt.Fprintf(w, `
			<div class="bg-zinc-800 p-4 rounded-xl flex justify-between">
				<div>%s</div>
				<div class="%s">₹%.2f</div>
			</div>
		`, name, color, amount); err != nil {
			log.Println("write error:", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("rows error:", err)
	}
}
