package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		pr       PullRequest
		expected error
	}{
		{
			name: "valid pr",
			pr: PullRequest{
				ID:       "pr-1",
				Name:     "Feature implementation",
				AuthorID: "u1",
				Status:   PRStatusOpen,
			},
			expected: nil,
		},
		{
			name: "empty ID",
			pr: PullRequest{
				ID:       "",
				Name:     "Feature implementation",
				AuthorID: "u1",
			},
			expected: ErrPullRequestIDEmpty,
		},
		{
			name: "empty name",
			pr: PullRequest{
				ID:       "pr-1",
				Name:     "",
				AuthorID: "u1",
			},
			expected: ErrPullRequestNameEmpty,
		},
		{
			name: "empty author ID",
			pr: PullRequest{
				ID:       "pr-1",
				Name:     "Feature implementation",
				AuthorID: "",
			},
			expected: ErrPullRequestAuthorIDEmpty,
		},
		{
			name: "empty",
			pr: PullRequest{
				ID:       "",
				Name:     "",
				AuthorID: "",
			},
			expected: ErrPullRequestIDEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pr.Validate()
			assert.Equal(t, tt.expected, err)
		})
	}
}
