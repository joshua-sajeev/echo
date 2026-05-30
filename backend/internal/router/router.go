// Package router
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/auth"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/transactions"
)

// Config holds the pre-constructed handlers passed from the app container.
type Config struct {
	AccountHandler     *accounts.AccountHandler
	JarHandler         *jars.JarHandler
	TransactionHandler *transactions.TransactionHandler

	AuthHandler *auth.Handler
}

// New constructs the root router and mounts all domain routes.
func New(cfg Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post(
			"/auth/login",
			cfg.AuthHandler.Login,
		)

		r.Post(
			"/auth/logout",
			cfg.AuthHandler.Logout,
		)
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(cfg.AuthHandler.Store))

			cfg.AccountHandler.RegisterRoutes(r)
			cfg.JarHandler.RegisterRoutes(r)
			cfg.TransactionHandler.RegisterRoutes(r)
		})

		r.Head("/", func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		})
	})

	return r
}
