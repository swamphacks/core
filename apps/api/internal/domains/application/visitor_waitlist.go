package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
)

var ErrVisitorWaitlistEligibility = errors.New(
	"only visitors with no completed application can use this form")

func registerVisitorWaitlistRoutes(h *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "submit-visitor-waitlist",
		Method:        http.MethodPost,
		Path:          "/visitor-waitlist",
		Summary:       "Submit visitor waitlist registration",
		Tags:          []string{"Application"},
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:        []int{400, 401, 409, 500},
		DefaultStatus: http.StatusOK,
	}, h.handleSubmitVisitorWaitlist)
}

func validateVisitorSubmission(data ApplicationSubmissionFields) error {
	if err := validator.New().StructExcept(
		data, "Essay1", "Essay2", "UfHackathonExp"); err != nil {
		return errors.New("Complete all required profile fields")
	}
	if !strings.HasSuffix(strings.ToLower(data.UniversityEmail), ".edu") {
		return errors.New("University email must end in .edu")
	}
	if !regexp.MustCompile("^[0-9]{10}$").MatchString(data.Phone) {
		return errors.New("Phone number must contain 10 digits")
	}
	if !regexp.MustCompile("^(0[1-9]|1[0-2])/[0-9]{4}$").MatchString(data.GraduationYear) {
		return errors.New("Graduation date must use MM/YYYY")
	}
	for _, consent := range []string{
		data.PictureConsent, data.InPersonAcknowledgement,
		data.AgreeToConduct, data.InfoShareAuthorization,
	} {
		if consent != "agree" {
			return errors.New("Agree to each required consent")
		}
	}
	if data.AgreeToMLHEmails != "" && data.AgreeToMLHEmails != "agree" {
		return errors.New("Invalid MLH email preference")
	}
	return nil
}

func (s *ApplicationService) SubmitVisitorWaitlist(
	ctx context.Context, data ApplicationSubmissionFields, resume []byte, userID uuid.UUID,
) (*time.Time, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)
	if err != nil {
		return nil, ErrGetHackathon
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	payload["registrationType"] = "visitor-waitlist"
	raw, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	submittedAt := time.Now()
	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		var role string
		var isFake bool
		if err := tx.QueryRow(ctx,
			"SELECT role::text, is_fake FROM users WHERE id=$1 FOR UPDATE",
			userID).Scan(&role, &isFake); err != nil {
			return err
		}

		var activeID string
		if err := tx.QueryRow(ctx,
			"SELECT id FROM hackathons WHERE id=$1 AND is_active FOR SHARE",
			hackathon.ID).Scan(&activeID); err != nil {
			return ErrGetHackathon
		}

		var id uuid.UUID
		var status string
		var existing []byte
		var previousSubmission *time.Time
		err := tx.QueryRow(ctx, `
            SELECT id, status::text, application, submitted_at
            FROM applications
            WHERE user_id=$1 AND hackathon_id=$2 FOR UPDATE
        `, userID, activeID).Scan(&id, &status, &existing, &previousSubmission)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if err == nil && status == "waitlisted" {
			var old map[string]any
			if json.Unmarshal(existing, &old) == nil &&
				old["registrationType"] == "visitor-waitlist" &&
				role == "applicant" && previousSubmission != nil {
				submittedAt = *previousSubmission
				return nil
			}
		}
		if role != "visitor" || (err == nil && status != "started") {
			return ErrVisitorWaitlistEligibility
		}
		if !admissionWaitlistIsOpen(time.Now()) {
			return ErrWaitlistClosed
		}

		if errors.Is(err, pgx.ErrNoRows) {
			err = tx.QueryRow(ctx, `
                INSERT INTO applications (user_id, hackathon_id, is_fake)
                VALUES ($1,$2,$3) RETURNING id
            `, userID, activeID, isFake).Scan(&id)
			if err != nil {
				return err
			}
		}

		contentType := "application/pdf"
		if err := s.storage.Store(ctx, s.buckets.ApplicationResumes,
			activeID+"/"+userID.String(), resume, &contentType); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
            UPDATE applications SET status='waitlisted', application=$2,
                submitted_at=$3, saved_at=$3, updated_at=$3
            WHERE id=$1
        `, id, raw, submittedAt); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, `
            INSERT INTO waitlist (hackathon_id,user_id) VALUES ($1,$2)
            ON CONFLICT (hackathon_id,user_id) DO NOTHING
        `, activeID, userID); err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
            UPDATE users SET role='applicant', preferred_email=$2, updated_at=now()
            WHERE id=$1
        `, userID, data.PreferredEmail)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &submittedAt, nil
}

func (h *handler) handleSubmitVisitorWaitlist(ctx context.Context, input *struct{}) (*SubmitApplicationOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	// TODO: refactor this to using Huma's request API instead of using the raw http package

	r := ctx.Value(middleware.RawRequestKey{}).(*http.Request)
	r.Body = http.MaxBytesReader(nil, r.Body, 12<<20)

	// Parse multipart form (10 MB max memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to parse form")
	}

	defer r.MultipartForm.RemoveAll()

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
	submission.Essay3 = r.FormValue("essay3")
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

	if _, err := io.Copy(resumeFileBuffer, io.LimitReader(resumeFile, (10<<20)+1)); err != nil {
		return nil, huma.Error500InternalServerError("Error while parsing resume")
	}

	if resumeFileBuffer.Len() > 10<<20 ||
		!bytes.HasPrefix(resumeFileBuffer.Bytes(), []byte("%PDF-")) {
		return nil, huma.Error400BadRequest("Upload a PDF resume under 10 MB")
	}
	submission.Essay1, submission.Essay2, submission.Essay3 = "", "", ""
	if err := validateVisitorSubmission(submission); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	submittedAt, err := h.applicationService.SubmitVisitorWaitlist(r.Context(), submission, resumeFileBuffer.Bytes(), userCtx.UserID)

	if err != nil {
		if errors.Is(err, ErrVisitorWaitlistEligibility) || errors.Is(err, ErrWaitlistClosed) {
			return nil, huma.Error409Conflict(err.Error())
		}
		if errors.Is(err, ErrApplicationNotOpened) {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &SubmitApplicationOutput{
		Body: SubmitApplicationResponseDto{
			SubmittedAt: submittedAt,
		},
	}, nil
}
