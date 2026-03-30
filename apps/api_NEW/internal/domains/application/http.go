package application

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/domains/bat"
)

func RegisterRoutes(applicationHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "get-application",
		Method:        http.MethodGet,
		Summary:       "Get Application",
		Description:   "Get the application of the current user",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetApplication)

	huma.Register(group, huma.Operation{
		OperationID:   "save-application",
		Method:        http.MethodPost,
		Summary:       "Save Application",
		Description:   "Save user's progress on the application. File/Upload fields are not saved (eg. resumes).",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/save",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleSaveApplication)

	huma.Register(group, huma.Operation{
		OperationID:   "submit-application",
		Method:        http.MethodPost,
		Summary:       "Submit Application",
		Description:   "Submit the application",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RawHTTPMiddlewareHuma, mw.Auth.RequireAuthHuma},
		Path:          "/submit",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleSubmitApplication)

	huma.Register(group, huma.Operation{
		OperationID:   "get-download-resume-url",
		Method:        http.MethodGet,
		Summary:       "Get Resume Download URL",
		Description:   "Returns a presigned S3 URL with GET permission for the user's specific object, which is their uploaded resume. The client can use this URL to download the object.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/resume",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetDownloadResumeURL)

	huma.Register(group, huma.Operation{
		OperationID:   "get-application-statistics",
		Method:        http.MethodGet,
		Summary:       "Get Application Statistics",
		Description:   "Aggregates applications by race, gender, age, majors, and schools",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/stats",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetApplicationStatistics)

	huma.Register(group, huma.Operation{
		OperationID:   "submit-application-review",
		Method:        http.MethodPost,
		Summary:       "Submit Application Review",
		Description:   "Handles ratings submissions from staff during the application review process",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/{applicantId}",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleSubmitApplicationReview)

	huma.Register(group, huma.Operation{
		OperationID:   "get-assigned-applications",
		Method:        http.MethodGet,
		Summary:       "Get Assigned Applications",
		Description:   "Returns assigned applications and their review progress for the authenticated reviewer",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/assigned",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetAssignedApplications)

	huma.Register(group, huma.Operation{
		OperationID:   "assign-application-reviewers",
		Method:        http.MethodPost,
		Summary:       "Assign Application Reviewers",
		Description:   "Assigns applications to reviewers for the application review process.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/assign",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleAssignApplicationReviewers)

	huma.Register(group, huma.Operation{
		OperationID:   "reset-application-reviews",
		Method:        http.MethodPost,
		Summary:       "Reset Application Reviews",
		Description:   "Resets all application reviews, clearing any existing reviewer assignments.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/reset",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleResetApplicationReviews)

	huma.Register(group, huma.Operation{
		OperationID:   "get-resume",
		Method:        http.MethodGet,
		Summary:       "Get Resume URL (for review process)",
		Description:   "Returns a presigned S3 URL with GET permission for a specific user's resume as an object. The client can use this URL to download the object temporarily for application review.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/{applicantId}/resume",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetResumePresignedUrl)

	huma.Register(group, huma.Operation{
		OperationID:   "join-waitlist",
		Method:        http.MethodPatch,
		Summary:       "Join Waitlist",
		Description:   "Adds a waitlist join time to application. Sets status to waitlisted",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/join-waitlist",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleJoinWaitlist)

	huma.Register(group, huma.Operation{
		OperationID:   "withdraw-acceptance",
		Method:        http.MethodPatch,
		Summary:       "Withdraw Acceptance",
		Description:   "Withdraw an acceptance after being accepted to an event. Sets application status from accepted to rejected.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/withdraw-acceptance",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleWithdrawAcceptance)

	huma.Register(group, huma.Operation{
		OperationID:   "withdraw-attendance",
		Method:        http.MethodPatch,
		Summary:       "Withdraw Attendance",
		Description:   "Withdraw attendance after accepting to go to the hackathon. Sets application status from accepted to withdrawn. Sets event role from attendee, back to applicant.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/withdraw-attendance",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleWithdrawAttendance)

	huma.Register(group, huma.Operation{
		OperationID:   "accept-application-acceptance",
		Method:        http.MethodPatch,
		Summary:       "Accept Application Acceptance",
		Description:   "Accept an acceptance after being accepted. Sets event role to attendee, from applicant.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/accept-acceptance",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleAcceptApplicationAcceptance)

	huma.Register(group, huma.Operation{
		OperationID:   "transition-waitlist",
		Method:        http.MethodPatch,
		Summary:       "Transition Waitlisted Applications",
		Description:   "Transitions all accepted users to waitlist, and accepts 50 from the waitlist. Sets application status from accepted to rejected.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/transition-waitlisted-applications",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleTransitionWaitlistedApplications)

	huma.Register(group, huma.Operation{
		OperationID:   "calculate-admissions-request",
		Method:        http.MethodPost,
		Summary:       "Submit Admissions Calculation Request",
		Description:   "Queues an admission calculation task to the BAT worker",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/calculate-admissions",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleCalculateAdmissionsRequest)

	huma.Register(group, huma.Operation{
		OperationID:   "release-decisions",
		Method:        http.MethodPost,
		Summary:       "Release Decisions",
		Description:   "Releases decisions that were calculated by the worker from a specific run id",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/release-decisions/{runId}",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleReleaseDecisions)
}

type handler struct {
	applicationService *ApplicationService
	batService         *bat.BatService
	config             *config.Config
	logger             zerolog.Logger
}

func NewHandler(
	applicationService *ApplicationService, batService *bat.BatService,
	config *config.Config, logger zerolog.Logger,
) *handler {
	return &handler{
		applicationService: applicationService,
		batService:         batService,
		config:             config,
		logger:             logger,
	}
}

type GetApplicationOutput struct {
	Body *sqlc.Application
}

func (h *handler) handleGetApplication(ctx context.Context, input *struct{}) (*GetApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	application, err := h.applicationService.GetApplicationByUserId(ctx, userCtx.UserID)
	if err != nil {
		if errors.Is(err, database.ErrApplicationNotFound) {
			newApplication, err := h.applicationService.CreateApplication(ctx, userCtx.UserID)

			if err != nil || newApplication == nil {
				return nil, huma.Error500InternalServerError("can't create application")
			}

			return &GetApplicationOutput{Body: newApplication}, nil
		}
		if errors.Is(err, ErrApplicationNotOpened) {
			return nil, huma.Error400BadRequest("application is unavailable")
		}

		return nil, huma.Error500InternalServerError("error retrieving application")
	}

	return &GetApplicationOutput{Body: application}, nil
}

type SaveApplicationOutput struct {
	Status int
}

func (h *handler) handleSaveApplication(ctx context.Context, input *struct {
	Body any
}) (*SaveApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	// TODO: validate input.Body, make sure that the data is the application
	err := h.applicationService.SaveApplication(ctx, input.Body, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to save application")
	}

	return &SaveApplicationOutput{Status: http.StatusOK}, nil
}

type SubmitApplicationOutput struct {
	Status int
}

func (h *handler) handleSubmitApplication(ctx context.Context, input *struct{}) (*SubmitApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	r := ctx.Value(middleware.RawRequestKey{}).(*http.Request)

	// Parse multipart form (10 MB max memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to parse form")
	}

	var submission ApplicationSubmissionFields

	// Map form values
	submission.FirstName = r.FormValue("firstName")
	submission.LastName = r.FormValue("lastName")

	if ageStr := r.FormValue("age"); ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil {
			submission.Age = age
		}
	}

	submission.Phone = r.FormValue("phone")
	submission.PreferredEmail = r.FormValue("preferredEmail")
	submission.UniversityEmail = r.FormValue("universityEmail")

	submission.Country = r.FormValue("country")
	submission.Gender = r.FormValue("gender")
	submission.GenderOther = r.FormValue("gender-other")
	submission.Pronouns = r.FormValue("pronouns")
	submission.Race = r.FormValue("race")
	submission.RaceOther = r.FormValue("race-other")
	submission.Orientation = r.FormValue("orientation")

	submission.Linkedin = r.FormValue("linkedin")
	submission.Github = r.FormValue("github")

	if ageCertStr := r.FormValue("ageCertification"); ageCertStr != "" {
		submission.AgeCertification = (ageCertStr == "true" || ageCertStr == "1")
	}

	submission.School = r.FormValue("school")
	submission.Level = r.FormValue("level")
	submission.LevelOther = r.FormValue("level-other")
	submission.Year = r.FormValue("year")
	submission.YearOther = r.FormValue("year-other")
	submission.GraduationYear = r.FormValue("graduationYear")
	submission.Majors = r.FormValue("majors")
	submission.Minors = r.FormValue("minors")
	submission.Experience = r.FormValue("experience")
	submission.UfHackathonExp = r.FormValue("ufHackathonExp")
	submission.ProjectExperience = r.FormValue("projectExperience")
	submission.ShirtSize = r.FormValue("shirtSize")
	submission.Diet = r.FormValue("diet")
	submission.Essay1 = r.FormValue("essay1")
	submission.Essay2 = r.FormValue("essay2")
	submission.Referral = r.FormValue("referral")
	submission.PictureConsent = r.FormValue("pictureConsent")
	submission.InPersonAcknowledgement = r.FormValue("inpersonAcknowledgement")
	submission.AgreeToConduct = r.FormValue("agreeToConduct")
	submission.InfoShareAuthorization = r.FormValue("infoShareAuthorization")
	submission.AgreeToMLHEmails = r.FormValue("agreeToMLHEmails")

	resumeFile, _, err := r.FormFile("resume[]")
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid resume file")
	}

	defer resumeFile.Close()

	resumeFileBuffer := bytes.NewBuffer(nil)

	if _, err := io.Copy(resumeFileBuffer, resumeFile); err != nil {
		return nil, huma.Error500InternalServerError("Error while parsing resume")
	}

	validate := validator.New()
	if err := validate.Struct(submission); err != nil {
		return nil, huma.Error400BadRequest("Unable to parse application submission")
	}

	err = h.applicationService.SubmitApplication(r.Context(), submission, resumeFileBuffer.Bytes(), userCtx.UserID)

	if err != nil {
		if errors.Is(err, ErrApplicationNotOpened) {
			return nil, huma.Error400BadRequest("Application is not opened")
		}

		return nil, huma.Error500InternalServerError("Fail to submit application")
	}

	return &SubmitApplicationOutput{Status: http.StatusOK}, nil
}

type GetDownloadResumeOutput struct {
	Body string
}

func (h *handler) handleGetDownloadResumeURL(ctx context.Context, input *struct{}) (*GetDownloadResumeOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	request, err := h.applicationService.GetDownloadResumeURL(ctx, userCtx.UserID, 60)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get url to download resume")
	}

	return &GetDownloadResumeOutput{Body: request.URL}, nil
}

type GetApplicationStatisticsOutput struct {
	Body *ApplicationStatistics
}

func (h *handler) handleGetApplicationStatistics(ctx context.Context, input *struct{}) (*GetApplicationStatisticsOutput, error) {
	stats, err := h.applicationService.GetApplicationStatistics(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get application statistics")
	}

	return &GetApplicationStatisticsOutput{Body: stats}, nil
}

type ReviewRatings struct {
	PassionRating    int `json:"passionRating" minLength:"1" maxLength:"5" required:"true"`
	ExperienceRating int `json:"experienceRating" minLength:"1" maxLength:"5" required:"true"`
}

type SubmitApplicationReviewOutput struct {
	Status int
}

func (h *handler) handleSubmitApplicationReview(ctx context.Context, input *struct {
	ApplicantId string `path:"applicantId"`
	Body        ReviewRatings
}) (*SubmitApplicationReviewOutput, error) {
	applicantId, err := uuid.Parse(input.ApplicantId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid applicant id")
	}

	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err = h.applicationService.SaveApplicationReview(ctx, userCtx.UserID, applicantId, input.Body.ExperienceRating, input.Body.PassionRating)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to save review")
	}

	return &SubmitApplicationReviewOutput{Status: http.StatusCreated}, nil
}

type GetAssignedApplicationsOutput struct {
	Body []AssignedApplication
}

func (h *handler) handleGetAssignedApplications(ctx context.Context, input *struct{}) (*GetAssignedApplicationsOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	applications, err := h.applicationService.GetAssignedApplicationsAndProgress(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get assigned applications")
	}

	return &GetAssignedApplicationsOutput{Body: applications}, nil
}

type AssignApplicationReviewersOutput struct {
	Status int
}

func (h *handler) handleAssignApplicationReviewers(ctx context.Context, input *struct {
	Body []ReviewerAssignment
}) (*AssignApplicationReviewersOutput, error) {
	err := h.applicationService.AssignReviewers(ctx, input.Body)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to assign reviewers")
	}

	return &AssignApplicationReviewersOutput{Status: http.StatusOK}, nil
}

type ResetApplicationReviewsOutput struct {
	Status int
}

func (h *handler) handleResetApplicationReviews(ctx context.Context, input *struct{}) (*ResetApplicationReviewsOutput, error) {
	err := h.applicationService.ResetApplicationReviews(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to reset application reviews")
	}

	return &ResetApplicationReviewsOutput{Status: http.StatusOK}, nil
}

type GetResumePresignedUrlOutput struct {
	Body string
}

func (h *handler) handleGetResumePresignedUrl(ctx context.Context, input *struct {
	ApplicantId string `path:"applicantId"`
}) (*GetResumePresignedUrlOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	applicantId, err := uuid.Parse(input.ApplicantId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid applicant id")
	}

	if userCtx.Role != sqlc.UserRoleStaff && userCtx.Role != sqlc.UserRoleAdmin && userCtx.UserID != applicantId {
		return nil, huma.Error400BadRequest("You are not allowed to see other ppls resumes :(")
	}

	request, err := h.applicationService.GetDownloadResumeURL(ctx, applicantId, 600)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to retrieve download url")
	}

	return &GetResumePresignedUrlOutput{Body: request.URL}, nil
}

type JoinWaitlistOutput struct {
	Status int
}

func (h *handler) handleJoinWaitlist(ctx context.Context, input *struct{}) (*JoinWaitlistOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.JoinWaitlist(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to join waitlist")
	}

	return &JoinWaitlistOutput{Status: http.StatusOK}, nil
}

type WithdrawAcceptanceOutput struct {
	Status int
}

func (h *handler) handleWithdrawAcceptance(ctx context.Context, input *struct{}) (*WithdrawAcceptanceOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.WithdrawAcceptance(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to withdraw acceptance")
	}

	return &WithdrawAcceptanceOutput{Status: http.StatusOK}, nil
}

type WithdrawAttendanceOutput struct {
	Status int
}

func (h *handler) handleWithdrawAttendance(ctx context.Context, input *struct{}) (*WithdrawAttendanceOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.WithdrawAttendance(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to withdraw attendance")
	}

	return &WithdrawAttendanceOutput{Status: http.StatusOK}, nil
}

type AcceptApplicationAcceptanceOutput struct {
	Status int
}

func (h *handler) handleAcceptApplicationAcceptance(ctx context.Context, input *struct{}) (*AcceptApplicationAcceptanceOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.AcceptApplicationAcceptance(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to withdraw attendance")
	}

	return &AcceptApplicationAcceptanceOutput{Status: http.StatusOK}, nil
}

type TransitionWaitlistedApplicationsOutput struct {
	Status int
}

func (h *handler) handleTransitionWaitlistedApplications(ctx context.Context, input *struct{}) (*TransitionWaitlistedApplicationsOutput, error) {
	// TODO: move these numbers elsewhere to a config file or something
	var acceptanceCount uint32 = 50
	var acceptanceQuota uint32 = 500
	err := h.applicationService.TransitionWaitlistedApplications(ctx, acceptanceCount, acceptanceQuota)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to transition waitlisted applications")
	}

	return &TransitionWaitlistedApplicationsOutput{Status: http.StatusOK}, nil
}

type CalculateAdmissionsRequestOutput struct {
	Status int
}

func (h *handler) handleCalculateAdmissionsRequest(ctx context.Context, input *struct{}) (*CalculateAdmissionsRequestOutput, error) {
	_, err := h.batService.QueueCalculateAdmissionsTask(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to handle admissions calculation request")
	}

	return &CalculateAdmissionsRequestOutput{Status: http.StatusOK}, nil
}

type ReleaseDecisionsOutput struct {
	Status int
}

func (h *handler) handleReleaseDecisions(ctx context.Context, input *struct {
	RunId string `path:"runId"`
}) (*ReleaseDecisionsOutput, error) {
	runId, err := uuid.Parse(input.RunId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid Run Id")
	}

	err = h.applicationService.ReleaseDecisions(ctx, runId)

	if err != nil {
		return nil, err
	}

	return &ReleaseDecisionsOutput{Status: http.StatusOK}, nil

}
