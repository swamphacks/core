package email

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type acceptanceEmailData struct {
	Name                 string
	ConfirmationDeadline string
}

func (s *EmailService) acceptanceData(
	ctx context.Context, userID uuid.UUID, name string,
) (acceptanceEmailData, error) {
	deadline, err := s.hackathonRepo.GetAcceptedApplicationConfirmationDeadline(ctx, userID)
	if err != nil {
		return acceptanceEmailData{}, err
	}
	if deadline == nil {
		return acceptanceEmailData{}, errors.New("accepted application has no confirmation deadline")
	}
	if !time.Now().Before(*deadline) {
		return acceptanceEmailData{}, errors.New("confirmation deadline has already passed")
	}
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return acceptanceEmailData{}, err
	}
	return acceptanceEmailData{
		Name: name,
		ConfirmationDeadline: deadline.In(location).Format(
			"Monday, January 2, 2006 at 3:04 PM MST",
		),
	}, nil
}
