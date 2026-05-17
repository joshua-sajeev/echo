// Package service
package service

import (
	"context"

	"github.com/joshu-sajeev/echo/internal/models"
)

type AccountServiceInterface interface {
	Create(ctx context.Context, name string) (int64, error)
	List(ctx context.Context) ([]models.Account, error)
	ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
	ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
	Rename(ctx context.Context, id int64, name string) error
	Archive(ctx context.Context, id int64) error
	Unarchive(ctx context.Context, id int64) error
}
