package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpErr "github.com/437d5/pr-review-manager/internal/application/http"
	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPRService struct {
	createFn   func(ctx context.Context, pr models.PullRequest) (models.PullRequest, error)
	mergeFn    func(ctx context.Context, id string) (models.PullRequest, error)
	reassignFn func(ctx context.Context, prID, oldReviewerID string) (models.PullRequest, string, error)
}

func (m *mockPRService) CreatePR(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
	return m.createFn(ctx, pr)
}

func (m *mockPRService) Merge(ctx context.Context, id string) (models.PullRequest, error) {
	return m.mergeFn(ctx, id)
}

func (m *mockPRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (models.PullRequest, string, error) {
	return m.reassignFn(ctx, prID, oldReviewerID)
}

// helper to read response body
func readBody(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(b)
}

func TestPRHandler_CreatePR_Success(t *testing.T) {
	svc := &mockPRService{
		createFn: func(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
			pr.Status = models.PRStatusOpen
			pr.AssignedReviewers = []string{"u2", "u3"}
			return pr, nil
		},
	}
	handler := &PRHandler{prService: svc}

	body := CreatePRRequest{
		ID:       "pr-1",
		Name:     "Feature",
		AuthorID: "u1",
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(data))
	rec := httptest.NewRecorder()

	handler.CreatePR(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resp PRResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, body.ID, resp.PullRequest.ID)
	assert.Equal(t, body.Name, resp.PullRequest.Name)
	assert.Equal(t, body.AuthorID, resp.PullRequest.AuthorID)
	assert.Equal(t, models.PRStatusOpen, resp.PullRequest.Status)
	assert.Len(t, resp.PullRequest.AssignedReviewers, 2)
}

func TestPRHandler_CreatePR_InvalidJSON(t *testing.T) {
	svc := &mockPRService{
		createFn: func(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
			t.Fatalf("CreatePR should not be called on invalid JSON")
			return models.PullRequest{}, nil
		},
	}
	handler := &PRHandler{prService: svc}

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{invalid json"))
	rec := httptest.NewRecorder()

	handler.CreatePR(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	body := readBody(t, res.Body)
	assert.Contains(t, body, "Invalid JSON")
}

func TestPRHandler_CreatePR_ErrorMapping_NotFound(t *testing.T) {
	svc := &mockPRService{
		createFn: func(ctx context.Context, pr models.PullRequest) (models.PullRequest, error) {
			return models.PullRequest{}, models.ErrUserNotFound
		},
	}
	handler := &PRHandler{prService: svc}

	body := CreatePRRequest{
		ID:       "pr-1",
		Name:     "Feature",
		AuthorID: "u-missing",
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(data))
	rec := httptest.NewRecorder()

	handler.CreatePR(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var errResp httpErr.ErrorResponse
	err = json.NewDecoder(res.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, httpErr.ErrNotFound.Error.Code, errResp.Error.Code)
}
