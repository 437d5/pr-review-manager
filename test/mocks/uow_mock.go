package mocks

import (
	"context"

	"github.com/437d5/pr-review-manager/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
)

type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Teams() repositories.TeamRepository {
	args := m.Called()
	return args.Get(0).(repositories.TeamRepository)
}

func (m *MockUnitOfWork) Users() repositories.UserRepository {
	args := m.Called()
	return args.Get(0).(repositories.UserRepository)
}

func (m *MockUnitOfWork) PR() repositories.PullRequestRepository {
	args := m.Called()
	return args.Get(0).(repositories.PullRequestRepository)
}

func (m *MockUnitOfWork) Begin(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUnitOfWork) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUnitOfWork) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUnitOfWork) Close() error {
	args := m.Called()
	return args.Error(0)
}
