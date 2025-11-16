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
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/ptr"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
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
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_team_id", "The team ID is missing from the URL!"))
		return
	}
	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_team_id", "The team ID is not a valid UUID"))
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

// Gets an events teams
//
//	@Summary		Get an event's teams
//	@Description	Gets all teams for a specific event.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string						true	"The authenticated session token/id"
//	@Param			event_id		path		string						true	"The ID of the event"
//	@Success		200				{array}		services.TeamWithMembers	"Teams successfully retrieved."
//	@Failure		400				{object}	response.ErrorResponse		"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse		"Unauthenticated: Requester is not currently authenticated."
//	@Failure		500				{object}	response.ErrorResponse		"Something went seriously wrong."
//
//	@Router			/events/{eventId}/teams [get]
func (h *TeamHandler) GetEventTeams(w http.ResponseWriter, r *http.Request) {
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

	// Parse limit and offset from query parameters
	query := r.URL.Query()
	limit, err := web.ParseParamInt32(query, "limit", ptr.Int32ToPtr(10))
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_params", "Malformed `limit` query parameter."))
		return
	}

	offset, err := web.ParseParamInt32(query, "offset", ptr.Int32ToPtr(0))
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_params", "Malformed `offset` query parameter."))
		return
	}

	teams, err := h.teamService.GetTeamsWithMembersByEvent(r.Context(), eventId, *limit, *offset)
	if err != nil {
		h.logger.Err(err).Msg("Failed to get teams for event")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusOK, teams)
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
//	@Produce		json
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

// Leave a team
//
//	@Summary		Leave a team
//	@Description	Leaves a team if the requester is on the team. Depends on cookies for user retrieval.
//	@Tags			Team
//	@Param			sh_session_id	cookie	string	true	"The authenticated session token/id"
//	@Param			team_id			path	string	true	"The ID of the team"
//	@Success		204				"Successfully left the team"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		500				{object}	response.ErrorResponse	"Something went seriously wrong."
//
//	@Router			/teams/{teamId}/members/me [delete]
func (h *TeamHandler) LeaveTeam(w http.ResponseWriter, r *http.Request) {
	teamIdStr := chi.URLParam(r, "teamId")
	if teamIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_team_id", "The team ID is missing from the URL!"))
		return
	}
	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_team_id", "The team ID is not a valid UUID"))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	err = h.teamService.LeaveTeam(r.Context(), *userId, teamId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))

	}
}

type CreateJoinRequest struct {
	Message *string `json:"message"`
}

// Request to join a team
//
//	@Summary		Request to join a team
//	@Description	Requests to join a team or fails if user is already on a team.
//	@Tags			Team
//	@Param			sh_session_id	cookie	string				true	"The authenticated session token/id"
//	@Param			team_id			path	string				true	"The ID of the team"
//	@Param			event_id		path	string				true	"The ID of the event"
//	@Param			request			body	CreateJoinRequest	true	"Team Creation Payload"
//	@Accept			json
//	@Produce		json
//	@Success		204	"Successfully left the team"
//	@Failure		400	{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401	{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//
//	@Failure		409	{object}	response.ErrorResponse	"Conflict: User is already on a team."
//
//	@Failure		500	{object}	response.ErrorResponse	"Something went seriously wrong."
//
//	@Router			/events/{eventId}/teams/{teamId}/join [post]
func (h *TeamHandler) RequestToJoinTeam(w http.ResponseWriter, r *http.Request) {
	teamId, err := web.PathParamToUUID(r, "teamId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_team_id", "The team ID is malformed/missing."))
		return
	}

	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_event_id", "The event ID is malformed/missing."))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	var body CreateJoinRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Your request is missing a field"))
		return
	}
	defer r.Body.Close()

	request, err := h.teamService.RequestToJoinTeam(r.Context(), eventId, teamId, *userId, body.Message)
	if err != nil {
		if errors.Is(err, services.ErrUserOnTeam) {
			res.SendError(w, http.StatusConflict, res.NewError("already_teamed", "You are already on a team."))
			return
		}

		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusOK, request)
}

// Get a team's pending join requets
//
//	@Summary		Get team's pending join requests
//	@Description	Retrieves a team's pending join requests. This is only allowed for the team's owner.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string											true	"The authenticated session token/id"
//	@Param			team_id			path		string											true	"The ID of the team"
//	@Success		200				{array}		sqlc.ListJoinRequestsByTeamAndStatusWithUserRow	"Successfully retrieved pending requests"
//	@Failure		400				{object}	response.ErrorResponse							"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse							"Unauthenticated: Requester is not currently authenticated."
//	@Failure		403				{object}	response.ErrorResponse							"Forbidden: Requester is not allowed to perform this action."
//	@Failure		500				{object}	response.ErrorResponse							"Something went seriously wrong."
//
//	@Router			/teams/{teamId}/pending-joins [get]
func (h *TeamHandler) GetPendingRequestsForTeam(w http.ResponseWriter, r *http.Request) {
	teamId, err := web.PathParamToUUID(r, "teamId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_team_id", "The team ID is malformed/missing."))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	requests, err := h.teamService.GetPendingJoinRequestForTeam(r.Context(), *userId, teamId)
	if err != nil {
		if errors.Is(err, services.ErrUserNotTeamOwner) {
			res.SendError(w, http.StatusForbidden, res.NewError("forbidden", "You do not have the permissions for this action."))
			return
		}

		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusOK, requests)
}

// Get your pending requests for an event
//
//	@Summary		Get your pending requests
//	@Description	Retrieves the current user's pending requests for a specific event's teams.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string					true	"The authenticated session token/id"
//	@Param			team_id			path		string					true	"The ID of the team"
//	@Success		200				{array}		sqlc.TeamJoinRequest	"Successfully retrieved pending requests"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		500				{object}	response.ErrorResponse	"Something went seriously wrong."
//
//	@Router			/events/{eventId}/teams/me/pending-joins [get]
func (h *TeamHandler) GetMyPendingRequests(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_event_id", "The event ID is malformed/missing."))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	requests, err := h.teamService.GetUserPendingJoinRequestsByEvent(r.Context(), *userId, eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusOK, requests)
}

// Accept/Approve a team join request
//
//	@Summary		Accept a team join request
//	@Description	Accepts a pending team join request. Only the team owner can perform this action.
//	@Tags			Team
//	@Param			sh_session_id	cookie	string	true	"The authenticated session token/id"
//	@Param			team_id			path	string	true	"The ID of the team"
//	@Param			request_id		path	string	true	"The ID of the join request"
//	@Success		204				"Successfully accepted the join request"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		403				{object}	response.ErrorResponse	"Forbidden: Requester is not allowed to perform this action."
//
//	@Failure		404				{object}	response.ErrorResponse	"Not Found: The join request does not exist."
//	@Failure		409				{object}	response.ErrorResponse	"Conflict: The join request has already been responded to."
//
//	@Failure		500				{object}	response.ErrorResponse	"Something went wrong."
//
//	@Router			/teams/join/{requestId}/accept [post]
func (h *TeamHandler) AcceptTeamJoinRequest(w http.ResponseWriter, r *http.Request) {
	requestId, err := web.PathParamToUUID(r, "requestId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_request_id", "The join request ID is malformed/missing."))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	err = h.teamService.RespondToJoinRequest(r.Context(), *userId, requestId, true)
	if err != nil {
		status, code, message := mapTeamServiceError(err)
		res.SendError(w, status, res.NewError(code, message))
		return
	}

	res.Send(w, http.StatusNoContent, nil)
}

// Reject a team join request
//
//	@Summary		Reject a team join request
//	@Description	Rejects a pending team join request. Only the team owner can perform this action.
//	@Tags			Team
//	@Param			sh_session_id	cookie	string	true	"The authenticated session token/id"
//	@Param			team_id			path	string	true	"The ID of the team"
//	@Param			request_id		path	string	true	"The ID of the join request"
//	@Success		204				"Successfully accepted the join request"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		403				{object}	response.ErrorResponse	"Forbidden: Requester is not allowed to perform this action."
//
//	@Failure		404				{object}	response.ErrorResponse	"Not Found: The join request does not exist."
//	@Failure		409				{object}	response.ErrorResponse	"Conflict: The join request has already been responded to."
//
//	@Failure		500				{object}	response.ErrorResponse	"Something went wrong."
//
//	@Router			/teams/join/{requestId}/reject [post]
func (h *TeamHandler) RejectTeamJoinRequest(w http.ResponseWriter, r *http.Request) {
	requestId, err := web.PathParamToUUID(r, "requestId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_request_id", "The join request ID is malformed/missing."))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	err = h.teamService.RespondToJoinRequest(r.Context(), *userId, requestId, false)
	if err != nil {
		status, code, message := mapTeamServiceError(err)
		res.SendError(w, status, res.NewError(code, message))
		return
	}

	res.Send(w, http.StatusNoContent, nil)
}

type InviteUserRequest struct {
	Email string `json:"email"`
}

// Invite user to team
//
//	@Summary		Invite user to team
//	@Description	Invites a user to a team by email. Only the team owner can invite users.
//	@Tags			Team
//	@Param			sh_session_id	cookie			string				true	"The authenticated session token/id"
//	@Param			teamId			path			string				true	"The ID of the team"
//	@Param			request			body			InviteUserRequest	true	"Invitation request payload"
//	@Accept			json
//	@Produce		json
//	@Success		204				"Successfully sent invitation"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		403				{object}	response.ErrorResponse	"Forbidden: User is not the team owner."
//	@Failure		404				{object}	response.ErrorResponse	"Not Found: Team not found."
//	@Failure		409				{object}	response.ErrorResponse	"Conflict: Invitation already exists."
//	@Failure		500				{object}	response.ErrorResponse	"Something went wrong."
//
//	@Router			/teams/{teamId}/invite [post]
func (h *TeamHandler) InviteUserToTeam(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	teamIdStr := chi.URLParam(r, "teamId")
	if teamIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_team_id", "The team ID is missing from the URL!"))
		return
	}

	teamId, err := uuid.Parse(teamIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_team_id", "The team ID is not a valid UUID"))
		return
	}

	var body InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Your request is missing a field"))
		return
	}
	defer r.Body.Close()

	if body.Email == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_email", "Email is required"))
		return
	}

	err = h.teamService.InviteUserToTeam(r.Context(), teamId, *userId, body.Email)
	if err != nil {
		if errors.Is(err, services.ErrUserNotTeamOwner) {
			res.SendError(w, http.StatusForbidden, res.NewError("forbidden", "You do not have permission to invite users to this team"))
			return
		}
		if errors.Is(err, services.ErrTeamNotFound) {
			res.SendError(w, http.StatusNotFound, res.NewError("team_not_found", "Team not found"))
			return
		}
		if errors.Is(err, services.ErrInvitationAlreadyExists) {
			res.SendError(w, http.StatusConflict, res.NewError("invitation_exists", "A pending invitation already exists for this email and team"))
			return
		}
		h.logger.Err(err).Msg("Failed to invite user to team")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusNoContent, nil)
}

// Get invitation details
//
//	@Summary		Get invitation details
//	@Description	Retrieves invitation details including team name, inviter name, and event name. Works for unauthenticated users.
//	@Tags			Team
//	@Param			invitationId	path		string						true	"The ID of the invitation"
//	@Success		200				{object}	services.InvitationDetails	"Invitation details successfully retrieved."
//	@Failure		404				{object}	response.ErrorResponse		"Invitation not found or expired."
//	@Failure		500				{object}	response.ErrorResponse		"Something went wrong."
//
//	@Router			/teams/invite/{invitationId} [get]
func (h *TeamHandler) GetInvitation(w http.ResponseWriter, r *http.Request) {
	invitationIdStr := chi.URLParam(r, "invitationId")
	if invitationIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_invitation_id", "The invitation ID is missing from the URL!"))
		return
	}

	invitationId, err := uuid.Parse(invitationIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_invitation_id", "The invitation ID is not a valid UUID"))
		return
	}

	details, err := h.teamService.GetInvitationDetails(r.Context(), invitationId)
	if err != nil {
		if errors.Is(err, services.ErrInvitationNotFound) || errors.Is(err, services.ErrInvitationExpired) || errors.Is(err, services.ErrInvitationAlreadyAccepted) {
			res.SendError(w, http.StatusNotFound, res.NewError("invitation_not_found", "Invitation not found, expired, or already accepted"))
			return
		}
		h.logger.Err(err).Msg("Failed to get invitation details")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	res.Send(w, http.StatusOK, details)
}

// Accept invitation
//
//	@Summary		Accept team invitation
//	@Description	Accepts a team invitation. Requires authentication and the user's email must match the invitation email. User must have a submitted application for the event.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string	true	"The authenticated session token/id"
//	@Param			invitationId	path		string	true	"The ID of the invitation"
//	@Success		204				"Successfully accepted the invitation"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		403				{object}	response.ErrorResponse	"Forbidden: Email mismatch or application required."
//	@Failure		404				{object}	response.ErrorResponse	"Not Found: Invitation not found or expired."
//	@Failure		409				{object}	response.ErrorResponse	"Conflict: Invitation already accepted."
//	@Failure		500				{object}	response.ErrorResponse	"Something went wrong."
//
//	@Router			/teams/invite/{invitationId}/accept [post]
func (h *TeamHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	invitationIdStr := chi.URLParam(r, "invitationId")
	if invitationIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_invitation_id", "The invitation ID is missing from the URL!"))
		return
	}

	invitationId, err := uuid.Parse(invitationIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_invitation_id", "The invitation ID is not a valid UUID"))
		return
	}

	err = h.teamService.AcceptInvitation(r.Context(), invitationId, *userId)
	if err != nil {
		status, code, message := mapInvitationServiceError(err)
		res.SendError(w, status, res.NewError(code, message))
		return
	}

	res.Send(w, http.StatusNoContent, nil)
}

// Reject invitation
//
//	@Summary		Reject team invitation
//	@Description	Rejects a team invitation. Requires authentication and the user's email must match the invitation email.
//	@Tags			Team
//	@Param			sh_session_id	cookie		string	true	"The authenticated session token/id"
//	@Param			invitationId	path		string	true	"The ID of the invitation"
//	@Success		204				"Successfully rejected the invitation"
//	@Failure		400				{object}	response.ErrorResponse	"Bad Request: Missing or malformed parameters."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		403				{object}	response.ErrorResponse	"Forbidden: Email mismatch."
//	@Failure		404				{object}	response.ErrorResponse	"Not Found: Invitation not found."
//	@Failure		500				{object}	response.ErrorResponse	"Something went wrong."
//
//	@Router			/teams/invite/{invitationId}/reject [post]
func (h *TeamHandler) RejectInvitation(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	invitationIdStr := chi.URLParam(r, "invitationId")
	if invitationIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_invitation_id", "The invitation ID is missing from the URL!"))
		return
	}

	invitationId, err := uuid.Parse(invitationIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_invitation_id", "The invitation ID is not a valid UUID"))
		return
	}

	err = h.teamService.RejectInvitation(r.Context(), invitationId, *userId)
	if err != nil {
		status, code, message := mapInvitationServiceError(err)
		res.SendError(w, status, res.NewError(code, message))
		return
	}

	res.Send(w, http.StatusNoContent, nil)
}

// Maps team service errors to HTTP status codes and messages
func mapTeamServiceError(err error) (status int, code, message string) {
	switch {
	case errors.Is(err, services.ErrUserNotTeamOwner):
		return http.StatusForbidden, "forbidden", "You do not have permission to perform this action."
	case errors.Is(err, services.ErrUserNotApplicantOrAttendee):
		return http.StatusForbidden, "invalid_user_role", "The user is not an applicant or attendee for this event."
	case errors.Is(err, services.ErrUserOnTeam):
		return http.StatusConflict, "user_on_team", "The user is already on a team for this event."
	case errors.Is(err, services.ErrTeamFull):
		return http.StatusConflict, "team_full", "The team is already full."
	case errors.Is(err, repository.ErrTeamNotFound):
		return http.StatusNotFound, "team_not_found", "Team resource was not found."
	default:
		return http.StatusInternalServerError, "internal_error", "Something went wrong."
	}
}

// Maps invitation service errors to HTTP status codes and messages
func mapInvitationServiceError(err error) (status int, code, message string) {
	switch {
	case errors.Is(err, services.ErrInvitationNotFound):
		return http.StatusNotFound, "invitation_not_found", "Invitation not found."
	case errors.Is(err, services.ErrInvitationExpired):
		return http.StatusNotFound, "invitation_expired", "Invitation has expired."
	case errors.Is(err, services.ErrInvitationAlreadyAccepted):
		return http.StatusConflict, "invitation_already_accepted", "Invitation has already been accepted."
	case errors.Is(err, services.ErrEmailMismatch):
		return http.StatusForbidden, "email_mismatch", "Your email does not match the invitation email."
	case errors.Is(err, services.ErrApplicationRequired):
		return http.StatusForbidden, "application_required", "Please complete your application first."
	default:
		return http.StatusInternalServerError, "internal_error", "Something went wrong."
	}
}
