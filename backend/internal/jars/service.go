package jars

import (
	"context"
	"fmt"
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
			return 0, err
		}

		if total+jar.Value > 100 {
			return 0, fmt.Errorf("total percentage exceeds 100")
		}
	}

	return s.repo.Create(ctx, jar)
}

func (s *JarService) ListJars(ctx context.Context) ([]Jar, error) {
	return s.repo.List(ctx)
}

func (s *JarService) UpdateJar(ctx context.Context, jar Jar) error {
	if jar.ID == 0 {
		return fmt.Errorf("invalid jar id")
	}

	if jar.Name == "" {
		return fmt.Errorf("jar name required")
	}

	// Optional: enforce same percentage rule on update
	if jar.AllocationType == AllocationPercentage {
		if jar.Value <= 0 {
			return fmt.Errorf("percentage must be positive")
		}

		jars, err := s.repo.List(ctx)
		if err != nil {
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

	return s.repo.Update(ctx, jar)
}

func (s *JarService) DeleteJar(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid jar id")
	}

	return s.repo.Delete(ctx, id)
}

func (s *JarService) totalPercentage(ctx context.Context) (float64, error) {
	jars, err := s.repo.List(ctx)
	if err != nil {
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
