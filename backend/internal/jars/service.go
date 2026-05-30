package jars

import (
	"context"
	"errors"

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
	if request.Name == "" {
		return 0, ErrJarNameRequired
	}

	if request.AllocationType == string(AllocationPercentage) {
		if request.Value <= 0 {
			return 0, ErrPercentageMustBePositive
		}

		total, err := s.totalPercentage(ctx)
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

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
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

	// validate name if provided
	if request.Name != nil && *request.Name == "" {
		return ErrJarNameRequired
	}

	// percentage logic only if provided
	if request.AllocationType != nil && *request.AllocationType == string(AllocationPercentage) {

		if request.Value != nil && *request.Value <= 0 {
			return ErrPercentageMustBePositive
		}

		jars, err := s.repo.List(ctx)
		if err != nil {
			utils.LogError(ctx, "JarService.UpdateJar (repo.List)", err)
			return err
		}

		var total int64
		for _, j := range jars {
			if j.ID != id && j.AllocationType == AllocationPercentage {
				total += j.Value
			}
		}

		if request.Value != nil && total+*request.Value > 100 {
			return ErrTotalPercentageExceeded
		}
	}

	// build update struct safely
	jar := Jar{ID: id}

	if request.Name != nil {
		jar.Name = *request.Name
	}

	if request.AllocationType != nil {
		jar.AllocationType = AllocationType(*request.AllocationType)
	}

	if request.Priority != nil {
		jar.Priority = *request.Priority
	}

	if request.Value != nil {
		jar.Value = *request.Value
	}

	return s.repo.Update(ctx, jar)
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

func (s *JarService) totalPercentage(ctx context.Context) (int64, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
		// Let the caller handle the LogError wrapper to prevent double logging
		return 0, err
	}

	var total int64
	for _, jar := range jars {
		if jar.AllocationType == AllocationPercentage {
			total += jar.Value
		}
	}

	return total, nil
}
