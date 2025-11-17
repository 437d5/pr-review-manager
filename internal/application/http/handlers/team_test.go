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

type mockTeamService struct {
	createFn func(ctx context.Context, team models.Team) (models.Team, error)
	getFn    func(ctx context.Context, name string) (models.Team, error)
}

func (m *mockTeamService) CreateTeam(ctx context.Context, team models.Team) (models.Team, error) {
	return m.createFn(ctx, team)
}

func (m *mockTeamService) GetTeam(ctx context.Context, name string) (models.Team, error) {
	return m.getFn(ctx, name)
}

func TestTeamHandler_CreateTeam_Success(t *testing.T) {
	service := &mockTeamService{
		createFn: func(ctx context.Context, team models.Team) (models.Team, error) {
			team.Members[0].TeamName = team.Name
			return team, nil
		},
	}

	handler := NewTeamHandler(service)

	reqBody := CreateTeamRequest{
		Name: "backend",
		Members: []models.User{
			{ID: "u1", Username: "Alice", IsActive: true},
		},
	}

	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.CreateTeam(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resp TeamResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, reqBody.Name, resp.Team.Name)
	require.Len(t, resp.Team.Members, 1)
	assert.Equal(t, "u1", resp.Team.Members[0].ID)
	assert.Empty(t, resp.Team.Members[0].TeamName, "handler should hide team name")
}

func TestTeamHandler_CreateTeam_InvalidJSON(t *testing.T) {
	service := &mockTeamService{
		createFn: func(ctx context.Context, team models.Team) (models.Team, error) {
			t.Fatal("CreateTeam should not be called on invalid JSON")
			return models.Team{}, nil
		},
	}

	handler := NewTeamHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBufferString("{invalid"))
	rec := httptest.NewRecorder()

	handler.CreateTeam(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestTeamHandler_CreateTeam_TeamExists(t *testing.T) {
	service := &mockTeamService{
		createFn: func(ctx context.Context, team models.Team) (models.Team, error) {
			return models.Team{}, models.ErrTeamExists
		},
	}

	handler := NewTeamHandler(service)

	reqBody := CreateTeamRequest{
		Name: "backend",
		Members: []models.User{
			{ID: "u1", Username: "Alice", IsActive: true},
		},
	}
	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.CreateTeam(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var errResp httpErr.ErrorResponse
	err = json.NewDecoder(res.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, httpErr.ErrTeamExists.Error.Code, errResp.Error.Code)
}

func TestTeamHandler_GetTeam_Success(t *testing.T) {
	service := &mockTeamService{
		getFn: func(ctx context.Context, name string) (models.Team, error) {
			return models.Team{
				Name: name,
				Members: []models.User{
					{ID: "u1", Username: "Alice", TeamName: name},
				},
			}, nil
		},
	}

	handler := NewTeamHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=backend", nil)
	rec := httptest.NewRecorder()

	handler.GetTeam(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp TeamResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "backend", resp.Team.Name)
	assert.Empty(t, resp.Team.Members[0].TeamName)
}

func TestTeamHandler_GetTeam_NotFound(t *testing.T) {
	service := &mockTeamService{
		getFn: func(ctx context.Context, name string) (models.Team, error) {
			return models.Team{}, models.ErrTeamNotFound
		},
	}

	handler := NewTeamHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=backend", nil)
	rec := httptest.NewRecorder()

	handler.GetTeam(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var errResp httpErr.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, httpErr.ErrNotFound.Error.Code, errResp.Error.Code)
}
