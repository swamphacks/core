package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
	logger      zerolog.Logger
}

func NewTeamHandler(teamService *services.TeamService, logger zerolog.Logger) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		logger:      logger.With().Str("handler", "TeamHandler").Str("component", "team").Logger(),
	}
}

// Get the authenticated user's team for this event, including its members.
//
//	@Summary		Get the authenticated user's team and its members for this specific event.
//	@Description	Retrieves the team information and the full list of team members for the currently authenticated user within a specified event.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string						true	"The authenticated session token/id"
//	@Param			event_id		path		int							true	"The ID of the event"
//	@Success		200				{object}	services.TeamWithMembers	"Team information and members successfully retrieved."
//	@Failure		401				{object}	response.ErrorResponse		"Unauthenticated: Requester is not currently authenticated."
//	@Failure		404				{object}	response.ErrorResponse		"Team not found for the user in this event."
//	@Failure		500				{object}	response.ErrorResponse		"Something went seriously wrong."
//
//	@Router			/events/{eventId}/teams/me [get]
func (h *TeamHandler) GetMyTeam(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	teamWithMembers, err := h.teamService.GetUserTeamWithMembers(r.Context(), *userId, eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	if teamWithMembers == nil {
		res.SendError(w, http.StatusNotFound, res.NewError("no_team", "user does not have a team for this event"))
		return
	}

	res.Send(w, http.StatusOK, teamWithMembers)
}

// Get team by ID
//
//	@Summary		Get a team and its members by team id.
//	@Description	Retrieves the team information and the full list of team members by a team id.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string						true	"The authenticated session token/id"
//	@Param			team_id			path		string						true	"The ID of the team"
//	@Success		200				{object}	services.TeamWithMembers	"Team information and members successfully retrieved."
//	@Failure		401				{object}	response.ErrorResponse		"Unauthenticated: Requester is not currently authenticated."
//	@Failure		404				{object}	response.ErrorResponse		"Team not found for the user in this event."
//	@Failure		500				{object}	response.ErrorResponse		"Something went seriously wrong."
//
//	@Router			/teams/{teamId} [get]
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamIdStr := chi.URLParam(r, "teamId")
	if teamIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}
	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	teamWithMembers, err := h.teamService.GetTeamWithMembers(r.Context(), teamId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	if teamWithMembers == nil {
		res.SendError(w, http.StatusNotFound, res.NewError("not_found", "could not find the specified team"))
		return
	}

	res.Send(w, http.StatusOK, teamWithMembers)
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

// Create a new team
//
//	@Summary		Create a new team
//	@Description	Creates a new team for a specific event and assigns the creator as the owner.
//	@Tags			Team
//
//	@Accept			json
//	@Product		json
//
//	@Param			sh_session_id	cookie		string					true	"The authenticated session token/id"
//	@Param			event_id		path		int						true	"The ID of the event"
//
//	@Param			request			body		CreateTeamRequest		true	"Team Creation Payload"
//
//	@Success		200				{object}	sqlc.Team				"A team object"
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//
//	@Failure		400				{object}	response.ErrorResponse	"Bad request: you had request parameters needed for this method."
//	@Failure		409				{object}	response.ErrorResponse	"Conflict: You already have a team."
//
//	@Failure		500				{object}	response.ErrorResponse	"Something went seriously wrong."
//
//	@Router			/events/{eventId}/teams [post]
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	// Extract from body
	var body CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Your request is missing a field"))
		return
	}
	defer r.Body.Close()

	team, err := h.teamService.CreateTeam(r.Context(), body.Name, eventId, *userId)
	if err != nil {
		if errors.Is(err, services.ErrTeamExists) {
			res.SendError(w, http.StatusConflict, res.NewError("team_exists", "the current user already has a team!"))
			return
		}

		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusOK, team)
}
