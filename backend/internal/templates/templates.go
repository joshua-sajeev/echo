package templates

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/httpresponse"
)

type Template struct {
	ID             int64     `json:"id"`
	TemplateName   string    `json:"template_name"`
	Type           string    `json:"type"`
	Amount         *int64    `json:"amount"`
	Name           string    `json:"name"`
	Category       *string   `json:"category"`
	AccountID      *int64    `json:"account_id"`
	FromAccountID  *int64    `json:"from_account_id"`
	ToAccountID    *int64    `json:"to_account_id"`
	JarID          *int64    `json:"jar_id"`
	IsMasterIncome bool      `json:"is_master_income"`
	CreatedAt      time.Time `json:"created_at"`
}

type TemplateRequest struct {
	TemplateName   string  `json:"template_name"`
	Type           string  `json:"type"`
	Amount         *int64  `json:"amount"`
	Name           string  `json:"name"`
	Category       *string `json:"category"`
	AccountID      *int64  `json:"account_id"`
	FromAccountID  *int64  `json:"from_account_id"`
	ToAccountID    *int64  `json:"to_account_id"`
	JarID          *int64  `json:"jar_id"`
	IsMasterIncome bool    `json:"is_master_income"`
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, t Template) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO transaction_templates (
			template_name, type, amount, name, category,
			account_id, from_account_id, to_account_id, jar_id, is_master_income
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, t.TemplateName, t.Type, t.Amount, t.Name, t.Category,
		t.AccountID, t.FromAccountID, t.ToAccountID, t.JarID, t.IsMasterIncome).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create template db: %w", err)
	}
	return id, nil
}

func (r *Repository) List(ctx context.Context) ([]Template, error) {
	templates := make([]Template, 0)
	rows, err := r.db.Query(ctx, `
		SELECT id, template_name, type, amount, name, category,
		       account_id, from_account_id, to_account_id, jar_id, is_master_income, created_at
		FROM transaction_templates
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list templates db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Template
		err := rows.Scan(&t.ID, &t.TemplateName, &t.Type, &t.Amount, &t.Name, &t.Category,
			&t.AccountID, &t.FromAccountID, &t.ToAccountID, &t.JarID, &t.IsMasterIncome, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan template db: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *Repository) Update(ctx context.Context, id int64, t Template) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE transaction_templates
		SET template_name = $1, type = $2, amount = $3, name = $4, category = $5,
		    account_id = $6, from_account_id = $7, to_account_id = $8, jar_id = $9, is_master_income = $10
		WHERE id = $11
	`, t.TemplateName, t.Type, t.Amount, t.Name, t.Category,
		t.AccountID, t.FromAccountID, t.ToAccountID, t.JarID, t.IsMasterIncome, id)
	if err != nil {
		return fmt.Errorf("update template db: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("template not found")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM transaction_templates WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete template db: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("template not found")
	}
	return nil
}

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/templates", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req TemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	t := Template{
		TemplateName:   req.TemplateName,
		Type:           req.Type,
		Amount:         req.Amount,
		Name:           req.Name,
		Category:       req.Category,
		AccountID:      req.AccountID,
		FromAccountID:  req.FromAccountID,
		ToAccountID:    req.ToAccountID,
		JarID:          req.JarID,
		IsMasterIncome: req.IsMasterIncome,
	}

	id, err := h.repo.Create(r.Context(), t)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error(), "", "DATABASE_ERROR")
		return
	}

	httpcallJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		httpcallJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid id", "id", "INVALID_ID")
		return
	}

	var req TemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body", "", "INVALID_REQUEST_BODY")
		return
	}

	t := Template{
		TemplateName:   req.TemplateName,
		Type:           req.Type,
		Amount:         req.Amount,
		Name:           req.Name,
		Category:       req.Category,
		AccountID:      req.AccountID,
		FromAccountID:  req.FromAccountID,
		ToAccountID:    req.ToAccountID,
		JarID:          req.JarID,
		IsMasterIncome: req.IsMasterIncome,
	}

	err = h.repo.Update(r.Context(), id, t)
	if err != nil {
		if err.Error() == "template not found" {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error(), "", "NOT_FOUND")
		} else {
			httpcallJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "invalid id", "id", "INVALID_ID")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "template not found" {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error(), "", "NOT_FOUND")
		} else {
			httpcallJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func httpcallJSON(w http.ResponseWriter, status int, data any) {
	httpresponse.WriteJSON(w, status, data)
}

func httpcallJSONError(w http.ResponseWriter, status int, msg string) {
	httpresponse.WriteError(w, status, msg, "", "SERVER_ERROR")
}
