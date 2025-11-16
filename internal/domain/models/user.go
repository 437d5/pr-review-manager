package models

type User struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name,omitempty"`
}

func (u User) Equals(other User) bool {
	return u.ID == other.ID &&
		u.Username == other.Username &&
		u.IsActive == other.IsActive &&
		u.TeamName == other.TeamName
}
