package dto

import (
	"database/sql"
	"testing"
	"time"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestPullRequestDTO_ToDomain(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)

	tests := []struct {
		name          string
		pr            PullRequestDTO
		expected      models.PullRequest
		expectedError error
	}{
		{
			name: "valid OPEN PR",
			pr: PullRequestDTO{
				ID:        "pr-1",
				Name:      "PR 1",
				AuthorID:  "author",
				Status:    "OPEN",
				CreatedAt: now,
				MergedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			expected: models.PullRequest{
				ID:       "pr-1",
				Name:     "PR 1",
				AuthorID: "author",
				Status:   models.PRStatusOpen,
				MergedAt: &nowStr,
			},
			expectedError: nil,
		},
		{
			name: "valid MERGED PR",
			pr: PullRequestDTO{
				ID:        "pr-1",
				Name:      "PR 1",
				AuthorID:  "author",
				Status:    "MERGED",
				CreatedAt: now,
				MergedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			expected: models.PullRequest{
				ID:       "pr-1",
				Name:     "PR 1",
				AuthorID: "author",
				Status:   models.PRStatusMerged,
				MergedAt: &nowStr,
			},
			expectedError: nil,
		},
		{
			name: "valid without merged at PR",
			pr: PullRequestDTO{
				ID:        "pr-1",
				Name:      "PR 1",
				AuthorID:  "author",
				Status:    "OPEN",
				CreatedAt: now,
				MergedAt: sql.NullTime{
					Time:  now,
					Valid: false,
				},
			},
			expected: models.PullRequest{
				ID:       "pr-1",
				Name:     "PR 1",
				AuthorID: "author",
				Status:   models.PRStatusOpen,
				MergedAt: nil,
			},
			expectedError: nil,
		},
		{
			name: "invalid status PR",
			pr: PullRequestDTO{
				ID:        "pr-1",
				Name:      "PR 1",
				AuthorID:  "author",
				Status:    "INVALID",
				CreatedAt: now,
				MergedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			expected: models.PullRequest{
				ID:       "pr-1",
				Name:     "PR 1",
				AuthorID: "author",
				Status:   models.PRStatusOpen,
				MergedAt: &nowStr,
			},
			expectedError: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr, err := tt.pr.ToDomain()
			assert.Equal(t, tt.expectedError, err)
			if err == nil {
				assert.Equal(t, tt.expected, pr)
			}
		})
	}
}
