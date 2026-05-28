package accounts

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	service AccountServiceInterface
}

// NewAccountHandler creates a new AccountHandler.
func NewAccountHandler(service AccountServiceInterface) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

// Create handles POST /accounts requests.
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, ErrInvalidAccountName) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print(err)
		http.Error(w, "failed to create account", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, CreateAccountResponse{
		ID: id,
	})
}

// List handles GET /accounts requests.
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.List(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list accounts", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// ListWithBalances handles GET /accounts/balances requests.
func (h *AccountHandler) ListWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListWithBalances(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list account balances", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// ListArchivedWithBalances handles GET /accounts/archived requests.
func (h *AccountHandler) ListArchivedWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListArchivedWithBalances(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list archived accounts", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// Rename handles PATCH /accounts/{id}/rename requests.
func (h *AccountHandler) Rename(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	var req RenameAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Print(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.Rename(r.Context(), id, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAccountID):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, ErrInvalidAccountName):
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			log.Print(err)
			http.Error(w, "failed to rename account", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive handles PATCH /accounts/{id}/archive requests.
func (h *AccountHandler) Archive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	err = h.service.Archive(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvalidAccountID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print(err)
		http.Error(w, "failed to archive account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Unarchive handles PATCH /accounts/{id}/unarchive requests.
func (h *AccountHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	err = h.service.Unarchive(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvalidAccountID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print(err)
		http.Error(w, "failed to unarchive account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Print(err)
	}
}
