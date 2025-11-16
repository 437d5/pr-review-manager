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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(httpErr.ErrTeamExists)
		case models.ErrTeamNameEmpty, models.ErrTeamMembersEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	createdTeam = hideTeamName(createdTeam)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TeamResponse{Team: createdTeam})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		switch err {
		case models.ErrTeamNameEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case models.ErrTeamNotFound:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(httpErr.ErrNotFound)
		default:
			slog.Error("failed get team", "error", err.Error(), "team_name", teamName)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	team = hideTeamName(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TeamResponse{Team: team})
}

func hideTeamName(team models.Team) models.Team {
	for i := range team.Members {
		team.Members[i].TeamName = ""
	}

	return team
}
