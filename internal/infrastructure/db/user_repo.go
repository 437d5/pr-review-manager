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

type UserRepository struct {
	db sqlx.ExtContext
}

func NewUserRepository(db sqlx.ExtContext) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user models.User,
	teamID int,
) error {
	const query = `
		INSERT INTO users (id, username, is_active, team_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.IsActive, teamID)
	if err != nil {
		slog.Error("cannot create user", "error", err.Error(), "user", user.Username, "id", user.ID)
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id string,
) (models.User, error) {
	const query = `
		SELECT 
			u.id,
			u.username,
			u.team_id,
			u.created_at,
			t.name as team_name
		FROM users u
		LEFT JOIN teams t ON u.team_id = t.id
		WHERE u.id = $1
	`

	var userDTO dto.User
	if err := sqlx.GetContext(ctx, r.db, &userDTO, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		slog.Error("cannot get user", "error", err.Error(), "id", id)
		return models.User{}, err
	}

	return userDTO.ToDomain(), nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	user models.User,
	teamID int,
) (models.User, error) {
	const query = `
		UPDATE users
		SET username = $1, is_active = $2, team_id = $3
		WHERE id = $4
		RETURNING
			id,
			username,
			is_active,
			created_at,
			(SELECT name FROM teams WHERE id = users.team_id) as team_name
	`

	var userDTO dto.User
	if err := sqlx.GetContext(ctx, r.db, &userDTO, query, user.Username, user.IsActive, teamID, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		slog.Error("cannot update user", "error", err.Error(), "user", user.Username, "id", user.ID)
		return models.User{}, err
	}

	return userDTO.ToDomain(), nil
}

func (r *UserRepository) SetIsActive(
	ctx context.Context,
	id string,
	isActive bool,
) (models.User, error) {
	const query = `
		UPDATE users
		SET is_active = $1
		WHERE id = $2
		RETURNING
			id,
			username,
			is_active,
			created_at,
			(SELECT name FROM teams WHERE id = users.team_id) as team_name
	`
	var userDTO dto.User
	if err := sqlx.GetContext(ctx, r.db, &userDTO, query, isActive, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		slog.Error("cannot update user", "error", err.Error(), "id", id, "is_active", isActive)
		return models.User{}, err
	}

	return userDTO.ToDomain(), nil
}

func (r *UserRepository) GetActiveTeammatesByUserID(ctx context.Context, userID string) ([]models.User, error) {
	const query = `
		SELECT 
			u.id,
			u.username,
			u.team_id,
			u.is_active,
			u.created_at,
			t.name as team_name
		FROM users u
		INNER JOIN users cu ON cu.team_id = u.team_id
		LEFT JOIN teams t ON u.team_id = t.id
		WHERE cu.id = $1
			AND u.id != $1
			AND u.is_active = true
		ORDER BY u.username
	`

	var userDTOs []dto.User
	err := sqlx.SelectContext(ctx, r.db, &userDTOs, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.User{}, nil
		}
		slog.Error("cannot get active teammates", "error", err.Error(), "user_id", userID)
		return []models.User{}, err
	}

	users := make([]models.User, len(userDTOs))
	for i, dto := range userDTOs {
		users[i] = dto.ToDomain()
	}

	return users, nil
}
