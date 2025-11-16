package models

import "errors"

var (
	ErrTeamExists       = errors.New("team_name already exists")
	ErrTeamNameEmpty    = errors.New("team_name cannot be empty")
	ErrTeamMembersEmpty = errors.New("team_members cannot be empty")
	ErrTeamNotFound     = errors.New("team not found")

	ErrUserNotFound = errors.New("user not found")
	ErrEmptyUserID  = errors.New("user id cannot be empty")

	ErrPullRequestIDEmpty       = errors.New("pull_request_id cannot be empty")
	ErrPullRequestNameEmpty     = errors.New("pull_request_name cannot be empty")
	ErrPullRequestAuthorIDEmpty = errors.New("pr author_id cannot be empty")
	ErrPullRequestNotFound      = errors.New("pr not found")
	ErrPullRequestExists        = errors.New("pr already exists")
	ErrPullRequestAlreadyMerged = errors.New("pr already merged")

	ErrNoCandidateToReassign = errors.New("no candidates to reassign pr")
	ErrUserNotReviewer       = errors.New("user is not a reviewer of pr")
)
