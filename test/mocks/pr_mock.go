package mocks

import (
	"context"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockPRRepository struct {
	mock.Mock
}

func (m *MockPRRepository) Create(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
	args := m.Called(ctx, pr)
	return args.Get(0).(models.PullRequest), args.Error(1)
}

func (m *MockPRRepository) Merge(ctx context.Context, id string) (models.PullRequest, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.PullRequest), args.Error(1)
}

func (m *MockPRRepository) Reassign(ctx context.Context, prID, oldReviewerID, newReviewerID string) (models.PullRequest, error) {
	args := m.Called(ctx, prID, oldReviewerID, newReviewerID)
	return args.Get(0).(models.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetPRs(ctx context.Context, userID string) ([]models.PullRequest, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetByID(ctx context.Context, id string) (models.PullRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return models.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(models.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetReviewers(ctx context.Context, prID string) ([]models.User, error) {
	args := m.Called(ctx, prID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockPRRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
