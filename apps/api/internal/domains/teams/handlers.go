package teams

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type GetMyTeamOutput struct {
	Body TeamDetailsDto
}

func parseMemberStatus(member *sqlc.GetTeamMembersRow) TeamMemberStatus {
	if member == nil {
		return TeamMemberStatusNotAccepted
	}

	if member.ApplicationStatus != nil {
		switch *member.ApplicationStatus {
		case sqlc.ApplicationStatusConfirmed:
			return TeamMemberStatusConfirmed
		case sqlc.ApplicationStatusAccepted:
			return TeamMemberStatusAccepted
		}
	} else if member.Role == sqlc.RoleAttendee {
		// I don't believe it is possible for someone to be an "attendee" without first
		// submitting an application, but we add this branch as a safety measure
		return TeamMemberStatusConfirmed
	}
	return TeamMemberStatusNotAccepted
}

func (h *handler) handleGetMyTeam(ctx context.Context, input *struct{}) (*GetMyTeamOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	team, err := h.teamService.GetTeamByUserId(ctx, userCtx.UserID)

	if err != nil {
		if errors.Is(err, ErrNotInTeam) {
			return nil, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	teamMembers, err := h.teamService.GetTeamMembers(ctx, team.ID)

	teamMembersDto := make([]TeamMemberDto, len(teamMembers))

	for i, val := range teamMembers {
		teamMembersDto[i] = TeamMemberDto{
			ID:     val.UserID,
			Name:   val.Name,
			Image:  val.Image,
			Status: parseMemberStatus(&val),
		}
	}

	return &GetMyTeamOutput{Body: TeamDetailsDto{
		ID:      team.ID,
		Name:    team.Name,
		OwnerID: team.OwnerID,
		Members: teamMembersDto,
	}}, nil
}

type GetTeamDetailsOutput struct {
	Body TeamDetailsDto
}

func (h *handler) handleGetTeamDetails(ctx context.Context, input *struct {
	TeamId uuid.UUID `path:"teamId"`
}) (*GetTeamDetailsOutput, error) {
	team, err := h.teamService.GetTeamById(ctx, input.TeamId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to get team details")
	}

	teamMembers, err := h.teamService.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	teamMembersDto := make([]TeamMemberDto, len(teamMembers))

	for i, val := range teamMembers {
		teamMembersDto[i] = TeamMemberDto{
			ID:     val.UserID,
			Name:   val.Name,
			Image:  val.Image,
			Status: parseMemberStatus(&val),
		}
	}

	return &GetTeamDetailsOutput{Body: TeamDetailsDto{
		ID:      team.ID,
		Name:    team.Name,
		OwnerID: team.OwnerID,
		// CreatedAt: teamDetails.CreatedAt,
		Members: teamMembersDto,
	}}, nil
}

type GetTeamByInvitationIdOutput struct {
	Body TeamDto
}

func (h *handler) handleGetTeamByInvitationId(ctx context.Context, input *struct {
	InviteId uuid.UUID `path:"inviteId"`
}) (*GetTeamByInvitationIdOutput, error) {
	team, err := h.teamService.GetTeamByInvitationId(ctx, input.InviteId)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &GetTeamByInvitationIdOutput{Body: TeamDto{
		ID:      team.ID,
		Name:    team.Name,
		OwnerID: team.OwnerID,
	}}, nil
}

type GetTeamMembersOutput struct {
	Body []TeamMemberDto `nullable:"false"`
}

func (h *handler) handleGetTeamMembers(ctx context.Context, input *struct {
	TeamId uuid.UUID `path:"teamId"`
}) (*GetTeamMembersOutput, error) {
	members, err := h.teamService.GetTeamMembers(ctx, input.TeamId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to get team members")
	}

	teamMembersDto := make([]TeamMemberDto, len(members))

	for i, val := range members {
		teamMembersDto[i] = TeamMemberDto{
			ID:     val.UserID,
			Name:   val.Name,
			Image:  val.Image,
			Status: parseMemberStatus(&val),
		}
	}

	return &GetTeamMembersOutput{
		Body: teamMembersDto,
	}, nil
}

type CreateTeamOutput struct {
	Body TeamDto
}

func (h *handler) handleCreateTeam(ctx context.Context, input *struct {
	Body CreateTeamRequestDto
}) (*CreateTeamOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	team, err := h.teamService.CreateTeam(ctx, input.Body.Name, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &CreateTeamOutput{Body: TeamDto{
		ID:        team.ID,
		Name:      team.Name,
		OwnerID:   team.OwnerID,
		CreatedAt: team.CreatedAt,
	}}, nil
}

type DeleteTeamOutput struct {
	Status int
}

func (h *handler) handleDeleteTeam(ctx context.Context, input *struct {
	Body DeleteTeamRequestDto
}) (*DeleteTeamOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.teamService.DeleteTeam(ctx, userCtx.UserID, input.Body.TeamID)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &DeleteTeamOutput{Status: http.StatusOK}, nil
}

type JoinTeamOutput struct {
	Status int
}

func (h *handler) handleJoinTeam(ctx context.Context, input *struct {
	InvitationId uuid.UUID `path:"id"`
}) (*JoinTeamOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	invitation, err := h.teamService.GetInvitation(ctx, input.InvitationId)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Check if team has been "accepted"
	// A team is treated as accepted if any one of its members has been accepted
	members, err := h.teamService.GetTeamMembers(ctx, invitation.TeamID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	teamIsAccepted := false
	for _, member := range members {
		status := parseMemberStatus(&member)
		if status == TeamMemberStatusConfirmed || status == TeamMemberStatusAccepted {
			teamIsAccepted = true
			break
		}
	}
	userIsAccepted := userCtx.Role == sqlc.RoleAttendee

	// A user cannot join a non-accepted team if they are already accepted
	// A user can join an accepted team if they too are accepted
	if userIsAccepted && !teamIsAccepted {
		return nil, huma.Error400BadRequest("You cannot join a non-accepted team if you are already accepted!")
	} else if !userIsAccepted && teamIsAccepted {
		return nil, huma.Error400BadRequest("You cannot join an accepted team if you are not yet accepted!")
	}

	err = h.teamService.JoinTeam(ctx, userCtx.UserID, invitation.TeamID)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &JoinTeamOutput{Status: http.StatusOK}, nil
}

type LeaveTeamOutput struct {
	Status int
}

func (h *handler) handleLeaveTeam(ctx context.Context, input *struct {
	TeamId uuid.UUID `path:"teamId"`
}) (*LeaveTeamOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.teamService.LeaveTeam(ctx, userCtx.UserID, input.TeamId)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &LeaveTeamOutput{Status: http.StatusOK}, nil
}

type KickMemberFromTeamOutput struct {
	Status int
}

func (h *handler) handleKickMember(ctx context.Context, input *struct {
	Body   KickMemberRequestDto
	TeamId uuid.UUID `path:"teamId"`
}) (*KickMemberFromTeamOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.teamService.KickMember(ctx, input.Body.MemberId, input.TeamId, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &KickMemberFromTeamOutput{Status: http.StatusOK}, nil
}

type CreateInvitationOutput struct {
	Body uuid.UUID
}

func (h *handler) handleCreateInvitation(ctx context.Context, input *struct{}) (*CreateInvitationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	invitation, err := h.teamService.CreateInvitation(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &CreateInvitationOutput{Body: *invitation}, nil
}

type GetInvitationOutput struct {
	Body uuid.UUID
}

func (h *handler) handleGetInvitation(ctx context.Context, input *struct {
	TeamID uuid.UUID `path:"teamId"`
}) (*GetInvitationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	invitation, err := h.teamService.GetInvitationByTeamID(ctx, input.TeamID, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &GetInvitationOutput{Body: invitation.ID}, nil
}

// type CreateJoinRequest struct {
// 	Message *string `json:"message"`
// }

// type RequestToJoinTeamOutput struct {
// 	Body *sqlc.TeamJoinRequest
// }

// func (h *handler) handleRequestToJoinTeam(ctx context.Context, input *struct {
// 	Body   CreateJoinRequest
// 	TeamId string `path:"teamId"`
// }) (*RequestToJoinTeamOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	teamId, err := uuid.Parse(input.TeamId)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid team id")
// 	}

// 	request, err := h.teamService.RequestToJoinTeam(ctx, teamId, userCtx.UserID, input.Body.Message)

// 	if err != nil {
// 		if errors.Is(err, ErrUserOnTeam) {
// 			return nil, huma.Error409Conflict("User is already on a team")
// 		}

// 		return nil, huma.Error500InternalServerError("Fail to create join team request")
// 	}

// 	return &RequestToJoinTeamOutput{Body: request}, nil
// }

// type GetPendingRequestsForTeamOutput struct {
// 	Body []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow
// }

// func (h *handler) handleGetPendingRequestsForTeam(ctx context.Context, input *struct {
// 	TeamId string `path:"teamId"`
// }) (*GetPendingRequestsForTeamOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	teamId, err := uuid.Parse(input.TeamId)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid team id")
// 	}

// 	requests, err := h.teamService.GetPendingJoinRequestForTeam(ctx, userCtx.UserID, teamId)

// 	if err != nil {
// 		if errors.Is(err, ErrUserNotTeamOwner) {
// 			return nil, huma.Error403Forbidden("Not authorized to perform this action")
// 		}

// 		return nil, huma.Error500InternalServerError("Fail to get pending requests for team")
// 	}

// 	return &GetPendingRequestsForTeamOutput{Body: requests}, nil
// }

// type GetMyPendingRequestsOutput struct {
// 	Body []sqlc.TeamJoinRequest
// }

// func (h *handler) handleGetMyPendingRequests(ctx context.Context, input *struct{}) (*GetMyPendingRequestsOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}
// 	requests, err := h.teamService.GetUserPendingJoinRequests(ctx, userCtx.UserID)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Fail to get user's pending requests")
// 	}

// 	return &GetMyPendingRequestsOutput{Body: requests}, nil
// }

// type AcceptTeamJoinRequestOutput struct {
// 	Status int
// }

// func (h *handler) handleAcceptTeamJoinRequest(ctx context.Context, input *struct {
// 	RequestId string `path:"requestId"`
// }) (*AcceptTeamJoinRequestOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	requestId, err := uuid.Parse(input.RequestId)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid join request id")
// 	}

// 	err = h.teamService.RespondToJoinRequest(ctx, userCtx.UserID, requestId, true)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Fail to accept team join request")
// 	}

// 	return &AcceptTeamJoinRequestOutput{Status: http.StatusOK}, nil
// }

// type RejectTeamJoinRequestOutput struct {
// 	Status int
// }

// func (h *handler) handleRejectTeamJoinRequest(ctx context.Context, input *struct {
// 	RequestId string `path:"requestId"`
// }) (*RejectTeamJoinRequestOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	requestId, err := uuid.Parse(input.RequestId)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid join request id")
// 	}

// 	err = h.teamService.RespondToJoinRequest(ctx, userCtx.UserID, requestId, false)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Fail to accept team join request")
// 	}

// 	return &RejectTeamJoinRequestOutput{Status: http.StatusOK}, nil
// }
