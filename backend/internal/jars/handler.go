package jars

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kittipat1413/go-common/framework/validator"
)

type JarHandler struct {
	service JarServiceInterface
	v       *validator.Validator
}

var (
	ErrInvalidAllocationType = errors.New("invalid or missing allocation type")
	ErrJarNotFound           = errors.New("jar not found")
	ErrJarValidation         = errors.New("jar validation failed")
)

func NewJarHandler(service JarServiceInterface) *JarHandler {
	v, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		log.Fatal("failed to initialize validator:", err)
	}
	return &JarHandler{service: service, v: v}
}

func (h *JarHandler) CreateJar(w http.ResponseWriter, r *http.Request) {
	var req CreateJarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.v.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateJar(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidAllocationType) || errors.Is(err, ErrJarValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

	var req UpdateJarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateJar(r.Context(), id, req); err != nil {
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
