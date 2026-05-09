// Package router
package router

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joshu-sajeev/echo/internal/handlers"
	"github.com/joshu-sajeev/echo/internal/repository"
)

func New(tmpl *template.Template, conn *pgx.Conn) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// repos
	accountRepo := repository.NewAccountRepository(conn)
	txRepo := repository.NewTransactionRepository(conn)
	jarRepo := repository.NewJarRepository(conn)

	// handlers
	accountH := handlers.NewAccountHandler(accountRepo, txRepo)
	txH := handlers.NewTransactionHandler(txRepo)
	optionsH := handlers.NewOptionsHandler(accountRepo, jarRepo)

	// pages
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// accounts
	r.Post("/accounts", accountH.Create)
	r.Get("/accounts", accountH.List)

	// transactions
	r.Post("/transactions", txH.Create)
	r.Get("/transactions", txH.List)

	// htmx option endpoints
	r.Get("/accounts/options", optionsH.AccountFields)
	r.Get("/jars/options", optionsH.JarFields)

	return r
}
