// Package router
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/auth"
	"github.com/joshu-sajeev/echo/internal/dashboard"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/transactions"
)

// Config holds the pre-constructed handlers passed from the app container.
type Config struct {
	AccountHandler     *accounts.AccountHandler
	JarHandler         *jars.JarHandler
	TransactionHandler *transactions.TransactionHandler
	DashboardHandler   *dashboard.Handler
	AuthHandler        *auth.Handler
}

// New constructs the root router and mounts all domain routes.
func New(cfg Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://10.122.147.88:5173",
			"http://192.168.0.112:5173",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		AllowCredentials: true,
	}))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post(
			"/auth/login",
			cfg.AuthHandler.Login,
		)

		r.Post(
			"/auth/logout",
			cfg.AuthHandler.Logout,
		)

		r.Get(
			"/auth/me",
			cfg.AuthHandler.Me,
		)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(cfg.AuthHandler.Store))

			r.Get(
				"/dashboard",
				cfg.DashboardHandler.GetDashboard,
			)
			cfg.AccountHandler.RegisterRoutes(r)
			cfg.JarHandler.RegisterRoutes(r)
			cfg.TransactionHandler.RegisterRoutes(r)
		})
	})

	return r
}
