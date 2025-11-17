package mocks

import (
	"context"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user models.User, teamID int) error {
	args := m.Called(ctx, user, teamID)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user models.User, teamID int) (models.User, error) {
	args := m.Called(ctx, user, teamID)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) (models.User, error) {
	args := m.Called(ctx, userID, isActive)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetActiveTeammatesByUserID(ctx context.Context, userID string) ([]models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.User), args.Error(1)
}
