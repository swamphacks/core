package teams

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
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
		OperationID:   "get-team-by-invite-id",
		Method:        http.MethodGet,
		Summary:       "Get Team By Invitation Id",
		Description:   "Returns the team information by invitation id",
		Tags:          []string{"Team"},
		Path:          "/invitation/{inviteId}/team",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetTeamByInvitationId)

	huma.Register(group, huma.Operation{
		OperationID:   "get-team-details",
		Method:        http.MethodGet,
		Summary:       "Get Team Details",
		Description:   "Returns the team information and the full list of team members by team id",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/details",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetTeamDetails)

	huma.Register(group, huma.Operation{
		OperationID:   "get-team-members",
		Method:        http.MethodGet,
		Summary:       "Get Team Members",
		Description:   "Returns the list of members of a team",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/members",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetTeamMembers)

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
		OperationID:   "delete-team",
		Method:        http.MethodDelete,
		Summary:       "Delete Team",
		Description:   "Delete a team. The user must be the owner of the team.",
		Tags:          []string{"Team"},
		Path:          "",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusConflict, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleDeleteTeam)

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
		OperationID:   "kick-member",
		Method:        http.MethodPost,
		Summary:       "Kick Team Member",
		Description:   "Kicks a member from a team. Only the team owner can perform this action.",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/kick",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleKickMember)

	huma.Register(group, huma.Operation{
		OperationID:   "join-team",
		Method:        http.MethodPost,
		Summary:       "Join Team",
		Description:   "Join a team through an invitation link.",
		Tags:          []string{"Team"},
		Path:          "/join/{id}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleJoinTeam)

	huma.Register(group, huma.Operation{
		OperationID:   "create-invitation",
		Method:        http.MethodPost,
		Summary:       "Create Invitation",
		Description:   "Create an invitation for other users to join the team",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/invitation",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleCreateInvitation)

	huma.Register(group, huma.Operation{
		OperationID:   "get-invitation",
		Method:        http.MethodGet,
		Summary:       "Get Invitation",
		Description:   "Get an invitation id, which can be used to construct a team join link",
		Tags:          []string{"Team"},
		Path:          "/{teamId}/invitation",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, teamHandler.handleGetInvitation)
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
