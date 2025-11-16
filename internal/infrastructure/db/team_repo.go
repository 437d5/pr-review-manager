package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/infrastructure/dto"
	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db sqlx.ExtContext
}

func NewTeamRepository(db sqlx.ExtContext) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(
	ctx context.Context,
	team models.Team,
) (int, error) {
	const query = `
		INSERT INTO teams (name)
		VALUES ($1)
		RETURNING id
	`

	var id int
	err := sqlx.GetContext(ctx, r.db, &id, query, team.Name)
	if err != nil {
		slog.Error("cannot insert team", "error", err.Error(), "team", team.Name)
		return 0, err
	}

	return id, nil
}

func (r *TeamRepository) Exists(
	ctx context.Context,
	name string,
) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`

	var exists bool
	err := sqlx.GetContext(ctx, r.db, &exists, query, name)
	if err != nil {
		slog.Error("cannot check if team exists", "error", err.Error(), "team", name, "query", query)
		return false, err
	}

	return exists, nil
}

func (r *TeamRepository) GetByName(
	ctx context.Context,
	name string,
) (models.Team, error) {
	const teamQuery = `
		SELECT id, name, created_at
		FROM teams
		WHERE name = $1
	`

	var teamDTO dto.Team
	err := sqlx.GetContext(ctx, r.db, &teamDTO, teamQuery, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Team{}, models.ErrTeamNotFound
		}
		slog.Error("cannot get team", "error", err.Error(), "team", name)
		return models.Team{}, err
	}

	const userQuery = `
		SELECT 
			u.id,
			u.username,
			u.is_active,
			u.team_id,
			t.name as team_name,
			u.created_at
		FROM users u
		LEFT JOIN teams t ON u.team_id = t.id
		WHERE u.team_id = $1
	`

	var usersDTO []dto.User
	err = sqlx.SelectContext(ctx, r.db, &usersDTO, userQuery, teamDTO.ID)
	if err != nil {
		slog.Error("cannot get team members", "error", err.Error(), "team", name)
		return models.Team{}, err
	}

	res := dto.TeamWithMembers{Team: teamDTO, Members: usersDTO}

	return res.ToDomain(), nil
}
