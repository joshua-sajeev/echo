package accounts

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

// Create handles POST /accounts
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	id, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, CreateAccountResponse{
		ID: id,
	})
}

// List handles GET /accounts
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.List(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, accounts)
}

// ListWithBalances handles GET /accounts/balances
func (h *AccountHandler) ListWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListWithBalances(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, accounts)
}

// ListArchivedWithBalances handles GET /accounts/archived
func (h *AccountHandler) ListArchivedWithBalances(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListArchivedWithBalances(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, accounts)
}

// Rename handles PATCH /accounts/{id}/rename
func (h *AccountHandler) Rename(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid account id", "id", "INVALID_ID")
		return
	}

	var req RenameAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	if err := h.service.Rename(r.Context(), id, req.Name); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive handles PATCH /accounts/{id}/archive
func (h *AccountHandler) Archive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid account id", "id", "INVALID_ID")
		return
	}

	if err := h.service.Archive(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Unarchive handles PATCH /accounts/{id}/unarchive
func (h *AccountHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid account id", "id", "INVALID_ID")
		return
	}

	if err := h.service.Unarchive(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidAccountID):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "id", "INVALID_ID")

	case errors.Is(err, ErrInvalidAccountName):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "name", "INVALID_ACCOUNT_NAME")

	case errors.Is(err, ErrAccountAlreadyExists):
		httpresponse.WriteError(w, http.StatusConflict, "account already exists", "name", "ACCOUNT_ALREADY_EXISTS")

	case errors.Is(err, ErrAccountNotFound):
		httpresponse.WriteError(w, http.StatusNotFound, "account not found", "id", "ACCOUNT_NOT_FOUND")

	case errors.Is(err, ErrAccountAlreadyArchived):
		httpresponse.WriteError(w, http.StatusConflict, err.Error(), "id", "ACCOUNT_ALREADY_ARCHIVED")

	case errors.Is(err, ErrAccountAlreadyActive):
		httpresponse.WriteError(w, http.StatusConflict, err.Error(), "id", "ACCOUNT_ALREADY_ACTIVE")

	default:
		log.Print(err)
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal server error", "", "INTERNAL_ERROR")
	}
}
