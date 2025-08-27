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
	"github.com/rs/zerolog/log"
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

func getEventID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		err := errors.New("missing_event_id")
		errMsg := "The event ID is missing from the URL!"

		log.Err(err).Msg(errMsg)

		res.SendError(w, http.StatusBadRequest,
			res.NewError(err.Error(), errMsg))

		return uuid.Nil, err
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		err = errors.New("invalid_event_id")
		errMsg := "The event ID is not a valid UUID"

		log.Err(err).Msg(errMsg)

		res.SendError(w, http.StatusBadRequest,
			res.NewError(err.Error(), errMsg))

		return uuid.Nil, err
	}

	return eventId, nil
}

func (h *ApplicationHandler) GetApplicationByUserAndEventID(w http.ResponseWriter, r *http.Request) {
	eventId, err := getEventID(w, r)

	if err != nil {
		return
	}

	userIdPtr := ctxutils.GetUserIdFromCtx(r.Context())

	if userIdPtr == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	userId := *userIdPtr

	params := sqlc.GetApplicationByUserAndEventIDParams{
		UserID:  userId,
		EventID: eventId,
	}

	application, err := h.appService.GetApplicationByUserAndEventID(r.Context(), params)

	if err != nil {
		if err == repository.ErrApplicationNotFound {
			params := sqlc.CreateApplicationParams{
				UserID:  userId,
				EventID: eventId,
			}
			newApplication, err := h.appService.CreateApplication(r.Context(), params)
			if err != nil {
				res.SendError(w, http.StatusBadRequest, res.NewError("create_application_error", "can't create application"))
			} else {
				res.Send(w, http.StatusOK, newApplication)
			}
			res.SendError(w, http.StatusBadRequest, res.NewError("application_not_found", "can't find application"))
		} else {
			res.SendError(w, http.StatusBadRequest, res.NewError("get_application_error", "error retrieving application"))
		}

		return
	}

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
	if err == nil {
		defer resumeFile.Close()
	}

	resumeFileBuffer := bytes.NewBuffer(nil)

	if _, err := io.Copy(resumeFileBuffer, resumeFile); err != nil {
		return
	}

	validate := validator.New()
	if err := validate.Struct(submission); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error()))
		return
	}

	eventId, err := getEventID(w, r)

	if err != nil {
		return
	}

	userIdPtr := ctxutils.GetUserIdFromCtx(r.Context())

	if userIdPtr == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	userId := *userIdPtr

	err = h.appService.SubmitApplication(r.Context(), submission, resumeFileBuffer.Bytes(), userId, eventId)

	if err != nil {
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

	eventId, err := getEventID(w, r)

	if err != nil {
		return
	}

	userIdPtr := ctxutils.GetUserIdFromCtx(r.Context())

	if userIdPtr == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "invalid user id"))
		return
	}

	userId := *userIdPtr

	err = h.appService.SaveApplication(r.Context(), data, sqlc.UpdateApplicationParams{
		UserID:  userId,
		EventID: eventId,
	})

	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("save_application_error", "Something went wrong while saving application"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
