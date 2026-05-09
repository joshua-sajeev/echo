package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)

	log.Println("Server running on :3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := map[string]string{
		"Name": "Joshua",
	}

	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}
