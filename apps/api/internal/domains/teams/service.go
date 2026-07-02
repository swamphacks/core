package teams

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type TeamService struct {
	db     *database.DB
	txm    *database.TransactionManager
	logger zerolog.Logger
}

func NewService(
	db *database.DB,
	txm *database.TransactionManager,
	logger zerolog.Logger) *TeamService {
	return &TeamService{
		db:     db,
		txm:    txm,
		logger: logger.With().Str("service", "TeamService").Str("component", "team").Logger(),
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string, userID uuid.UUID) (*sqlc.Team, error) {
	var newTeam sqlc.Team

	// Transactionally create a new team and assign the user as the owner.
	if err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		team, err := txDB.Query.CreateTeam(ctx, sqlc.CreateTeamParams{
			Name:    name,
			OwnerID: userID,
		})

		if err != nil {
			return err
		}

		if _, err = txDB.Query.AddUserToTeam(ctx, sqlc.AddUserToTeamParams{
			TeamID: team.ID,
			UserID: userID,
		}); err != nil {
			return err
		}

		_, err = txDB.Query.CreateInvitation(ctx, sqlc.CreateInvitationParams{
			TeamID:    team.ID,
			InviterID: userID,
			ExpiresAt: nil,
		})

		if err != nil {
			return err
		}

		newTeam = team
		return nil
	}); err != nil {
		s.logger.Err(err).Msg("CreateTeam fail")
		return nil, ErrCreateTeam
	}

	return &newTeam, nil
}

func (s *TeamService) GetTeamByUserId(ctx context.Context, userID uuid.UUID) (*sqlc.GetTeamByUserIdRow, error) {
	team, err := s.db.Query.GetTeamByUserId(ctx, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotInTeam
		}
		s.logger.Err(err).Msg("GetTeamByUserId fail")
		return nil, ErrGetTeam
	}

	return &team, nil
}

func (s *TeamService) GetTeamByInvitationId(ctx context.Context, inviteId uuid.UUID) (*sqlc.Team, error) {
	team, err := s.db.Query.GetTeamByInvitationId(ctx, inviteId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoTeamFound
		}
		s.logger.Err(err).Msg("fail to get team by invitation id")
		return nil, ErrGetTeam
	}

	return &team, nil
}

func (s *TeamService) GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]sqlc.GetTeamMembersRow, error) {
	members, err := s.db.Query.GetTeamMembers(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("fail to get team members")
		return nil, ErrGetTeamMembers
	}

	return members, nil
}

func (s *TeamService) GetTeamDetails(ctx context.Context, teamID uuid.UUID) (*sqlc.GetTeamDetailsRow, error) {
	teamDetails, err := s.db.Query.GetTeamDetails(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("GetTeam fail")
		return nil, ErrGetTeamDetails
	}

	return &teamDetails, nil
}

func (s *TeamService) JoinTeam(ctx context.Context, userID, teamID uuid.UUID) error {
	members, err := s.GetTeamMembers(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("fail to get team members")
		return ErrJoinTeam
	}

	if len(members) >= 4 {
		return ErrMembersLimitReached
	}

	team, err := s.db.Query.GetTeamByUserId(ctx, userID)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// User is not on a team, continue joining.

	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		s.logger.Err(err).Msg("JoinTeam failed: GetTeamByUserId")
		return ErrJoinTeam

	case team.ID == teamID:
		return ErrJoinSameTeam

	default:
		return ErrAlreadyHasTeam
	}

	_, err = s.db.Query.AddUserToTeam(ctx, sqlc.AddUserToTeamParams{
		TeamID: teamID,
		UserID: userID,
	})

	if err != nil {
		s.logger.Err(err).Msg("Join fail")
		return ErrJoinTeam
	}

	return nil
}

func (s *TeamService) KickMember(ctx context.Context, memberID, teamID, userID uuid.UUID) error {
	team, err := s.db.Query.GetTeamById(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("KickMember fail, unable to get team info by id")
		return errors.New("Unable to kick member")
	}

	if team.OwnerID != userID {
		return ErrUserNotTeamOwner
	}

	if memberID == userID {
		return ErrKickSelf
	}

	err = s.db.Query.RemoveUserFromTeam(ctx, sqlc.RemoveUserFromTeamParams{
		TeamID: teamID,
		UserID: memberID,
	})

	if err != nil {
		s.logger.Err(err).Msg("KickMember fail, unable to remove user from team")
		return ErrKickError
	}

	return nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, userID, teamID uuid.UUID) error {
	team, err := s.db.Query.GetTeamById(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("DeleteTeam fail, unable to get team info by id")
		return ErrDeleteTeam
	}

	if team.OwnerID != userID {
		return ErrUserNotTeamOwner
	}

	err = s.db.Query.DeleteTeamById(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("DeleteTeam fail")
		return ErrDeleteTeam
	}

	return nil
}

func (s *TeamService) LeaveTeam(ctx context.Context, userID, teamID uuid.UUID) error {
	team, err := s.db.Query.GetTeamById(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("LeaveTeam fail, unable to get team info by id")
		return ErrLeaveTeam
	}

	if team.OwnerID == userID {
		return ErrOwnerLeaveTeam
	}

	err = s.db.Query.RemoveUserFromTeam(ctx, sqlc.RemoveUserFromTeamParams{
		TeamID: teamID,
		UserID: userID,
	})

	if err != nil {
		s.logger.Err(err).Msg("LeaveTeam fail")
		return ErrLeaveTeam
	}

	return nil
}

func (s *TeamService) CreateInvitation(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	team, err := s.GetTeamByUserId(ctx, userID)

	if err != nil {
		if errors.Is(err, ErrNotInTeam) {
			return nil, ErrNotInTeam
		}
		s.logger.Err(err).Msg("fail to create invitation")
		return nil, ErrGetTeam
	}

	if team.OwnerID != userID {
		return nil, ErrUserNotTeamOwner
	}

	invitation, err := s.db.Query.CreateInvitation(ctx, sqlc.CreateInvitationParams{
		TeamID:    team.ID,
		InviterID: userID,
		ExpiresAt: nil,
	})

	if err != nil {
		s.logger.Err(err).Msg("fail to create invitation")
		return nil, ErrCreateInvitation
	}

	return &invitation.ID, nil
}

func (s *TeamService) GetInvitation(ctx context.Context, id uuid.UUID) (*sqlc.TeamInvitation, error) {
	invitation, err := s.db.Query.GetInvitationByID(ctx, id)

	if err != nil {
		s.logger.Err(err).Msg("fail to get invitation")
		return nil, ErrGetInvitation
	}

	return &invitation, nil
}

func (s *TeamService) GetInvitationByTeamID(ctx context.Context, teamID, userID uuid.UUID) (*sqlc.TeamInvitation, error) {
	team, err := s.GetTeamByUserId(ctx, userID)

	if err != nil {
		s.logger.Err(err).Msg("fail to get team by user id")
		return nil, ErrGetInvitation
	}

	if team.OwnerID != userID {
		return nil, ErrUserNotTeamOwner
	}

	invitation, err := s.db.Query.GetInvitationByTeamID(ctx, teamID)

	if err != nil {
		s.logger.Err(err).Msg("fail to get invitation by team id")
		return nil, ErrGetInvitation
	}

	return &invitation, nil
}

// type MemberWithUserInfo struct {
// 	UserID   uuid.UUID `json:"userID"`
// 	Email    *string   `json:"email"`
// 	Image    *string   `json:"image"`
// 	Name     string    `json:"name"`
// 	JoinedAt time.Time `json:"joinedAt"`
// }

// type TeamWithMembers struct {
// 	ID      uuid.UUID            `json:"id"`
// 	OwnerId *uuid.UUID           `json:"ownerId"`
// 	Name    string               `json:"name"`
// 	Members []MemberWithUserInfo `json:"members"`
// }

// func (s *TeamService) GetUserTeamWithMembers(ctx context.Context, userID uuid.UUID) (*TeamWithMembers, error) {
// 	team, err := s.teamRepo.GetTeamByMember(ctx, userID)
// 	if err != nil {
// 		// If no team, just return nil
// 		if errors.Is(err, repository.ErrTeamNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var parsedMembers []MemberWithUserInfo
// 	for _, member := range members {
// 		parsedMembers = append(parsedMembers, MemberWithUserInfo{
// 			UserID:   member.UserID,
// 			Email:    member.Email,
// 			Image:    member.Image,
// 			Name:     member.Name,
// 			JoinedAt: member.JoinedAt,
// 		})
// 	}

// 	teamWithMembers := TeamWithMembers{
// 		ID:      team.ID,
// 		OwnerId: team.OwnerID,
// 		Name:    team.Name,
// 		Members: parsedMembers,
// 	}

// 	return &teamWithMembers, nil
// }

// func (s *TeamService) GetTeamsWithMembers(ctx context.Context, limit, offset int32) ([]TeamWithMembers, error) {
// 	teams, err := s.teamRepo.GetTeamsWithMembers(ctx, sqlc.ListTeamsWithMembersParams{
// 		Limit:  limit,
// 		Offset: offset,
// 	})

// 	if err != nil {
// 		return []TeamWithMembers{}, err
// 	}

// 	var result []TeamWithMembers

// 	for _, team := range teams {
// 		var parsedMembers []MemberWithUserInfo

// 		if err := json.Unmarshal(team.Members, &parsedMembers); err != nil {
// 			return []TeamWithMembers{}, err
// 		}

// 		result = append(result, TeamWithMembers{
// 			ID:      team.ID,
// 			OwnerId: team.OwnerID,
// 			Name:    team.Name,
// 			Members: parsedMembers,
// 		})
// 	}

// 	return result, nil
// }

// func (s *TeamService) GetTeamWithMembersByTeamId(ctx context.Context, teamID uuid.UUID) (*TeamWithMembers, error) {
// 	team, err := s.teamRepo.GetByID(ctx, teamID)
// 	if err != nil {
// 		// If no team, just return nil
// 		if errors.Is(err, repository.ErrTeamNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var parsedMember []MemberWithUserInfo
// 	for _, member := range members {
// 		parsedMember = append(parsedMember, MemberWithUserInfo{
// 			UserID:   member.UserID,
// 			Email:    member.Email,
// 			Image:    member.Image,
// 			Name:     member.Name,
// 			JoinedAt: member.JoinedAt,
// 		})
// 	}

// 	teamWithMembers := TeamWithMembers{
// 		ID:      team.ID,
// 		OwnerId: team.OwnerID,
// 		Name:    team.Name,
// 		Members: parsedMember,
// 	}

// 	return &teamWithMembers, nil
// }

// func (s *TeamService) LeaveTeam(ctx context.Context, userID, teamID uuid.UUID) error {
// 	team, err := s.teamRepo.GetByID(ctx, teamID)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrTeamNotFound) {
// 			return ErrTeamNotFound
// 		}
// 		return err
// 	}

// 	// Check if user is NOT the owner
// 	if team.OwnerID == nil || *team.OwnerID != userID {
// 		return s.teamMemberRepo.Delete(ctx, sqlc.RemoveTeamMemberParams{
// 			TeamID: teamID,
// 			UserID: userID,
// 		})
// 	}

// 	// User IS the owner
// 	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
// 	if err != nil {
// 		return err
// 	}

// 	// Confirm that team has members. This should not run unless there are inconsistencies.
// 	if len(members) == 0 {
// 		s.logger.Warn().Msg("Team has no members but owner exists — inconsistent state.")
// 		return s.teamRepo.Delete(ctx, team.ID)
// 	}

// 	// If current user is the last member → delete both team + membership
// 	if len(members) == 1 {
// 		return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
// 			txTeamRepo := s.teamRepo.NewTx(tx)
// 			txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

// 			if err := txTeamMemberRepo.Delete(ctx, sqlc.RemoveTeamMemberParams{
// 				TeamID: teamID,
// 				UserID: userID,
// 			}); err != nil {
// 				return err
// 			}
// 			return txTeamRepo.Delete(ctx, team.ID)
// 		})
// 	}

// 	// Choose the next owner deterministically
// 	var nextOwner sqlc.GetTeamMembersRow
// 	for _, m := range members {
// 		if m.UserID != userID {
// 			nextOwner = m
// 			break
// 		}
// 	}

// 	if nextOwner.UserID == uuid.Nil {
// 		// Should not happen unless data is inconsistent
// 		s.logger.Error().Msg("No eligible next owner found.")
// 		return ErrNoEligibleNextOwner
// 	}

// 	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
// 		txTeamRepo := s.teamRepo.NewTx(tx)
// 		txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

// 		if err := txTeamMemberRepo.Delete(ctx, sqlc.RemoveTeamMemberParams{
// 			TeamID: teamID,
// 			UserID: userID,
// 		}); err != nil {
// 			return err
// 		}

// 		// _, err := txTeamRepo.Update(ctx, team.ID, nil, &nextOwner.UserID)
// 		_, err := txTeamRepo.Update(ctx, sqlc.UpdateTeamByIdParams{
// 			ID:              team.ID,
// 			OwnerIDDoUpdate: true,
// 			OwnerID:         &nextOwner.UserID,
// 		})
// 		return err
// 	})
// }

// func (s *TeamService) RequestToJoinTeam(ctx context.Context, teamID, userID uuid.UUID, message *string) (*sqlc.TeamJoinRequest, error) {
// 	// Ensure user is not already on a team (User's can't request to join a team if they are already on a team)
// 	_, err := s.teamMemberRepo.GetTeamMemberByUser(ctx, userID)
// 	if err != nil && !errors.Is(err, repository.ErrTeamMemberNotFound) {
// 		return nil, err
// 	}

// 	if err == nil {
// 		return nil, ErrUserOnTeam
// 	}

// 	request, err := s.teamJoinRequestRepo.Create(ctx, sqlc.CreateTeamJoinRequestParams{
// 		TeamID:         teamID,
// 		UserID:         userID,
// 		RequestMessage: message,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return request, nil
// }

// func (s *TeamService) GetPendingJoinRequestForTeam(ctx context.Context, userID, teamID uuid.UUID) ([]sqlc.ListJoinRequestsByTeamAndStatusWithUserRow, error) {
// 	// Get team and check if user is the owner of the tema
// 	team, err := s.teamRepo.GetByID(ctx, teamID)
// 	if err != nil {
// 		return []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow{}, err
// 	}

// 	if *team.OwnerID != userID {
// 		return []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow{}, ErrUserNotTeamOwner
// 	}

// 	requests, err := s.teamJoinRequestRepo.ListJoinRequestsByTeamWithUser(ctx, sqlc.ListJoinRequestsByTeamAndStatusWithUserParams{
// 		TeamID: teamID,
// 		Status: sqlc.TeamJoinRequestStatusPending,
// 	})

// 	if err != nil {
// 		return []sqlc.ListJoinRequestsByTeamAndStatusWithUserRow{}, err
// 	}

// 	return requests, nil
// }

// func (s *TeamService) GetUserPendingJoinRequests(ctx context.Context, userID uuid.UUID) ([]sqlc.TeamJoinRequest, error) {
// 	requests, err := s.teamJoinRequestRepo.ListJoinRequestsByUserAndStatus(ctx, sqlc.ListTeamJoinRequestsByUserAndStatusParams{
// 		UserID: userID,
// 		Status: sqlc.TeamJoinRequestStatusPending,
// 	})

// 	if err != nil {
// 		return []sqlc.TeamJoinRequest{}, err
// 	}

// 	return requests, nil
// }

// func (s *TeamService) RespondToJoinRequest(ctx context.Context, ownerID, requestID uuid.UUID, accept bool) error {
// 	// Retrieve the join request
// 	oldRequest, err := s.teamJoinRequestRepo.GetById(ctx, requestID)
// 	if err != nil {
// 		return err
// 	}

// 	// Get team and check if user is the owner of the team
// 	team, err := s.teamRepo.GetByID(ctx, oldRequest.TeamID)
// 	if err != nil {
// 		return err
// 	}
// 	if *team.OwnerID != ownerID {
// 		return ErrUserNotTeamOwner
// 	}

// 	// Also ensure user is actually an applicant or attendee
// 	user, err := s.userRepo.GetUserByID(ctx, oldRequest.UserID)
// 	if err != nil {
// 		if errors.Is(err, database.ErrEntityNotFound) {
// 			return ErrUserNotApplicantOrAttendee
// 		}

// 		return err
// 	}
// 	if user == nil {
// 		return ErrUserNotApplicantOrAttendee
// 	}
// 	if user.Role != sqlc.UserRoleAttendee && user.Role != sqlc.UserRoleApplicant {
// 		return ErrUserNotApplicantOrAttendee
// 	}

// 	// Ensure user is not already on a team
// 	_, err = s.teamMemberRepo.GetTeamMemberByUser(ctx, oldRequest.UserID)
// 	if err != nil && !errors.Is(err, repository.ErrTeamMemberNotFound) {
// 		return err
// 	}

// 	if err == nil {
// 		return ErrUserOnTeam
// 	}

// 	if accept {
// 		// Check if team is already full
// 		members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
// 		if err != nil {
// 			return err
// 		}

// 		if len(members) >= 4 {
// 			return ErrTeamFull
// 		}

// 		// Accepting the request: add user to team and update request status
// 		return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
// 			txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)
// 			txTeamJoinRequestRepo := s.teamJoinRequestRepo.NewTx(tx)

// 			if _, err := txTeamMemberRepo.Create(ctx, sqlc.CreateTeamMemberParams{
// 				TeamID: oldRequest.TeamID,
// 				UserID: oldRequest.UserID,
// 			}); err != nil {
// 				return err
// 			}

// 			if _, err := txTeamJoinRequestRepo.UpdateStatus(ctx, sqlc.UpdateTeamJoinRequestParams{
// 				ID:             requestID,
// 				StatusDoUpdate: true,
// 				Status:         sqlc.TeamJoinRequestStatusApproved,
// 			}); err != nil {
// 				return err
// 			}

// 			// Delete all other pending requests by the user for this event
// 			return txTeamJoinRequestRepo.DeleteByUserAndStatus(ctx, sqlc.DeleteJoinRequestsByUserAndStatusParams{
// 				UserID: oldRequest.UserID,
// 				Status: sqlc.TeamJoinRequestStatusPending,
// 			})
// 		})
// 	} else {
// 		// Rejecting the request: just update request status
// 		_, err := s.teamJoinRequestRepo.UpdateStatus(ctx, sqlc.UpdateTeamJoinRequestParams{
// 			ID:             requestID,
// 			StatusDoUpdate: true,
// 			Status:         sqlc.TeamJoinRequestStatusApproved,
// 		})
// 		return err
// 	}
// }
