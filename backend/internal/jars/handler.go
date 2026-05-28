package jars

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JarHandler struct {
	service JarServiceInterface
}

// Global Sentinel Errors used across layers
var (
	ErrInvalidAllocationType = errors.New("invalid or missing allocation type")
	ErrJarNotFound           = errors.New("jar not found")
	ErrJarValidation         = errors.New("jar validation failed")
)

func NewJarHandler(service JarServiceInterface) *JarHandler {
	return &JarHandler{service: service}
}

func (h *JarHandler) CreateJar(w http.ResponseWriter, r *http.Request) {
	var jar Jar
	if err := json.NewDecoder(r.Body).Decode(&jar); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateJar(r.Context(), jar)
	if err != nil {
		// Business validation errors -> 400 Bad Request
		if errors.Is(err, ErrInvalidAllocationType) || errors.Is(err, ErrJarValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Database structural errors are already logged in service layer via dbutil.LogError.
		// Just report clean 500 to client without double-logging here.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *JarHandler) ListJars(w http.ResponseWriter, r *http.Request) {
	jars, err := h.service.ListJars(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
		switch {
		case errors.Is(err, ErrJarNotFound):
			http.Error(w, "jar not found", http.StatusNotFound)
		case errors.Is(err, ErrInvalidAllocationType) || errors.Is(err, ErrJarValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
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
		if errors.Is(err, ErrJarNotFound) {
			http.Error(w, "jar not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *JarHandler) RegisterRoutes(r chi.Router) {
	r.Route("/jars", func(r chi.Router) {
		r.Post("/", h.CreateJar)
		r.Get("/", h.ListJars)
		r.Put("/{id}", h.UpdateJar)
		r.Delete("/{id}", h.DeleteJar)
	})
}
