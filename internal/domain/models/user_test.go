// internal/domain/models/user_test.go
package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Equals(t *testing.T) {
	baseUser := User{
		ID:       "user-1",
		Username: "john_doe",
		IsActive: true,
		TeamName: "backend-team",
	}

	tests := []struct {
		name     string
		user1    User
		user2    User
		expected bool
	}{
		{
			name:     "identical users",
			user1:    baseUser,
			user2:    baseUser,
			expected: true,
		},
		{
			name:  "different ID",
			user1: baseUser,
			user2: User{
				ID:       "user-2", // diff
				Username: "john_doe",
				IsActive: true,
				TeamName: "backend-team",
			},
			expected: false,
		},
		{
			name:  "different username",
			user1: baseUser,
			user2: User{
				ID:       "user-1",
				Username: "john_smith", // diff
				IsActive: true,
				TeamName: "backend-team",
			},
			expected: false,
		},
		{
			name:  "different active status",
			user1: baseUser,
			user2: User{
				ID:       "user-1",
				Username: "john_doe",
				IsActive: false, // diff
				TeamName: "backend-team",
			},
			expected: false,
		},
		{
			name:  "different team name",
			user1: baseUser,
			user2: User{
				ID:       "user-1",
				Username: "john_doe",
				IsActive: true,
				TeamName: "frontend-team", // diff
			},
			expected: false,
		},
		{
			name: "empty team vs filled",
			user1: User{
				ID:       "user-1",
				Username: "john_doe",
				IsActive: true,
				TeamName: "",
			},
			user2: User{
				ID:       "user-1",
				Username: "john_doe",
				IsActive: true,
				TeamName: "backend-team",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user1.Equals(tt.user2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
