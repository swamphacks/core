package team

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func RegisterRoutes(teamHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "get-my-team",
		Method:        http.MethodGet,
		Summary:       "Get My Team",
		Description:   "Returns the team information and the full list of team members for the currently authenticated user",
		Tags:          []string{"Team"},
		Path:          "/me",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetMyTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "get-team",
		Method:        http.MethodGet,
		Summary:       "Get Team",
		Description:   "Returns the team information and the full list of team members by team id",
		Tags:          []string{"Team"},
		Path:          "/{teamId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "create-team",
		Method:        http.MethodPost,
		Summary:       "Create Team",
		Description:   "Creates a new team and assigns the user as the owner. Returns the team.",
		Tags:          []string{"Team"},
		Path:          "",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusConflict, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleCreateTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "leave-team",
		Method:        http.MethodPost,
		Summary:       "Leave Team",
		Description:   "Leaves a team if the user is on the team.",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/leave",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleLeaveTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "create-join-team-request",
		Method:        http.MethodPost,
		Summary:       "Request to Join Team",
		Description:   "Requests to join a team or fails if user is already on a team.",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/join",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusConflict, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleRequestToJoinTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "get-pending-join-team-requests",
		Method:        http.MethodGet,
		Summary:       "Get Pending Join Requests for Team",
		Description:   "Returns a team's pending join requests. This is only allowed for the team's owner.",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/pending-joins",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetPendingRequestsForTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "get-my-pending-join-requests",
		Method:        http.MethodGet,
		Summary:       "Get User's Pending Join Requests",
		Description:   "Returns the current user's pending requests for teams.",
		Tags:          []string{"Team"},
		Path:          "/me/pending-joins",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetMyPendingRequests)

	huma.Register(group, huma.Operation{
		OperationID:   "accept-team-join-request",
		Method:        http.MethodPost,
		Summary:       "Accept Team Join Request",
		Description:   "Accepts a pending team join request. Only the team owner can perform this action.",
		Tags:          []string{"Team"},
		Path:          "/{requestId}/accept",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleAcceptTeamJoinRequest)

	huma.Register(group, huma.Operation{
		OperationID:   "reject-team-join-request",
		Method:        http.MethodPost,
		Summary:       "Reject Team Join Request",
		Description:   "Rejects a pending team join request. Only the team owner can perform this action.",
		Tags:          []string{"Team"},
		Path:          "/{requestId}/reject",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleRejectTeamJoinRequest)

	huma.Register(group, huma.Operation{
		OperationID:   "kick-member-from-team",
		Method:        http.MethodPost,
		Summary:       "Kick Team Member",
		Description:   "Kicks a member from a team. Only the team owner can perform this action.",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/kick/{memberId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleKickMemberFromTeam)
}

type handler struct {
	teamService *TeamService
	logger      zerolog.Logger
}

func NewHandler(teamService *TeamService, logger zerolog.Logger) *handler {
	return &handler{
		teamService: teamService,
		logger:      logger.With().Str("handler", "TeamHandler").Str("domain", "team").Logger(),
	}
}

type GetMyTeamOutput struct {
	Body *TeamWithMembers
}

func (h *handler) handleGetMyTeam(ctx context.Context, input *struct{}) (*GetMyTeamOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	team, err := h.teamService.GetUserTeamWithMembers(ctx, *userId)
	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to get my team")
	}

	if team == nil {
		return nil, huma.Error404NotFound("user does not have a team")
	}

	return &GetMyTeamOutput{Body: team}, nil
}

type GetTeamOutput struct {
	Body *TeamWithMembers
}

func (h *handler) handleGetTeam(ctx context.Context, input *struct {
	TeamId string `path:"teamId"`
}) (*GetTeamOutput, error) {
	teamId, err := uuid.Parse(input.TeamId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}

	team, err := h.teamService.GetTeamWithMembersByTeamId(ctx, teamId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to get team by id")
	}

	if team == nil {
		return nil, huma.Error404NotFound("Could not find team associated with this team id")
	}

	return &GetTeamOutput{Body: team}, nil
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type CreateTeamOutput struct {
	Body *sqlc.Team
}

func (h *handler) handleCreateTeam(ctx context.Context, input *struct {
	Body CreateTeamRequest
}) (*CreateTeamOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	team, err := h.teamService.CreateTeam(ctx, input.Body.Name, *userId)

	if err != nil {
		if errors.Is(err, ErrTeamExists) {
			return nil, huma.Error409Conflict("User already in a team")
		}

		return nil, huma.Error500InternalServerError("Fail to create team")
	}

	return &CreateTeamOutput{Body: team}, nil
}

type LeaveTeamOutput struct {
	Status int
}

func (h *handler) handleLeaveTeam(ctx context.Context, input *struct {
	TeamId string `path:"teamId"`
}) (*LeaveTeamOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	teamId, err := uuid.Parse(input.TeamId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}

	err = h.teamService.LeaveTeam(ctx, *userId, teamId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to leave team")
	}

	return &LeaveTeamOutput{Status: http.StatusOK}, nil
}

type CreateJoinRequest struct {
	Message *string `json:"message"`
}

type RequestToJoinTeamOutput struct {
	Body *sqlc.TeamJoinRequest
}

func (h *handler) handleRequestToJoinTeam(ctx context.Context, input *struct {
	Body   CreateJoinRequest
	TeamId string `path:"teamId"`
}) (*RequestToJoinTeamOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	teamId, err := uuid.Parse(input.TeamId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}

	request, err := h.teamService.RequestToJoinTeam(ctx, teamId, *userId, input.Body.Message)

	if err != nil {
		if errors.Is(err, ErrUserOnTeam) {
			return nil, huma.Error409Conflict("User is already on a team")
		}

		return nil, huma.Error500InternalServerError("Fail to create join team request")
	}

	return &RequestToJoinTeamOutput{Body: request}, nil
}

type GetPendingRequestsForTeamOutput struct {
	Body []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow
}

func (h *handler) handleGetPendingRequestsForTeam(ctx context.Context, input *struct {
	TeamId string `path:"teamId"`
}) (*GetPendingRequestsForTeamOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	teamId, err := uuid.Parse(input.TeamId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}

	requests, err := h.teamService.GetPendingJoinRequestForTeam(ctx, *userId, teamId)

	if err != nil {
		if errors.Is(err, ErrUserNotTeamOwner) {
			return nil, huma.Error403Forbidden("Not authorized to perform this action")
		}

		return nil, huma.Error500InternalServerError("Fail to get pending requests for team")
	}

	return &GetPendingRequestsForTeamOutput{Body: requests}, nil
}

type GetMyPendingRequestsOutput struct {
	Body []sqlc.TeamJoinRequest
}

func (h *handler) handleGetMyPendingRequests(ctx context.Context, input *struct{}) (*GetMyPendingRequestsOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	requests, err := h.teamService.GetUserPendingJoinRequests(ctx, *userId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to get user's pending requests")
	}

	return &GetMyPendingRequestsOutput{Body: requests}, nil
}

type AcceptTeamJoinRequestOutput struct {
	Status int
}

func (h *handler) handleAcceptTeamJoinRequest(ctx context.Context, input *struct {
	RequestId string `path:"requestId"`
}) (*AcceptTeamJoinRequestOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	requestId, err := uuid.Parse(input.RequestId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid join request id")
	}

	err = h.teamService.RespondToJoinRequest(ctx, *userId, requestId, true)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to accept team join request")
	}

	return &AcceptTeamJoinRequestOutput{Status: http.StatusOK}, nil
}

type RejectTeamJoinRequestOutput struct {
	Status int
}

func (h *handler) handleRejectTeamJoinRequest(ctx context.Context, input *struct {
	RequestId string `path:"requestId"`
}) (*RejectTeamJoinRequestOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	requestId, err := uuid.Parse(input.RequestId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid join request id")
	}

	err = h.teamService.RespondToJoinRequest(ctx, *userId, requestId, false)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to accept team join request")
	}

	return &RejectTeamJoinRequestOutput{Status: http.StatusOK}, nil
}

type KickMemberFromTeamOutput struct {
	Status int
}

func (h *handler) handleKickMemberFromTeam(ctx context.Context, input *struct {
	MemberId string `path:"memberId"`
	TeamId   string `path:"teamId"`
}) (*KickMemberFromTeamOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	memberId, err := uuid.Parse(input.MemberId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid memberid")
	}

	teamId, err := uuid.Parse(input.TeamId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}

	err = h.teamService.KickMemberFromTeam(ctx, memberId, teamId, *userId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to kick member from team")
	}

	return &KickMemberFromTeamOutput{Status: http.StatusOK}, nil
}
