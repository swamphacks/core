package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Errors
var (
	ErrEntityNotFound          = errors.New("Entity not found")
	ErrUnexpectedFileType      = errors.New("did not expect this file type")
	ErrFailedToUploadBanner    = errors.New("failed to upload banner")
	ErrFailedToUpdateHackathon = errors.New("failed to update hackathon")

	// Account
	ErrAccountNotFound = errors.New("account not found")

	// Application
	ErrCreateApplication      = errors.New("unable to create application")
	ErrSaveApplication        = errors.New("unable to save application")
	ErrSubmitApplication      = errors.New("unable to submit application")
	ErrInvalidApplicationData = errors.New("unable to parse application data")
	ErrGetApplication         = errors.New("unable to get application for user")
	ErrApplicationNotFound    = errors.New("can not find application for user")

	// Bat runs
	ErrDuplicateRun        = errors.New("Run already exists in the database")
	ErrRunNotFound         = errors.New("Run not found")
	ErrNoRunsDeleted       = errors.New("No Runs deleted")
	ErrMultipleRunsDeleted = errors.New("Multiple Runs affected by delete query expecting to delete one")

	// Emails
	ErrDuplicateEmails = errors.New("email already exists in the database")
)

func IsUniqueViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}
	return false
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
