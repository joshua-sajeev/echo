package jars

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupRouter(h *JarHandler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func TestJarHandlerCreate(t *testing.T) {
	mock := &MockJarService{
		CreateJarFunc: func(ctx context.Context, jar CreateJarRequest) (int64, error) {
			return 10, nil
		},
	}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	body := `{"name":"Savings","allocation_type":"percentage","value":20}`
	req := httptest.NewRequest(http.MethodPost, "/jars/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rec.Code)
	}
}

func TestJarHandlerCreateValidationFail(t *testing.T) {
	mock := &MockJarService{}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	body := `{"name":"Savings","allocation_type":"fixed_amount","value":20}`
	req := httptest.NewRequest(http.MethodPost, "/jars/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestJarHandlerList(t *testing.T) {
	mock := &MockJarService{
		ListJarsFunc: func(ctx context.Context) ([]Jar, error) {
			return []Jar{
				{Name: "A"},
				{Name: "B"},
			}, nil
		},
	}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/jars/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestJarHandlerUpdate(t *testing.T) {
	mock := &MockJarService{
		UpdateJarFunc: func(ctx context.Context, id int64, jar UpdateJarRequest) error {
			return nil
		},
	}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	body := `{"name":"Updated","allocation_type":"percentage","value":30}`
	req := httptest.NewRequest(http.MethodPut, "/jars/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestJarHandlerDelete(t *testing.T) {
	mock := &MockJarService{
		DeleteJarFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/jars/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", rec.Code)
	}
}
