package allocations

import (
	"context"

	"github.com/joshu-sajeev/echo/internal/goals"
)

type AllocationServiceInterface interface {
	Run(
		ctx context.Context,
		amount int64,
	) error
}

type AllocationService struct {
	repo      AllocationRepositoryInterface
	goalsRepo goals.GoalRepositoryInterface
}

func NewAllocationService(
	repo AllocationRepositoryInterface,
	goalRepo goals.GoalRepositoryInterface,
) *AllocationService {
	return &AllocationService{
		repo:      repo,
		goalsRepo: goalRepo,
	}
}

func (s *AllocationService) Run(
	ctx context.Context,
	amount int64,
) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	goalsList, err := s.goalsRepo.List(ctx)
	if err != nil {
		return err
	}

	if len(goalsList) == 0 {
		return ErrNoGoalsConfigured
	}

	var totalPercentage int64

	for _, goal := range goalsList {
		totalPercentage += goal.AllocationPercentage
	}

	if totalPercentage != 100 {
		return ErrInvalidAllocationPercentages
	}

	return s.repo.RunAllocation(
		ctx,
		amount,
		goalsList,
	)
}
