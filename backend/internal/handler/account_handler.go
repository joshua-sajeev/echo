// internal/handler/account_handler.go
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/joshu-sajeev/echo/internal/dto"
	"github.com/joshu-sajeev/echo/internal/service"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

// POST /accounts
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAccountName) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print(err)
		http.Error(w, "failed to create account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(dto.CreateAccountResponse{
		ID: id,
	})
}

// GET /accounts
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.List(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// GET /accounts/balances
func (h *AccountHandler) ListWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListWithBalances(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list account balances", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// GET /accounts/archived
func (h *AccountHandler) ListArchivedWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListArchivedWithBalances(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list archived accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// PATCH /accounts/{id}/rename
func (h *AccountHandler) Rename(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	log.Println("idParam:", idParam)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	var req dto.RenameAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Print(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.Rename(r.Context(), id, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAccountID):
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, service.ErrInvalidAccountName):
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			log.Print(err)
			http.Error(w, "failed to rename account", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /accounts/{id}/archive
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
		if errors.Is(err, service.ErrInvalidAccountID) {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print(err)
		http.Error(w, "failed to archive account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /accounts/{id}/unarchive
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
		if errors.Is(err, service.ErrInvalidAccountID) {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "failed to unarchive account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
