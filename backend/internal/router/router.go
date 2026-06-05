// Package router
package router

import (
	"net/http"
	"os"
	"path/filepath"

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
			"http://10.174.66.88:8080",
			"http://localhost:5173",
			"https://echo-ui.onrender.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Head("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", cfg.AuthHandler.Login)
		r.Post("/auth/logout", cfg.AuthHandler.Logout)
		r.Get("/auth/me", cfg.AuthHandler.Me)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(cfg.AuthHandler.Store))

			r.Get("/dashboard", cfg.DashboardHandler.GetDashboard)
			cfg.AccountHandler.RegisterRoutes(r)
			cfg.JarHandler.RegisterRoutes(r)
			cfg.TransactionHandler.RegisterRoutes(r)
		})
	})

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	distPath := filepath.Join(cwd, "../frontend/dist")
	fs := http.FileServer(http.Dir(distPath))

	r.Handle("/assets/*", fs)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := "../frontend/dist" + r.URL.Path

		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, "../frontend/dist/index.html")
	})

	return r
}
