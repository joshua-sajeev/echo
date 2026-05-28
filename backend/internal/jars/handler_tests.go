package jars

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupRouter(h *JarHandler) http.Handler {
	r := chi.NewRouter()
	r.Mount("/", h.Routes())
	return r
}

func TestJarHandler_Create(t *testing.T) {
	mock := &MockJarService{
		CreateJarFunc: func(ctx context.Context, jar Jar) (int64, error) {
			return 10, nil
		},
	}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/jars",
		bytes.NewBufferString(`{"name":"Savings","allocation_type":"percentage","value":20}`))

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rec.Code)
	}
}

func TestJarHandler_List(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/jars", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestJarHandler_Update(t *testing.T) {
	mock := &MockJarService{
		UpdateJarFunc: func(ctx context.Context, jar Jar) error {
			return nil
		},
	}

	handler := NewJarHandler(mock)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/jars/1",
		bytes.NewBufferString(`{"name":"Updated","allocation_type":"fixed_amount","value":100}`))

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestJarHandler_Delete(t *testing.T) {
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
