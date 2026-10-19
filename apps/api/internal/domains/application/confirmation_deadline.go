package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *ApplicationService) GetApplicationConfirmationDeadline(
	ctx context.Context, applicationID uuid.UUID,
) (*time.Time, error) {
	return s.db.GetApplicationConfirmationDeadline(ctx, applicationID)
}
