// Package router
package router

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joshu-sajeev/echo/internal/handlers"
)

func New(tmpl *template.Template) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// pages
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "base.html", nil)
	})

	// accounts
	r.Post("/accounts", handlers.CreateAccountHandler)
	r.Get("/accounts", handlers.ListAccountsHandler)

	// transactions
	r.Post("/transactions", handlers.CreateTransactionHandler)
	r.Get("/transactions", handlers.ListTransactionsHandler)

	return r
}
