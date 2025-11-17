package mocks

import (
	"context"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(ctx context.Context, team models.Team) (int, error) {
	args := m.Called(ctx, team)
	return args.Int(0), args.Error(1)
}

func (m *MockTeamRepository) GetByName(ctx context.Context, name string) (models.Team, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(models.Team), args.Error(1)
}

func (m *MockTeamRepository) Exists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}
