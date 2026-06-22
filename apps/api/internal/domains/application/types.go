package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrApplicationNotOpened      = errors.New("Application not opened")
	ErrFailedToCreateApplication = errors.New("Failed to create application")
	ErrFailedToGetHackathon      = errors.New("Failed to get hackathon information")
	ErrCannotReplaceResume       = errors.New("cannot replace resume before the application has been submitted")
)

type AppUser struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"name"`
	Image    *string   `json:"image"`
	Email    *string   `json:"email"`
}

type MyApplicationResponseDto struct {
	ID          uuid.UUID  `json:"id"`
	Status      string     `json:"status"`
	Application []byte     `json:"application"`
	CreatedAt   time.Time  `json:"createdAt"`
	SubmittedAt *time.Time `json:"submittedAt" required:"false"`
	UserID      uuid.UUID  `json:"userId"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	SavedAt     time.Time  `json:"savedAt"`
	HackathonID string     `json:"hackathonId"`
}

type ExtendedApplicationResponseDto struct {
	ID                  uuid.UUID               `json:"id"`
	Status              string                  `json:"status"`
	Application         []byte                  `json:"application"`
	CreatedAt           time.Time               `json:"createdAt"`
	SubmittedAt         *time.Time              `json:"submittedAt" required:"false"`
	UpdatedAt           time.Time               `json:"updatedAt"`
	User                AppUser                 `json:"user"`
	IsEarly             bool                    `json:"isEarly"`
	Review              *ReviewDto              `json:"review" required:"false"`
	AutoDecisionRequest *AutoDecisionRequestDto `json:"autoDecisionRequest" required:"false"`
	ResumeURL           string                  `json:"resumeUrl"`
}

type SearchApplicationsResponseDto struct {
	Applications []ApplicationWithUserInfoDto `json:"applications"`
	Count        int64                        `json:"count"`
}

type SubmitApplicationResponseDto struct {
	SubmittedAt *time.Time `json:"submittedAt"`
}

type UpdateApplicationRequestDto struct {
	ApplicationID uuid.UUID               `json:"applicationId"`
	Status        *sqlc.ApplicationStatus `json:"status"`
}

type ApplicationWithUserInfoDto struct {
	ID          uuid.UUID  `json:"id"`
	Status      string     `json:"status"`
	Application []byte     `json:"application"`
	CreatedAt   time.Time  `json:"createdAt"`
	SubmittedAt *time.Time `json:"submittedAt" required:"false"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	User        AppUser    `json:"user"`
	IsEarly     bool       `json:"isEarly"`
}

type ApplicationReviewResponseDto struct {
	ID                  uuid.UUID               `json:"id"`
	ExperienceRating    *int32                  `json:"experienceRating"`
	PassionRating       *int32                  `json:"passionRating"`
	Notes               *string                 `json:"notes"`
	ReviewUpdatedAt     time.Time               `json:"reviewUpdatedAt"`
	ReviewUpdatedBy     *uuid.UUID              `json:"reviewUpdatedBy"`
	Application         []byte                  `json:"application"`
	ResumeURL           string                  `json:"resumeUrl"`
	AutoDecisionRequest *AutoDecisionRequestDto `json:"autoDecisionRequest" required:"false"`
}

type ReviewAssignmentDto struct {
	ReviewID      uuid.UUID               `json:"reviewId"`
	ApplicationID uuid.UUID               `json:"applicationId"`
	UserID        uuid.UUID               `json:"userId"`
	Status        ApplicationReviewStatus `json:"status"`
}

type ReviewDto struct {
	ID               uuid.UUID  `json:"id"`
	ExperienceRating *int32     `json:"experienceRating"`
	PassionRating    *int32     `json:"passionRating"`
	Notes            *string    `json:"notes"`
	ReviewUpdatedAt  time.Time  `json:"reviewUpdatedAt"`
	ReviewUpdatedBy  *uuid.UUID `json:"reviewUpdatedBy"`
	Reviewer         AppUser    `json:"reviewer"`
}

type ReviewerAssignmentRequestDto struct {
	ID     uuid.UUID `json:"userId"` // User/Reviewer ID
	Amount *int      `json:"amount"` // Number of applications assigned (nil if autoassign)
}

type ReviewerProgressResponseDto struct {
	ID             *uuid.UUID `json:"id"`
	Name           *string    `json:"name"`
	Image          *string    `json:"image"`
	TotalAssigned  int64      `json:"totalAssigned"`
	CompletedCount int64      `json:"completedCount"`
}

type SaveReviewRequestDto struct {
	ReviewId         uuid.UUID `json:"id"`
	ApplicationId    uuid.UUID `json:"applicationId"`
	ExperienceRating int       `json:"experienceRating"`
	PassionRating    int       `json:"passionRating"`
	Notes            string    `json:"notes"`
}

type UpdateReviewRequestDto struct {
	ReviewId         uuid.UUID `json:"id"`
	PassionRating    *int      `json:"passionRating"`
	ExperienceRating *int      `json:"experienceRating"`
	Notes            *string   `json:"notes"`
}

type CreateAutoDecisionRequestDto struct {
	ApplicationID     uuid.UUID `json:"applicationId"`
	RequestedDecision string    `json:"decision"`
	Justification     *string   `json:"justification"`
}

type UpdateAutoDecisionRequestDto struct {
	RequestID uuid.UUID `json:"requestId"`
	Approved  bool      `json:"approved"`
}

type AutoDecisionRequestDto struct {
	ID                   uuid.UUID  `json:"id"`
	ApplicationID        uuid.UUID  `json:"applicationId"`
	Justification        *string    `json:"justification" required:"false"`
	CreatedAt            time.Time  `json:"createdAt"`
	RequestedDecision    string     `json:"decision"`
	AutoDecisionApproved bool       `json:"approved"`
	DecidedBy            *uuid.UUID `json:"decidedBy" required:"false"`
}

type ExtendedAutoDecisionRequestDto struct {
	ID                uuid.UUID `json:"id"`
	ApplicationID     uuid.UUID `json:"applicationId"`
	Justification     *string   `json:"justification" required:"false"`
	CreatedAt         time.Time `json:"createdAt"`
	RequestedDecision string    `json:"requestedDecision"`
	Approved          *bool     `json:"approved" required:"false"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Reviewer          AppUser   `json:"reviewer"`
	DecidedBy         *AppUser  `json:"decidedBy" required:"false"`
	User              AppUser   `json:"user"`
}

type SearchAutoDecisionRequestsDto struct {
	Search   string `query:"search"`
	Approved string `query:"approved"`
	Decision string `query:"decision"`
	Offset   int32  `query:"offset"`
	Limit    int32  `query:"limit"`
}

type SearchAutoDecisionRequestsResponseDto struct {
	Requests []ExtendedAutoDecisionRequestDto `json:"autoDecisionRequests" nullable:"false"`
	Count    int64                            `json:"count"`
}

// TODO: figure out a way to create the submission fields dynamically using the json form files with proper validation.
// these fields are only applicable to swamphacks xi, not other events
type ApplicationSubmissionFields struct {
	FirstName               string `json:"firstName" validate:"required,max=50"`
	LastName                string `json:"lastName" validate:"required,max=50"`
	Age                     int    `json:"age" validate:"required,min=0,max=99"`
	Phone                   string `json:"phone" validate:"required,len=10"`
	PreferredEmail          string `json:"preferredEmail" validate:"required,email"`
	UniversityEmail         string `json:"universityEmail" validate:"required,email"`
	Country                 string `json:"country" validate:"required"`
	Gender                  string `json:"gender"`
	GenderOther             string `json:"gender-other"`
	Pronouns                string `json:"pronouns"`
	Race                    string `json:"race"`
	RaceOther               string `json:"race-other"`
	Orientation             string `json:"orientation"`
	Linkedin                string `json:"linkedin" validate:"required,url"`
	Github                  string `json:"github" validate:"required,url"`
	AgeCertification        bool   `json:"ageCertification" validate:"required,boolean"`
	School                  string `json:"school" validate:"required"`
	Level                   string `json:"level" validate:"required"`
	LevelOther              string `json:"level-other"`
	Year                    string `json:"year" validate:"required"`
	YearOther               string `json:"year-other"`
	GraduationYear          string `json:"graduationYear" validate:"required"`
	Majors                  string `json:"majors" validate:"required"`
	Minors                  string `json:"minors"`
	Experience              string `json:"experience" validate:"required"`
	UfHackathonExp          string `json:"ufHackathonExp" validate:"required"`
	ProjectExperience       string `json:"projectExperience" validate:"required"`
	ShirtSize               string `json:"shirtSize" validate:"required"`
	Diet                    string `json:"diet"`
	Essay1                  string `json:"essay1" validate:"required"`
	Essay2                  string `json:"essay2" validate:"required"`
	Referral                string `json:"referral" validate:"required"`
	PictureConsent          string `json:"pictureConsent" validate:"required"`
	InPersonAcknowledgement string `json:"inpersonAcknowledgement" validate:"required"`
	AgreeToConduct          string `json:"agreeToConduct" validate:"required"`
	InfoShareAuthorization  string `json:"infoShareAuthorization" validate:"required"`
	AgreeToMLHEmails        string `json:"agreeToMLHEmails"`
}

type ApplicationStatisticsDto struct {
	GenderStatistics sqlc.GetSubmittedApplicationGendersRow   `json:"genderStats"`
	AgeStatistics    sqlc.GetSubmittedApplicationAgesRow      `json:"ageStats"`
	RaceStatistics   []sqlc.GetSubmittedApplicationRacesRow   `json:"raceStats"`
	MajorStatistics  []sqlc.GetSubmittedApplicationMajorsRow  `json:"majorStats"`
	SchoolStatistics []sqlc.GetSubmittedApplicationSchoolsRow `json:"schoolStats"`
	StatusStatistics sqlc.GetApplicationStatusesRow           `json:"statusStats"`
}
