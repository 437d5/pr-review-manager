package dto

import (
	"time"

	"github.com/437d5/pr-review-manager/internal/domain/models"
)

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	IsActive  bool      `db:"is_active"`
	TeamID    int       `db:"team_id"`
	TeamName  string    `db:"team_name"`
	CreatedAt time.Time `db:"created_at"`
}

func (u User) ToDomain() models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Username,
		IsActive: u.IsActive,
		TeamName: u.TeamName,
	}
}
