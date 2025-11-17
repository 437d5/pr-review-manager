package dto

import (
	"database/sql"
	"errors"
	"time"

	"github.com/437d5/pr-review-manager/internal/domain/models"
)

var (
	ErrInvalidStatus = errors.New("invalid PR status")
)

type PullRequestDTO struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	AuthorID  string       `db:"author_id"`
	Status    string       `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	MergedAt  sql.NullTime `db:"merged_at"`
}

func (pr PullRequestDTO) ToDomain() (models.PullRequest, error) {
	var status models.PRStatus
	switch pr.Status {
	case string(models.PRStatusOpen):
		status = models.PRStatusOpen
	case string(models.PRStatusMerged):
		status = models.PRStatusMerged
	default:
		return models.PullRequest{}, ErrInvalidStatus
	}

	domainPR := models.PullRequest{
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
		Status:   status,
	}

	if pr.MergedAt.Valid {
		mergedAtStr := pr.MergedAt.Time.Format(time.RFC3339)
		domainPR.MergedAt = &mergedAtStr
	} else {
		domainPR.MergedAt = nil
	}

	return domainPR, nil
}
