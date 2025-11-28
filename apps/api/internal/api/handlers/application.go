package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
)

type ApplicationHandler struct {
	appService *services.ApplicationService
}

func NewApplicationHandler(appService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}

// Get Application By User and Event ID
//
//	@Summary		Get Application By User and Event ID
//	@Description	Get the current user's application progress for an event. If this is their first time filling out the application, a new application will be created.
//	@Tags			Application
//	@Accept			json
//	@Produce		json
//	@Param			eventId		path		string					true	"Event ID"
//	@Param			sh_session	cookie		string					true	"The authenticated session token/id"
//	@Success		200			{object}	sqlc.Application		"OK: An application was found"
//	@Success		200			{object}	map[string]any			"OK: An application was found"
//	@Failure		400			{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500			{object}	response.ErrorResponse	"Server Error: error retrieving application"\
//	@Router			/events/{eventId}/application [get]
func (h *ApplicationHandler) GetApplicationByUserAndEventID(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")

	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())

	if userId == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	params := sqlc.GetApplicationByUserAndEventIDParams{
		UserID:  *userId,
		EventID: eventId,
	}

	application, err := h.appService.GetApplicationByUserAndEventID(r.Context(), params)
	if err != nil {
		if errors.Is(err, repository.ErrApplicationNotFound) {
			params := sqlc.CreateApplicationParams{
				UserID:  *userId,
				EventID: eventId,
			}

			newApplication, err := h.appService.CreateApplication(r.Context(), params)

			if err != nil {
				res.SendError(w, http.StatusBadRequest, res.NewError("create_application_error", "can't create application"))
				return
			}

			if newApplication == nil {
				res.SendError(w, http.StatusBadRequest, res.NewError("create_application_error", "can't create application"))
				return
			}

			res.Send(w, http.StatusOK, newApplication)
			return
		}
		if errors.Is(err, services.ErrApplicationUnavailable) {
			res.SendError(w, http.StatusBadRequest, res.NewError("get_application_error", "the application is unavailable"))
			return
		}

		res.SendError(w, http.StatusBadRequest, res.NewError("get_application_error", "error retrieving application"))
		return
	}

	res.Send(w, http.StatusOK, application)
}

// Submit Application
//
//	@Summary		Submit Application
//	@Description	Submit the application for an event.
//	@Tags			Application
//	@Accept			json
//	@Produce		json
//	@Param			formBody	formData	any		true	"Submission form data"
//	@Param			eventId		path		string	true	"Event ID"	Format(uuid)
//	@Success		200
//	@Failure		400	{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500	{object}	response.ErrorResponse	"Server Error: error submitting application"
//	@Router			/events/{eventId}/application/submit [post]
func (h *ApplicationHandler) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (10 MB max memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("parse_form_invalid", "Failed to parse form: "+err.Error()))
		return
	}

	var submission services.ApplicationSubmissionFields

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
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "invalid resume file"))
		return
	}

	defer resumeFile.Close()

	resumeFileBuffer := bytes.NewBuffer(nil)

	if _, err := io.Copy(resumeFileBuffer, resumeFile); err != nil {
		return
	}

	validate := validator.New()
	if err := validate.Struct(submission); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error()))
		return
	}

	eventIdStr := chi.URLParam(r, "eventId")

	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())

	if userId == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	err = h.appService.SubmitApplication(r.Context(), submission, resumeFileBuffer.Bytes(), *userId, eventId)

	if err != nil {
		if errors.Is(err, services.ErrApplicationDeadlinePassed) {
			res.SendError(w, http.StatusInternalServerError, res.NewError("submit_application_error", services.ErrApplicationDeadlinePassed.Error()))
			return
		}

		res.SendError(w, http.StatusInternalServerError, res.NewError("submit_application_error", "Something went wrong while submitting application"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Save Application
//
//	@Summary		Save Application
//	@Description	Save user's progress on the application. File/Upload fields are not saved.
//	@Tags			Application
//	@Accept			json
//	@Produce		json
//	@Param			data	body	any		true	"Form data"
//	@Param			eventId	path	string	true	"Event ID"	Format(uuid)
//	@Success		200
//	@Failure		400	{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500	{object}	response.ErrorResponse	"Server Error: error saving application"
//	@Router			/events/{eventId}/application/save [post]
func (h *ApplicationHandler) SaveApplication(w http.ResponseWriter, r *http.Request) {
	var data any

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_form_data", "Something went wrong while parsing form submission"))
		return
	}

	eventIdStr := chi.URLParam(r, "eventId")

	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())

	if userId == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	err = h.appService.SaveApplication(r.Context(), data, *userId, eventId)

	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("save_application_error", "Something went wrong while saving application"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Download Resume
//
//	@Summary		Download the user's uploaded resume from their event application
//	@Description	This handler creates a presigned S3 URL with GET permission for the user's specific object, which is their uploaded resume. The client can use this URL to download the object.
//	@Tags			Application
//	@Produce		json
//	@Param			eventId	path		string	true	"Event ID"	Format(uuid)
//	@Success		200		{object}	string
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: error handling download resume request"
//	@Router			/events/{eventId}/application/download-resume [get]
func (h *ApplicationHandler) DownloadResume(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")

	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	userId := ctxutils.GetUserIdFromCtx(r.Context())

	if userId == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	request, err := h.appService.DownloadResume(r.Context(), *userId, eventId, 60)

	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("resume_download_error", "unable to retrieve resume download url"))
		return
	}

	res.Send(w, http.StatusOK, request.URL)
}

// Get Application Statistics
//
//	@Summary		Gets an event's submitted application statistics
//	@Description	This aggregates applications by race, gender, age, majors, and schools. This route is only available to event staff and admins.
//	@Tags			Application
//	@Produce		json
//	@Param			eventId	path		string	true	"Event ID"	Format(uuid)
//	@Success		200		{object}	services.ApplicationStatistics
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: error getting statistics"
//	@Router			/events/{eventId}/application/stats [get]
func (h *ApplicationHandler) GetApplicationStatistics(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID."))
		return
	}

	appStats, err := h.appService.GetApplicationStatistics(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("application_stats_err", "Something went wrong while aggregating application statistics."))
		return
	}

	res.Send(w, http.StatusOK, appStats)
}

// Get Assigned Application For Review
//
//	@Summary		Get a staff's assigned application for reviewing
//	@Description	Search through the `applications` and retrieves the one with matching assigned reviewer user id to the current user id.
//	@Tags			Application
//	@Produce		json
//	@Param			userId		query		string					true	"The user ID"
//	@Param			eventId		path		string					true	"Event ID"
//	@Param			sh_session	cookie		string					true	"The authenticated session token/id"
//	@Success		200			{object}	sqlc.Application		"OK: An application was found"
//	@Success		200			{object}	map[string]any			"OK: An application was found"
//	@Failure		400			{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500			{object}	response.ErrorResponse	"Server Error: error retrieving assigned application"
//	@Router			/events/{eventId}/application/assigned [get]
func (h *ApplicationHandler) GetAssignedApplication(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID."))
		return
	}

	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_user_id", "The user ID is missing from the URL query!"))
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "The user ID is not a valid UUID."))
		return
	}

	params := sqlc.GetAssignedApplicationByUserAndEventIdParams{
		UserID:  userId,
		EventID: eventId,
	}

	application, err := h.appService.GetAssignedApplicationByUserAndEventID(r.Context(), params)

	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("get_assigned_application_error", "error retrieving assigned application"))
		return
	}

	res.Send(w, http.StatusOK, application)
}

// Submit application review
//
//	@Summary		Submit application review
//	@Description	Handles ratings submissions from staff during the application review process.
//	@Tags			Application
//	@Produce		json
//	@Param			reviewData	body	any	true	"An object containing the passion and experience ratings"
//	@Success		200
//	@Failure		400	{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500	{object}	response.ErrorResponse	"Server Error: error submitting application review"
//	@Router			/events/{eventId}/application/submit-review [post]
func (h *ApplicationHandler) SubmitApplicationReview(w http.ResponseWriter, r *http.Request) {
	// eventIdStr := chi.URLParam(r, "eventId")
	// if eventIdStr == "" {
	// 	res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
	// 	return
	// }

	// eventId, err := uuid.Parse(eventIdStr)
	// if err != nil {
	// 	res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID."))
	// 	return
	// }

	// userIdStr := r.URL.Query().Get("userId")
	// if userIdStr == "" {
	// 	res.SendError(w, http.StatusBadRequest, res.NewError("missing_user_id", "The user ID is missing from the URL query!"))
	// 	return
	// }

	// userId, err := uuid.Parse(userIdStr)
	// if err != nil {
	// 	res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "The user ID is not a valid UUID."))
	// 	return
	// }

	// var data any

	// if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
	// 	res.SendError(w, http.StatusBadRequest, res.NewError("invalid_review_data", "Something went wrong while parsing review submission"))
	// 	return
	// }

	// params := sqlc.GetAssignedApplicationByUserAndEventIdParams{
	// 	UserID:  userId,
	// 	EventID: eventId,
	// }

	// application, err := h.appService.GetAssignedApplicationByUserAndEventID(r.Context(), params)

	// if err != nil {
	// 	res.SendError(w, http.StatusBadRequest, res.NewError("get_assigned_application_error", "error retrieving assigned application"))
	// 	return
	// }

	// res.Send(w, http.StatusOK)
}

// Assign application to reviewers
//
//	@Summary		Assign application to reviewers
//	@Description	Assigns applications for an event to reviewers for the application review process.
//	@Tags			Application
//	@Accept			json
//	@Param			eventId	path	string							true	"Event ID"	Format(uuid)
//	@Param			request	body	[]services.ReviewerAssignment	true	"Reviewer assignmnet payload"
//	@Success		201		"Reviewers assigned"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: error assigning reviewers"
//	@Router			/events/{eventId}/application/assign-reviewers [post]
func (h *ApplicationHandler) AssignApplicationReviewers(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid."))
		return
	}

	var payload []services.ReviewerAssignment
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Failed to parse request body: "+err.Error()))
		return
	}

	// Process assignments
	err = h.appService.AssignReviewers(r.Context(), eventId, payload)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("assign_reviewers_error", "Something went wrong while assigning reviewers to applications."))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Reset application reviews
//
//	@Summary		Reset application reviews
//	@Description	Resets all application reviews for a given event, clearing any existing reviewer assignments.
//	@Tags			Application
//
//	@Param			eventId	path	string	true	"ID of the event to reset reviews for"
//	@Success		200		"Application reviews reset successfully"
//	@Failure		400		{object}	res.ErrorResponse	"Bad request: invalid event ID"
//	@Failure		500		{object}	res.ErrorResponse	"Server error: failed to reset application reviews"
//	@Router			/events/{eventId}/application/reset-reviews [post]
func (h *ApplicationHandler) ResetApplicationReviews(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid."))
		return
	}

	err = h.appService.ResetApplicationReviews(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("reset_reviews_error", "Something went wrong while resetting application reviews."))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Get resume for application review
//
//	@Summary		Get resume for application review
//	@Description	This handler creates a presigned S3 URL with GET permission for a specific user's resume as an object. The client can use this URL to download the object temporarily for application review.
//	@Tags			Application
//	@Produce		json
//	@Param			eventId	path		string	true	"Event ID"	Format(uuid)
//	@Param			applicationId path  string  true	"The application ID (userId of applicant)" Format(uuid)
//	@Success		200		{object}	string
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request."
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: error handling download resume request"
//	@Router			/events/{eventId}/application/{applicationId}/download-reviewable-resume [get]
func (h *ApplicationHandler) DownloadReviewableResume(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")

	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	userIdStr := chi.URLParam(r, "userId")

	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_user_id", "The user ID is missing from the URL!"))
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "The user ID is not a valid UUID"))
		return
	}

	request, err := h.appService.DownloadResume(r.Context(), userId, eventId, 600)

	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("resume_download_error", "unable to retrieve resume download url"))
		return
	}

	res.Send(w, http.StatusOK, request.URL)
}
