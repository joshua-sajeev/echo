package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joshu-sajeev/echo/internal/db"
)

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	opening := r.FormValue("opening_balance")

	if name == "" {
		http.Error(w, "name required", 400)
		return
	}

	// 1. insert account
	var accountID int64

	err := db.Conn.QueryRow(
		context.Background(),
		`INSERT INTO accounts (name)
		 VALUES ($1)
		 RETURNING id`,
		name,
	).Scan(&accountID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 2. if opening balance exists → create transaction
	if opening != "" && opening != "0" {

		_, err := db.Conn.Exec(
			context.Background(),
			`INSERT INTO transactions
			(type, amount, name, to_account_id)
			VALUES ('income', $1, 'Opening Balance', $2)`,
			opening,
			accountID,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Write([]byte(`
		<div class="bg-green-500/10 border border-green-500/30 text-green-400 p-3 rounded-xl">
			Account created
		</div>
	`))
}

func ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Conn.Query(
		context.Background(),
		`SELECT id, name FROM accounts ORDER BY id DESC`,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string

		rows.Scan(&id, &name)

		fmt.Fprintf(w, `
			<div class="bg-zinc-800 rounded-xl p-4 flex justify-between">
				<div>
					<p class="font-medium">%s</p>
					<p class="text-xs text-zinc-500">Account</p>
				</div>
			</div>
		`, name)
	}
}
