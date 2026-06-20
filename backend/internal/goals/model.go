package goals

import "time"

type Goal struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	TargetAmount         int64      `json:"target_amount"`
	SavedAmount          int64      `json:"saved_amount"`
	Deadline             *time.Time `json:"deadline"`
	AllocationPercentage int64      `json:"allocation_percentage"`
	IsArchived           bool       `json:"is_archived"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type GoalWithProgress struct {
	Goal
	Status   string `json:"status"`   // "not_started", "in_progress", "completed"
	Progress int64  `json:"progress"` // percentage 0-100
}

// CalculateStatus determines goal status based on saved vs target amount
func (g *Goal) CalculateStatus() string {
	if g.SavedAmount >= g.TargetAmount {
		return "completed"
	}
	if g.SavedAmount > 0 {
		return "in_progress"
	}
	return "not_started"
}

// CalculateProgress returns the progress percentage (0-100)
func (g *Goal) CalculateProgress() int64 {
	if g.TargetAmount <= 0 {
		return 0
	}

	progress := (g.SavedAmount * 100) / g.TargetAmount

	if progress > 100 {
		return 100
	}

	return progress
}

// NewGoalWithProgress creates a GoalWithProgress with calculated values
func NewGoalWithProgress(g Goal) GoalWithProgress {
	return GoalWithProgress{
		Goal:     g,
		Status:   g.CalculateStatus(),
		Progress: g.CalculateProgress(),
	}
}
