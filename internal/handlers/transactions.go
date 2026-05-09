package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/db"
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	txType := r.FormValue("type")
	name := r.FormValue("name")
	amount := r.FormValue("amount")

	if txType == "" || amount == "" {
		http.Error(w, "invalid input", 400)
		return
	}

	// convert amount
	var amt float64
	fmt.Sscanf(amount, "%f", &amt)

	var fromID, toID string

	switch txType {

	case "expense":
		fromID = r.FormValue("account")
		toID = "NULL"

		db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, from_account_id)
			VALUES ('expense', $1, $2, $3)`,
			amt, name, fromID,
		)

	case "income":
		toID = r.FormValue("account")

		db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, to_account_id)
			VALUES ('income', $1, $2, $3)`,
			amt, name, toID,
		)

	case "transfer":
		fromID = r.FormValue("from")
		toID = r.FormValue("to")

		db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, from_account_id, to_account_id)
			VALUES ('transfer', $1, $2, $3, $4)`,
			amt, name, fromID, toID,
		)
	}

	w.Write([]byte(`
		<div class="bg-zinc-800 p-4 rounded-xl">
			Transaction added
		</div>
	`))
}

func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Conn.Query(
		context.Background(),
		`SELECT name, amount, type FROM transactions ORDER BY id DESC`,
	)
	defer rows.Close()

	for rows.Next() {
		var name, txType string
		var amount float64

		rows.Scan(&name, &amount, &txType)

		color := "text-red-400"
		if txType == "income" {
			color = "text-green-400"
		}

		fmt.Fprintf(w, `
			<div class="bg-zinc-800 p-4 rounded-xl flex justify-between">
				<div>%s</div>
				<div class="%s">₹%.2f</div>
			</div>
		`, name, color, amount)
	}
}
