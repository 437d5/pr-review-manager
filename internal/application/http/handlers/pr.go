package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	httpErr "github.com/437d5/pr-review-manager/internal/application/http"
	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/services"
)

type PRHandler struct {
	prService *services.PRService
}

func NewPRHandler(prService *services.PRService) *PRHandler {
	return &PRHandler{
		prService: prService,
	}
}

type CreatePRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type PRResponse struct {
	PullRequest models.PullRequest `json:"pr"`
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req CreatePRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pr := models.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}

	createdPR, err := h.prService.CreatePR(r.Context(), pr)
	if err != nil {
		switch err {
		case models.ErrUserNotFound, models.ErrTeamNotFound:
			httpErr.WriteError(w, http.StatusNotFound, httpErr.ErrNotFound)
		case models.ErrPullRequestExists:
			httpErr.WriteError(w, http.StatusConflict, httpErr.ErrPRExists)
		default:
			httpErr.WriteInernalError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(PRResponse{PullRequest: createdPR})
	if err != nil {
		slog.Error("failed to encode response", "error", err.Error())
	}
}

type MergeRequest struct {
	ID string `json:"pull_request_id"`
}

func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err.Error())
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mergedPR, err := h.prService.Merge(r.Context(), req.ID)
	if err != nil {
		if errors.Is(err, models.ErrPullRequestNotFound) {
			httpErr.WriteError(w, http.StatusNotFound, httpErr.ErrNotFound)
		} else if !errors.Is(err, models.ErrPullRequestAlreadyMerged) {
			httpErr.WriteInernalError(w, err)
		}
		return
	}

	res := PRResponse{PullRequest: mergedPR}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("failed to encode response", "error", err.Error(), "response", res)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

type ReassignRequest struct {
	ID            string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignResponse struct {
	PRResponse
	NewReviewerID string `json:"replaced_by"`
}

func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err.Error())
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedPR, newReviewerID, err := h.prService.ReassignReviewer(r.Context(), req.ID, req.OldReviewerID)
	if err != nil {
		switch err {
		case models.ErrEmptyUserID, models.ErrPullRequestIDEmpty:
			http.Error(w, "Invalid request", http.StatusBadRequest)
		case models.ErrPullRequestNotFound, models.ErrUserNotFound:
			httpErr.WriteError(w, http.StatusNotFound, httpErr.ErrNotFound)
		case models.ErrPullRequestAlreadyMerged:
			httpErr.WriteError(w, http.StatusConflict, httpErr.ErrCannotChangeAfterMerge)
		case models.ErrUserNotReviewer:
			httpErr.WriteError(w, http.StatusConflict, httpErr.ErrUserWasNotAssigned)
		case models.ErrNoCandidateToReassign:
			httpErr.WriteError(w, http.StatusConflict, httpErr.ErrNoCandidate)
		default:
			httpErr.WriteInernalError(w, err)
		}
		return
	}

	res := ReassignResponse{
		PRResponse:    PRResponse{updatedPR},
		NewReviewerID: newReviewerID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("failed to encode response", "error", err.Error(), "response", res)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
