// Package router provides HTTP route registration for the API server.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/handler"
	"github.com/joshu-sajeev/echo/internal/repository"
	"github.com/joshu-sajeev/echo/internal/service"
)

func New(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Head("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	accountRepo := repository.NewAccountRepository(pool)

	accountService := service.NewAccountService(accountRepo)

	accountHandler := handler.NewAccountHandler(accountService)

	RegisterAccountRoutes(r, accountHandler)
	return r
}

func RegisterAccountRoutes(r chi.Router, h *handler.AccountHandler) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.Create)

		r.Get("/", h.List)
		r.Get("/balances", h.ListWithBalances)
		r.Get("/archived", h.ListArchivedWithBalances)

		r.Patch("/{id}/rename", h.Rename)
		r.Patch("/{id}/archive", h.Archive)
		r.Patch("/{id}/unarchive", h.Unarchive)
	})
}
