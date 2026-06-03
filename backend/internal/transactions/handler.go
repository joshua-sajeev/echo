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

	return &TransactionHandler{
		service: service,
		v:       v,
	}
}

func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", h.CreateTransaction)
		r.Get("/", h.ListTransactions)
		r.Get("/{id}", h.GetTransaction)
		r.Put("/{id}", h.UpdateTransaction)
		r.Delete("/{id}", h.DeleteTransaction)
	})
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
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

	httpresponse.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.List(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, data)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid id", "id", "INVALID_ID")
		return
	}

	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	if err := h.v.ValidateStruct(req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "", "VALIDATION_FAILED")
		return
	}

	err = h.service.Update(r.Context(), id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpresponse.WriteError(w,
			http.StatusBadRequest,
			"invalid id",
			"id",
			"INVALID_ID")
		return
	}

	tx, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, tx)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid id", "id", "INVALID_ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTransactionNotFound):
		httpresponse.WriteError(w, http.StatusNotFound, "transaction not found", "id", "TRANSACTION_NOT_FOUND")

	case errors.Is(err, ErrInvalidTransactionID):
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid transaction id", "id", "INVALID_ID")

	case errors.Is(err, ErrTransactionNameRequired):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "name", "TRANSACTION_NAME_REQUIRED")

	case errors.Is(err, ErrTransactionTypeRequired):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "type", "TRANSACTION_TYPE_REQUIRED")

	case errors.Is(err, ErrTransactionAmountInvalid):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "amount", "INVALID_AMOUNT")

	case errors.Is(err, ErrTransactionSameAccount):
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error(), "from_account_id", "SAME_ACCOUNT")

	case errors.Is(err, ErrJarNotFound):
		httpresponse.WriteError(w, http.StatusNotFound, "jar not found", "jar_id", "JAR_NOT_FOUND")

	case errors.Is(err, ErrAccountNotFound):
		httpresponse.WriteError(w, http.StatusNotFound, "account not found", "from_account_id", "ACCOUNT_NOT_FOUND")

	default:
		httpresponse.WriteError(w, http.StatusInternalServerError, "internal server error", "", "INTERNAL_ERROR")
	}
}
