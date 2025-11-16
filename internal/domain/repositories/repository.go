package repositories

import (
	"context"

	"github.com/437d5/pr-review-manager/internal/domain/models"
)

type PullRequestRepository interface {
	Create(context.Context, models.PullRequest) (models.PullRequest, error)
	Merge(context.Context, string) (models.PullRequest, error)
	Reassign(context.Context, string, string, string) (models.PullRequest, error)
	GetPRs(context.Context, string) ([]models.PullRequest, error)
	GetByID(context.Context, string) (models.PullRequest, error)
	GetReviewers(context.Context, string) ([]models.User, error)
}

type TeamRepository interface {
	Create(context.Context, models.Team) (int, error)
	GetByName(context.Context, string) (models.Team, error)
	Exists(context.Context, string) (bool, error)
}

type UserRepository interface {
	Create(context.Context, models.User, int) error
	GetByID(context.Context, string) (models.User, error)
	Update(context.Context, models.User, int) (models.User, error)
	SetIsActive(context.Context, string, bool) (models.User, error)
	GetActiveTeammatesByUserID(context.Context, string) ([]models.User, error)
}

type UnitOfWork interface {
	Teams() TeamRepository
	Users() UserRepository
	PR() PullRequestRepository

	// Work with transactions
	Begin(context.Context) error
	Commit() error
	Rollback() error
	Close() error
}
