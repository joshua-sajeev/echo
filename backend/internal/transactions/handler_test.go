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
	r.Put("/transactions/{id}", h.UpdateTransaction)
	r.Delete("/transactions/{id}", h.DeleteTransaction)

	return r
}

func TestCreateTransactionHandler(t *testing.T) {
	mockService := &MockTransactionService{
		CreateFunc: func(ctx context.Context, request CreateTransactionRequest) (int64, error) {
			return 101, nil
		},
	}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	body := `{
  "name":"Salary",
  "type":"income",
  "amount":5000,
  "date":"2026-05-30T10:00:00Z"
}`

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var resp struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.ID != 101 {
		t.Fatalf("expected id 101, got %d", resp.ID)
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
	name := "Updated"
	amount := int64(999)

	mockService := &MockTransactionService{
		UpdateFunc: func(
			ctx context.Context,
			id int64,
			req UpdateTransactionRequest,
		) error {
			if id != 1 {
				t.Fatalf("expected id 1, got %d", id)
			}
			return nil
		},
	}

	handler := NewTransactionHandler(mockService)
	router := setupRouter(handler)

	body, _ := json.Marshal(UpdateTransactionRequest{
		Name:   &name,
		Amount: &amount,
	})

	req := httptest.NewRequest(
		http.MethodPut,
		"/transactions/1",
		bytes.NewBuffer(body),
	)
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
