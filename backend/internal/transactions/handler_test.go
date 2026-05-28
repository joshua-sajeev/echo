package transactions

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupRouter(h *TransactionHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/transactions", h.CreateTransaction)
	r.Get("/transactions", h.ListTransactions)
	r.Put("/transactions", h.UpdateTransaction)
	r.Delete("/transactions/{id}", h.DeleteTransaction)

	return r
}

func TestCreateTransactionHandler(t *testing.T) {
	mockService := &MockTransactionService{
		CreateFunc: func(ctx context.Context, tx Transaction) (int64, error) {
			return 101, nil
		},
	}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	body := `{
		"name": "Salary",
		"type": "income",
		"amount": 5000
	}`

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var resp map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	if resp["id"].(int64) != 101 {
		t.Fatalf("expected id 101, got %v", resp["id"])
	}
}

func TestCreateTransaction_InvalidBody(t *testing.T) {
	mockService := &MockTransactionService{}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{invalid-json`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestListTransactionsHandler(t *testing.T) {
	mockService := &MockTransactionService{
		ListFunc: func(ctx context.Context) ([]Transaction, error) {
			return []Transaction{
				{ID: 1, Name: "A", Amount: 100},
			}, nil
		},
	}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUpdateTransactionHandler(t *testing.T) {
	mockService := &MockTransactionService{
		UpdateFunc: func(ctx context.Context, tx Transaction) error {
			return nil
		},
	}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	body := `{
		"id": 1,
		"name": "Updated",
		"amount": 999
	}`

	req := httptest.NewRequest(http.MethodPut, "/transactions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestDeleteTransactionHandler(t *testing.T) {
	mockService := &MockTransactionService{
		DeleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/transactions/10", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
