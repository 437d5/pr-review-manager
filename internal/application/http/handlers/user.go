package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httpErr "github.com/437d5/pr-review-manager/internal/application/http"
	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User models.User `json:"user"`
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userService.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		switch err {
		case models.ErrEmptyUserID:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case models.ErrUserNotFound:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(httpErr.ErrNotFound)
		default:
			slog.Error("failed to set user active status", "error", err.Error(), "user_id", req.UserID, "is_active", req.IsActive)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{User: updatedUser})
}

type GetPRsResponse struct {
	UserID string               `json:"user_id"`
	PRs    []models.PullRequest `json:"pull_requests"`
}

func (h *UserHandler) GetPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	prs, err := h.userService.GetPRs(r.Context(), userID)
	if err != nil {
		slog.Error("failed to get PRs", "error", err.Error(), "user_id", userID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	prs = hideMergeAt(prs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetPRsResponse{
		UserID: userID,
		PRs:    prs,
	})
}

func hideMergeAt(prs []models.PullRequest) []models.PullRequest {
	for i := range prs {
		prs[i].MergedAt = nil
	}
	return prs
}
