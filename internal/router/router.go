package router

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/handlers"
	"github.com/joshu-sajeev/echo/internal/repository"
)

func New(tmpl *template.Template, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	accountRepo := repository.NewAccountRepository(pool)
	txRepo := repository.NewTransactionRepository(pool)
	jarRepo := repository.NewJarRepository(pool)

	accountH := handlers.NewAccountHandler(accountRepo, txRepo)
	txH := handlers.NewTransactionHandler(txRepo)
	optionsH := handlers.NewOptionsHandler(accountRepo, jarRepo)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Post("/accounts", accountH.Create)
	r.Get("/accounts", accountH.List)
	r.Get("/accounts/archived", accountH.ListArchived)
	r.Patch("/accounts/{id}/rename", accountH.Rename)
	r.Delete("/accounts/{id}", accountH.Archive)
	r.Patch("/accounts/{id}/unarchive", accountH.Unarchive)

	r.Post("/transactions", txH.Create)
	r.Get("/transactions", txH.List)

	r.Get("/accounts/options", optionsH.AccountFields)
	r.Get("/jars/options", optionsH.JarFields)

	return r
}
