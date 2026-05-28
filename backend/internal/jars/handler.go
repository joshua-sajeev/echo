package jars

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JarHandler struct {
	service JarServiceInterface
}

func NewJarHandler(service JarServiceInterface) *JarHandler {
	return &JarHandler{service: service}
}

func (h *JarHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/jars", h.CreateJar)
	r.Get("/jars", h.ListJars)
	r.Put("/jars/{id}", h.UpdateJar)
	r.Delete("/jars/{id}", h.DeleteJar)

	return r
}

func (h *JarHandler) CreateJar(w http.ResponseWriter, r *http.Request) {
	var jar Jar

	if err := json.NewDecoder(r.Body).Decode(&jar); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateJar(r.Context(), jar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *JarHandler) ListJars(w http.ResponseWriter, r *http.Request) {
	jars, err := h.service.ListJars(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(jars)
}

func (h *JarHandler) UpdateJar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var jar Jar
	if err := json.NewDecoder(r.Body).Decode(&jar); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	jar.ID = id

	if err := h.service.UpdateJar(r.Context(), jar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *JarHandler) DeleteJar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteJar(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
