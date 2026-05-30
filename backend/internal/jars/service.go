package jars

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joshu-sajeev/echo/internal/utils"
)

type JarService struct {
	repo JarRepositoryInterface
}

type JarServiceInterface interface {
	CreateJar(ctx context.Context, jar CreateJarRequest) (int64, error)
	ListJars(ctx context.Context) ([]Jar, error)
	UpdateJar(ctx context.Context, id int64, jar UpdateJarRequest) error
	DeleteJar(ctx context.Context, id int64) error
}

var _ JarServiceInterface = (*JarService)(nil)

func NewJarService(repo JarRepositoryInterface) *JarService {
	return &JarService{repo: repo}
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
			utils.LogError(ctx, "JarService.CreateJar (totalPercentage)", err)
			return 0, err
		}

		if total+request.Value > 100 {
			return 0, ErrTotalPercentageExceeded
		}
	}

	jar := Jar{
		Name:           request.Name,
		AllocationType: AllocationType(request.AllocationType),
		Value:          request.Value,
		Priority:       request.Priority,
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

	if request.Priority != nil {
		current.Priority = *request.Priority
	}

	if current.AllocationType == AllocationPercentage {
		if current.Value <= 0 {
			return ErrPercentageMustBePositive
		}

		total, err := s.totalPercentage(ctx, id)
		if err != nil {
			utils.LogError(ctx, "JarService.UpdateJar (totalPercentage)", err)
			return err
		}

		if total+current.Value > 100 {
			return ErrTotalPercentageExceeded
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

func (s *JarService) totalPercentage(ctx context.Context, excludeID int64) (int64, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, jar := range jars {
		// Only apply target exclusion if a genuine database tracking ID is requested
		if excludeID > 0 && jar.ID == excludeID {
			continue
		}
		if jar.AllocationType == AllocationPercentage {
			total += jar.Value
		}
	}

	return total, nil
}
