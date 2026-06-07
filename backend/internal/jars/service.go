package jars

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joshu-sajeev/echo/internal/transactions"
	"github.com/joshu-sajeev/echo/internal/utils"
)

type JarService struct {
	repo   JarRepositoryInterface
	txRepo transactions.TransactionRepositoryInterface
}

type JarServiceInterface interface {
	CreateJar(ctx context.Context, jar CreateJarRequest) (int64, error)
	ListJars(ctx context.Context) ([]Jar, error)
	ListJarAllocations(ctx context.Context) ([]JarWithAllocation, error)
	UpdateJar(ctx context.Context, id int64, jar UpdateJarRequest) error
	DeleteJar(ctx context.Context, id int64) error
}

var _ JarServiceInterface = (*JarService)(nil)

func NewJarService(
	repo JarRepositoryInterface,
	txRepo transactions.TransactionRepositoryInterface,
) *JarService {
	return &JarService{
		repo:   repo,
		txRepo: txRepo,
	}
}

func (s *JarService) CreateJar(ctx context.Context, request CreateJarRequest) (int64, error) {
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" {
		return 0, ErrJarNameRequired
	}

	if request.AllocationType == string(AllocationPercentage) {
		if request.Value <= 0 {
			return 0, ErrPercentageMustBePositive
		}

		total, err := s.totalPercentage(ctx, 0)
		if err != nil {
			return 0, err
		}

		if total+request.Value > 100 {
			return 0, ErrTotalPercentageExceeded
		}
	}

	if request.AllocationType == string(AllocationRemainder) {
		if request.Value != 0 {
			return 0, ErrRemainderMustBeZero
		}

		exists, err := s.hasRemainderJar(ctx, 0)
		if err != nil {
			return 0, err
		}

		if exists {
			return 0, ErrRemainderJarAlreadyExists
		}
	}

	jar := Jar{
		Name:           request.Name,
		AllocationType: AllocationType(request.AllocationType),
		Value:          request.Value,
	}

	id, err := s.repo.Create(ctx, jar)
	if err != nil {
		utils.LogError(ctx, "JarService.CreateJar (repo.Create)", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "jars_name_key" {
				return 0, ErrJarNameAlreadyExists
			}
		}
		return 0, err
	}

	return id, nil
}

func (s *JarService) ListJars(ctx context.Context) ([]Jar, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
		utils.LogError(ctx, "JarService.ListJars", err)
		return nil, err
	}
	return jars, nil
}

func (s *JarService) UpdateJar(ctx context.Context, id int64, request UpdateJarRequest) error {
	if id <= 0 {
		return ErrInvalidJarID
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrJarNotFound) {
			return ErrJarNotFound
		}
		utils.LogError(ctx, "JarService.UpdateJar (repo.GetByID)", err)
		return err
	}

	if request.Name != nil {
		trimmedName := strings.TrimSpace(*request.Name)
		if trimmedName == "" {
			return ErrJarNameRequired
		}
		current.Name = trimmedName
	}
	if request.AllocationType != nil {
		current.AllocationType = AllocationType(*request.AllocationType)
	}

	if request.Value != nil {
		current.Value = *request.Value
	}
	if current.AllocationType == AllocationPercentage {
		if current.Value <= 0 {
			return ErrPercentageMustBePositive
		}

		total, err := s.totalPercentage(ctx, id)
		if err != nil {
			return err
		}

		if total+current.Value > 100 {
			return ErrTotalPercentageExceeded
		}
	}

	if current.AllocationType == AllocationRemainder {
		if current.Value != 0 {
			return ErrRemainderMustBeZero
		}

		exists, err := s.hasRemainderJar(ctx, id)
		if err != nil {
			return err
		}

		if exists {
			return ErrRemainderJarAlreadyExists
		}
	}
	err = s.repo.Update(ctx, current)
	if err != nil {
		utils.LogError(ctx, "JarService.UpdateJar (repo.Update)", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "jars_name_key" {
				return ErrJarNameAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (s *JarService) DeleteJar(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidJarID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		utils.LogError(ctx, "JarService.DeleteJar", err)
		return err
	}

	return nil
}

// ListJarAllocations returns each jar with:
//   - AllocatedAmount: this month's share of master income
//   - Balance: all-time running balance
//   - SpentThisMonth: expenses charged to this jar in the current month
func (s *JarService) ListJarAllocations(ctx context.Context) ([]JarWithAllocation, error) {
	// ── current-month master income (for AllocatedAmount preview) ──────────
	monthIncome, err := s.txRepo.GetCurrentMonthIncome(ctx)
	if err != nil {
		utils.LogError(ctx, "JarService.ListJarAllocations (GetCurrentMonthIncome)", err)
		return nil, err
	}

	// ── all jars ────────────────────────────────────────────────────────────
	jars, err := s.repo.List(ctx)
	if err != nil {
		utils.LogError(ctx, "JarService.ListJarAllocations (repo.List)", err)
		return nil, err
	}

	// ── all-time running balances per jar ───────────────────────────────────
	balances, err := s.repo.GetAllJarBalances(ctx)
	if err != nil {
		utils.LogError(ctx, "JarService.ListJarAllocations (GetAllJarBalances)", err)
		return nil, err
	}

	// ── this month's spending per jar ───────────────────────────────────────
	spentThisMonth, err := s.repo.GetSpentThisMonthPerJar(ctx)
	if err != nil {
		utils.LogError(ctx, "JarService.ListJarAllocations (GetSpentThisMonthPerJar)", err)
		return nil, err
	}

	// ── build result ────────────────────────────────────────────────────────
	result := make([]JarWithAllocation, 0, len(jars))

	remaining := monthIncome
	remainderIndex := -1

	for _, jar := range jars {
		var allocated int64

		switch jar.AllocationType {
		case AllocationPercentage:
			allocated = monthIncome * jar.Value / 100
			remaining -= allocated

		case AllocationRemainder:
			// fill in after the loop
			remainderIndex = len(result)
		}

		result = append(result, JarWithAllocation{
			Jar:             jar,
			AllocatedAmount: allocated,
			Balance:         balances[jar.ID],
			SpentThisMonth:  spentThisMonth[jar.ID],
		})
	}

	if remainderIndex >= 0 {
		result[remainderIndex].AllocatedAmount = remaining
	}

	return result, nil
}

func (s *JarService) totalPercentage(ctx context.Context, excludeID int64) (int64, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, jar := range jars {
		if excludeID > 0 && jar.ID == excludeID {
			continue
		}
		if jar.AllocationType == AllocationPercentage {
			total += jar.Value
		}
	}

	return total, nil
}

func (s *JarService) hasRemainderJar(ctx context.Context, excludeID int64) (bool, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
		return false, err
	}

	for _, jar := range jars {
		if excludeID > 0 && jar.ID == excludeID {
			continue
		}

		if jar.AllocationType == AllocationRemainder {
			return true, nil
		}
	}

	return false, nil
}
