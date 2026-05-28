package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// ---------------- TEST HELPERS ----------------
func setupRouter(h *AccountHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/accounts", h.Create)
	r.Get("/accounts", h.List)
	r.Get("/accounts/balances", h.ListWithBalances)
	r.Get("/accounts/archived", h.ListArchivedWithBalances)
	r.Patch("/accounts/{id}/rename", h.Rename)
	r.Patch("/accounts/{id}/archive", h.Archive)
	r.Patch("/accounts/{id}/unarchive", h.Unarchive)

	return r
}

func newRouter(mock *MockAccountService) *chi.Mux {
	h := NewAccountHandler(mock)
	return setupRouter(h)
}

func executeRequest(router *chi.Mux, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d got %d", want, got)
	}
}

// ---------------- CREATE ----------------

func TestAccountHandler_Create_Success(t *testing.T) {
	mock := &MockAccountService{
		CreateFn: func(ctx context.Context, name string) (int64, error) {
			return 1, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPost, "/accounts",
		bytes.NewBufferString(`{"name":"test"}`))

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusCreated)

	var resp CreateAccountResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid response json: %v", err)
	}

	if resp.ID != 1 {
		t.Fatalf("expected id=1 got=%d", resp.ID)
	}
}

func TestAccountHandler_Create_InvalidJSON(t *testing.T) {
	mock := &MockAccountService{}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPost, "/accounts",
		bytes.NewBufferString("{bad json"))

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}

func TestAccountHandler_Create_ServiceError(t *testing.T) {
	mock := &MockAccountService{
		CreateFn: func(ctx context.Context, name string) (int64, error) {
			return 0, ErrInvalidAccountName
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPost, "/accounts",
		bytes.NewBufferString(`{"name":"test"}`))

	rec := executeRequest(router, req)

	// since ErrInvalidAccountName maps to 400 in handler
	assertStatus(t, rec.Code, http.StatusBadRequest)
}

// ---------------- LIST ----------------

func TestAccountHandler_List(t *testing.T) {
	mock := &MockAccountService{
		ListFn: func(ctx context.Context) ([]Account, error) {
			return []Account{
				{ID: 1, Name: "A"},
			}, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusOK)

	var res []Account
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if len(res) != 1 || res[0].Name != "A" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestAccountHandler_ListWithBalances(t *testing.T) {
	mock := &MockAccountService{
		ListWithBalancesFn: func(ctx context.Context) ([]AccountWithBalance, error) {
			return []AccountWithBalance{
				{Account: Account{ID: 1, Name: "A"}, Balance: 100},
			}, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/accounts/balances", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusOK)
}

func TestAccountHandler_ListArchivedWithBalances(t *testing.T) {
	mock := &MockAccountService{
		ListArchivedWithBalancesFn: func(ctx context.Context) ([]AccountWithBalance, error) {
			return []AccountWithBalance{
				{Account: Account{ID: 2, Name: "B"}, Balance: 0},
			}, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/accounts/archived", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusOK)
}

// ---------------- RENAME ----------------

func TestAccountHandler_Rename_Success(t *testing.T) {
	mock := &MockAccountService{
		RenameFn: func(ctx context.Context, id int64, name string) error {
			if id != 1 || name != "new" {
				t.Fatalf("unexpected args id=%d name=%s", id, name)
			}
			return nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch,
		"/accounts/1/rename",
		bytes.NewBufferString(`{"name":"new"}`))

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusNoContent)
}

func TestAccountHandler_Rename_InvalidID(t *testing.T) {
	mock := &MockAccountService{}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch,
		"/accounts/abc/rename",
		bytes.NewBufferString(`{"name":"new"}`))

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}

func TestAccountHandler_Rename_InvalidJSON(t *testing.T) {
	mock := &MockAccountService{}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch,
		"/accounts/1/rename",
		bytes.NewBufferString("{bad"))

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}

func TestAccountHandler_Rename_InvalidAccountError(t *testing.T) {
	mock := &MockAccountService{
		RenameFn: func(ctx context.Context, id int64, name string) error {
			return ErrInvalidAccountName
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch,
		"/accounts/1/rename",
		bytes.NewBufferString(`{"name":"new"}`))

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}

// ---------------- ARCHIVE ----------------

func TestAccountHandler_Archive_Success(t *testing.T) {
	mock := &MockAccountService{
		ArchiveFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch, "/accounts/1/archive", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusNoContent)
}

func TestAccountHandler_Archive_InvalidID(t *testing.T) {
	mock := &MockAccountService{
		ArchiveFn: func(ctx context.Context, id int64) error {
			return ErrInvalidAccountID
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch, "/accounts/1/archive", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}

// ---------------- UNARCHIVE ----------------

func TestAccountHandler_Unarchive_Success(t *testing.T) {
	mock := &MockAccountService{
		UnarchiveFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch, "/accounts/1/unarchive", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusNoContent)
}

func TestAccountHandler_Unarchive_InvalidID(t *testing.T) {
	mock := &MockAccountService{
		UnarchiveFn: func(ctx context.Context, id int64) error {
			return ErrInvalidAccountID
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(http.MethodPatch, "/accounts/abc/unarchive", nil)
	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}
