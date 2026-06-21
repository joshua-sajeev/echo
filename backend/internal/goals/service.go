package goals

import (
	"context"
	"strings"
	"time"
)

type GoalService struct {
	repo GoalRepositoryInterface
}

type GoalServiceInterface interface {
	Create(ctx context.Context, request CreateGoalRequest) (int64, error)
	CreateWithRebalance(ctx context.Context, request CreateGoalWithRebalanceRequest) (*RebalanceSummary, error)
	List(ctx context.Context) ([]GoalWithProgress, error)
	GetByID(ctx context.Context, id int64) (*GoalWithProgress, error)
	Update(ctx context.Context, id int64, request UpdateGoalRequest) error
	AddProgress(ctx context.Context, id int64, amount int64) error
	Archive(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
}

var _ GoalServiceInterface = (*GoalService)(nil)

// NewGoalService creates a new goal service
func NewGoalService(repo GoalRepositoryInterface) *GoalService {
	return &GoalService{repo: repo}
}

// Create creates a new goal (without rebalancing)
// This allows the first goal to be 100%
func (s *GoalService) Create(ctx context.Context, request CreateGoalRequest) (int64, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return 0, ErrGoalNameRequired
	}

	if request.TargetAmount <= 0 {
		return 0, ErrTargetAmountInvalid
	}

	if request.AllocationPercentage < 0 || request.AllocationPercentage > 100 {
		return 0, ErrTargetAmountInvalid
	}

	// Validate deadline is in the future
	now := time.Now().UTC()

	if request.Deadline != nil && request.Deadline.UTC().Before(now) {
		return 0, ErrDeadlinePassed
	}

	goal := Goal{
		Name:                 name,
		TargetAmount:         request.TargetAmount,
		SavedAmount:          0,
		Deadline:             request.Deadline,
		AllocationPercentage: request.AllocationPercentage,
	}

	return s.repo.Create(ctx, goal)
}

// CreateWithRebalance creates a new goal and rebalances all others to total 100%
// IMPORTANT: Can only be used when there are existing goals
// If there are no existing goals, use Create() instead and the new goal will be 100%
func (s *GoalService) CreateWithRebalance(ctx context.Context, request CreateGoalWithRebalanceRequest) (*RebalanceSummary, error) {
	// Get all existing goals
	existingGoals, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// If no existing goals, this endpoint shouldn't be used
	// Use Create() instead for the first goal
	if len(existingGoals) == 0 {
		return nil, ErrInvalidAllocationCount // "must create first goal with POST /goals, then use /goals/rebalance for subsequent goals"
	}

	// Now we require at least the existing goals + 1 new goal
	if len(request.Goals) < 2 {
		return nil, ErrInvalidAllocationCount
	}

	// Validate total allocation equals 100
	totalAllocation := int64(0)
	for _, item := range request.Goals {
		totalAllocation += item.AllocationPercentage
	}

	if totalAllocation != 100 {
		return nil, ErrAllocationPercentageTotal
	}

	// Validate all items
	newGoalIndex := -1
	for i, item := range request.Goals {
		if item.Name == "" {
			return nil, ErrGoalNameRequired
		}
		if item.TargetAmount <= 0 {
			return nil, ErrTargetAmountInvalid
		}
		if item.AllocationPercentage < 0 || item.AllocationPercentage > 100 {
			return nil, ErrTargetAmountInvalid
		}

		// Validate deadline is in the future
		now := time.Now().UTC()
		if item.Deadline != nil && item.Deadline.UTC().Before(now) {
			return nil, ErrDeadlinePassed
		}

		// Find the new goal (no ID)
		if item.ID == nil {
			if newGoalIndex != -1 {
				return nil, ErrInvalidAllocationCount // Only one new goal allowed
			}
			newGoalIndex = i
		}
	}

	if newGoalIndex == -1 {
		return nil, ErrInvalidAllocationCount // Must have at least one new goal
	}

	summary := &RebalanceSummary{
		Changes: make([]AllocationChange, 0),
		Message: "Goals rebalanced successfully to total 100%",
	}

	// Process existing goals (update their allocations)
	for i, item := range request.Goals {
		if i == newGoalIndex {
			continue // Skip new goal for now
		}

		// This is an existing goal - must have an ID
		if item.ID == nil {
			return nil, ErrGoalNotInRebalance
		}

		existingGoal, err := s.repo.GetByID(ctx, *item.ID)
		if err != nil {
			return nil, err
		}

		oldPercentage := existingGoal.AllocationPercentage
		newPercentage := item.AllocationPercentage

		// Record the change
		summary.Changes = append(summary.Changes, AllocationChange{
			ID:            existingGoal.ID,
			Name:          existingGoal.Name,
			OldPercentage: oldPercentage,
			NewPercentage: newPercentage,
			IsNew:         false,
		})

		// Update the goal
		existingGoal.AllocationPercentage = newPercentage
		existingGoal.Name = item.Name
		existingGoal.TargetAmount = item.TargetAmount
		existingGoal.Deadline = item.Deadline

		if err := s.repo.Update(ctx, *existingGoal); err != nil {
			return nil, err
		}
	}

	// Create the new goal
	newItem := request.Goals[newGoalIndex]
	newGoal := Goal{
		Name:                 strings.TrimSpace(newItem.Name),
		TargetAmount:         newItem.TargetAmount,
		SavedAmount:          0,
		Deadline:             newItem.Deadline,
		AllocationPercentage: newItem.AllocationPercentage,
	}

	newID, err := s.repo.Create(ctx, newGoal)
	if err != nil {
		return nil, err
	}

	// Record the new goal in summary
	summary.Changes = append(summary.Changes, AllocationChange{
		ID:            newID,
		Name:          newGoal.Name,
		OldPercentage: 0,
		NewPercentage: newGoal.AllocationPercentage,
		IsNew:         true,
	})

	return summary, nil
}

// List returns all goals with their progress calculated
func (s *GoalService) List(ctx context.Context) ([]GoalWithProgress, error) {
	goals, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]GoalWithProgress, len(goals))
	for i, goal := range goals {
		result[i] = NewGoalWithProgress(goal)
	}

	return result, nil
}

// GetByID retrieves a goal by ID with progress calculated
func (s *GoalService) GetByID(ctx context.Context, id int64) (*GoalWithProgress, error) {
	if id <= 0 {
		return nil, ErrInvalidGoalID
	}

	goal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	goalWithProgress := NewGoalWithProgress(*goal)
	return &goalWithProgress, nil
}

// Update updates an existing goal
func (s *GoalService) Update(ctx context.Context, id int64, request UpdateGoalRequest) error {
	if id <= 0 {
		return ErrInvalidGoalID
	}

	goal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields if provided
	if request.Name != nil {
		trimmedName := strings.TrimSpace(*request.Name)
		if trimmedName == "" {
			return ErrGoalNameRequired
		}
		goal.Name = trimmedName
	}

	if request.TargetAmount != nil {
		if *request.TargetAmount <= 0 {
			return ErrTargetAmountInvalid
		}
		goal.TargetAmount = *request.TargetAmount
	}

	if request.AllocationPercentage != nil {
		if *request.AllocationPercentage < 0 || *request.AllocationPercentage > 100 {
			return ErrTargetAmountInvalid
		}
		goal.AllocationPercentage = *request.AllocationPercentage
	}

	if request.Deadline != nil {
		// Validate new deadline is in the future
		if request.Deadline.Before(time.Now()) {
			return ErrDeadlinePassed
		}
		goal.Deadline = request.Deadline
	}

	return s.repo.Update(ctx, *goal)
}

// AddProgress adds progress towards a goal
func (s *GoalService) AddProgress(ctx context.Context, id int64, amount int64) error {
	if id <= 0 {
		return ErrInvalidGoalID
	}

	if amount <= 0 {
		return ErrProgressAmountInvalid
	}

	goal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if goal is already completed
	if goal.SavedAmount >= goal.TargetAmount {
		return ErrGoalAlreadyCompleted
	}

	return s.repo.AddProgress(ctx, id, amount)
}

// Archive removes a goal
func (s *GoalService) Archive(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidGoalID
	}

	return s.repo.Archive(ctx, id)
}

// Unarchives a goal
func (s *GoalService) Restore(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidGoalID
	}

	return s.repo.Restore(ctx, id)
}
