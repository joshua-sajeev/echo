package transactions

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

type TransactionHandler struct {
	service TransactionServiceInterface
	v       *validator.Validator
}

func NewTransactionHandler(service TransactionServiceInterface) *TransactionHandler {
	v, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		log.Fatal("failed to initialize validator:", err)
	}
	return &TransactionHandler{service: service, v: v}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.v.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTransactionNameRequired):
			httpresponse.WriteError(w, 400, err.Error(), "name", "TRANSACTION_NAME_REQUIRED")
		case errors.Is(err, ErrTransactionTypeRequired):
			httpresponse.WriteError(w, 400, err.Error(), "type", "TRANSACTION_TYPE_REQUIRED")
		case errors.Is(err, ErrTransactionAmountInvalid):
			httpresponse.WriteError(w, 400, err.Error(), "amount", "INVALID_AMOUNT")
		case errors.Is(err, ErrTransactionSameAccount):
			httpresponse.WriteError(w, 400, err.Error(), "from_account_id", "SAME_ACCOUNT")
		default:
			httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.List(r.Context())
	if err != nil {
		httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpresponse.WriteError(w, 400, "invalid id", "id", "INVALID_ID")
		return
	}
	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.v.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Update(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, ErrTransactionNotFound):
			httpresponse.WriteError(w, 404, "transaction not found", "id", "TRANSACTION_NOT_FOUND")
		case errors.Is(err, ErrInvalidTransactionID):
			httpresponse.WriteError(w, 400, "invalid transaction id", "id", "INVALID_ID")
		case errors.Is(err, ErrTransactionNameRequired):
			httpresponse.WriteError(w, 400, err.Error(), "name", "TRANSACTION_NAME_REQUIRED")
		case errors.Is(err, ErrTransactionAmountInvalid):
			httpresponse.WriteError(w, 400, err.Error(), "amount", "INVALID_AMOUNT")
		case errors.Is(err, ErrTransactionSameAccount):
			httpresponse.WriteError(w, 400, err.Error(), "from_account_id", "SAME_ACCOUNT")
		default:
			httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpresponse.WriteError(w, 400, "invalid id", "id", "INVALID_ID")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrTransactionNotFound):
			httpresponse.WriteError(w, 404, "transaction not found", "id", "TRANSACTION_NOT_FOUND")
		case errors.Is(err, ErrInvalidTransactionID):
			httpresponse.WriteError(w, 400, "invalid transaction id", "id", "INVALID_ID")
		default:
			httpresponse.WriteError(w, 500, "internal server error", "", "INTERNAL_ERROR")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", h.CreateTransaction)
		r.Get("/", h.ListTransactions)
		r.Put("/{id}", h.UpdateTransaction)
		r.Delete("/{id}", h.DeleteTransaction)
	})
}
