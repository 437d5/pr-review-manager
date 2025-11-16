package dto

import (
	"time"

	"github.com/437d5/pr-review-manager/internal/domain/models"
)

type Team struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type TeamWithMembers struct {
	Team
	Members []User
}

func (t TeamWithMembers) ToDomain() models.Team {
	members := make([]models.User, len(t.Members))
	for i, m := range t.Members {
		members[i] = m.ToDomain()
	}

	return models.Team{
		Name:    t.Name,
		Members: members,
	}
}
