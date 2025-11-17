package dto

import (
	"testing"

	"github.com/437d5/pr-review-manager/internal/domain/models"
)

func TestUserDTO_ToDomain(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected models.User
	}{
		{
			name: "valid user",
			user: User{
				ID:       "user-1",
				Username: "User 1",
				IsActive: true,
				TeamName: "backend-team",
			},
			expected: models.User{
				ID:       "user-1",
				Username: "User 1",
				IsActive: true,
				TeamName: "backend-team"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.ToDomain()
			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
