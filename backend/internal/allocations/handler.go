package allocations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joshu-sajeev/echo/internal/httpresponse"
)

type AllocationHandler struct {
	service AllocationServiceInterface
}

func NewAllocationHandler(
	service AllocationServiceInterface,
) *AllocationHandler {
	return &AllocationHandler{
		service: service,
	}
}

func (h *AllocationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/allocation", func(r chi.Router) {
		r.Post("/run", h.RunManual)                  // Manual: to specific goal
		r.Post("/distribute", h.Distribute)          // Automatic: by percentage
		r.Get("/leisure-leftover", h.GetLeftover)   // Get last month's leisure leftover
	})
}

// RunManual handles POST /allocation/run (manual allocation to specific goal)
func (h *AllocationHandler) RunManual(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req RunAllocationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			"invalid request body",
			"",
			"INVALID_REQUEST_BODY",
		)
		return
	}

	if err := h.service.RunManual(
		r.Context(),
		req.GoalID,
		req.Amount,
	); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Distribute handles POST /allocation/distribute (automatic allocation by percentage)
func (h *AllocationHandler) Distribute(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req DistributeAllocationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			"invalid request body",
			"",
			"INVALID_REQUEST_BODY",
		)
		return
	}

	if err := h.service.DistributeAutomatic(
		r.Context(),
		req.Type,
		req.Amount,
	); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetLeftover handles GET /allocation/leisure-leftover
func (h *AllocationHandler) GetLeftover(
	w http.ResponseWriter,
	r *http.Request,
) {
	leftover, alreadyAllocated, err := h.service.GetLeisureLeftover(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"amount":            leftover,
		"already_allocated": alreadyAllocated,
	})
}

func (h *AllocationHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidAmount):
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			"amount",
			"INVALID_AMOUNT",
		)

	case errors.Is(err, ErrInvalidGoalID):
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			"goal_id is required and must be greater than 0",
			"goal_id",
			"INVALID_GOAL_ID",
		)

	case errors.Is(err, ErrGoalNotFound):
		httpresponse.WriteError(
			w,
			http.StatusNotFound,
			"goal not found",
			"goal_id",
			"GOAL_NOT_FOUND",
		)

	case errors.Is(err, ErrGoalArchived):
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			"cannot allocate to archived goal",
			"goal_id",
			"GOAL_ARCHIVED",
		)

	case errors.Is(err, ErrNoGoalsConfigured):
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			"",
			"NO_GOALS_CONFIGURED",
		)

	case errors.Is(err, ErrInvalidAllocationPercentages):
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			"",
			"INVALID_PERCENTAGES",
		)

	case errors.Is(err, ErrLeisureJarNotFound):
		httpresponse.WriteError(
			w,
			http.StatusNotFound,
			err.Error(),
			"type",
			"LEISURE_JAR_NOT_FOUND",
		)

	case errors.Is(err, ErrInvalidAllocationType):
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			"type",
			"INVALID_ALLOCATION_TYPE",
		)

	default:
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			"allocation failed",
			"",
			"ALLOCATION_ERROR",
		)
	}
}
