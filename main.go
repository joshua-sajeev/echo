package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/db"
	"github.com/joshu-sajeev/echo/internal/repository"
	"github.com/joshu-sajeev/echo/internal/router"
)

func main() {
	_ = godotenv.Load()

	db.Connect(os.Getenv("DATABASE_URL"))

	// ensure system jars exist (idempotent)
	jarRepo := repository.NewJarRepository(db.Pool)
	if err := jarRepo.EnsureDefaults(context.Background()); err != nil {
		log.Fatal("failed to seed default jars:", err)
	}

	templates := loadTemplates()

	r := router.New(templates, db.Pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func loadTemplates() *template.Template {
	t := template.New("")
	template.Must(t.ParseFiles(
		"templates/base.html",
		"templates/index.html",
	))
	template.Must(t.ParseGlob("templates/partials/*.html"))
	return t
}
