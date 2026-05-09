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
	txH := handlers.NewTransactionHandler(txRepo, accountRepo, jarRepo, tmpl)
	optionsH := handlers.NewOptionsHandler(accountRepo, jarRepo)
	pageH := handlers.NewPageHandler(tmpl, txRepo)

	r.Get("/", pageH.Index)

	r.Post("/accounts", accountH.Create)
	r.Get("/accounts", accountH.List)
	r.Get("/accounts/archived", accountH.ListArchived)
	r.Patch("/accounts/{id}/rename", accountH.Rename)
	r.Delete("/accounts/{id}", accountH.Archive)
	r.Patch("/accounts/{id}/unarchive", accountH.Unarchive)

	r.Post("/transactions", txH.Create)
	r.Get("/transactions/recent", txH.List)
	r.Get("/transactions/all", txH.ListAll)
	r.Get("/transactions/new", txH.NewForm)
	r.Get("/transactions/page", txH.AllPage)
	r.Get("/transactions/fields", txH.Fields)
	r.Get("/transactions/filter-options", txH.FilterOptions)
	r.Patch("/transactions/{id}", txH.Update)
	r.Delete("/transactions/{id}", txH.Delete)

	r.Get("/accounts/options", optionsH.AccountFields)
	r.Get("/jars/options", optionsH.JarFields)

	return r
}
