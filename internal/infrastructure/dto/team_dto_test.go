package dto

import (
	"testing"
	"time"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestTeamDTO_ToDomain(t *testing.T) {
	tests := []struct {
		name     string
		team     TeamWithMembers
		expected models.Team
	}{
		{
			name: "valid team with members",
			team: TeamWithMembers{
				Team: Team{
					ID:        1,
					Name:      "backend-team",
					CreatedAt: time.Now(),
				},
				Members: []User{
					{ID: "user-1", Username: "User 1", IsActive: true, TeamName: "backend-team"},
				},
			},
			expected: models.Team{
				Name: "backend-team",
				Members: []models.User{
					{ID: "user-1", Username: "User 1", IsActive: true},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team := tt.team.ToDomain()
			assert.Equal(t, tt.expected.Name, team.Name)
			assert.Len(t, team.Members, len(tt.expected.Members))
			for i := range tt.expected.Members {
				assert.Equal(t, tt.expected.Members[i].ID, team.Members[i].ID)
				assert.Equal(t, tt.expected.Members[i].Username, team.Members[i].Username)
				assert.Equal(t, tt.expected.Members[i].IsActive, team.Members[i].IsActive)
			}
		})
	}
}
