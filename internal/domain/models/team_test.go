package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeam_Validate(t *testing.T) {
	validUser := User{
		ID:       "u1",
		Username: "John Doe",
		IsActive: true,
	}

	tests := []struct {
		name     string
		team     Team
		expected error
	}{
		{
			name: "valid team",
			team: Team{
				Name:    "backend",
				Members: []User{validUser},
			},
			expected: nil,
		},
		{
			name: "empty team name",
			team: Team{
				Name:    "",
				Members: []User{validUser},
			},
			expected: ErrTeamNameEmpty,
		},
		{
			name: "empty members",
			team: Team{
				Name:    "backend",
				Members: []User{},
			},
			expected: ErrTeamMembersEmpty,
		},
		{
			name: "nil team members",
			team: Team{
				Name:    "backend",
				Members: nil,
			},
			expected: ErrTeamMembersEmpty,
		},
		{
			name: "valid team with multiple members",
			team: Team{
				Name: "backend",
				Members: []User{
					validUser,
					{
						ID:       "u2",
						Username: "Jane Doe",
						IsActive: true,
					},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.Validate()
			assert.Equal(t, tt.expected, err)
		})
	}
}
