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

func setupRouter(h *AccountHandler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func newRouter(mock *MockAccountService) *chi.Mux {
	return setupRouter(NewAccountHandler(mock))
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

func TestAccountHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		createFn     func(context.Context, string) (int64, error)
		wantStatus   int
		wantResponse *CreateAccountResponse
	}{
		{
			name: "success",
			body: `{"name":"Savings"}`,
			createFn: func(ctx context.Context, name string) (int64, error) {
				return 1, nil
			},
			wantStatus: http.StatusCreated,
			wantResponse: &CreateAccountResponse{
				ID: 1,
			},
		},
		{
			name:       "invalid json",
			body:       `{bad json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid account name",
			body: `{"name":"Savings"}`,
			createFn: func(ctx context.Context, name string) (int64, error) {
				return 0, ErrInvalidAccountName
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "account already exists",
			body: `{"name":"Savings"}`,
			createFn: func(ctx context.Context, name string) (int64, error) {
				return 0, ErrAccountAlreadyExists
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			body: `{"name":"Savings"}`,
			createFn: func(ctx context.Context, name string) (int64, error) {
				return 0, context.DeadlineExceeded
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockAccountService{
				CreateFn: tt.createFn,
			}

			router := newRouter(mock)

			req := httptest.NewRequest(
				http.MethodPost,
				"/accounts",
				bytes.NewBufferString(tt.body),
			)

			rec := executeRequest(router, req)

			assertStatus(t, rec.Code, tt.wantStatus)

			if tt.wantResponse != nil {
				var resp CreateAccountResponse

				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}

				if resp.ID != tt.wantResponse.ID {
					t.Fatalf("expected id=%d got=%d", tt.wantResponse.ID, resp.ID)
				}
			}
		})
	}
}

func TestAccountHandler_List(t *testing.T) {
	mock := &MockAccountService{
		ListFn: func(ctx context.Context) ([]Account, error) {
			return []Account{
				{
					ID:   1,
					Name: "Savings",
				},
			}, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts",
		nil,
	)

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusOK)

	var accounts []Account

	if err := json.NewDecoder(rec.Body).Decode(&accounts); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(accounts) != 1 {
		t.Fatalf("expected 1 account got %d", len(accounts))
	}
}

func TestAccountHandler_ListWithBalances(t *testing.T) {
	mock := &MockAccountService{
		ListWithBalancesFn: func(ctx context.Context) ([]AccountWithBalance, error) {
			return []AccountWithBalance{
				{
					Account: Account{
						ID:   1,
						Name: "Savings",
					},
					Balance: 100,
				},
			}, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts/balances",
		nil,
	)

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusOK)
}

func TestAccountHandler_ListArchivedWithBalances(t *testing.T) {
	mock := &MockAccountService{
		ListArchivedWithBalancesFn: func(ctx context.Context) ([]AccountWithBalance, error) {
			return []AccountWithBalance{
				{
					Account: Account{
						ID:   1,
						Name: "Archived",
					},
					Balance: 0,
				},
			}, nil
		},
	}

	router := newRouter(mock)

	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts/archived",
		nil,
	)

	rec := executeRequest(router, req)

	assertStatus(t, rec.Code, http.StatusOK)
}

func TestAccountHandler_Rename(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		body       string
		renameFn   func(context.Context, int64, string) error
		wantStatus int
	}{
		{
			name: "success",
			path: "/accounts/1/rename",
			body: `{"name":"Emergency Fund"}`,
			renameFn: func(ctx context.Context, id int64, name string) error {
				return nil
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			path:       "/accounts/abc/rename",
			body:       `{"name":"Emergency Fund"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			path:       "/accounts/1/rename",
			body:       `{bad`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid account name",
			path: "/accounts/1/rename",
			body: `{"name":"Emergency Fund"}`,
			renameFn: func(ctx context.Context, id int64, name string) error {
				return ErrInvalidAccountName
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate name",
			path: "/accounts/1/rename",
			body: `{"name":"Emergency Fund"}`,
			renameFn: func(ctx context.Context, id int64, name string) error {
				return ErrAccountAlreadyExists
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "not found",
			path: "/accounts/1/rename",
			body: `{"name":"Emergency Fund"}`,
			renameFn: func(ctx context.Context, id int64, name string) error {
				return ErrAccountNotFound
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockAccountService{
				RenameFn: tt.renameFn,
			}

			router := newRouter(mock)

			req := httptest.NewRequest(
				http.MethodPatch,
				tt.path,
				bytes.NewBufferString(tt.body),
			)

			rec := executeRequest(router, req)

			assertStatus(t, rec.Code, tt.wantStatus)
		})
	}
}

func TestAccountHandler_Archive(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		serviceError error
		wantStatus   int
	}{
		{
			name:       "success",
			path:       "/accounts/1/archive",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			path:       "/accounts/abc/archive",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "not found",
			path:         "/accounts/1/archive",
			serviceError: ErrAccountNotFound,
			wantStatus:   http.StatusNotFound,
		},
		{
			name:         "already archived",
			path:         "/accounts/1/archive",
			serviceError: ErrAccountAlreadyArchived,
			wantStatus:   http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockAccountService{
				ArchiveFn: func(ctx context.Context, id int64) error {
					return tt.serviceError
				},
			}

			router := newRouter(mock)

			req := httptest.NewRequest(
				http.MethodPatch,
				tt.path,
				nil,
			)

			rec := executeRequest(router, req)

			assertStatus(t, rec.Code, tt.wantStatus)
		})
	}
}

func TestAccountHandler_Unarchive(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		serviceError error
		wantStatus   int
	}{
		{
			name:       "success",
			path:       "/accounts/1/unarchive",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			path:       "/accounts/abc/unarchive",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "not found",
			path:         "/accounts/1/unarchive",
			serviceError: ErrAccountNotFound,
			wantStatus:   http.StatusNotFound,
		},
		{
			name:         "already active",
			path:         "/accounts/1/unarchive",
			serviceError: ErrAccountAlreadyActive,
			wantStatus:   http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockAccountService{
				UnarchiveFn: func(ctx context.Context, id int64) error {
					return tt.serviceError
				},
			}

			router := newRouter(mock)

			req := httptest.NewRequest(
				http.MethodPatch,
				tt.path,
				nil,
			)

			rec := executeRequest(router, req)

			assertStatus(t, rec.Code, tt.wantStatus)
		})
	}
}
