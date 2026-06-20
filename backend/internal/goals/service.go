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

// Create creates a new goal
func (s *GoalService) Create(ctx context.Context, request CreateGoalRequest) (int64, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return 0, ErrGoalNameRequired
	}

	if request.TargetAmount <= 0 {
		return 0, ErrTargetAmountInvalid
	}

	// Validate deadline is in the future
	now := time.Now().UTC()

	if request.Deadline != nil && request.Deadline.UTC().Before(now) {
		return 0, ErrDeadlinePassed
	}

	goal := Goal{
		Name:         name,
		TargetAmount: request.TargetAmount,
		SavedAmount:  0,
		Deadline:     request.Deadline,
	}

	return s.repo.Create(ctx, goal)
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
