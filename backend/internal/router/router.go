// Package router
package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/allocations"
	"github.com/joshu-sajeev/echo/internal/auth"
	"github.com/joshu-sajeev/echo/internal/dashboard"
	"github.com/joshu-sajeev/echo/internal/goals"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/transactions"
)

// Config holds the pre-constructed handlers passed from the app container.
type Config struct {
	DB                 *pgxpool.Pool
	AccountHandler     *accounts.AccountHandler
	JarHandler         *jars.JarHandler
	TransactionHandler *transactions.TransactionHandler
	DashboardHandler   *dashboard.Handler
	AuthHandler        *auth.Handler
	GoalsHandler       *goals.GoalHandler
	AllocationsHandler *allocations.AllocationHandler
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

		r.Post("/demo/reset", func(w http.ResponseWriter, r *http.Request) {
			if os.Getenv("DEMO_MODE") != "true" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"demo mode is not enabled"}`))
				return
			}

			_, err := cfg.DB.Exec(r.Context(), `
				BEGIN;
				TRUNCATE TABLE goal_transactions, goals, transactions, jars, accounts RESTART IDENTITY CASCADE;

				INSERT INTO accounts (name) VALUES ('HDFC Savings'), ('Cash Wallet');

				INSERT INTO jars (name, allocation_type, value) VALUES
				('Necessities', 'percentage', 55),
				('Leisure', 'percentage', 20),
				('Investments', 'percentage', 15),
				('Giving', 'remainder', 0);

				INSERT INTO transactions (type, amount, name, date, to_account_id, category, is_master_income) VALUES
				('income', 85000, 'June Salary', '2026-06-01', 1, 'Income', TRUE);

				INSERT INTO transactions (type, amount, name, date, from_account_id, to_account_id, category, jar_id) VALUES
				('transfer', 46750, 'Necessities Allocation', '2026-06-01', 1, 1, 'Transfers', 1),
				('transfer', 17000, 'Leisure Allocation', '2026-06-01', 1, 1, 'Transfers', 2),
				('transfer', 12750, 'Investment Allocation', '2026-06-01', 1, 1, 'Transfers', 3),
				('transfer', 8500, 'Giving Allocation', '2026-06-01', 1, 1, 'Transfers', 4),
				('transfer', 4000, 'Cash Withdrawal', '2026-06-03', 1, 2, 'Transfers', NULL);

				INSERT INTO transactions (type, amount, name, date, from_account_id, category, jar_id) VALUES
				('expense', 18000, 'House Rent', '2026-06-02', 1, 'Housing', 1),
				('expense', 3200, 'Groceries', '2026-06-05', 1, 'Food', 1),
				('expense', 500, 'Coffee', '2026-06-08', 2, 'Food', 2),
				('expense', 5000, 'Nifty 50 SIP', '2026-06-15', 1, 'Investment', 3),
				('expense', 1000, 'Church Offering', '2026-06-20', 1, 'Donations', 4);

				INSERT INTO transactions (type, amount, name, date, to_account_id, category, is_master_income) VALUES
				('income', 90000, 'July Salary', '2026-07-01', 1, 'Income', TRUE);

				INSERT INTO transactions (type, amount, name, date, from_account_id, to_account_id, category, jar_id) VALUES
				('transfer', 49500, 'Necessities Allocation', '2026-07-01', 1, 1, 'Transfers', 1),
				('transfer', 18000, 'Leisure Allocation', '2026-07-01', 1, 1, 'Transfers', 2),
				('transfer', 13500, 'Investment Allocation', '2026-07-01', 1, 1, 'Transfers', 3),
				('transfer', 9000, 'Giving Allocation', '2026-07-01', 1, 1, 'Transfers', 4),
				('transfer', 5000, 'Cash Withdrawal', '2026-07-03', 1, 2, 'Transfers', NULL);

				INSERT INTO transactions (type, amount, name, date, from_account_id, category, jar_id) VALUES
				('expense', 18000, 'House Rent', '2026-07-02', 1, 'Housing', 1),
				('expense', 3500, 'Groceries', '2026-07-04', 1, 'Food', 1),
				('expense', 1200, 'Fuel', '2026-07-06', 1, 'Transport', 1),
				('expense', 450, 'Cafe Coffee Day', '2026-07-08', 2, 'Food', 2),
				('expense', 1200, 'Movie Night', '2026-07-12', 2, 'Entertainment', 2),
				('expense', 5000, 'Nifty 50 SIP', '2026-07-15', 1, 'Investment', 3),
				('expense', 1000, 'Church Offering', '2026-07-20', 1, 'Donations', 4);

				INSERT INTO goals (name, target_amount, saved_amount, deadline, allocation_percentage) VALUES
				('Emergency Fund', 300000, 35000, '2027-12-31', 20),
				('Japan Trip', 150000, 10000, '2027-06-01', 10);

				INSERT INTO goal_transactions (goal_id, amount, transaction_type, notes) VALUES
				(1, 17000, 'allocation', 'June allocation'),
				(1, 18000, 'allocation', 'July allocation'),
				(2, 10000, 'manual_contribution', 'Initial contribution');
				COMMIT;
			`)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"failed to reset demo database: ` + err.Error() + `"}`))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		})

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(cfg.AuthHandler.Store))

			r.Get("/dashboard", cfg.DashboardHandler.GetDashboard)
			cfg.AccountHandler.RegisterRoutes(r)
			cfg.JarHandler.RegisterRoutes(r)
			cfg.TransactionHandler.RegisterRoutes(r)
			cfg.GoalsHandler.RegisterRoutes(r)
			cfg.AllocationsHandler.RegisterRoutes(r)
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
