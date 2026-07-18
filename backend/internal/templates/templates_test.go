package templates

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/utils"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB = utils.GetTestDB()

	code := m.Run()

	utils.CleanupTestDB()

	os.Exit(code)
}

func resetTemplatesTable(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(context.Background(), "TRUNCATE TABLE transaction_templates RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate transaction_templates: %v", err)
	}
}

func TestRepository_CRUD(t *testing.T) {
	resetTemplatesTable(t)

	ctx := context.Background()
	repo := NewRepository(testDB)

	amount := int64(1000)
	category := "Salary"
	tpl := Template{
		TemplateName:   "Monthly Salary",
		Type:           "income",
		Amount:         &amount,
		Name:           "Employer",
		Category:       &category,
		IsMasterIncome: true,
	}

	// 1. Create
	id, err := repo.Create(ctx, tpl)
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero id")
	}

	// 2. List
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list templates: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 template, got %d", len(list))
	}
	if list[0].TemplateName != tpl.TemplateName {
		t.Fatalf("expected template name %q, got %q", tpl.TemplateName, list[0].TemplateName)
	}

	// 3. Update
	newAmount := int64(2000)
	tpl.Amount = &newAmount
	tpl.TemplateName = "Updated Monthly Salary"
	err = repo.Update(ctx, id, tpl)
	if err != nil {
		t.Fatalf("failed to update template: %v", err)
	}

	list2, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list templates: %v", err)
	}
	if *list2[0].Amount != newAmount {
		t.Fatalf("expected amount %d, got %d", newAmount, *list2[0].Amount)
	}
	if list2[0].TemplateName != "Updated Monthly Salary" {
		t.Fatalf("expected template name 'Updated Monthly Salary', got %q", list2[0].TemplateName)
	}

	// 4. Delete
	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete template: %v", err)
	}

	list3, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list templates: %v", err)
	}
	if len(list3) != 0 {
		t.Fatalf("expected 0 templates, got %d", len(list3))
	}
}

func TestHandler_CRUD(t *testing.T) {
	resetTemplatesTable(t)

	repo := NewRepository(testDB)
	handler := NewHandler(repo)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	amount := int64(500)
	reqBody := TemplateRequest{
		TemplateName:   "Coffee",
		Type:           "expense",
		Amount:         &amount,
		Name:           "Starbucks",
		IsMasterIncome: false,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// 1. POST /templates (Create)
	req := httptest.NewRequest(http.MethodPost, "/templates", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", rec.Code)
	}

	var createResp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}
	idVal, ok := createResp["id"].(float64)
	if !ok {
		t.Fatalf("expected template id in response")
	}
	id := int64(idVal)

	// 2. GET /templates (List)
	req = httptest.NewRequest(http.MethodGet, "/templates", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rec.Code)
	}

	var list []Template
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("failed to unmarshal list response: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 template, got %d", len(list))
	}
	if list[0].TemplateName != "Coffee" {
		t.Fatalf("expected template name 'Coffee', got %q", list[0].TemplateName)
	}

	// 3. PUT /templates/{id} (Update)
	newAmount := int64(600)
	reqBody.Amount = &newAmount
	reqBody.TemplateName = "Updated Coffee"
	bodyBytes2, _ := json.Marshal(reqBody)

	req = httptest.NewRequest(http.MethodPut, "/templates/"+strconvFormat(id), bytes.NewReader(bodyBytes2))
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content, got %d", rec.Code)
	}

	// Verify update
	req = httptest.NewRequest(http.MethodGet, "/templates", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if *list[0].Amount != 600 || list[0].TemplateName != "Updated Coffee" {
		t.Fatalf("update was not saved correctly")
	}

	// 4. DELETE /templates/{id} (Delete)
	req = httptest.NewRequest(http.MethodDelete, "/templates/"+strconvFormat(id), nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content, got %d", rec.Code)
	}

	// Verify delete
	req = httptest.NewRequest(http.MethodGet, "/templates", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Fatalf("expected 0 templates after delete, got %d", len(list))
	}
}

func strconvFormat(v int64) string {
	return javaString(v)
}

func javaString(v int64) string {
	return bytes.NewBufferString(jsonString(v)).String()
}

func jsonString(v int64) string {
	b, _ := json.Marshal(v)
	return string(b)
}
