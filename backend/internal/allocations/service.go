package allocations

import (
	"context"

	"github.com/joshu-sajeev/echo/internal/goals"
)

type AllocationServiceInterface interface {
	// Manual: allocate to one specific goal
	RunManual(
		ctx context.Context,
		goalID int64,
		amount int64,
	) error

	// Automatic: distribute across all goals by percentage
	DistributeAutomatic(
		ctx context.Context,
		allocType string,
		amount *int64,
	) error

	// Get last month's leisure leftover amount and already allocated status
	GetLeisureLeftover(ctx context.Context) (int64, bool, error)
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

// GetLeisureLeftover retrieves the calculated leftover for the leisure jar from last month and checks if already allocated
func (s *AllocationService) GetLeisureLeftover(ctx context.Context) (int64, bool, error) {
	leftover, err := s.repo.GetLastMonthLeisureLeftover(ctx)
	if err != nil {
		return 0, false, err
	}

	alreadyAllocated, err := s.repo.IsAlreadyAllocatedThisMonth(ctx)
	if err != nil {
		return leftover, false, err
	}

	return leftover, alreadyAllocated, nil
}

// RunManual allocates amount to a specific goal
func (s *AllocationService) RunManual(
	ctx context.Context,
	goalID int64,
	amount int64,
) error {
	// Validate inputs
	if goalID <= 0 {
		return ErrInvalidGoalID
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}

	// Verify goal exists and not archived
	goal, err := s.goalsRepo.GetByID(ctx, goalID)
	if err != nil {
		return err
	}

	if goal == nil {
		return ErrGoalNotFound
	}

	if goal.IsArchived {
		return ErrGoalArchived
	}

	// Allocate to specific goal
	return s.repo.RunManual(ctx, goalID, amount)
}

// DistributeAutomatic distributes amount across all goals by their allocation_percentage
func (s *AllocationService) DistributeAutomatic(
	ctx context.Context,
	allocType string,
	amount *int64,
) error {
	var targetAmount int64

	switch allocType {
	case "automatic_splitting":
		if amount == nil || *amount <= 0 {
			return ErrInvalidAmount
		}
		targetAmount = *amount

	case "leisure_leftover":
		leftover, err := s.repo.GetLastMonthLeisureLeftover(ctx)
		if err != nil {
			return err
		}
		if leftover <= 0 {
			return ErrInvalidAmount // No leftover to allocate
		}
		targetAmount = leftover

	default:
		return ErrInvalidAllocationType
	}

	// Get all non-archived goals
	goalsList, err := s.goalsRepo.List(ctx)
	if err != nil {
		return err
	}

	if len(goalsList) == 0 {
		return ErrNoGoalsConfigured
	}

	// Validate percentages sum to 100%
	var totalPercentage int64
	for _, goal := range goalsList {
		totalPercentage += goal.AllocationPercentage
	}

	if totalPercentage != 100 {
		return ErrInvalidAllocationPercentages
	}

	// Distribute by percentage
	return s.repo.DistributeAutomatic(ctx, targetAmount, goalsList)
}
