package application

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
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

	// huma.Register(group, huma.Operation{
	// 	OperationID:   "get-application-by-user-id",
	// 	Method:        http.MethodGet,
	// 	Summary:       "Get Application By User ID",
	// 	Description:   "Get the application of the specified user",
	// 	Tags:          []string{"Application"},
	// 	Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
	// 	Path:          "/{userId}",
	// 	Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
	// 	Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
	// 	DefaultStatus: http.StatusOK,
	// }, applicationHandler.handleGetApplicationByUserId)

	huma.Register(group, huma.Operation{
		OperationID:   "get-extended-application-by-id",
		Method:        http.MethodGet,
		Summary:       "Get Extended Application By ID",
		Description:   "Get the full details of an application, including reviews, applicant and other information",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/extended/{applicationId}",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetExtendedApplicationById)

	huma.Register(group, huma.Operation{
		OperationID:   "update-application-by-id",
		Method:        http.MethodPatch,
		Summary:       "Update Application By ID",
		Description:   "Update an application",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleUpdateApplicationById)

	huma.Register(group, huma.Operation{
		OperationID:   "search-applications",
		Method:        http.MethodGet,
		Summary:       "Search Applications",
		Description:   "Search applications for the current hackthaton. If params are empty, all applications are returned with pagination.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/search",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleSearchApplications)

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
		Method:      http.MethodPatch,
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
		OperationID:   "withdraw-application",
		Method:        http.MethodPatch,
		Summary:       "Withdraw Application",
		Description:   "Withdraw application after being accepted to the hackthon. Sets application status from accepted to withdrawn.",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Path:          "/withdraw",
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
		Path:          "/confirm",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleConfirmAttendance)

	// REVIEWS

	huma.Register(group, huma.Operation{
		OperationID:   "get-assigned-applications",
		Method:        http.MethodGet,
		Summary:       "Get Assigned Applications",
		Description:   "Returns assigned applications and their review progress for the authenticated reviewer",
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/assigned",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetReviewAssignments)

	huma.Register(group, huma.Operation{
		OperationID:   "get-review-by-id",
		Method:        http.MethodGet,
		Summary:       "Get Review By Id",
		Description:   "Get an application review detail including ratings, the application json, and resume",
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetReviewById)

	huma.Register(group, huma.Operation{
		OperationID:   "get-reviewers-and-progress",
		Method:        http.MethodGet,
		Summary:       "Get All Reviewers and Progress",
		Description:   "Get all reviewers and their progress",
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/progress",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetAllReviewersAndProgress)

	huma.Register(group, huma.Operation{
		OperationID:   "submit-application-review",
		Method:        http.MethodPost,
		Summary:       "Submit Application Review",
		Description:   "Handles ratings submissions from staff during the application review process",
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Path:          "/review",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleSubmitApplicationReview)

	huma.Register(group, huma.Operation{
		OperationID:   "update-application-review",
		Method:        http.MethodPatch,
		Summary:       "Update Application Review",
		Description:   "Update application review",
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleUpdateApplicationReview)

	huma.Register(group, huma.Operation{
		OperationID:   "update-application-review-status",
		Method:        http.MethodPost,
		Summary:       "Update Application Review Status",
		Description:   "Update the application review status for the current hackathon",
		Tags:          []string{"Application Review"},
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
		Tags:          []string{"Application Review"},
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
		Tags:          []string{"Application Review"},
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
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/all-auto-decision-requests",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleGetAutoDecisionRequests)

	huma.Register(group, huma.Operation{
		OperationID:   "request-auto-decision",
		Method:        http.MethodPost,
		Summary:       "Request Auto Decision",
		Description:   "Create a request to auto accept or auto reject applications.",
		Tags:          []string{"Application Review"},
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
		Tags:          []string{"Application Review"},
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
		Tags:          []string{"Application Review"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Path:          "/review/auto-decision",
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, applicationHandler.handleUpdateAutoDecision)

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
