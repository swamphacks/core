package application

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

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
)

func RegisterRoutes(applicationHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "get-application",
		Method:        http.MethodGet,
		Summary:       "Get My Application",
		Description:   "Get the application of the current user",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetMyApplication)

	huma.Register(group, huma.Operation{
		OperationID:   "get-application-by-user-id",
		Method:        http.MethodGet,
		Summary:       "Get Application By User ID",
		Description:   "Get the application of the specified user",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/{userId}",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetApplicationByUserId)

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
	}, applicationHandler.handleGetApplicationResumeURL)

	huma.Register(group, huma.Operation{
		OperationID: "replace-resume",
		Method:      http.MethodPut,
		Summary:     "Replace Resume",
		Description: "Replaces the resume of an already-submitted application without modifying any question responses.",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:        "/resume",
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, applicationHandler.handleReplaceResume)

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
		OperationID:   "get-application-for-review",
		Method:        http.MethodGet,
		Summary:       "Get Application For Review",
		Description:   "Get application details including ratings, application json, and resume for review",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/{applicationId}",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetApplicationForReview)

	huma.Register(group, huma.Operation{
		OperationID:   "get-resume",
		Method:        http.MethodGet,
		Summary:       "Get Resume URL By Application Id (for review process)",
		Description:   "Returns a presigned S3 URL with GET permission for a specific user's resume as an object. The client can use this URL to download the object temporarily for application review.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/{applicationId}/resume",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetResumePresignedUrlByApplicationId)

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
		OperationID:   "submit-application-review",
		Method:        http.MethodPost,
		Summary:       "Submit Application Review",
		Description:   "Handles ratings submissions from staff during the application review process",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/{applicationId}",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleSubmitApplicationReview)

	huma.Register(group, huma.Operation{
		OperationID:   "update-application-review-status",
		Method:        http.MethodPost,
		Summary:       "Update Application Review Status",
		Description:   "Update the application review status for the current hackathon",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/update-status",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleUpdateApplicationReviewStatusForHackathon)

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
		OperationID:   "get-auto-decision-requests",
		Method:        http.MethodGet,
		Summary:       "Get Auto Decision Requests",
		Description:   "Get all auto deicision requests created by staff",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/auto-decision",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetAutoDecisionRequests)

	huma.Register(group, huma.Operation{
		OperationID:   "request-auto-decision",
		Method:        http.MethodPost,
		Summary:       "Request Auto Decision",
		Description:   "Create a request to auto accept or auto reject applications",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/auto-decision",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleRequestAutoDecision)

	huma.Register(group, huma.Operation{
		OperationID:   "delete-auto-decision",
		Method:        http.MethodDelete,
		Summary:       "Delete Auto Decision",
		Description:   "Delete an existing auto decision made by current reviewer",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review/auto-decision",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleDeleteAutoDecision)

	huma.Register(group, huma.Operation{
		OperationID:   "update-auto-decision-request",
		Method:        http.MethodPatch,
		Summary:       "Update Auto Decision Request",
		Description:   "Update an auto decision request",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/auto-decision",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleUpdateAutoDecision)

	huma.Register(group, huma.Operation{
		OperationID:   "withdraw-application",
		Method:        http.MethodPatch,
		Summary:       "Withdraw Application",
		Description:   "Withdraw application after being accepted to the hackthon. Sets application status from accepted to withdrawn.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/withdraw-application",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleWithdrawApplication)

	huma.Register(group, huma.Operation{
		OperationID:   "confirm-attendance",
		Method:        http.MethodPatch,
		Summary:       "Confirm Attendance",
		Description:   "Confirm attendance after being accepted. Sets event role to attendee from applicant.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/confirm-attendance",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleConfirmAttendance)

	// huma.Register(group, huma.Operation{
	// 	OperationID:   "join-waitlist",
	// 	Method:        http.MethodPatch,
	// 	Summary:       "Join Waitlist",
	// 	Description:   "Adds a waitlist join time to application. Sets status to waitlisted",
	// 	Tags:          []string{"Application"},
	// 	Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
	// 	Path:          "/join-waitlist",
	// 	Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
	// 	Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
	// 	DefaultStatus: http.StatusOK,
	// }, applicationHandler.handleJoinWaitlist)

	// huma.Register(group, huma.Operation{
	// 	OperationID:   "withdraw-acceptance",
	// 	Method:        http.MethodPatch,
	// 	Summary:       "Withdraw Acceptance",
	// 	Description:   "Withdraw an acceptance after being accepted to an event. Sets application status from accepted to rejected.",
	// 	Tags:          []string{"Application"},
	// 	Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
	// 	Path:          "/withdraw-acceptance",
	// 	Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
	// 	Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
	// 	DefaultStatus: http.StatusOK,
	// }, applicationHandler.handleWithdrawAcceptance)

	// huma.Register(group, huma.Operation{
	// 	OperationID:   "transition-waitlist",
	// 	Method:        http.MethodPatch,
	// 	Summary:       "Transition Waitlisted Applications",
	// 	Description:   "Transitions all accepted users to waitlist, and accepts 50 from the waitlist. Sets application status from accepted to rejected.",
	// 	Tags:          []string{"Application"},
	// 	Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
	// 	Path:          "/transition-waitlisted-applications",
	// 	Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
	// 	Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
	// 	DefaultStatus: http.StatusOK,
	// }, applicationHandler.handleTransitionWaitlistedApplications)

	// huma.Register(group, huma.Operation{
	// 	OperationID:   "calculate-admissions-request",
	// 	Method:        http.MethodPost,
	// 	Summary:       "Submit Admissions Calculation Request",
	// 	Description:   "Queues an admission calculation task to the BAT worker",
	// 	Tags:          []string{"Application"},
	// 	Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
	// 	Path:          "/calculate-admissions",
	// 	Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
	// 	Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
	// 	DefaultStatus: http.StatusOK,
	// }, applicationHandler.handleCalculateAdmissionsRequest)

	// huma.Register(group, huma.Operation{
	// 	OperationID:   "release-decisions",
	// 	Method:        http.MethodPost,
	// 	Summary:       "Release Decisions",
	// 	Description:   "Releases decisions that were calculated by the worker from a specific run id",
	// 	Tags:          []string{"Application"},
	// 	Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
	// 	Path:          "/release-decisions/{runId}",
	// 	Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
	// 	Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
	// 	DefaultStatus: http.StatusOK,
	// }, applicationHandler.handleReleaseDecisions)
}

type handler struct {
	applicationService *ApplicationService
	config             *config.Config
	logger             zerolog.Logger
}

func NewHandler(
	applicationService *ApplicationService,
	config *config.Config, logger zerolog.Logger,
) *handler {
	return &handler{
		applicationService: applicationService,
		config:             config,
		logger:             logger,
	}
}

type HackerApplication struct {
	UserID      uuid.UUID              `json:"userId"`
	Status      sqlc.ApplicationStatus `json:"status"`
	Application []byte                 `json:"application"`
	CreatedAt   time.Time              `json:"createdAt"`
	SavedAt     time.Time              `json:"savedAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	SubmittedAt *time.Time             `json:"submittedAt"`
	HackathonID string                 `json:"hackathonId"`
}

type GetApplicationOutput struct {
	Body HackerApplication
}

func (h *handler) handleGetMyApplication(ctx context.Context, input *struct{}) (*GetApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	application, err := h.applicationService.GetApplicationByUserId(ctx, userCtx.UserID)
	if err != nil {
		if errors.Is(err, database.ErrApplicationNotFound) {
			newApplication, err := h.applicationService.CreateApplication(ctx, userCtx.UserID)

			if err != nil || newApplication == nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			return &GetApplicationOutput{Body: HackerApplication{
				UserID:      newApplication.UserID,
				Status:      newApplication.Status,
				Application: newApplication.Application,
				CreatedAt:   newApplication.CreatedAt,
				SavedAt:     newApplication.SavedAt,
				UpdatedAt:   newApplication.UpdatedAt,
				SubmittedAt: newApplication.SubmittedAt,
				HackathonID: newApplication.HackathonID,
			}}, nil
		}

		if errors.Is(err, ErrApplicationNotOpened) {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, huma.Error500InternalServerError("error retrieving application")
	}

	return &GetApplicationOutput{Body: HackerApplication{
		UserID:      application.UserID,
		Status:      application.Status,
		Application: application.Application,
		CreatedAt:   application.CreatedAt,
		SavedAt:     application.SavedAt,
		UpdatedAt:   application.UpdatedAt,
		SubmittedAt: application.SubmittedAt,
		HackathonID: application.HackathonID,
	}}, nil
}

type GetApplicationByUserIdOutput struct {
	Body HackerApplication
}

func (h *handler) handleGetApplicationByUserId(ctx context.Context, input *struct {
	UserID string `path:"userId"`
}) (*GetApplicationByUserIdOutput, error) {
	userID, err := uuid.Parse(input.UserID)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	application, err := h.applicationService.GetApplicationByUserId(ctx, userID)

	if err != nil {
		if errors.Is(err, database.ErrApplicationNotFound) {
			return nil, huma.Error404NotFound("Application not found for user")
		}

		return nil, huma.Error500InternalServerError("error retrieving application")
	}

	return &GetApplicationByUserIdOutput{Body: HackerApplication{
		UserID:      application.UserID,
		Status:      application.Status,
		Application: application.Application,
		SubmittedAt: application.SubmittedAt,
		HackathonID: application.HackathonID,
	}}, nil
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
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &SaveApplicationOutput{Status: http.StatusOK}, nil
}

type SubmissionResult struct {
	SubmittedAt *time.Time `json:"submittedAt"`
}

type SubmitApplicationOutput struct {
	Body SubmissionResult
}

func (h *handler) handleSubmitApplication(ctx context.Context, input *struct{}) (*SubmitApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	// TODO: refactor this to using Huma's request API instead of using the raw http package

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

	submittedAt, err := h.applicationService.SubmitApplication(r.Context(), submission, resumeFileBuffer.Bytes(), userCtx.UserID)

	if err != nil {
		if errors.Is(err, ErrApplicationNotOpened) {
			return nil, huma.Error400BadRequest("Application is not opened")
		}

		return nil, huma.Error500InternalServerError("Fail to submit application")
	}

	return &SubmitApplicationOutput{
		Body: SubmissionResult{
			SubmittedAt: submittedAt,
		},
	}, nil
}

type GetDownloadResumeOutput struct {
	Body string
}

func (h *handler) handleGetApplicationResumeURL(ctx context.Context, input *struct{}) (*GetDownloadResumeOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	request, err := h.applicationService.GetApplicationResumeURL(ctx, userCtx.UserID, 60)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get url to download resume")
	}

	return &GetDownloadResumeOutput{Body: request.URL}, nil
}

type ReplaceResumeOutput struct {
	Status int
}

func (h *handler) handleReplaceResume(ctx context.Context, input *struct {
	RawBody huma.MultipartFormFiles[struct {
		Resume huma.FormFile `form:"resume" contentType:"application/pdf" required:"true"`
	}]
}) (*ReplaceResumeOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	fileHeaders := input.RawBody.Form.File["resume"]
	if len(fileHeaders) == 0 {
		return nil, huma.Error400BadRequest("Invalid resume file")
	}
	fileHeader := fileHeaders[0]

	if fileHeader.Size > 10*1024*1024 { // 10 MiB
		return nil, huma.Error400BadRequest("File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to parse uploaded resume")
	}
	defer file.Close()

	resumeBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(resumeBuffer, file); err != nil {
		return nil, huma.Error500InternalServerError("Error while parsing resume")
	}

	if err := h.applicationService.ReplaceResume(ctx, userCtx.UserID, resumeBuffer.Bytes()); err != nil {
		if errors.Is(err, database.ErrApplicationNotFound) {
			return nil, huma.Error400BadRequest("No application found to replace resume for")
		}
		if errors.Is(err, ErrCannotReplaceResume) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, huma.Error500InternalServerError("Unable to replace resume")
	}

	return &ReplaceResumeOutput{Status: http.StatusNoContent}, nil
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

type UpdateApplicationReviewStatusForHackathonOutput struct {
	Status int
}

type ReviewStatus struct {
	Started bool `json:"started"  required:"true"`
}

func (h *handler) handleUpdateApplicationReviewStatusForHackathon(ctx context.Context, input *struct {
	Body ReviewStatus
}) (*UpdateApplicationReviewStatusForHackathonOutput, error) {
	err := h.applicationService.UpdateApplicationReviewStatusForHackathon(ctx, input.Body.Started)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to update review status for hackathon")
	}

	return &UpdateApplicationReviewStatusForHackathonOutput{Status: http.StatusOK}, nil
}

type GetApplicationForReviewOutput struct {
	Body ApplicationReviewDetails
}

func (h *handler) handleGetApplicationForReview(ctx context.Context, input *struct {
	ApplicationId string `path:"applicationId"`
}) (*GetApplicationForReviewOutput, error) {
	applicationId, err := uuid.Parse(input.ApplicationId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid applicant id")
	}

	applicationReviewDetails, err := h.applicationService.GetApplicationReviewDetails(ctx, applicationId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get application review details")
	}

	return &GetApplicationForReviewOutput{Body: *applicationReviewDetails}, nil
}

type ReviewRatings struct {
	PassionRating    int `json:"passionRating" minLength:"1" maxLength:"100" required:"true"`
	ExperienceRating int `json:"experienceRating" minLength:"1" maxLength:"100" required:"true"`
}

type SubmitApplicationReviewOutput struct {
	Status int
}

func (h *handler) handleSubmitApplicationReview(ctx context.Context, input *struct {
	ApplicationId string `path:"applicationId"`
	Body          ReviewRatings
}) (*SubmitApplicationReviewOutput, error) {
	applicationId, err := uuid.Parse(input.ApplicationId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid applicant id")
	}

	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err = h.applicationService.SaveApplicationReview(ctx, userCtx.UserID, applicationId, input.Body.ExperienceRating, input.Body.PassionRating)

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

	applications, err := h.applicationService.GetAssignedApplicationsForReviewer(ctx, userCtx.UserID)

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
	err := h.applicationService.AssignReviewerToApplications(ctx, input.Body)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to assign reviewers")
	}

	return &AssignApplicationReviewersOutput{Status: http.StatusOK}, nil
}

type GetAutoDecisionRequestsOutput struct {
	Body []sqlc.ListAutoDecisionRequestsRow
}

func (h *handler) handleGetAutoDecisionRequests(ctx context.Context, input *struct{}) (*GetAutoDecisionRequestsOutput, error) {
	requests, err := h.applicationService.GetAllAutoDecisionRequests(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get auto decision requests")
	}

	return &GetAutoDecisionRequestsOutput{Body: requests}, nil
}

type RequestAutoDecisionOutput struct {
	Status int
}

func (h *handler) handleRequestAutoDecision(ctx context.Context, input *struct {
	Body AutoDecisionRequest
}) (*RequestAutoDecisionOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.RequestAutoDecision(ctx, input.Body, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to request decision")
	}

	return &RequestAutoDecisionOutput{Status: http.StatusOK}, nil
}

type DeleteAutoDecisionRequest struct {
	RequestId uuid.UUID `json:"requestId"`
}

type DeleteAutoDecisionOutput struct {
	Status int
}

func (h *handler) handleDeleteAutoDecision(ctx context.Context, input *struct {
	Body DeleteAutoDecisionRequest
}) (*DeleteAutoDecisionOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.DeleteAutoDecisionRequest(ctx, input.Body.RequestId, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to delete auto decision request")
	}

	return &DeleteAutoDecisionOutput{Status: http.StatusOK}, nil
}

type UpdateAutoDecisionRequestBody struct {
	RequestId uuid.UUID `json:"requestId"`
	Approved  bool      `json:"approved"`
}

type UpdateAutoDecisionOutput struct {
	Status int
}

func (h *handler) handleUpdateAutoDecision(ctx context.Context, input *struct {
	Body UpdateAutoDecisionRequestBody
}) (*UpdateAutoDecisionOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.UpdateAutoDecisionRequest(
		ctx,
		input.Body.RequestId,
		userCtx.UserID,
		input.Body.Approved,
	)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to update auto decision request")
	}

	return &UpdateAutoDecisionOutput{Status: http.StatusOK}, nil
}

type ResetApplicationReviewsOutput struct {
	Status int
}

func (h *handler) handleResetApplicationReviews(ctx context.Context, input *struct{}) (*ResetApplicationReviewsOutput, error) {
	err := h.applicationService.DeleteAllApplicationReviews(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to reset application reviews")
	}

	return &ResetApplicationReviewsOutput{Status: http.StatusOK}, nil
}

type GetResumePresignedUrlOutput struct {
	Body string
}

func (h *handler) handleGetResumePresignedUrlByApplicationId(ctx context.Context, input *struct {
	ApplicationId string `path:"applicationId"`
}) (*GetResumePresignedUrlOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	applicationId, err := uuid.Parse(input.ApplicationId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid applicationId")
	}

	if userCtx.Role != sqlc.UserRoleStaff && userCtx.Role != sqlc.UserRoleAdmin && userCtx.UserID != applicationId {
		return nil, huma.Error400BadRequest("You are not allowed to see other ppls resumes :(")
	}

	request, err := h.applicationService.GetApplicationResumeURL(ctx, applicationId, 600)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to retrieve download url")
	}

	return &GetResumePresignedUrlOutput{Body: request.URL}, nil
}

type WithdrawApplicationOutput struct {
	Status int
}

func (h *handler) handleWithdrawApplication(ctx context.Context, input *struct{}) (*WithdrawApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	if userCtx.Role != sqlc.UserRoleApplicant {
		return nil, huma.Error400BadRequest("Not an applicant")
	}

	err := h.applicationService.WithdrawApplication(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to withdraw attendance")
	}

	return &WithdrawApplicationOutput{Status: http.StatusOK}, nil
}

type ConfirmAttendanceOutput struct {
	Status int
}

func (h *handler) handleConfirmAttendance(ctx context.Context, input *struct{}) (*ConfirmAttendanceOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	if userCtx.Role != sqlc.UserRoleApplicant {
		return nil, huma.Error400BadRequest("Not an applicant")
	}

	err := h.applicationService.ConfirmAttendance(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to withdraw attendance")
	}

	return &ConfirmAttendanceOutput{Status: http.StatusOK}, nil
}

// type JoinWaitlistOutput struct {
// 	Status int
// }

// func (h *handler) handleJoinWaitlist(ctx context.Context, input *struct{}) (*JoinWaitlistOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	err := h.applicationService.JoinWaitlist(ctx, userCtx.UserID)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Unable to join waitlist")
// 	}

// 	return &JoinWaitlistOutput{Status: http.StatusOK}, nil
// }

// type WithdrawAcceptanceOutput struct {
// 	Status int
// }

// func (h *handler) handleWithdrawAcceptance(ctx context.Context, input *struct{}) (*WithdrawAcceptanceOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	err := h.applicationService.WithdrawAcceptance(ctx, userCtx.UserID)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Unable to withdraw acceptance")
// 	}

// 	return &WithdrawAcceptanceOutput{Status: http.StatusOK}, nil
// }

// type TransitionWaitlistedApplicationsOutput struct {
// 	Status int
// }

// func (h *handler) handleTransitionWaitlistedApplications(ctx context.Context, input *struct{}) (*TransitionWaitlistedApplicationsOutput, error) {
// 	// TODO: move these numbers elsewhere to a config file or something
// 	var acceptanceCount uint32 = 50
// 	var acceptanceQuota uint32 = 500
// 	err := h.applicationService.TransitionWaitlistedApplications(ctx, acceptanceCount, acceptanceQuota)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Fail to transition waitlisted applications")
// 	}

// 	return &TransitionWaitlistedApplicationsOutput{Status: http.StatusOK}, nil
// }

// type CalculateAdmissionsRequestOutput struct {
// 	Status int
// }

// func (h *handler) handleCalculateAdmissionsRequest(ctx context.Context, input *struct{}) (*CalculateAdmissionsRequestOutput, error) {
// 	_, err := h.batService.QueueCalculateAdmissionsTask(ctx)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Failed to handle admissions calculation request")
// 	}

// 	return &CalculateAdmissionsRequestOutput{Status: http.StatusOK}, nil
// }

// type ReleaseDecisionsOutput struct {
// 	Status int
// }

// func (h *handler) handleReleaseDecisions(ctx context.Context, input *struct {
// 	RunId string `path:"runId"`
// }) (*ReleaseDecisionsOutput, error) {
// 	runId, err := uuid.Parse(input.RunId)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid Run Id")
// 	}

// 	err = h.applicationService.ReleaseDecisions(ctx, runId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ReleaseDecisionsOutput{Status: http.StatusOK}, nil

// }
