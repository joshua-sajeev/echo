package jars

import (
	"context"
	"fmt"

	"github.com/joshu-sajeev/echo/internal/utils"
)

type JarService struct {
	repo JarRepositoryInterface
}

type JarServiceInterface interface {
	CreateJar(ctx context.Context, jar Jar) (int64, error)
	ListJars(ctx context.Context) ([]Jar, error)
	UpdateJar(ctx context.Context, jar Jar) error
	DeleteJar(ctx context.Context, id int64) error
}

var _ JarServiceInterface = (*JarService)(nil)

func NewJarService(repo JarRepositoryInterface) *JarService {
	return &JarService{repo: repo}
}

func (s *JarService) CreateJar(ctx context.Context, jar Jar) (int64, error) {
	if jar.Name == "" {
		return 0, fmt.Errorf("jar name required")
	}

	if jar.AllocationType == AllocationPercentage {
		if jar.Value <= 0 {
			return 0, fmt.Errorf("percentage must be positive")
		}

		total, err := s.totalPercentage(ctx)
		if err != nil {
			// Log here because s.totalPercentage calls the DB repo layer
			utils.LogError(ctx, "JarService.CreateJar (totalPercentage)", err)
			return 0, err
		}

		if total+jar.Value > 100 {
			return 0, fmt.Errorf("total percentage exceeds 100")
		}
	}

	id, err := s.repo.Create(ctx, jar)
	if err != nil {
		// Catches Postgres errors (e.g. duplicate name unique key constraint)
		utils.LogError(ctx, "JarService.CreateJar (repo.Create)", err)
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

func (s *JarService) UpdateJar(ctx context.Context, jar Jar) error {
	if jar.ID == 0 {
		return fmt.Errorf("invalid jar id")
	}

	if jar.Name == "" {
		return fmt.Errorf("jar name required")
	}

	if jar.AllocationType == AllocationPercentage {
		if jar.Value <= 0 {
			return fmt.Errorf("percentage must be positive")
		}

		jars, err := s.repo.List(ctx)
		if err != nil {
			utils.LogError(ctx, "JarService.UpdateJar (repo.List)", err)
			return err
		}

		var total float64
		for _, j := range jars {
			if j.ID != jar.ID && j.AllocationType == AllocationPercentage {
				total += j.Value
			}
		}

		if total+jar.Value > 100 {
			return fmt.Errorf("total percentage exceeds 100")
		}
	}

	err := s.repo.Update(ctx, jar)
	if err != nil {
		// Catches database runtime errors or "jar not found" string match

		utils.LogError(ctx, "JarService.UpdateJar (repo.Update)", err)
		return err
	}

	return nil
}

func (s *JarService) DeleteJar(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid jar id")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		utils.LogError(ctx, "JarService.DeleteJar", err)
		return err
	}

	return nil
}

func (s *JarService) totalPercentage(ctx context.Context) (float64, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
		// Let the caller handle the LogError wrapper to prevent double logging
		return 0, err
	}

	var total float64
	for _, jar := range jars {
		if jar.AllocationType == AllocationPercentage {
			total += jar.Value
		}
	}

	return total, nil
}
