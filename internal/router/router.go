// Package router
package router

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joshu-sajeev/echo/internal/handler"
	"github.com/joshu-sajeev/echo/internal/repository"
)

func New(tmpl *template.Template, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	_ = repository.NewAccountRepository(pool)

	pageH := handler.NewPageHandler(tmpl)
	fs := http.FileServer(http.Dir("./ui/static"))

	r.Get("/", pageH.Index)
	r.Handle(
		"/static/*",
		http.StripPrefix("/static/", fs),
	)
	return r
}
