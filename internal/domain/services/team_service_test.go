package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
	"github.com/437d5/pr-review-manager/test/mocks"
)

func TestTeamService_CreateTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("successful team creation", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockTeams := &mocks.MockTeamRepository{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewTeamService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		team := models.Team{
			Name: "backend-team",
			Members: []models.User{
				{ID: "user-1", Username: "User 1", IsActive: true, TeamName: "backend-team"},
				{ID: "user-2", Username: "User 2", IsActive: true, TeamName: "backend-team"},
			},
		}

		createdTeam := team

		mockUOW.On("Begin", ctx).Return(nil).Once()
		mockUOW.On("Commit").Return(nil).Once()
		mockUOW.On("Close").Return(nil).Once()
		mockUOW.On("Teams").Return(mockTeams)
		mockUOW.On("Users").Return(mockUsers)

		mockTeams.On("Exists", ctx, "backend-team").Return(false, nil).Once()
		mockTeams.On("Create", ctx, team).Return(1, nil).Once()
		mockUsers.On("GetByID", ctx, "user-1").Return(models.User{}, models.ErrUserNotFound).Once()
		mockUsers.On("Create", ctx, team.Members[0], 1).Return(nil).Once()
		mockUsers.On("GetByID", ctx, "user-2").Return(models.User{}, models.ErrUserNotFound).Once()
		mockUsers.On("Create", ctx, team.Members[1], 1).Return(nil).Once()
		mockTeams.On("GetByName", ctx, "backend-team").Return(createdTeam, nil).Once()

		result, err := service.CreateTeam(ctx, team)

		require.NoError(t, err)
		assert.Equal(t, createdTeam.Name, result.Name)
		assert.Len(t, result.Members, 2)
		for i := range result.Members {
			assert.Equal(t, createdTeam.Members[i].ID, result.Members[i].ID)
			assert.Equal(t, createdTeam.Members[i].Username, result.Members[i].Username)
			assert.Equal(t, createdTeam.Members[i].IsActive, result.Members[i].IsActive)
		}
		mockUOW.AssertExpectations(t)
		mockTeams.AssertExpectations(t)
		mockUsers.AssertExpectations(t)
	})

	t.Run("team already exists", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockTeams := &mocks.MockTeamRepository{}

		service := NewTeamService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		team := models.Team{
			Name: "backend-team",
			Members: []models.User{
				{ID: "user-1", Username: "User 1", IsActive: true},
			},
		}

		mockUOW.On("Begin", ctx).Return(nil).Once()
		mockUOW.On("Rollback").Return(nil).Once()
		mockUOW.On("Close").Return(nil).Once()
		mockUOW.On("Teams").Return(mockTeams)

		mockTeams.On("Exists", ctx, "backend-team").Return(true, nil).Once()

		result, err := service.CreateTeam(ctx, team)

		assert.Error(t, err)
		assert.Equal(t, models.ErrTeamExists, err)
		assert.Equal(t, models.Team{}, result)

		mockUOW.AssertExpectations(t)
		mockTeams.AssertExpectations(t)
	})
}

func TestTeamService_GetTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get team", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockTeams := &mocks.MockTeamRepository{}

		service := NewTeamService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		expectedTeam := models.Team{
			Name: "backend-team",
			Members: []models.User{
				{ID: "user-1", Username: "User 1", IsActive: true},
			},
		}

		mockUOW.On("Close").Return(nil)
		mockUOW.On("Teams").Return(mockTeams)

		mockTeams.On("GetByName", ctx, "backend-team").Return(expectedTeam, nil)

		result, err := service.GetTeam(ctx, "backend-team")

		require.NoError(t, err)
		assert.Equal(t, expectedTeam, result)
	})

	t.Run("team not found", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockTeams := &mocks.MockTeamRepository{}

		service := NewTeamService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		mockUOW.On("Close").Return(nil)
		mockUOW.On("Teams").Return(mockTeams)

		mockTeams.On("GetByName", ctx, "non-existent-team").Return(models.Team{}, models.ErrTeamNotFound)

		result, err := service.GetTeam(ctx, "non-existent-team")

		assert.Error(t, err)
		assert.Equal(t, models.ErrTeamNotFound, err)
		assert.Equal(t, models.Team{}, result)
	})
}
