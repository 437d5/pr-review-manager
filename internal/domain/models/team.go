package models

type Team struct {
	Name    string `json:"team_name"`
	Members []User `json:"members"`
}

func (t Team) Validate() error {
	if t.Name == "" {
		return ErrTeamNameEmpty
	}

	if len(t.Members) == 0 {
		return ErrTeamMembersEmpty
	}

	return nil
}
