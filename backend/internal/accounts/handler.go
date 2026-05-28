package accounts

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kittipat1413/go-common/framework/validator"
)

type AccountHandler struct {
	service AccountServiceInterface
	v       *validator.Validator
}

// NewAccountHandler creates a new AccountHandler.
func NewAccountHandler(service AccountServiceInterface) *AccountHandler {
	v, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		log.Fatal("failed to initialize validator:", err)
	}

	return &AccountHandler{
		service: service,
		v:       v,
	}
}

// Create handles POST /accounts
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAccountName):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, ErrAccountAlreadyExists):
			http.Error(w, "account already exists", http.StatusConflict)

		default:
			log.Print(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusCreated, CreateAccountResponse{
		ID: id,
	})
}

// List handles GET /accounts
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.List(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list accounts", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// ListWithBalances handles GET /accounts/balances
func (h *AccountHandler) ListWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListWithBalances(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list account balances", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// ListArchivedWithBalances handles GET /accounts/archived
func (h *AccountHandler) ListArchivedWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListArchivedWithBalances(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to list archived accounts", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// Rename handles PATCH /accounts/{id}/rename
func (h *AccountHandler) Rename(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	var req RenameAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

// Archive handles PATCH /accounts/{id}/archive
func (h *AccountHandler) Archive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	err = h.service.Archive(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAccountID):
			http.Error(w, "invalid account id", http.StatusBadRequest)

		case errors.Is(err, ErrAccountAlreadyState):
			http.Error(w, "account not found or already archived", http.StatusConflict)

		default:
			log.Print(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Unarchive handles PATCH /accounts/{id}/unarchive
func (h *AccountHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	err = h.service.Unarchive(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAccountID):
			http.Error(w, "invalid account id", http.StatusBadRequest)

		case errors.Is(err, ErrAccountAlreadyState):
			http.Error(w, "account not found or not archived", http.StatusConflict)

		default:
			log.Print(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers all account routes
func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)

		r.Get("/balances", h.ListWithBalances)
		r.Get("/archived", h.ListArchivedWithBalances)

		r.Patch("/{id}/rename", h.Rename)
		r.Patch("/{id}/archive", h.Archive)
		r.Patch("/{id}/unarchive", h.Unarchive)
	})
}

// writeJSON writes JSON response
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Print(err)
	}
}
