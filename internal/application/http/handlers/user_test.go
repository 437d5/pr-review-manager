package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpErr "github.com/437d5/pr-review-manager/internal/application/http"
	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	setActiveFn func(ctx context.Context, userID string, isActive bool) (models.User, error)
	getPRsFn    func(ctx context.Context, userID string) ([]models.PullRequest, error)
}

func (m *mockUserService) SetIsActive(ctx context.Context, userID string, isActive bool) (models.User, error) {
	return m.setActiveFn(ctx, userID, isActive)
}

func (m *mockUserService) GetPRs(ctx context.Context, userID string) ([]models.PullRequest, error) {
	return m.getPRsFn(ctx, userID)
}

func TestUserHandler_SetIsActive_Success(t *testing.T) {
	service := &mockUserService{
		setActiveFn: func(ctx context.Context, userID string, isActive bool) (models.User, error) {
			return models.User{ID: userID, Username: "Bob", TeamName: "backend", IsActive: isActive}, nil
		},
	}

	handler := NewUserHandler(service)

	reqBody := SetIsActiveRequest{
		UserID:   "u1",
		IsActive: false,
	}
	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.SetIsActive(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp UserResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, reqBody.UserID, resp.User.ID)
	assert.False(t, resp.User.IsActive)
}

func TestUserHandler_SetIsActive_InvalidJSON(t *testing.T) {
	service := &mockUserService{
		setActiveFn: func(ctx context.Context, userID string, isActive bool) (models.User, error) {
			t.Fatal("SetIsActive should not run for invalid JSON")
			return models.User{}, nil
		},
	}

	handler := NewUserHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBufferString("{invalid"))
	rec := httptest.NewRecorder()

	handler.SetIsActive(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserHandler_SetIsActive_NotFound(t *testing.T) {
	service := &mockUserService{
		setActiveFn: func(ctx context.Context, userID string, isActive bool) (models.User, error) {
			return models.User{}, models.ErrUserNotFound
		},
	}

	handler := NewUserHandler(service)

	reqBody := SetIsActiveRequest{
		UserID:   "missing",
		IsActive: true,
	}
	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.SetIsActive(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var errResp httpErr.ErrorResponse
	err = json.NewDecoder(res.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, httpErr.ErrNotFound.Error.Code, errResp.Error.Code)
}

func TestUserHandler_GetPRs_Success(t *testing.T) {
	service := &mockUserService{
		getPRsFn: func(ctx context.Context, userID string) ([]models.PullRequest, error) {
			return []models.PullRequest{
				{ID: "pr-1", Name: "Feature", AuthorID: "u2", Status: models.PRStatusOpen},
			}, nil
		},
	}

	handler := NewUserHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u1", nil)
	rec := httptest.NewRecorder()

	handler.GetPRs(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp GetPRsResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "u1", resp.UserID)
	assert.Len(t, resp.PRs, 1)
}

func TestUserHandler_GetPRs_InternalError(t *testing.T) {
	service := &mockUserService{
		getPRsFn: func(ctx context.Context, userID string) ([]models.PullRequest, error) {
			return nil, assert.AnError
		},
	}

	handler := NewUserHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u1", nil)
	rec := httptest.NewRecorder()

	handler.GetPRs(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
