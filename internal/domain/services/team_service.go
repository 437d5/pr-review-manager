package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
)

type TeamService struct {
	uowFactory func(context.Context) (repositories.UnitOfWork, error)
}

func NewTeamService(uowFactory func(ctx context.Context) (repositories.UnitOfWork, error)) *TeamService {
	return &TeamService{
		uowFactory: uowFactory,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team models.Team) (models.Team, error) {
	if err := team.Validate(); err != nil {
		slog.Error("invalid team", "error", err.Error(), "team", team.Name)
		return models.Team{}, err
	}

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return models.Team{}, err
	}
	defer uow.Close()

	if err := uow.Begin(ctx); err != nil {
		slog.Error("cannot begin transaction", "error", err.Error())
		return models.Team{}, err
	}

	var resTeam models.Team
	err = func() error {
		exists, err := uow.Teams().Exists(ctx, team.Name)
		if err != nil {
			slog.Error("cannot check team existence", "error", err.Error(), "team", team.Name)
			return err
		}
		if exists {
			return models.ErrTeamExists
		}

		createdTeamID, err := uow.Teams().Create(ctx, team)
		if err != nil {
			slog.Error("cannot create team", "error", err.Error(), "team", team.Name)
			return err
		}

		for _, member := range team.Members {
			member.TeamName = team.Name
			existingUser, err := uow.Users().GetByID(ctx, member.ID)
			if err != nil {
				// if user not found, create it
				if errors.Is(err, models.ErrUserNotFound) {
					slog.Info("creating new user", "user_id", member.ID, "name", member.Username)

					if err := uow.Users().Create(ctx, member, createdTeamID); err != nil {
						slog.Error("cannot create user", "error", err.Error(), "user_id", member.ID)
						return fmt.Errorf("failed to create user %s: %w", member.Username, err)
					}
					// create user and go to next member
					continue
				}
				slog.Error("cannot get user", "error", err.Error(), "user_id", member.ID)
				return err
			}

			// check if we need to update existing user
			if !member.Equals(existingUser) {
				slog.Info("updating user data", "user_id", member.ID, "name", member.Username)

				if _, err := uow.Users().Update(ctx, member, createdTeamID); err != nil {
					slog.Error("cannot update user", "error", err.Error(), "user_id", member.ID)
					return err
				}
			}
		}

		resTeam, err = uow.Teams().GetByName(ctx, team.Name)
		if err != nil {
			slog.Error("cannot get team", "error", err.Error(), "team", team.Name)
			return err
		}

		return nil
	}()

	// rollback if any error occurred
	if err != nil {
		if err := uow.Rollback(); err != nil {
			slog.Error("cannot rollback transaction", "error", err.Error())
			return models.Team{}, fmt.Errorf("rollback failed: %w", err)
		}
		return models.Team{}, err
	}

	// commit if no errors
	if err := uow.Commit(); err != nil {
		slog.Error("cannot commit transaction", "error", err.Error())
		return models.Team{}, fmt.Errorf("failed to commit team creation: %w", err)
	}

	slog.Info("team created successfully",
		"name", resTeam.Name,
		"member_count", len(resTeam.Members),
	)

	return resTeam, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (models.Team, error) {
	if name == "" {
		return models.Team{}, models.ErrTeamNameEmpty
	}

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return models.Team{}, err
	}
	defer uow.Close()

	team, err := uow.Teams().GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, models.ErrTeamNotFound) {
			return models.Team{}, models.ErrTeamNotFound
		}
		slog.Error("cannot get team", "error", err.Error(), "team", team.Name)
		return models.Team{}, err
	}

	return team, nil
}
