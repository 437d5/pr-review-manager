package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httpErr "github.com/437d5/pr-review-manager/internal/application/http"
	"github.com/437d5/pr-review-manager/internal/domain/models"
	"github.com/437d5/pr-review-manager/internal/domain/services"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

type CreateTeamRequest struct {
	Name    string        `json:"team_name"`
	Members []models.User `json:"members"`
}

type TeamResponse struct {
	Team models.Team `json:"team"`
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	team := models.Team{
		Name:    req.Name,
		Members: req.Members,
	}

	createdTeam, err := h.teamService.CreateTeam(r.Context(), team)
	if err != nil {
		switch err {
		case models.ErrTeamExists:
			httpErr.WriteError(w, http.StatusBadRequest, httpErr.ErrTeamExists)
		case models.ErrTeamNameEmpty, models.ErrTeamMembersEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			httpErr.WriteInernalError(w, err)
		}
		return
	}

	createdTeam = hideTeamName(createdTeam)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(TeamResponse{Team: createdTeam})
	if err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		switch err {
		case models.ErrTeamNameEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case models.ErrTeamNotFound:
			httpErr.WriteError(w, http.StatusNotFound, httpErr.ErrNotFound)
		default:
			httpErr.WriteInernalError(w, err)
		}
		return
	}

	team = hideTeamName(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(TeamResponse{Team: team})
	if err != nil {
		slog.Error("cannot encode response", "error", err, "team", team.Name)
	}
}

func hideTeamName(team models.Team) models.Team {
	for i := range team.Members {
		team.Members[i].TeamName = ""
	}

	return team
}
