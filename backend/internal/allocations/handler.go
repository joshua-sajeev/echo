package allocations

import (
	"encoding/json"
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
		r.Post("/run", h.Run)
	})
}

func (h *AllocationHandler) Run(
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

	if err := h.service.Run(
		r.Context(),
		req.Amount,
	); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			"",
			"ALLOCATION_FAILED",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
