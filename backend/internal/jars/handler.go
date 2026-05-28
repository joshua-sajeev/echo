package jars

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JarHandler struct {
	service JarServiceInterface
}

var (
	ErrInvalidAllocationType = errors.New("invalid or missing allocation type")
	ErrJarNotFound           = errors.New("jar not found")
	// Add other user-facing errors here
)

// NewJarHandler creates a new JarHandler.
func NewJarHandler(service JarServiceInterface) *JarHandler {
	return &JarHandler{service: service}
}

// CreateJar handles POST /api/v1/jars requests.
func (h *JarHandler) CreateJar(w http.ResponseWriter, r *http.Request) {
	var jar Jar

	if err := json.NewDecoder(r.Body).Decode(&jar); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateJar(r.Context(), jar)
	if err != nil {
		if errors.Is(err, ErrInvalidAllocationType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("[ERROR] failed to create jar: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

// ListJars handles GET /api/v1/jars requests.
func (h *JarHandler) ListJars(w http.ResponseWriter, r *http.Request) {
	jars, err := h.service.ListJars(r.Context())
	if err != nil {
		log.Printf("[ERROR] failed to list jars: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jars)
}

// UpdateJar handles PUT /api/v1/jars/{id} requests.
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
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, ErrInvalidAllocationType):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("[ERROR] failed to update jar %d: %v", id, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteJar handles DELETE /api/v1/jars/{id} requests.
func (h *JarHandler) DeleteJar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteJar(r.Context(), id); err != nil {
		if errors.Is(err, ErrJarNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("[ERROR] failed to delete jar %d: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes mounts this domain's endpoints onto the provided router.
func (h *JarHandler) RegisterRoutes(r chi.Router) {
	r.Route("/jars", func(r chi.Router) {
		r.Post("/", h.CreateJar)
		r.Get("/", h.ListJars)
		r.Put("/{id}", h.UpdateJar)
		r.Delete("/{id}", h.DeleteJar)
	})
}
