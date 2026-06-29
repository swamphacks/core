package application

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type GetMyApplicationOutput struct {
	Body MyApplicationResponseDto
}

func (h *handler) handleGetMyApplication(ctx context.Context, input *struct{}) (*GetMyApplicationOutput, error) {
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

			return &GetMyApplicationOutput{Body: MyApplicationResponseDto{
				ID:          newApplication.ID,
				UserID:      newApplication.UserID,
				Status:      string(newApplication.Status),
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

	if application == nil {
		return nil, huma.Error500InternalServerError("Application is null")
	}

	return &GetMyApplicationOutput{Body: MyApplicationResponseDto{
		ID:          application.ID,
		UserID:      application.UserID,
		Status:      string(application.Status),
		Application: application.Application,
		CreatedAt:   application.CreatedAt,
		SavedAt:     application.SavedAt,
		UpdatedAt:   application.UpdatedAt,
		SubmittedAt: application.SubmittedAt,
		HackathonID: application.HackathonID,
	}}, nil
}

type UpdateApplicationByIdOutput struct {
	Status int
}

func (h *handler) handleUpdateApplicationById(ctx context.Context, input *struct {
	Body UpdateApplicationRequestDto
}) (*UpdateApplicationByIdOutput, error) {
	err := h.applicationService.UpdateApplicationById(ctx, input.Body)

	if err != nil {
		return nil, huma.Error500InternalServerError("error updating application")
	}

	return &UpdateApplicationByIdOutput{Status: http.StatusNoContent}, nil
}

type GetExtendedApplicationOutput struct {
	Body ExtendedApplicationResponseDto
}

func (h *handler) handleGetExtendedApplicationById(ctx context.Context, input *struct {
	ID uuid.UUID `path:"applicationId"`
}) (*GetExtendedApplicationOutput, error) {
	application, err := h.applicationService.GetExtendedApplicationById(ctx, input.ID)

	if err != nil {
		return nil, huma.Error500InternalServerError("error retrieving extended application")
	}

	resumeRequest, err := h.applicationService.GetApplicationResumeURL(ctx, application.UserID, 600)

	if err != nil {
		h.logger.Err(err).Str("ApplicationId", input.ID.String()).Msg(err.Error())
		return nil, huma.Error500InternalServerError("unable to get application resume")
	}

	var review *ReviewDto
	if application.ReviewID != nil {
		review = &ReviewDto{
			ID:               *application.ReviewID,
			ExperienceRating: application.ExperienceRating,
			PassionRating:    application.PassionRating,
			Notes:            application.Notes,
			ReviewUpdatedAt:  *application.ReviewUpdatedAt,
			ReviewUpdatedBy:  application.ReviewUpdatedBy,
			Reviewer: AppUser{
				ID:       *application.ReviewerID,
				UserName: *application.ReviewerName,
				Image:    application.ReviewerImage,
			},
		}
	}

	var autoDecisionRequest *AutoDecisionRequestDto
	if application.AutoDecisionRequestID != nil && application.RequestedDecision.Valid {
		autoDecisionRequest = &AutoDecisionRequestDto{
			ID:                   *application.AutoDecisionRequestID,
			ApplicationID:        application.ID,
			RequestedDecision:    string(application.RequestedDecision.ApplicationAutoDecisionType),
			Justification:        application.DecisionJustification,
			AutoDecisionApproved: *application.DecisionApproved,
			CreatedAt:            *application.DecisionRequestCreatedAt,
			DecidedBy:            application.DecidedBy,
		}
	}

	return &GetExtendedApplicationOutput{Body: ExtendedApplicationResponseDto{
		ID: application.ID,
		User: AppUser{
			ID:       application.UserID,
			UserName: application.UserName,
			Image:    application.UserImage,
			Email:    application.UserEmail,
		},
		Status:              string(application.Status),
		Application:         application.Application,
		CreatedAt:           application.CreatedAt,
		UpdatedAt:           application.UpdatedAt,
		SubmittedAt:         application.SubmittedAt,
		IsEarly:             application.IsEarly,
		Review:              review,
		AutoDecisionRequest: autoDecisionRequest,
		ResumeURL:           resumeRequest.URL,
	}}, nil
}

type SearchApplicationsOutput struct {
	Body SearchApplicationsResponseDto
}

func (h *handler) handleSearchApplications(ctx context.Context, input *struct {
	Limit  int32  `query:"limit"`
	Offset int32  `query:"offset"`
	Search string `query:"search"`
}) (*SearchApplicationsOutput, error) {
	count, applications, err := h.applicationService.SearchApplications(ctx, input.Limit, input.Offset, input.Search)

	if err != nil {
		return nil, huma.Error500InternalServerError("error retrieving applications")
	}

	searchApplicationResponses := make([]ApplicationWithUserInfoDto, len(applications))

	for i, val := range applications {
		searchApplicationResponses[i] = ApplicationWithUserInfoDto{
			ID: val.ID,
			User: AppUser{
				ID:       val.UserID,
				UserName: val.Name,
				Image:    val.Image,
				Email:    val.Email,
			},
			Status:      string(val.Status),
			Application: val.Application,
			CreatedAt:   val.CreatedAt,
			SubmittedAt: val.SubmittedAt,
			UpdatedAt:   val.CreatedAt,
			IsEarly:     val.IsEarly,
		}
	}

	return &SearchApplicationsOutput{Body: SearchApplicationsResponseDto{
		Count:        *count,
		Applications: searchApplicationResponses,
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

type SubmitApplicationOutput struct {
	Body SubmitApplicationResponseDto
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

	submission.Phone = r.FormValue("phone")
	submission.PreferredEmail = r.FormValue("preferredEmail")
	submission.UniversityEmail = r.FormValue("universityEmail")
	submission.Age = r.FormValue("age")
	submission.Country = r.FormValue("country")
	submission.Gender = r.FormValue("gender")
	submission.GenderOther = r.FormValue("gender-other")
	submission.Pronouns = r.FormValue("pronouns")
	submission.Race = r.FormValue("race")
	submission.RaceOther = r.FormValue("race-other")
	submission.Orientation = r.FormValue("orientation")

	submission.Linkedin = r.FormValue("linkedin")
	// submission.Github = r.FormValue("github")

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
		Body: SubmitApplicationResponseDto{
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
	Body ApplicationStatisticsDto
}

func (h *handler) handleGetApplicationStatistics(ctx context.Context, input *struct{}) (*GetApplicationStatisticsOutput, error) {
	stats, err := h.applicationService.GetApplicationStatistics(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get application statistics")
	}

	return &GetApplicationStatisticsOutput{Body: *stats}, nil
}

type UpdateApplicationReviewStatusForHackathonOutput struct {
	Status int
}

func (h *handler) handleUpdateApplicationReviewStatusForHackathon(ctx context.Context, input *struct {
	Body struct {
		Started bool `json:"started"  required:"true"`
	}
}) (*UpdateApplicationReviewStatusForHackathonOutput, error) {
	err := h.applicationService.UpdateApplicationReviewStatusForHackathon(ctx, input.Body.Started)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to update review status for hackathon")
	}

	return &UpdateApplicationReviewStatusForHackathonOutput{Status: http.StatusOK}, nil
}

type GetApplicationReviewDetailsOutput struct {
	Body ApplicationReviewResponseDto
}

func (h *handler) handleGetReviewById(ctx context.Context, input *struct {
	ID uuid.UUID `path:"reviewId"`
}) (*GetApplicationReviewDetailsOutput, error) {
	review, resume, err := h.applicationService.GetReviewById(ctx, input.ID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get application review details")
	}

	var autoDecisionRequest *AutoDecisionRequestDto
	if review.DecisionRequestID != nil && review.RequestedDecision.Valid {
		autoDecisionRequest = &AutoDecisionRequestDto{
			ID:                   *review.DecisionRequestID,
			ApplicationID:        review.ID,
			RequestedDecision:    string(review.RequestedDecision.ApplicationAutoDecisionType),
			Justification:        review.DecisionJustification,
			AutoDecisionApproved: *review.DecisionApproved,
			CreatedAt:            *review.DecisionRequestCreatedAt,
			DecidedBy:            review.DecisionDecidedBy,
		}
	}

	return &GetApplicationReviewDetailsOutput{Body: ApplicationReviewResponseDto{
		ID:                  review.ID,
		PassionRating:       review.PassionRating,
		ExperienceRating:    review.ExperienceRating,
		Notes:               review.Notes,
		Application:         review.Application,
		ResumeURL:           resume.URL,
		AutoDecisionRequest: autoDecisionRequest,
	}}, nil
}

type SubmitApplicationReviewOutput struct {
	Status int
}

func (h *handler) handleSubmitApplicationReview(ctx context.Context, input *struct {
	Body SaveReviewRequestDto
}) (*SubmitApplicationReviewOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.SaveApplicationReview(ctx, input.Body, userCtx.UserID, userCtx.Role)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to save review")
	}

	return &SubmitApplicationReviewOutput{Status: http.StatusCreated}, nil
}

type UpdateApplicationReviewOutput struct {
	Status int
}

func (h *handler) handleUpdateApplicationReview(ctx context.Context, input *struct {
	Body UpdateReviewRequestDto
}) (*UpdateApplicationReviewOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.UpdateApplicationReview(ctx, input.Body, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to update review")
	}

	return &UpdateApplicationReviewOutput{Status: http.StatusOK}, nil
}

type GetReviewAssignmentsOutput struct {
	Body []ReviewAssignmentDto `json:"body" nullable:"false"`
}

func (h *handler) handleGetReviewAssignments(ctx context.Context, input *struct{}) (*GetReviewAssignmentsOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	reviews, err := h.applicationService.GetReviewsForReviewer(ctx, userCtx.UserID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get assigned applications")
	}

	var assignments []ReviewAssignmentDto
	for _, review := range reviews {
		status := ApplicationReviewStatusInProgress
		if review.ExperienceRating != nil && review.PassionRating != nil {
			status = ApplicationReviewStatusCompleted
		}

		assignments = append(assignments, ReviewAssignmentDto{
			ReviewID:      review.ID,
			UserID:        review.UserID,
			ApplicationID: review.ApplicationID,
			Status:        status,
		})
	}

	return &GetReviewAssignmentsOutput{Body: assignments}, nil
}

type AssignApplicationReviewersOutput struct {
	Status int
}

func (h *handler) handleAssignApplicationReviewers(ctx context.Context, input *struct {
	Body []ReviewerAssignmentRequestDto
}) (*AssignApplicationReviewersOutput, error) {
	err := h.applicationService.AssignReviewersToApplications(ctx, input.Body)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to assign reviewers")
	}

	return &AssignApplicationReviewersOutput{Status: http.StatusOK}, nil
}

type GetAllReviewersAndProgressOutput struct {
	Body []ReviewerProgressResponseDto
}

func (h *handler) handleGetAllReviewersAndProgress(ctx context.Context, input *struct{}) (*GetAllReviewersAndProgressOutput, error) {
	results, err := h.applicationService.GetAllReviewersAndProgress(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get all reviewers and their progress")
	}

	reviewersProgress := make([]ReviewerProgressResponseDto, len(results))

	for i, val := range results {
		reviewersProgress[i] = ReviewerProgressResponseDto{
			ID:             val.ID,
			Name:           val.Name,
			Image:          val.Image,
			TotalAssigned:  val.TotalAssigned,
			CompletedCount: val.CompletedCount,
		}
	}

	return &GetAllReviewersAndProgressOutput{
		Body: reviewersProgress,
	}, nil
}

type SearchAutoDecisionRequestsOutput struct {
	Body SearchAutoDecisionRequestsResponseDto
}

func (h *handler) handleSearchAutoDecisionRequests(ctx context.Context, input *SearchAutoDecisionRequestsDto) (*SearchAutoDecisionRequestsOutput, error) {
	count, requests, err := h.applicationService.SearchAutoDecisionRequests(ctx, *input)

	if err != nil {
		return nil, huma.Error500InternalServerError("error retrieving requests")
	}

	searchAutoDecisionRequestsResponses := make([]ExtendedAutoDecisionRequestDto, len(requests))

	for i, val := range requests {
		var decidedBy *AppUser
		if val.DecidedBy != nil {
			decidedBy = &AppUser{
				ID:       *val.DecidedBy,
				UserName: *val.ApproverName,
				Image:    val.ApproverImage,
			}
		}

		searchAutoDecisionRequestsResponses[i] = ExtendedAutoDecisionRequestDto{
			ID: val.ID,
			User: AppUser{
				ID:       val.UserID,
				UserName: val.UserName,
				Image:    val.UserImage,
			},
			ApplicationID: val.ApplicationID,
			Reviewer: AppUser{
				ID:       val.ReviewerID,
				UserName: val.ReviewerName,
				Image:    val.ReviewerImage,
			},
			DecidedBy:         decidedBy,
			RequestedDecision: string(val.RequestedDecision),
			Justification:     val.Justification,
			Approved:          val.Approved,
			UpdatedAt:         val.UpdatedAt,
			CreatedAt:         val.CreatedAt,
		}
	}

	return &SearchAutoDecisionRequestsOutput{Body: SearchAutoDecisionRequestsResponseDto{
		Count:    *count,
		Requests: searchAutoDecisionRequestsResponses,
	}}, nil
}

type GetAutoDecisionRequestsOutput struct {
	Body []ExtendedAutoDecisionRequestDto `json:"body" nullable:"false"`
}

func (h *handler) handleGetAutoDecisionRequests(ctx context.Context, input *struct{}) (*GetAutoDecisionRequestsOutput, error) {
	requests, err := h.applicationService.GetAllAutoDecisionRequests(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to get auto decision requests")
	}

	extendedRequests := make([]ExtendedAutoDecisionRequestDto, len(requests))

	for i, val := range requests {
		var decidedBy *AppUser
		if val.DecidedBy != nil {
			decidedBy = &AppUser{
				ID:       *val.DecidedBy,
				UserName: *val.ApproverName,
				Image:    val.ApproverImage,
			}
		}

		extendedRequests[i] = ExtendedAutoDecisionRequestDto{
			ID:            val.ID,
			ApplicationID: val.ApplicationID,
			User: AppUser{
				ID:       val.UserID,
				UserName: val.UserName,
				Image:    val.UserImage,
			},
			Reviewer: AppUser{
				ID:       val.ReviewerID,
				UserName: val.ReviewerName,
				Image:    val.ReviewerImage,
			},
			DecidedBy:         decidedBy,
			RequestedDecision: string(val.RequestedDecision),
			Justification:     val.Justification,
			Approved:          val.Approved,
			UpdatedAt:         val.UpdatedAt,
			CreatedAt:         val.CreatedAt,
		}
	}

	return &GetAutoDecisionRequestsOutput{Body: extendedRequests}, nil
}

type RequestAutoDecisionOutput struct {
	Body AutoDecisionRequestDto
}

func (h *handler) handleRequestAutoDecision(ctx context.Context, input *struct {
	Body CreateAutoDecisionRequestDto
}) (*RequestAutoDecisionOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	req, err := h.applicationService.RequestAutoDecision(ctx, input.Body, userCtx.UserID, userCtx.Role)

	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to request decision")
	}

	return &RequestAutoDecisionOutput{Body: AutoDecisionRequestDto{
		ID:                   req.ID,
		ApplicationID:        req.ApplicationID,
		RequestedDecision:    string(req.RequestedDecision),
		Justification:        req.Justification,
		AutoDecisionApproved: *req.Approved,
		CreatedAt:            req.CreatedAt,
		DecidedBy:            req.DecidedBy,
	}}, nil
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

type UpdateAutoDecisionOutput struct {
	Status int
}

func (h *handler) handleUpdateAutoDecision(ctx context.Context, input *struct {
	Body UpdateAutoDecisionRequestDto
}) (*UpdateAutoDecisionOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err := h.applicationService.UpdateAutoDecisionRequest(
		ctx,
		input.Body,
		userCtx.UserID,
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

// type GetApplicationByUserIdOutput struct {
// 	Body Application
// }

// // TODO: Return StaffApplicationResponse instead
// func (h *handler) handleGetApplicationByUserId(ctx context.Context, input *struct {
// 	UserID string `path:"userId"`
// }) (*GetApplicationByUserIdOutput, error) {
// 	userID, err := uuid.Parse(input.UserID)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid user id")
// 	}

// 	application, err := h.applicationService.GetApplicationByUserId(ctx, userID)

// 	if err != nil {
// 		if errors.Is(err, database.ErrApplicationNotFound) {
// 			return nil, huma.Error404NotFound("Application not found for user")
// 		}

// 		return nil, huma.Error500InternalServerError("error retrieving application")
// 	}

// 	return &GetApplicationByUserIdOutput{Body: Application{
// 		ID:          application.ID,
// 		UserID:      application.UserID,
// 		Status:      string(application.Status),
// 		Application: application.Application,
// 		CreatedAt:   application.CreatedAt,
// 		SavedAt:     application.SavedAt,
// 		UpdatedAt:   application.UpdatedAt,
// 		SubmittedAt: application.SubmittedAt,
// 		HackathonID: application.HackathonID,
// 	}}, nil
// }

// type GetResumePresignedUrlOutput struct {
// 	Body string
// }

// func (h *handler) handleGetResumePresignedUrlByApplicationId(ctx context.Context, input *struct {
// 	ApplicationId string `path:"applicationId"`
// }) (*GetResumePresignedUrlOutput, error) {
// 	userCtx := ctxutils.GetUserFromCtx(ctx)

// 	if userCtx == nil {
// 		return nil, huma.Error400BadRequest("Failed to get current user info")
// 	}

// 	applicationId, err := uuid.Parse(input.ApplicationId)

// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Invalid applicationId")
// 	}

// 	if userCtx.Role != sqlc.UserRoleStaff && userCtx.Role != sqlc.UserRoleAdmin && userCtx.UserID != applicationId {
// 		return nil, huma.Error400BadRequest("You are not allowed to see other ppls resumes :(")
// 	}

// 	request, err := h.applicationService.GetApplicationResumeURL(ctx, applicationId, 600)

// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Unable to retrieve download url")
// 	}

// 	return &GetResumePresignedUrlOutput{Body: request.URL}, nil
// }

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
