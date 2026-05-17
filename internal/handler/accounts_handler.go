// Package handler contains HTTP handlers.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joshu-sajeev/echo/internal/dto"
	"github.com/joshu-sajeev/echo/internal/service"
)

// AccountHandler handles account HTTP requests.
type AccountHandler struct {
	service service.AccountServiceInterface
}

// NewAccountHandler creates a new AccountHandler.
func NewAccountHandler(
	service service.AccountServiceInterface,
) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

// Create handles account creation.
func (h *AccountHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(
		r.Context(),
		req.Name,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := writeJSON(
		w,
		http.StatusCreated,
		map[string]any{
			"id": id,
		},
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}

// List returns all active accounts.
func (h *AccountHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	accounts, err := h.service.List(r.Context())
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := make(
		[]dto.AccountResponse,
		0,
		len(accounts),
	)

	for _, account := range accounts {
		response = append(
			response,
			dto.ToAccountResponse(account),
		)
	}

	if err := writeJSON(
		w,
		http.StatusOK,
		response,
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}

// ListWithBalances returns active accounts with balances.
func (h *AccountHandler) ListWithBalances(
	w http.ResponseWriter,
	r *http.Request,
) {
	accounts, err := h.service.ListWithBalances(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := make(
		[]dto.AccountWithBalanceResponse,
		0,
		len(accounts),
	)

	for _, account := range accounts {
		response = append(
			response,
			dto.ToAccountWithBalanceResponse(account),
		)
	}

	if err := writeJSON(
		w,
		http.StatusOK,
		response,
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}

// ListArchivedWithBalances returns archived accounts with balances.
func (h *AccountHandler) ListArchivedWithBalances(
	w http.ResponseWriter,
	r *http.Request,
) {
	accounts, err := h.service.ListArchivedWithBalances(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := make(
		[]dto.AccountWithBalanceResponse,
		0,
		len(accounts),
	)

	for _, account := range accounts {
		response = append(
			response,
			dto.ToAccountWithBalanceResponse(account),
		)
	}

	if err := writeJSON(
		w,
		http.StatusOK,
		response,
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}

// Rename updates an account name.
func (h *AccountHandler) Rename(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(
		idParam,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid account id",
			http.StatusBadRequest,
		)
		return
	}

	var req dto.RenameAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Rename(
		r.Context(),
		id,
		req.Name,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive archives an account.
func (h *AccountHandler) Archive(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(
		idParam,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid account id",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Archive(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Unarchive restores an archived account.
func (h *AccountHandler) Unarchive(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(
		idParam,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid account id",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Unarchive(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) error {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
