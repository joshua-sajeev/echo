package jars

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joshu-sajeev/echo/internal/httpresponse"
	"github.com/kittipat1413/go-common/framework/validator"
)

type JarHandler struct {
	service JarServiceInterface
	v       *validator.Validator
}

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
		switch {
		case errors.Is(err, ErrJarNameRequired):
			httpresponse.WriteError(w, 400, err.Error(), "name", "JAR_NAME_REQUIRED")

		case errors.Is(err, ErrPercentageMustBePositive):
			httpresponse.WriteError(w, 400, err.Error(), "value", "INVALID_PERCENTAGE")

		case errors.Is(err, ErrTotalPercentageExceeded):
			httpresponse.WriteError(w, 400, err.Error(), "value", "PERCENTAGE_LIMIT_EXCEEDED")
		case errors.Is(err, ErrJarNameAlreadyExists):
			httpresponse.WriteError(w, 409, "jar name already exists", "name", "JAR_NAME_ALREADY_EXISTS")
		case errors.Is(err, ErrJarNotFound):
			httpresponse.WriteError(w, 404, "jar not found", "id", "JAR_NOT_FOUND")
		default:
			httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
		}
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

	err = h.service.UpdateJar(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrJarNotFound):
			httpresponse.WriteError(w, 404, "jar not found", "id", "JAR_NOT_FOUND")

		case errors.Is(err, ErrInvalidJarID):
			httpresponse.WriteError(w, 400, "invalid jar id", "id", "INVALID_ID")

		case errors.Is(err, ErrJarNameRequired):
			httpresponse.WriteError(w, 400, "jar name required", "name", "JAR_NAME_REQUIRED")

		case errors.Is(err, ErrPercentageMustBePositive):
			httpresponse.WriteError(w, 400, "percentage must be positive", "value", "INVALID_PERCENTAGE")

		case errors.Is(err, ErrTotalPercentageExceeded):
			httpresponse.WriteError(w, 400, "total percentage exceeds 100", "value", "PERCENTAGE_LIMIT_EXCEEDED")

		default:
			httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
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

	err = h.service.DeleteJar(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrJarNotFound):
			httpresponse.WriteError(w, 404, "jar not found", "id", "JAR_NOT_FOUND")

		case errors.Is(err, ErrInvalidJarID):
			httpresponse.WriteError(w, 400, "invalid jar id", "id", "INVALID_ID")

		default:
			httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
		}
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
