// Package dto contains request and response DTOs.
package dto

import (
	"time"

	"github.com/joshu-sajeev/echo/internal/models"
)

type CreateAccountRequest struct {
	Name string `json:"name"`
}

type RenameAccountRequest struct {
	Name string `json:"name"`
}

type AccountResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountWithBalanceResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func ToAccountResponse(a models.Account) AccountResponse {
	return AccountResponse{
		ID:        a.ID,
		Name:      a.Name,
		CreatedAt: a.CreatedAt,
	}
}

func ToAccountWithBalanceResponse(a models.AccountWithBalance) AccountWithBalanceResponse {
	return AccountWithBalanceResponse{
		ID:        a.ID,
		Name:      a.Name,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
	}
}
