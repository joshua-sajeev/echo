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

func (h *JarHandler) RegisterRoutes(r chi.Router) {
	r.Route("/jars", func(r chi.Router) {
		r.Post("/", h.CreateJar)
		r.Get("/", h.ListJars)
		r.Put("/{id}", h.UpdateJar)
		r.Delete("/{id}", h.DeleteJar)
	})
}

func (h *JarHandler) CreateJar(w http.ResponseWriter, r *http.Request) {
	var req CreateJarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	id, err := h.service.CreateJar(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *JarHandler) ListJars(w http.ResponseWriter, r *http.Request) {
	jars, err := h.service.ListJars(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, jars)
}

func (h *JarHandler) UpdateJar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid id", "id", "INVALID_ID")
		return
	}

	var req UpdateJarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	err = h.service.UpdateJar(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *JarHandler) DeleteJar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid id", "id", "INVALID_ID")
		return
	}

	err = h.service.DeleteJar(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrJarNameRequired):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "name", "JAR_NAME_REQUIRED")

	case errors.Is(err, ErrPercentageMustBePositive):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "value", "INVALID_PERCENTAGE")

	case errors.Is(err, ErrTotalPercentageExceeded):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "value", "PERCENTAGE_LIMIT_EXCEEDED")

	case errors.Is(err, ErrJarNameAlreadyExists):
		httpresponse.WriteError(w, http.StatusConflict, "jar name already exists", "name", "JAR_NAME_ALREADY_EXISTS")

	case errors.Is(err, ErrJarNotFound):
		httpresponse.WriteError(w, http.StatusNotFound, "jar not found", "id", "JAR_NOT_FOUND")

	case errors.Is(err, ErrInvalidJarID):
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid jar id", "id", "INVALID_ID")

	default:
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal server error", "", "INTERNAL_ERROR")
	}
}
