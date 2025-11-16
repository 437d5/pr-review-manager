package services

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
)

const maxReviewers = 2

type PRService struct {
	uowFactory func(context.Context) (repositories.UnitOfWork, error)
}

func NewPRService(uowFactory func(ctx context.Context) (repositories.UnitOfWork, error)) *PRService {
	return &PRService{
		uowFactory: uowFactory,
	}
}

func (s *PRService) CreatePR(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
	if err := pr.Validate(); err != nil {
		return models.PullRequest{}, err
	}
	pr.Status = models.PRStatusOpen

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return models.PullRequest{}, err
	}
	defer uow.Close()

	if err := uow.Begin(ctx); err != nil {
		slog.Error("cannot begin transaction", "error", err.Error())
		return models.PullRequest{}, err
	}

	var createdPR models.PullRequest
	err = func() error {
		author, err := uow.Users().GetByID(ctx, pr.AuthorID)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				slog.Error("author not found", "author_id", pr.AuthorID, "error", err.Error(), "pr", pr.Name)
				return models.ErrUserNotFound
			}
			slog.Error("cannot get author", "error", err.Error(), "author_id", pr.AuthorID)
			return err
		}

		if author.TeamName == "" {
			slog.Error("author has no team", "author_id", pr.AuthorID)
			return models.ErrTeamNotFound
		}

		existingPR, err := uow.PR().GetByID(ctx, pr.ID)
		if err != nil && !errors.Is(err, models.ErrPullRequestNotFound) {
			slog.Error("cannot get PR", "error", err.Error(), "pr_id", pr.ID)
			return err
		}
		if existingPR.ID != "" {
			slog.Warn("PR already exists", "pr_id", pr.ID)
			return models.ErrPullRequestExists
		}

		teammates, err := uow.Users().GetActiveTeammatesByUserID(ctx, pr.AuthorID)
		if err != nil {
			slog.Error("cannot get teammates", "error", err.Error())
			return err
		}

		reviewers := s.selectRandomReviewers(teammates, maxReviewers)
		for _, reviewer := range reviewers {
			pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer.ID)
		}

		createdPR, err = uow.PR().Create(ctx, pr)
		if err != nil {
			slog.Error("cannot create PR", "error", err.Error(), "pr", pr.Name)
			return err
		}

		slog.Info("PR created succefully",
			"pr_id", createdPR.ID,
			"author", pr.AuthorID,
			"reviewers_count", len(reviewers),
		)

		return nil
	}()

	if err != nil {
		if err := uow.Rollback(); err != nil {
			slog.Error("cannot rollback transaction", "error", err.Error())
			return models.PullRequest{}, err
		}
		return models.PullRequest{}, err
	}

	if err := uow.Commit(); err != nil {
		slog.Error("cannot commit transaction", "error", err.Error())
		return models.PullRequest{}, err
	}

	return createdPR, nil
}

func (s *PRService) selectRandomReviewers(users []models.User, max int) []models.User {
	if len(users) == 0 {
		return []models.User{}
	}

	if len(users) <= max {
		return users
	}

	shuffled := make([]models.User, len(users))
	copy(shuffled, users)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:max]
}

func (s *PRService) Merge(ctx context.Context, ID string) (models.PullRequest, error) {
	if ID == "" {
		return models.PullRequest{}, models.ErrPullRequestIDEmpty
	}

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return models.PullRequest{}, err
	}
	defer uow.Close()

	if err := uow.Begin(ctx); err != nil {
		slog.Error("cannot begin transaction", "error", err.Error())
		return models.PullRequest{}, err
	}

	var mergedPR models.PullRequest
	err = func() error {
		existingPR, err := uow.PR().GetByID(ctx, ID)
		if err != nil {
			if errors.Is(err, models.ErrPullRequestNotFound) {
				return models.ErrPullRequestNotFound
			}
			slog.Error("cannot get PR", "error", err.Error(), "pr_id", ID)
			return err
		}

		if existingPR.Status == models.PRStatusMerged {
			mergedPR = existingPR
			return models.ErrPullRequestAlreadyMerged
		}

		mergedPR, err = uow.PR().Merge(ctx, ID)
		if err != nil {
			slog.Error("cannot merge PR", "error", err.Error(), "pr_id", ID)
			return err
		}

		slog.Info("PR merged successfully", "pr_id", ID)
		return nil
	}()

	if err != nil {
		if err := uow.Rollback(); err != nil {
			slog.Error("cannot rollback transaction", "error", err.Error())
			return models.PullRequest{}, err
		}
		return mergedPR, err
	}

	if err := uow.Commit(); err != nil {
		slog.Error("cannot commit transaction", "error", err.Error())
		return models.PullRequest{}, err
	}

	return mergedPR, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (models.PullRequest, string, error) {
	if prID == "" {
		return models.PullRequest{}, "", models.ErrPullRequestIDEmpty
	}
	if oldReviewerID == "" {
		return models.PullRequest{}, "", models.ErrEmptyUserID
	}

	uow, err := s.uowFactory(ctx)
	if err != nil {
		return models.PullRequest{}, "", err
	}
	defer uow.Close()

	if err := uow.Begin(ctx); err != nil {
		slog.Error("cannot begin transaction", "error", err.Error())
		return models.PullRequest{}, "", err
	}

	var updatedPR models.PullRequest
	var newReviewerID string

	err = func() error {
		// check if pr exists
		pr, err := uow.PR().GetByID(ctx, prID)
		if err != nil {
			if errors.Is(err, models.ErrPullRequestNotFound) {
				return err
			}
			slog.Error("cannot get PR", "error", err.Error(), "pr_id", prID)
			return err
		}

		_, err = uow.Users().GetByID(ctx, oldReviewerID)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				return err
			}
			slog.Error("cannot get user", "error", err.Error(), "user_id", oldReviewerID)
			return err
		}

		// check if pr is not merged
		if pr.Status == models.PRStatusMerged {
			return models.ErrPullRequestAlreadyMerged
		}

		// check if old reviewer in reviewers
		reviewers, err := uow.PR().GetReviewers(ctx, prID)
		if err != nil {
			slog.Error("cannot get reviewers", "error", err.Error(), "pr_id", prID)
			return err
		}

		var isReviewer bool
		for _, reviewer := range reviewers {
			if reviewer.ID == oldReviewerID {
				isReviewer = true
				break
			}
		}

		if !isReviewer {
			return models.ErrUserNotReviewer
		}

		// get teammates except old reviewer
		candidates, err := uow.Users().GetActiveTeammatesByUserID(ctx, oldReviewerID)
		if err != nil {
			slog.Error("cannot get teammates", "error", err.Error(), "user_id", oldReviewerID)
			return err
		}

		// exlude old reviewer and author from candidates
		candidates = s.filterCandidates(candidates, pr, oldReviewerID)
		// select new random reviewer
		newReviewer := s.selectRandomReviewers(candidates, 1)

		if len(newReviewer) == 0 {
			return models.ErrNoCandidateToReassign
		}

		updatedPR, err = uow.PR().Reassign(ctx, prID, oldReviewerID, newReviewer[0].ID)
		if err != nil {
			slog.Error("cannot reassign reviewer", "error", err.Error(), "pr_id", prID, "old_reviewer_id",
				oldReviewerID, "new_reviewer_id", newReviewer[0].ID)
			return err
		}

		newReviewerID = newReviewer[0].ID

		slog.Info("reviewer reassigned successfully",
			"pr_id", prID,
			"old_reviewer", oldReviewerID,
			"new_reviewer", newReviewerID,
		)
		return nil
	}()

	if err != nil {
		if err := uow.Rollback(); err != nil {
			slog.Error("cannot rollback transaction", "error", err.Error())
			return models.PullRequest{}, "", err
		}
		return models.PullRequest{}, "", err
	}

	if err := uow.Commit(); err != nil {
		slog.Error("cannot commit transaction", "error", err.Error())
		return models.PullRequest{}, "", err
	}

	return updatedPR, newReviewerID, nil
}

// exclude old reviewer and author from candidates
func (s *PRService) filterCandidates(candidates []models.User, pr models.PullRequest, oldReviewerID string) []models.User {
	badSet := make(map[string]struct{}, len(pr.AssignedReviewers)+1) // 1 is author id
	badSet[pr.AuthorID] = struct{}{}
	badSet[oldReviewerID] = struct{}{}
	for _, reviewer := range pr.AssignedReviewers {
		badSet[reviewer] = struct{}{}
	}

	res := make([]models.User, 0, len(candidates))
	for _, candidate := range candidates {
		if _, ok := badSet[candidate.ID]; ok {
			continue
		} else {
			res = append(res, candidate)
		}
	}

	return res
}
