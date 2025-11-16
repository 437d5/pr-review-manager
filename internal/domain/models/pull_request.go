package models

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            PRStatus `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

func (p PullRequest) Validate() error {
	switch {
	case p.ID == "":
		return ErrPullRequestIDEmpty
	case p.Name == "":
		return ErrPullRequestNameEmpty
	case p.AuthorID == "":
		return ErrPullRequestAuthorIDEmpty
	}
	return nil
}
