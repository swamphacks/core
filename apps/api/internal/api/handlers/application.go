package handlers

import (
	"bytes"
	"encoding/json"
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
)

type ApplicationHandler struct {
	appService *services.ApplicationService
}

func NewApplicationHandler(appService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}

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
		if err == repository.ErrApplicationNotFound {
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

		if err == services.ErrApplicationUnavailable {
			res.SendError(w, http.StatusBadRequest, res.NewError("get_application_error", "the application is unavailable"))
			return
		}

		res.SendError(w, http.StatusBadRequest, res.NewError("get_application_error", "error retrieving application"))
		return
	}

	// If the application status is not "started", then it means the user has submitted the application
	if application.Status.ApplicationStatus != sqlc.ApplicationStatusStarted {
		res.Send(w, http.StatusOK, map[string]any{"submitted": true})
		return
	}

	res.Send(w, http.StatusOK, application)
}

func (h *ApplicationHandler) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (10 MB max memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
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
	submission.Linkedin = r.FormValue("linkedin")
	submission.Github = r.FormValue("github")

	if ageCertStr := r.FormValue("ageCertification"); ageCertStr != "" {
		submission.AgeCertification = (ageCertStr == "true" || ageCertStr == "1")
	}

	submission.School = r.FormValue("school")
	submission.Level = r.FormValue("level")
	submission.Year = r.FormValue("year")
	submission.GraduationYear = r.FormValue("graduationYear")
	submission.Majors = r.FormValue("majors")
	submission.Minors = r.FormValue("minors")
	submission.Experience = r.FormValue("experience")
	submission.ProjectExperience = r.FormValue("projectExperience")
	submission.ShirtSize = r.FormValue("shirtSize")
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
		if err == services.ErrApplicationDeadlinePassed {
			res.SendError(w, http.StatusInternalServerError, res.NewError("submit_application_error", services.ErrApplicationDeadlinePassed.Error()))
			return
		}

		res.SendError(w, http.StatusInternalServerError, res.NewError("submit_application_error", "Something went wrong while submitting application"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
