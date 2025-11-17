package services

import (
	"context"
	"testing"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
	"github.com/437d5/pr-review-manager/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPRService_CreatePR(t *testing.T) {
	ctx := context.Background()

	t.Run("successful PR creation", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockUsers := &mocks.MockUserRepository{}
		mockPR := &mocks.MockPRRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		author := models.User{
			ID:       "user-1",
			Username: "author",
			IsActive: true,
			TeamName: "backend",
		}

		teammates := []models.User{
			{ID: "user-2", Username: "reviewer1", IsActive: true},
			{ID: "user-3", Username: "reviewer2", IsActive: true},
			{ID: "user-4", Username: "reviewer3", IsActive: true},
		}

		pr := models.PullRequest{
			ID:       "pr-1",
			Name:     "Test PR",
			AuthorID: "user-1",
		}

		expectedPR := pr
		expectedPR.Status = models.PRStatusOpen
		expectedPR.AssignedReviewers = []string{"user-2", "user-3"}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Commit").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("Users").Return(mockUsers)
		mockUOW.On("PR").Return(mockPR)

		mockUsers.On("GetByID", ctx, "user-1").Return(author, nil)
		mockPR.On("GetByID", ctx, "pr-1").Return(models.PullRequest{}, models.ErrPullRequestNotFound)
		mockUsers.On("GetActiveTeammatesByUserID", ctx, "user-1").Return(teammates, nil)
		mockPR.On("Create", ctx, mock.AnythingOfType("models.PullRequest")).Return(expectedPR, nil)

		result, err := service.CreatePR(ctx, pr)

		require.NoError(t, err)
		assert.Equal(t, expectedPR.ID, result.ID)
		assert.Equal(t, models.PRStatusOpen, result.Status)
		assert.Len(t, result.AssignedReviewers, 2)
		mockUOW.AssertExpectations(t)
	})

	t.Run("author not found", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		pr := models.PullRequest{
			ID:       "pr-1",
			Name:     "Test PR",
			AuthorID: "user-1",
		}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Rollback").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("Users").Return(mockUsers)

		mockUsers.On("GetByID", ctx, "user-1").Return(models.User{}, models.ErrUserNotFound)

		result, err := service.CreatePR(ctx, pr)

		assert.Error(t, err)
		assert.Equal(t, models.ErrUserNotFound, err)
		assert.Equal(t, models.PullRequest{}, result)
		mockUOW.AssertExpectations(t)
	})

	t.Run("author has no team", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		author := models.User{
			ID:       "user-1",
			Username: "author",
			IsActive: true,
			TeamName: "",
		}

		pr := models.PullRequest{
			ID:       "pr-1",
			Name:     "Test PR",
			AuthorID: "user-1",
		}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Rollback").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("Users").Return(mockUsers)

		mockUsers.On("GetByID", ctx, "user-1").Return(author, nil)

		result, err := service.CreatePR(ctx, pr)

		assert.Error(t, err)
		assert.Equal(t, models.ErrTeamNotFound, err)
		assert.Equal(t, models.PullRequest{}, result)
		mockUOW.AssertExpectations(t)
	})

	t.Run("pr already exists", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockUsers := &mocks.MockUserRepository{}
		mockPR := &mocks.MockPRRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		author := models.User{
			ID:       "user-1",
			Username: "author",
			IsActive: true,
			TeamName: "backend-team",
		}

		existingPR := models.PullRequest{
			ID:   "pr-1",
			Name: "Existing PR",
		}

		pr := models.PullRequest{
			ID:       "pr-1",
			Name:     "Test PR",
			AuthorID: "user-1",
		}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Rollback").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("Users").Return(mockUsers)
		mockUOW.On("PR").Return(mockPR)

		mockUsers.On("GetByID", ctx, "user-1").Return(author, nil)
		mockPR.On("GetByID", ctx, "pr-1").Return(existingPR, nil)

		result, err := service.CreatePR(ctx, pr)

		assert.Error(t, err)
		assert.Equal(t, models.ErrPullRequestExists, err)
		assert.Equal(t, models.PullRequest{}, result)
		mockUOW.AssertExpectations(t)
	})
}

func TestPRService_Merge(t *testing.T) {
	ctx := context.Background()

	t.Run("successful merge", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockPR := &mocks.MockPRRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		existingPR := models.PullRequest{
			ID:     "pr-1",
			Name:   "Test PR",
			Status: models.PRStatusOpen,
		}

		mergedPR := existingPR
		mergedPR.Status = models.PRStatusMerged

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Commit").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("PR").Return(mockPR)

		mockPR.On("GetByID", ctx, "pr-1").Return(existingPR, nil)
		mockPR.On("Merge", ctx, "pr-1").Return(mergedPR, nil)

		result, err := service.Merge(ctx, "pr-1")

		require.NoError(t, err)
		assert.Equal(t, models.PRStatusMerged, result.Status)
		mockUOW.AssertExpectations(t)
	})

	t.Run("pr not found", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockPR := &mocks.MockPRRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Rollback").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("PR").Return(mockPR)

		mockPR.On("GetByID", ctx, "pr-1").Return(models.PullRequest{}, models.ErrPullRequestNotFound)

		result, err := service.Merge(ctx, "pr-1")

		assert.Error(t, err)
		assert.Equal(t, models.ErrPullRequestNotFound, err)
		assert.Equal(t, models.PullRequest{}, result)
		mockUOW.AssertExpectations(t)
	})

	t.Run("PR already merged", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockPR := &mocks.MockPRRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		mergedPR := models.PullRequest{
			ID:     "pr-1",
			Name:   "Test PR",
			Status: models.PRStatusMerged,
		}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Rollback").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("PR").Return(mockPR)

		mockPR.On("GetByID", ctx, "pr-1").Return(mergedPR, nil)

		result, err := service.Merge(ctx, "pr-1")

		assert.Error(t, err)
		assert.Equal(t, models.ErrPullRequestAlreadyMerged, err)
		assert.Equal(t, mergedPR, result)
		mockUOW.AssertExpectations(t)
	})
}

func TestPRService_ReassignReviewer(t *testing.T) {
	ctx := context.Background()

	t.Run("successful reassignment", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockPR := &mocks.MockPRRepository{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		pr := models.PullRequest{
			ID:                "pr-1",
			Name:              "Test PR",
			AuthorID:          "author-1",
			Status:            models.PRStatusOpen,
			AssignedReviewers: []string{"old-reviewer-1", "reviewer-2"},
		}

		reviewers := []models.User{
			{ID: "old-reviewer-1", Username: "Old Reviewer", IsActive: true},
			{ID: "reviewer-2", Username: "Reviewer 2", IsActive: true},
		}

		teammates := []models.User{
			{ID: "new-reviewer-1", Username: "New Reviewer 1", IsActive: true},
			{ID: "new-reviewer-2", Username: "New Reviewer 2", IsActive: true},
		}

		updatedPR := pr
		updatedPR.AssignedReviewers = []string{"new-reviewer-1", "reviewer-2"}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Commit").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("PR").Return(mockPR)
		mockUOW.On("Users").Return(mockUsers)

		mockPR.On("GetByID", ctx, "pr-1").Return(pr, nil)
		mockUsers.On("GetByID", ctx, "old-reviewer-1").Return(models.User{ID: "old-reviewer-1"}, nil)
		mockPR.On("GetReviewers", ctx, "pr-1").Return(reviewers, nil)
		mockUsers.On("GetActiveTeammatesByUserID", ctx, "old-reviewer-1").Return(teammates, nil)
		mockPR.On("Reassign", ctx, "pr-1", "old-reviewer-1", mock.AnythingOfType("string")).Return(updatedPR, nil)

		result, newReviewerID, err := service.ReassignReviewer(ctx, "pr-1", "old-reviewer-1")

		require.NoError(t, err)
		assert.Equal(t, updatedPR, result)
		assert.Contains(t, []string{"new-reviewer-1", "new-reviewer-2"}, newReviewerID)
		mockUOW.AssertExpectations(t)
	})

	t.Run("user is not a reviewer", func(t *testing.T) {
		mockUOW := &mocks.MockUnitOfWork{}
		mockPR := &mocks.MockPRRepository{}
		mockUsers := &mocks.MockUserRepository{}

		service := NewPRService(func(ctx context.Context) (repositories.UnitOfWork, error) {
			return mockUOW, nil
		})

		pr := models.PullRequest{
			ID:                "pr-1",
			Name:              "Test PR",
			AuthorID:          "author-1",
			Status:            models.PRStatusOpen,
			AssignedReviewers: []string{"reviewer-1", "reviewer-2"},
		}

		reviewers := []models.User{
			{ID: "reviewer-1", Username: "Reviewer 1", IsActive: true},
			{ID: "reviewer-2", Username: "Reviewer 2", IsActive: true},
		}

		mockUOW.On("Begin", ctx).Return(nil)
		mockUOW.On("Rollback").Return(nil)
		mockUOW.On("Close").Return(nil)
		mockUOW.On("PR").Return(mockPR)
		mockUOW.On("Users").Return(mockUsers)

		mockPR.On("GetByID", ctx, "pr-1").Return(pr, nil)
		mockUsers.On("GetByID", ctx, "not-a-reviewer").Return(models.User{ID: "not-a-reviewer"}, nil)
		mockPR.On("GetReviewers", ctx, "pr-1").Return(reviewers, nil)

		result, newReviewerID, err := service.ReassignReviewer(ctx, "pr-1", "not-a-reviewer")

		assert.Error(t, err)
		assert.Equal(t, models.ErrUserNotReviewer, err)
		assert.Equal(t, models.PullRequest{}, result)
		assert.Equal(t, "", newReviewerID)
	})
}

func TestPRService_selectRandomReviewers(t *testing.T) {
	service := &PRService{}

	t.Run("empty users", func(t *testing.T) {
		users := []models.User{}
		result := service.selectRandomReviewers(users, 2)
		assert.Empty(t, result)
	})

	t.Run("users less than max", func(t *testing.T) {
		users := []models.User{
			{ID: "user-1", Username: "User 1"},
			{ID: "user-2", Username: "User 2"},
		}
		result := service.selectRandomReviewers(users, 3)
		assert.Len(t, result, 2)
	})

	t.Run("users more than max", func(t *testing.T) {
		users := []models.User{
			{ID: "user-1", Username: "User 1"},
			{ID: "user-2", Username: "User 2"},
			{ID: "user-3", Username: "User 3"},
			{ID: "user-4", Username: "User 4"},
		}
		result := service.selectRandomReviewers(users, 2)
		assert.Len(t, result, 2)
	})
}

func TestPRService_filterCandidates(t *testing.T) {
	service := &PRService{}

	pr := models.PullRequest{
		ID:                "pr-1",
		AuthorID:          "author-1",
		AssignedReviewers: []string{"reviewer-1", "reviewer-2"},
	}

	candidates := []models.User{
		{ID: "author-1"},
		{ID: "old-reviewer-1"},
		{ID: "reviewer-1"},
		{ID: "reviewer-2"},
		{ID: "candidate-1"},
		{ID: "candidate-2"},
	}

	result := service.filterCandidates(candidates, pr, "old-reviewer-1")

	assert.Len(t, result, 2)
	assert.Equal(t, "candidate-1", result[0].ID)
	assert.Equal(t, "candidate-2", result[1].ID)
}
