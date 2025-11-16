package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

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

func WriteError(w http.ResponseWriter, status int, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		slog.Error("failed to write error response", "error", err.Error())
		return
	}
}

func WriteInernalError(w http.ResponseWriter, err error) {
	slog.Error("internal server error", "error", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
