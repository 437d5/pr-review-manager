package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
)

type UserService struct {
	uowFactory func(context.Context) (repositories.UnitOfWork, error)
}

func NewUserService(uowFactory func(ctx context.Context) (repositories.UnitOfWork, error)) *UserService {
	return &UserService{
		uowFactory: uowFactory,
	}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (models.User, error) {
	if userID == "" {
		return models.User{}, models.ErrEmptyUserID
	}

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer uow.Close()

	user, err := uow.Users().SetIsActive(ctx, userID, isActive)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.User{}, models.ErrUserNotFound
		}
		slog.Error("cannot set isActive for user", "error", err.Error(), "id", userID, "is_active", isActive)
		return models.User{}, err
	}

	return user, nil
}

func (s *UserService) GetPRs(ctx context.Context, userID string) ([]models.PullRequest, error) {
	if userID == "" {
		return []models.PullRequest{}, models.ErrEmptyUserID
	}

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return []models.PullRequest{}, err
	}
	defer uow.Close()

	prs, err := uow.PR().GetPRs(ctx, userID)
	if err != nil {
		slog.Error("cannot get PRs for user", "error", err.Error(), "id", userID)
		return []models.PullRequest{}, err
	}

	return prs, nil
}
