package teams

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCreateTeam          = errors.New("unable to create team")
	ErrGetTeamDetails      = errors.New("unable to get team details")
	ErrGetTeamMembers      = errors.New("unable to get team membes")
	ErrGetTeam             = errors.New("unable to get team")
	ErrNotInTeam           = errors.New("user is not part of a team")
	ErrUserNotTeamOwner    = errors.New("user is not the team owner")
	ErrDeleteTeam          = errors.New("unable to delete team")
	ErrKickSelf            = errors.New("cannot kick yourself")
	ErrKickError           = errors.New("unable to kick user")
	ErrJoinTeam            = errors.New("unable to join team")
	ErrJoinSameTeam        = errors.New("cannot join the same team again")
	ErrLeaveTeam           = errors.New("unable to leave team")
	ErrOwnerLeaveTeam      = errors.New("owner cannot leave team")
	ErrCreateInvitation    = errors.New("unable to create invitation")
	ErrGetInvitation       = errors.New("unable to get invitation")
	ErrMembersLimitReached = errors.New("members limit exceeded")
	ErrNoTeamFound         = errors.New("no team found")
	ErrAlreadyHasTeam      = errors.New("user is already in a team")
)

type TeamDto struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
}

type TeamMemberDto struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image *string   `json:"image"`
	// JoinedAt time.Time `json:"joinedAt"`
}

type TeamDetailsDto struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"ownerId"`
	// CreatedAt time.Time       `json:"createdAt"`
	Members []TeamMemberDto `json:"members"`
}

type CreateTeamRequestDto struct {
	Name string `json:"name"`
}

type DeleteTeamRequestDto struct {
	TeamID uuid.UUID `json:"teamId"`
}

type KickMemberRequestDto struct {
	MemberId uuid.UUID `json:"memberId"`
}
