package services

import (
	"context"
	"testing"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
	"github.com/437d5/pr-review-manager/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_SetIsActive(t *testing.T) {
	ctx := context.Background()

	t.Run("successful set is active", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewUserService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		user := models.User{
			ID:       "user-1",
			Username: "user-1",
			IsActive: true,
			TeamName: "team-1",
		}

		mockUOW.On("Close").Return(nil)
		mockUOW.On("Users").Return(mockUsers)

		mockUsers.On("SetIsActive", ctx, "user-1", true).Return(user, nil)

		result, err := service.SetIsActive(ctx, "user-1", true)
		require.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("empty user id", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewUserService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		mockUOW.On("Close").Return(nil)
		mockUOW.On("Users").Return(mockUsers)

		result, err := service.SetIsActive(ctx, "", true)
		require.Error(t, err)
		assert.Equal(t, models.ErrEmptyUserID, err)
		assert.Equal(t, models.User{}, result)
	})
}

func TestUserService_GetPRs(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get PRs", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockPR := &mocks.MockPRRepository{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewUserService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		prs := []models.PullRequest{
			{
				ID:                "pr-1",
				Name:              "PR 1",
				AuthorID:          "author",
				Status:            models.PRStatusOpen,
				AssignedReviewers: []string{"user-1", "user-2"},
			},
		}

		mockUOW.On("Close").Return(nil)
		mockUOW.On("PR").Return(mockPR)
		mockUOW.On("Users").Return(mockUsers)

		mockPR.On("GetPRs", ctx, "user-1").Return(prs, nil)

		result, err := service.GetPRs(ctx, "user-1")
		require.NoError(t, err)
		assert.Equal(t, prs, result)
	})

	t.Run("empty user id", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}

		service := NewUserService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		mockUOW.On("Close").Return(nil)

		result, err := service.GetPRs(ctx, "")
		require.Error(t, models.ErrEmptyUserID, err)
		assert.Equal(t, []models.PullRequest{}, result)
	})
}
