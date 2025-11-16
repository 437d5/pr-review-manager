package http

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

var (
	ErrTeamExists = ErrorResponse{
		Error: Error{
			Code:    "TEAM_EXISTS",
			Message: "team_name already exists",
		},
	}

	ErrNotFound = ErrorResponse{
		Error: Error{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		},
	}

	ErrPRExists = ErrorResponse{
		Error: Error{
			Code:    "PR_EXISTS",
			Message: "PR id already exists",
		},
	}

	ErrCannotChangeAfterMerge = ErrorResponse{
		Error: Error{
			Code:    "PR_MERGED",
			Message: "cannot reassign on merged PR",
		},
	}

	ErrUserWasNotAssigned = ErrorResponse{
		Error: Error{
			Code:    "NOT_ASSIGNED",
			Message: "reviewer is not assigned to this PR",
		},
	}

	ErrNoCandidate = ErrorResponse{
		Error: Error{
			Code:    "NO_CANDIDATE",
			Message: "no active replacement candidate in team",
		},
	}
)
