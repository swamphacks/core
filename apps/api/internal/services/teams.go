package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrTeamExists                 = errors.New("team already exists")
	ErrTeamNotFound               = errors.New("team does not exist")
	ErrNoEligibleNextOwner        = errors.New("no eligible next owner for team")
	ErrTeamMembersNotBytes        = errors.New("team members aren't bytes")
	ErrUserOnTeam                 = errors.New("user already on a team")
	ErrUserNotTeamOwner           = errors.New("user is not the team owner")
	ErrTeamFull                   = errors.New("team is full")
	ErrUserNotApplicantOrAttendee = errors.New("user is not an applicant or attendee for the event")
	ErrInvitationNotFound         = errors.New("invitation not found")
	ErrInvitationExpired          = errors.New("invitation expired")
	ErrInvitationAlreadyAccepted  = errors.New("invitation already accepted")
	ErrEmailMismatch              = errors.New("email mismatch")
	ErrApplicationRequired        = errors.New("application required")
	ErrInvitationAlreadyExists    = errors.New("pending invitation already exists for this email and team")
	ErrKickOwnerSelf              = errors.New("team owner cannot kick themselves")
)

type TeamService struct {
	teamRepo            *repository.TeamRepository
	teamMemberRepo      *repository.TeamMemberRepository
	teamJoinRequestRepo *repository.TeamJoinRequestRepository
	teamInvitationRepo  *repository.TeamInvitationRepository
	userRepo            *repository.UserRepository
	applicationRepo     *repository.ApplicationRepository
	eventRepo           *repository.EventRepository
	emailService        *EmailService
	txm                 *db.TransactionManager
	cfg                 *config.Config
	logger              zerolog.Logger
}

func NewTeamService(
	teamRepo *repository.TeamRepository,
	teamMemberRepo *repository.TeamMemberRepository,
	teamJoinRequestRepo *repository.TeamJoinRequestRepository,
	teamInvitationRepo *repository.TeamInvitationRepository,
	userRepo *repository.UserRepository,
	applicationRepo *repository.ApplicationRepository,
	eventRepo *repository.EventRepository,
	emailService *EmailService,
	txm *db.TransactionManager,
	cfg *config.Config,
	logger zerolog.Logger) *TeamService {
	return &TeamService{
		teamRepo:            teamRepo,
		teamMemberRepo:      teamMemberRepo,
		teamJoinRequestRepo: teamJoinRequestRepo,
		teamInvitationRepo:  teamInvitationRepo,
		userRepo:            userRepo,
		applicationRepo:     applicationRepo,
		eventRepo:           eventRepo,
		emailService:        emailService,
		txm:                 txm,
		cfg:                 cfg,
		logger:              logger.With().Str("service", "TeamService").Str("component", "team").Logger(),
	}
}

// TODO: Remove all join requests after creating a team
func (s *TeamService) CreateTeam(ctx context.Context, name string, eventId, userId uuid.UUID) (*sqlc.Team, error) {
	// Check if user already has a team for this event.
	member, err := s.teamMemberRepo.GetTeamMemberByUserAndEvent(ctx, userId, eventId)
	if err == nil && member != nil {
		// User already has a team
		return nil, ErrTeamExists
	}
	if err != nil && !errors.Is(err, repository.ErrTeamMemberNotFound) {
		return nil, err
	}

	var newTeam sqlc.Team

	// Transactionally create a new team and assign the user as the owner.
	if err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txTeamRepo := s.teamRepo.NewTx(tx)
		txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)
		txTeamJoinRequestRepo := s.teamJoinRequestRepo.NewTx(tx)

		// Delete any pending join requests by the user for this event
		if err := txTeamJoinRequestRepo.DeleteByUserAndEventAndStatus(ctx, userId, eventId, sqlc.JoinRequestStatusPENDING); err != nil {
			return err
		}

		team, err := txTeamRepo.Create(ctx, name, userId, eventId)
		if err != nil {
			return err
		}

		if _, err = txTeamMemberRepo.Create(ctx, team.ID, userId); err != nil {
			return err
		}

		newTeam = *team
		return nil
	}); err != nil {
		return nil, err
	}

	return &newTeam, nil
}

type MemberWithUserInfo struct {
	UserID   uuid.UUID  `json:"user_id"`
	Email    *string    `json:"email"`
	Image    *string    `json:"image"`
	Name     string     `json:"name"`
	JoinedAt *time.Time `json:"joined_at"`
}

type TeamWithMembers struct {
	ID      uuid.UUID            `json:"id"`
	EventId *uuid.UUID           `json:"event_id"`
	OwnerId *uuid.UUID           `json:"owner_id"`
	Name    string               `json:"name"`
	Members []MemberWithUserInfo `json:"members"`
}

func (s *TeamService) GetUserTeamWithMembers(ctx context.Context, userId, eventId uuid.UUID) (*TeamWithMembers, error) {
	team, err := s.teamRepo.GetTeamByMemberAndEvent(ctx, userId, eventId)
	if err != nil {
		// If no team, just return nil
		if errors.Is(err, repository.ErrTeamNotFound) {
			return nil, nil
		}
		return nil, err
	}

	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	var parsedMembers []MemberWithUserInfo
	for _, member := range members {
		parsedMembers = append(parsedMembers, MemberWithUserInfo{
			UserID:   member.UserID,
			Email:    member.Email,
			Image:    member.Image,
			Name:     member.Name,
			JoinedAt: member.JoinedAt,
		})
	}

	teamWithMembers := TeamWithMembers{
		ID:      team.ID,
		EventId: team.EventID,
		OwnerId: team.OwnerID,
		Name:    team.Name,
		Members: parsedMembers,
	}

	return &teamWithMembers, nil
}

func (s *TeamService) GetTeamsWithMembersByEvent(ctx context.Context, eventId uuid.UUID, limit, offset int32) ([]TeamWithMembers, error) {
	teams, err := s.teamRepo.GetTeamsWithMembersByEvent(ctx, eventId, limit, offset)
	if err != nil {
		return []TeamWithMembers{}, err
	}

	var result []TeamWithMembers

	for _, team := range teams {
		var parsedMembers []MemberWithUserInfo

		if err := json.Unmarshal(team.Members, &parsedMembers); err != nil {
			return []TeamWithMembers{}, err
		}

		result = append(result, TeamWithMembers{
			ID:      team.ID,
			EventId: team.EventID,
			OwnerId: team.OwnerID,
			Name:    team.Name,
			Members: parsedMembers,
		})
	}

	return result, nil
}

func (s *TeamService) GetTeamWithMembers(ctx context.Context, teamId uuid.UUID) (*TeamWithMembers, error) {
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		// If no team, just return nil
		if errors.Is(err, repository.ErrTeamNotFound) {
			return nil, nil
		}
		return nil, err
	}

	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	var parsedMember []MemberWithUserInfo
	for _, member := range members {
		parsedMember = append(parsedMember, MemberWithUserInfo{
			UserID:   member.UserID,
			Email:    member.Email,
			Image:    member.Image,
			Name:     member.Name,
			JoinedAt: member.JoinedAt,
		})
	}

	teamWithMembers := TeamWithMembers{
		ID:      team.ID,
		EventId: team.EventID,
		OwnerId: team.OwnerID,
		Name:    team.Name,
		Members: parsedMember,
	}

	return &teamWithMembers, nil
}

func (s *TeamService) JoinTeam(ctx context.Context, userId, teamId uuid.UUID) error {
	_, err := s.teamMemberRepo.Create(ctx, teamId, userId)
	return err
}

func (s *TeamService) InviteUserToTeam(ctx context.Context, teamId, inviterId uuid.UUID, inviteeEmail string) error {
	// 1. Validate team exists and inviter is team owner/leader
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	// Check if inviter is the team owner
	if team.OwnerID == nil || *team.OwnerID != inviterId {
		return ErrUserNotTeamOwner
	}

	// 2. Check for existing pending invitation for this email/team (prevent duplicates)
	_, err = s.teamInvitationRepo.GetByEmailAndTeam(ctx, inviteeEmail, teamId)
	if err == nil {
		// Invitation already exists
		return ErrInvitationAlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrInvitationNotFound) {
		return err
	}

	// 3. Check if user exists by email
	var invitedUserId *uuid.UUID
	user, err := s.userRepo.GetByEmail(ctx, inviteeEmail)
	if err == nil {
		invitedUserId = &user.ID
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	// 4. Create invitation record with expiration (7 days from now)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	invitation, err := s.teamInvitationRepo.Create(ctx, teamId, inviterId, inviteeEmail, &expiresAt)
	if err != nil {
		return err
	}

	// If user exists, update invitation with user ID
	if invitedUserId != nil {
		_, err = s.teamInvitationRepo.Update(ctx, invitation.ID, invitedUserId, nil)
		if err != nil {
			s.logger.Err(err).Msg("Failed to update invitation with user ID")
			// Non-fatal, continue with email sending
		}
	}

	// 5. Get team details (name, event) and inviter name for email
	event, err := s.eventRepo.GetEventByID(ctx, *team.EventID)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get event for invitation email")
		return err
	}

	inviter, err := s.userRepo.GetByID(ctx, inviterId)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get inviter for invitation email")
		return err
	}

	// 6. Generate invitation link
	baseURL := s.cfg.ClientUrl
	if baseURL == "" {
		baseURL = "https://app.swamphacks.com" // Fallback to production URL
	}
	inviteLink := fmt.Sprintf("%s/teams/invite/%s", baseURL, invitation.ID.String())

	// 7. Call emailService.QueueSendTeamInvitation with dynamic link
	_, err = s.emailService.QueueSendTeamInvitation(
		inviteeEmail,
		team.Name,
		inviter.Name,
		event.Name,
		inviteLink,
	)
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue team invitation email")
		return err
	}

	return nil
}

type InvitationDetails struct {
	ID              uuid.UUID          `json:"id"`
	TeamName        string             `json:"team_name"`
	InviterName     string             `json:"inviter_name"`
	EventName       string             `json:"event_name"`
	EventID         uuid.UUID          `json:"event_id"`
	InvitedEmail    string             `json:"invited_email"`
	Status          string             `json:"status"`
	ExpiresAt       *time.Time         `json:"expires_at"`
	CreatedAt       time.Time          `json:"created_at"`
	TeamMembers     []MemberWithUserInfo `json:"team_members"`
}

func (s *TeamService) GetInvitationDetails(ctx context.Context, invitationId uuid.UUID) (*InvitationDetails, error) {
	invitation, err := s.teamInvitationRepo.GetByID(ctx, invitationId)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	// Check if invitation is expired
	if invitation.ExpiresAt != nil && invitation.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvitationExpired
	}

	// Check if already accepted
	if invitation.Status == sqlc.InvitationStatusACCEPTED {
		return nil, ErrInvitationAlreadyAccepted
	}

	// Get team details
	team, err := s.teamRepo.GetByID(ctx, invitation.TeamID)
	if err != nil {
		return nil, err
	}

	// Get event details
	if team.EventID == nil {
		return nil, ErrTeamNotFound
	}
	event, err := s.eventRepo.GetEventByID(ctx, *team.EventID)
	if err != nil {
		return nil, err
	}

	// Get inviter details
	inviter, err := s.userRepo.GetByID(ctx, invitation.InvitedByUserID)
	if err != nil {
		return nil, err
	}

	// Get team members with their profile pictures
	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	var parsedMembers []MemberWithUserInfo
	for _, member := range members {
		parsedMembers = append(parsedMembers, MemberWithUserInfo{
			UserID:   member.UserID,
			Email:    member.Email,
			Image:    member.Image,
			Name:     member.Name,
			JoinedAt: member.JoinedAt,
		})
	}

	return &InvitationDetails{
		ID:           invitation.ID,
		TeamName:     team.Name,
		InviterName:  inviter.Name,
		EventName:    event.Name,
		EventID:      event.ID,
		InvitedEmail: invitation.InvitedEmail,
		Status:       string(invitation.Status),
		ExpiresAt:    invitation.ExpiresAt,
		CreatedAt:    invitation.CreatedAt,
		TeamMembers:  parsedMembers,
	}, nil
}

// LinkUserToInvitation links a user ID to an invitation if the invitation is still valid.
// This should always happen whether the user has access to the event or not.
// Returns true if the invitation is expired or already claimed.
func (s *TeamService) LinkUserToInvitation(ctx context.Context, invitationId, userId uuid.UUID) (bool, error) {
	// Get invitation
	invitation, err := s.teamInvitationRepo.GetByID(ctx, invitationId)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return true, ErrInvitationNotFound
		}
		return true, err
	}

	// Check if already accepted or rejected (this takes precedence over expiration check)
	if invitation.Status != sqlc.InvitationStatusPENDING {
		return true, nil // Already claimed (accepted/rejected)
	}

	// Check if invitation is expired
	if invitation.ExpiresAt != nil {
		now := time.Now()
		if invitation.ExpiresAt.Before(now) {
			s.logger.Debug().
				Str("invitation_id", invitationId.String()).
				Time("expires_at", *invitation.ExpiresAt).
				Time("now", now).
				Msg("Invitation expired")
			return true, ErrInvitationExpired
		}
	}

	// Link user ID to invitation if not already linked
	if invitation.InvitedUserID == nil {
		_, err = s.teamInvitationRepo.Update(ctx, invitationId, &userId, nil)
		if err != nil {
			s.logger.Err(err).
				Str("invitation_id", invitationId.String()).
				Str("user_id", userId.String()).
				Msg("Failed to update invitation with user ID")
			return false, err
		}
	}

	return false, nil // Not expired/claimed
}

func (s *TeamService) AcceptInvitation(ctx context.Context, invitationId, userId uuid.UUID) error {
	// Get invitation
	invitation, err := s.teamInvitationRepo.GetByID(ctx, invitationId)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return ErrInvitationNotFound
		}
		return err
	}

	// Validate invitation (pending, not expired)
	if invitation.Status != sqlc.InvitationStatusPENDING {
		if invitation.Status == sqlc.InvitationStatusACCEPTED {
			return ErrInvitationAlreadyAccepted
		}
		return ErrInvitationNotFound
	}

	if invitation.ExpiresAt != nil && invitation.ExpiresAt.Before(time.Now()) {
		return ErrInvitationExpired
	}

	// Get user to verify email
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	// Verify authenticated user's email matches invitation email
	if user.Email == nil || *user.Email != invitation.InvitedEmail {
		return ErrEmailMismatch
	}

	// Get team's event ID from team record
	team, err := s.teamRepo.GetByID(ctx, invitation.TeamID)
	if err != nil {
		return err
	}

	if team.EventID == nil {
		return ErrTeamNotFound
	}

	// Check if user has submitted application for event
	application, err := s.applicationRepo.GetApplicationByUserAndEventID(ctx, sqlc.GetApplicationByUserAndEventIDParams{
		UserID:  userId,
		EventID: *team.EventID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrApplicationNotFound) {
			return ErrApplicationRequired
		}
		return err
	}

	// If application exists, check if it's submitted
	if application.Status.ApplicationStatus != sqlc.ApplicationStatusSubmitted {
		return ErrApplicationRequired
	}

	members, err := s.teamMemberRepo.GetTeamMembers(ctx, invitation.TeamID)
	if err != nil {
		return err
	}

	if len(members) >= 4 {
		return ErrTeamFull
	}

	// Add user to team and mark invitation as accepted
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)
		txTeamInvitationRepo := s.teamInvitationRepo.NewTx(tx)

		// Add user to team
		if _, err := txTeamMemberRepo.Create(ctx, invitation.TeamID, userId); err != nil {
			return err
		}

		// Mark invitation as accepted
		_, err := txTeamInvitationRepo.AcceptInvitation(ctx, invitationId, userId)
		return err
	})
}

func (s *TeamService) RejectInvitation(ctx context.Context, invitationId, userId uuid.UUID) error {
	// Get invitation
	invitation, err := s.teamInvitationRepo.GetByID(ctx, invitationId)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return ErrInvitationNotFound
		}
		return err
	}

	// Get user to verify email
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	// Verify authenticated user's email matches invitation email
	if user.Email == nil || *user.Email != invitation.InvitedEmail {
		return ErrEmailMismatch
	}

	// Reject invitation
	_, err = s.teamInvitationRepo.RejectInvitation(ctx, invitationId)
	return err
}

func (s *TeamService) LeaveTeam(ctx context.Context, userId, teamId uuid.UUID) error {
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	// Check if user is NOT the owner
	if team.OwnerID == nil || *team.OwnerID != userId {
		return s.teamMemberRepo.Delete(ctx, team.ID, userId)
	}

	// User IS the owner
	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return err
	}

	// Confirm that team has members. This should not run unless there are inconsistencies.
	if len(members) == 0 {
		s.logger.Warn().Msg("Team has no members but owner exists — inconsistent state.")
		return s.teamRepo.Delete(ctx, team.ID)
	}

	// If current user is the last member → delete both team + membership
	if len(members) == 1 {
		return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
			txTeamRepo := s.teamRepo.NewTx(tx)
			txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

			if err := txTeamMemberRepo.Delete(ctx, team.ID, userId); err != nil {
				return err
			}
			return txTeamRepo.Delete(ctx, team.ID)
		})
	}

	// Choose the next owner deterministically
	var nextOwner sqlc.GetTeamMembersRow
	for _, m := range members {
		if m.UserID != userId {
			nextOwner = m
			break
		}
	}

	if nextOwner.UserID == uuid.Nil {
		// Should not happen unless data is inconsistent
		s.logger.Error().Msg("No eligible next owner found.")
		return ErrNoEligibleNextOwner
	}

	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txTeamRepo := s.teamRepo.NewTx(tx)
		txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

		if err := txTeamMemberRepo.Delete(ctx, team.ID, userId); err != nil {
			return err
		}

		_, err := txTeamRepo.Update(ctx, team.ID, nil, &nextOwner.UserID)
		return err
	})
}

func (s *TeamService) RequestToJoinTeam(ctx context.Context, eventId, teamId, userId uuid.UUID, message *string) (*sqlc.TeamJoinRequest, error) {
	// Ensure user is not already on a team (User's can't request to join a team if they are already on a team)
	_, err := s.teamMemberRepo.GetTeamMemberByUserAndEvent(ctx, userId, eventId)
	if err != nil && !errors.Is(err, repository.ErrTeamMemberNotFound) {
		return nil, err
	}

	if err == nil {
		return nil, ErrUserOnTeam
	}

	request, err := s.teamJoinRequestRepo.Create(ctx, teamId, userId, message)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *TeamService) GetPendingJoinRequestForTeam(ctx context.Context, userId, teamId uuid.UUID) ([]sqlc.ListJoinRequestsByTeamAndStatusWithUserRow, error) {
	// Get team and check if user is the owner of the tema
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		return []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow{}, err
	}

	if *team.OwnerID != userId {
		return []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow{}, ErrUserNotTeamOwner
	}

	requests, err := s.teamJoinRequestRepo.ListJoinRequestsByTeamWithUser(ctx, teamId, sqlc.JoinRequestStatusPENDING)
	if err != nil {
		return []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow{}, err
	}

	return requests, nil
}

func (s *TeamService) GetUserPendingJoinRequestsByEvent(ctx context.Context, userId, eventId uuid.UUID) ([]sqlc.TeamJoinRequest, error) {
	requests, err := s.teamJoinRequestRepo.ListJoinRequestsByUserAndEvent(ctx, userId, eventId, sqlc.JoinRequestStatusPENDING)
	if err != nil {
		return []sqlc.TeamJoinRequest{}, err
	}

	return requests, nil
}

func (s *TeamService) RespondToJoinRequest(ctx context.Context, ownerId, requestId uuid.UUID, accept bool) error {
	// Retrieve the join request
	oldRequest, err := s.teamJoinRequestRepo.GetById(ctx, requestId)
	if err != nil {
		return err
	}

	// Get team and check if user is the owner of the team
	team, err := s.teamRepo.GetByID(ctx, oldRequest.TeamID)
	if err != nil {
		return err
	}
	if *team.OwnerID != ownerId {
		return ErrUserNotTeamOwner
	}

	// Also ensure user is actually an applicant or attendee
	role, err := s.eventRepo.GetEventRoleByIds(ctx, oldRequest.UserID, *team.EventID)
	if err != nil {
		if errors.Is(err, repository.ErrEventRoleNotFound) {
			return ErrUserNotApplicantOrAttendee
		}

		return err
	}
	if role.Role != sqlc.EventRoleTypeAttendee && role.Role != sqlc.EventRoleTypeApplicant {
		return ErrUserNotApplicantOrAttendee
	}

	// Ensure user is not already on a team
	_, err = s.teamMemberRepo.GetTeamMemberByUserAndEvent(ctx, oldRequest.UserID, *team.EventID)
	if err != nil && !errors.Is(err, repository.ErrTeamMemberNotFound) {
		return err
	}

	if err == nil {
		return ErrUserOnTeam
	}

	if accept {
		// Check if team is already full
		members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
		if err != nil {
			return err
		}

		if len(members) >= 4 {
			return ErrTeamFull
		}

		// Accepting the request: add user to team and update request status
		return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
			txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)
			txTeamJoinRequestRepo := s.teamJoinRequestRepo.NewTx(tx)

			if _, err := txTeamMemberRepo.Create(ctx, oldRequest.TeamID, oldRequest.UserID); err != nil {
				return err
			}

			if _, err := txTeamJoinRequestRepo.UpdateStatus(ctx, requestId, sqlc.JoinRequestStatusAPPROVED); err != nil {
				return err
			}

			// Delete all other pending requests by the user for this event
			return txTeamJoinRequestRepo.DeleteByUserAndEventAndStatus(ctx, oldRequest.UserID, *team.EventID, sqlc.JoinRequestStatusPENDING)
		})
	} else {
		// Rejecting the request: just update request status
		_, err := s.teamJoinRequestRepo.UpdateStatus(ctx, requestId, sqlc.JoinRequestStatusREJECTED)
		return err
	}
}

func (s *TeamService) KickMemberFromTeam(ctx context.Context, memberId, teamId, userId uuid.UUID) error {
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		return err
	}

	if *team.OwnerID != userId {
		s.logger.Warn().Str("team_owner", team.OwnerID.String()).Str("user", userId.String()).Msg("User attempting to kick member is not the team owner")
		return ErrUserNotTeamOwner
	}

	// Prevent owner from kicking themselves
	if memberId == userId {
		s.logger.Warn().Str("team_owner", team.OwnerID.String()).Msg("Team owner attempting to kick themselves")
		return ErrKickOwnerSelf
	}

	return s.teamMemberRepo.Delete(ctx, team.ID, memberId)
}
