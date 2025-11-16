package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/infrastructure/dto"
	"github.com/jmoiron/sqlx"
)

type PullRequestRepository struct {
	db sqlx.ExtContext
}

func NewPullRequestRepository(db sqlx.ExtContext) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) Create(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
	const query = `
		INSERT INTO pull_requests (id, name, author_id, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, pr.ID, pr.Name, pr.AuthorID, pr.Status)
	if err != nil {
		slog.Error("cannot create pull request", "error", err, "pr_name", pr.Name, "pr_id", pr.ID)
		return models.PullRequest{}, err
	}

	const reviewerQuery = `
		INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	for _, reviewerID := range pr.AssignedReviewers {
		if _, err := r.db.ExecContext(ctx, reviewerQuery, pr.ID, reviewerID); err != nil {
			slog.Error("cannot assign reviewer to pull request", "error", err, "pr_id", pr.ID, "reviewer_id", reviewerID)
			return models.PullRequest{}, err
		}
	}

	return pr, nil
}

func (r *PullRequestRepository) Merge(ctx context.Context, ID string) (models.PullRequest, error) {
	const query = `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = $1
		WHERE id = $2
		RETURNING id, name, author_id, status, created_at, merged_at
	`

	now := time.Now()

	var prDTO dto.PullRequestDTO
	err := sqlx.GetContext(ctx, r.db, &prDTO, query, now, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequest{}, models.ErrPullRequestNotFound
		}
		slog.Error("cannot merge pull request", "error", err.Error(), "id", ID)
		return models.PullRequest{}, err
	}

	// need to refactor
	reviewers, err := r.GetReviewers(ctx, ID)
	if err != nil {
		slog.Error("cannot get reviewers", "error", err.Error(), "id", ID)
	}

	pr := prDTO.ToDomain()

	pr.AssignedReviewers = make([]string, len(reviewers))
	for i, reviewer := range reviewers {
		pr.AssignedReviewers[i] = reviewer.ID
	}

	slog.Info("PR merged successfully", "pr_id", ID, "merged_at", now)
	return pr, nil
}

func (r *PullRequestRepository) GetReviewers(ctx context.Context, ID string) ([]models.User, error) {
	const query = `
		SELECT 
			u.id, 
			u.username,
			u.is_active,
			u.team_id,
			t.name as team_name,
			u.created_at
		FROM users u
		INNER JOIN pull_requests_reviewers prr ON u.id = prr.reviewer_id
		LEFT JOIN teams t ON u.team_id = t.id
		WHERE prr.pull_request_id = $1
	`

	var userDTOs []dto.User
	err := sqlx.SelectContext(ctx, r.db, &userDTOs, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.User{}, nil
		} else {
			slog.Error("cannot get PR reviewers", "error", err.Error(), "id", ID)
			return []models.User{}, err
		}
	}

	users := make([]models.User, len(userDTOs))
	for i, dto := range userDTOs {
		users[i] = dto.ToDomain()
	}

	return users, nil
}

func (r *PullRequestRepository) Reassign(ctx context.Context, prID, oldReviewerID, newReviewerID string) (models.PullRequest, error) {
	const deleteQuery = `
		DELETE FROM pull_requests_reviewers
		WHERE pull_request_id = $1 AND reviewer_id = $2
	`

	_, err := r.db.ExecContext(ctx, deleteQuery, prID, oldReviewerID)
	if err != nil {
		slog.Error("cannot remove old reviewer", "error", err.Error(), "pr_id", prID, "reviewer_id", oldReviewerID)
		return models.PullRequest{}, err
	}

	const insertQuery = `
		INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	_, err = r.db.ExecContext(ctx, insertQuery, prID, newReviewerID)
	if err != nil {
		slog.Error("cannot assign new reviewer", "error", err.Error(), "pr_id", prID,
			"reviewer_id", newReviewerID, "new_reviewer_id", newReviewerID, "old_reviewer_id", oldReviewerID)
		return models.PullRequest{}, err
	}

	return r.GetByID(ctx, prID)
}

func (r *PullRequestRepository) GetPRs(ctx context.Context, userID string) ([]models.PullRequest, error) {
	const query = `
		SELECT
			pr.id,
			pr.name,
			pr.author_id,
			pr.status,
			pr.created_at,
			pr.merged_at
		FROM pull_requests pr
		INNER JOIN pull_requests_reviewers prr ON pr.id = prr.pull_request_id
		WHERE prr.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	var prDTOs []dto.PullRequestDTO
	err := sqlx.SelectContext(ctx, r.db, &prDTOs, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PullRequest{}, nil
		}
		slog.Error("cannot get pull requests for user", "error", err.Error(), "user_id", userID)
		return []models.PullRequest{}, err
	}

	prs := make([]models.PullRequest, len(prDTOs))
	for i, dto := range prDTOs {
		prs[i] = dto.ToDomain()
	}

	return prs, nil
}

func (r *PullRequestRepository) GetByID(ctx context.Context, ID string) (models.PullRequest, error) {
	const query = `
		SELECT
			pr.id,
			pr.name,
			pr.author_id,
			pr.status,
			pr.created_at,
			pr.merged_at
		FROM pull_requests pr
		WHERE pr.id = $1
	`

	var prDTO dto.PullRequestDTO
	err := sqlx.GetContext(ctx, r.db, &prDTO, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequest{}, models.ErrPullRequestNotFound
		}
		slog.Error("cannot get pull request", "error", err.Error(), "pr_id", ID)
		return models.PullRequest{}, err
	}

	reviewers, err := r.GetReviewers(ctx, ID)
	if err != nil {
		slog.Error("cannot get reviewers", "error", err.Error(), "id", ID)
	}

	pr := prDTO.ToDomain()

	pr.AssignedReviewers = make([]string, len(reviewers))
	for i, reviewer := range reviewers {
		pr.AssignedReviewers[i] = reviewer.ID
	}

	return pr, nil
}
