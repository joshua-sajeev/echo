package router

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/auth"
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
	pageH := handlers.NewPageHandler(tmpl, txRepo, jarRepo)
	jarH := handlers.NewJarHandler(jarRepo, txRepo)
	loginH := handlers.NewLoginHandler()

	// ── public routes (no auth required) ────────────────────────
	r.Get("/login", loginH.Page)
	r.Post("/login", loginH.Submit)
	r.Get("/logout", loginH.Logout)

	// ── protected routes ─────────────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware)

		// pages
		r.Get("/", pageH.Index)

		// accounts
		r.Post("/accounts", accountH.Create)
		r.Get("/accounts", accountH.List)
		r.Get("/accounts/archived", accountH.ListArchived)
		r.Patch("/accounts/{id}/rename", accountH.Rename)
		r.Delete("/accounts/{id}", accountH.Archive)
		r.Patch("/accounts/{id}/unarchive", accountH.Unarchive)

		// transactions
		r.Post("/transactions", txH.Create)
		r.Get("/transactions/recent", txH.List)
		r.Get("/transactions/all", txH.ListAll)
		r.Get("/transactions/new", txH.NewForm)
		r.Get("/transactions/page", txH.AllPage)
		r.Get("/transactions/fields", txH.Fields)
		r.Get("/transactions/filter-options", txH.FilterOptions)
		r.Patch("/transactions/{id}", txH.Update)
		r.Delete("/transactions/{id}", txH.Delete)

		// transaction edit page
		r.Get("/transactions/{id}/edit", txH.EditPage)
		r.Get("/transactions/{id}/fields", txH.PrefilledFields)
		r.Get("/transactions/{id}/data", txH.Data)

		// jars
		r.Get("/jars/page", jarH.Page)
		r.Get("/jars", jarH.List)
		r.Post("/jars", jarH.Create)
		r.Patch("/jars/{id}", jarH.Update)
		r.Delete("/jars/{id}", jarH.Delete)

		// options (for transaction form dropdowns)
		r.Get("/accounts/options", optionsH.AccountFields)
		r.Get("/jars/options", optionsH.JarFields)
	})

	r.Head("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return r
}
