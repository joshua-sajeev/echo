package goals

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

type GoalHandler struct {
	service GoalServiceInterface
	v       *validator.Validator
}

// NewGoalHandler creates a new GoalHandler
func NewGoalHandler(service GoalServiceInterface) *GoalHandler {
	v, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		log.Fatal("failed to initialize validator:", err)
	}

	return &GoalHandler{
		service: service,
		v:       v,
	}
}

// RegisterRoutes registers all goal routes
func (h *GoalHandler) RegisterRoutes(r chi.Router) {
	r.Route("/goals", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}/archive", h.Archive)
		r.Patch("/{id}/restore", h.Restore)
		r.Post("/{id}/progress", h.AddProgress)
	})
}

// Create handles POST /goals
func (h *GoalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateGoalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	id, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// List handles GET /goals
func (h *GoalHandler) List(w http.ResponseWriter, r *http.Request) {
	goals, err := h.service.List(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, goals)
}

// GetByID handles GET /goals/{id}
func (h *GoalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid goal id", "id", "INVALID_ID")
		return
	}

	goal, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, goal)
}

// Update handles PUT /goals/{id}
func (h *GoalHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid goal id", "id", "INVALID_ID")
		return
	}

	var req UpdateGoalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	if err := h.service.Update(r.Context(), id, req); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive Patch  /goals/{id}/archive
func (h *GoalHandler) Archive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid goal id", "id", "INVALID_ID")
		return
	}

	if err := h.service.Archive(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Restore handles Patch /goals/{id}/restore
func (h *GoalHandler) Restore(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			"invalid goal id",
			"id",
			"INVALID_ID",
		)
		return
	}

	if err := h.service.Restore(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddProgress handles POST /goals/{id}/progress
func (h *GoalHandler) AddProgress(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid goal id", "id", "INVALID_ID")
		return
	}

	var req AddProgressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	if err := h.service.AddProgress(r.Context(), id, req.Amount); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GoalHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidGoalID):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "id", "INVALID_ID")

	case errors.Is(err, ErrGoalNameRequired):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "name", "GOAL_NAME_REQUIRED")

	case errors.Is(err, ErrTargetAmountInvalid):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "target_amount", "INVALID_TARGET_AMOUNT")

	case errors.Is(err, ErrProgressAmountInvalid):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "amount", "INVALID_PROGRESS_AMOUNT")

	case errors.Is(err, ErrGoalNotFound):
		httpresponse.WriteError(w, http.StatusNotFound, "goal not found", "id", "GOAL_NOT_FOUND")

	case errors.Is(err, ErrGoalAlreadyCompleted):
		httpresponse.WriteError(w, http.StatusConflict, err.Error(), "id", "GOAL_ALREADY_COMPLETED")

	case errors.Is(err, ErrDeadlinePassed):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "deadline", "DEADLINE_PASSED")

	default:
		log.Print(err)
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal server error", "", "INTERNAL_ERROR")
	}
}
