package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/swamphacks/core/apps/api/internal/api/handlers"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

type API struct {
	Router     *chi.Mux
	Logger     *zerolog.Logger
	Handlers   *handlers.Handlers
	Middleware *mw.Middleware
}

func NewAPI(logger *zerolog.Logger, handlers *handlers.Handlers, middleware *mw.Middleware) *API {
	api := &API{
		Router:     chi.NewRouter(),
		Logger:     logger,
		Handlers:   handlers,
		Middleware: middleware,
	}

	api.setupRoutes(api.Middleware)

	return api
}

func (api *API) setupRoutes(mw *mw.Middleware) {
	var (
		// Both requireXXRole functions automatically allow superusers
		ensureSuperuser  = mw.Auth.RequirePlatformRole([]sqlc.AuthUserRole{sqlc.AuthUserRoleSuperuser})
		ensureEventAdmin = mw.Event.RequireEventRole([]sqlc.EventRoleType{sqlc.EventRoleTypeAdmin})
		// Event Admins are technically Staff...
		ensureEventStaff = mw.Event.RequireEventRole([]sqlc.EventRoleType{sqlc.EventRoleTypeAdmin, sqlc.EventRoleTypeStaff})
	)

	AllowedOrigins := config.Load().AllowedOrigins

	api.Router.Use(middleware.Logger)
	api.Router.Use(middleware.RealIP)
	api.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api.Router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "SwampHacks API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

	// Health check
	api.Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		api.Logger.Trace().Str("method", r.Method).Str("path", r.URL.Path).Msg("Received ping.")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "6")
		if _, err := w.Write([]byte("pong!\n")); err != nil {
			log.Err(err)
		}
	})

	// --- Auth routes ---
	api.Router.Route("/auth", func(r chi.Router) {
		r.Get("/callback", api.Handlers.Auth.OAuthCallback)

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequireAuth)
			r.Get("/me", api.Handlers.Auth.GetMe)
			r.Post("/logout", api.Handlers.Auth.Logout)
		})
	})

	// --- User routes ---
	api.Router.Route("/users", func(r chi.Router) {
		r.Use(mw.Auth.RequireAuth)
		r.Get("/", api.Handlers.User.GetUsers)
		r.Get("/me", api.Handlers.User.GetProfile)
		r.Patch("/me", api.Handlers.User.UpdateUser)
		r.Patch("/me/email-consent", api.Handlers.User.UpdateEmailConsent)
		r.Patch("/me/onboarding", api.Handlers.User.CompleteOnboarding)
	})

	// --- Team routes (non Event specific) ---
	api.Router.Route("/teams", func(r chi.Router) {
		r.Use(mw.Auth.RequireAuth)
		r.Get("/{teamId}", api.Handlers.Teams.GetTeam)
		r.Get("/{teamId}/pending-joins", api.Handlers.Teams.GetPendingRequestsForTeam)
		r.Delete("/{teamId}/members/me", api.Handlers.Teams.LeaveTeam)
		r.Delete("/{teamId}/members/{userId}", api.Handlers.Teams.KickMemberFromTeam)
		r.Post("/join/{requestId}/accept", api.Handlers.Teams.AcceptTeamJoinRequest)
		r.Post("/join/{requestId}/reject", api.Handlers.Teams.RejectTeamJoinRequest)
	})

	// --- Discord routes (for Discord bot) ---
	api.Router.Route("/discord", func(r chi.Router) {
		r.Use(mw.Auth.RequireAuth)
		r.Get("/event/{event_id}/attendees", api.Handlers.Discord.GetEventAttendeesWithDiscord)
	})

	// --- Event routes ---
	api.Router.Route("/events", func(r chi.Router) {
		// Superuser-only
		r.With(mw.Auth.RequireAuth, ensureSuperuser).Post("/", api.Handlers.Event.CreateEvent)

		// Authenticated
		r.With(mw.Auth.RequireAuth).Get("/", api.Handlers.Event.GetEvents)

		r.Post("/{eventId}/interest", api.Handlers.EventInterest.AddEmailToEvent) // Unprotected

		// Event-specific routes
		r.Route("/{eventId}", func(r chi.Router) {
			r.Use(mw.Auth.RequireAuth) // routes below this are protected

			r.Get("/", api.Handlers.Event.GetEventByID)
			r.Get("/role", api.Handlers.Event.GetEventRole)
			r.Get("/discord/{discordId}", api.Handlers.Discord.GetUserEventRoleByDiscordIDAndEventId)

			r.With(ensureEventStaff).Get("/overview", api.Handlers.Event.GetEventOverview)

			// Check in and event day of routes (STAFF+)
			r.With(ensureEventStaff).Post("/checkin", api.Handlers.Admission.HandleEventCheckIn)
			// Used to fetch user info for checking in
			r.With(ensureEventStaff).Get("/users/{userId}", api.Handlers.Event.GetUserForEvent)
			// Get user ID by RFID
			r.With(ensureEventStaff).Get("/users/by-rfid/{rfid}", api.Handlers.Event.GetUserByRFID)

			// Admin-only
			r.With(ensureEventAdmin).Post("/queue-confirmation-email", api.Handlers.Email.QueueConfirmationEmail)
			r.With(ensureEventAdmin).Post("/queue-welcome-email", api.Handlers.Email.QueueWelcomeEmail)
			r.With(ensureEventAdmin).Post("/send-welcome-emails", api.Handlers.Bat.SendWelcomeEmails)
			r.With(ensureEventAdmin).Post("/calc-admissions", api.Handlers.Admission.HandleCalculateAdmissionsRequest)
			r.With(ensureEventAdmin).Patch("/transition-waitlisted-applications", api.Handlers.Application.TransitionWaitlistedApplications)
			r.With(ensureEventAdmin).Post("/begin-waitlist-transition", api.Handlers.Bat.QueueScheduleWaitlistTransitionTask)
			r.With(ensureEventAdmin).Post("/shutdown-waitlist-scheduler", api.Handlers.Bat.QueueShutdownWaitlistSchedulerTask)
			r.With(ensureEventAdmin).Post("/reviews/bat-runs/{runId}/release", api.Handlers.Admission.ReleaseDecisions)
			r.With(ensureEventAdmin).Patch("/", api.Handlers.Event.UpdateEventById)
			r.With(ensureEventAdmin).Post("/banner", api.Handlers.Event.UploadEventBanner)
			r.With(ensureEventAdmin).Delete("/banner", api.Handlers.Event.DeleteBanner)
			r.With(ensureEventAdmin).Get("/staff", api.Handlers.Event.GetEventStaffUsers)
			r.With(ensureEventStaff).Get("/users", api.Handlers.Event.GetEventUsers)
			r.With(ensureEventAdmin).Post("/roles", api.Handlers.Event.AssignEventRole)
			r.With(ensureEventAdmin).Delete("/roles/{userId}", api.Handlers.Event.RevokeEventRole)
			r.With(ensureEventAdmin).Post("/roles/batch", api.Handlers.Event.BatchAssignEventRoles)
			r.With(ensureEventAdmin).Get("/bat-runs", api.Handlers.Bat.GetRunsByEventId)
			r.With(ensureEventAdmin).Delete("/bat-runs", api.Handlers.Bat.GetRunsByEventId)
			r.With(ensureEventAdmin).Get("/review-status", api.Handlers.Bat.CheckApplicationReviewsComplete)
			r.With(ensureEventAdmin).Post("/reviews/bat-runs", api.Handlers.Admission.HandleCalculateAdmissionsRequest)
			r.With(ensureEventAdmin).Post("/reviews/bat-runs/{runId}/release", api.Handlers.Admission.ReleaseDecisions)

			// Superuser-only
			r.With(ensureSuperuser).Delete("/", api.Handlers.Event.DeleteEventById)

			// Application routes
			r.Route("/application", func(r chi.Router) {
				r.Use(mw.Auth.RequireAuth)
				r.Get("/", api.Handlers.Application.GetMyApplication)
				r.Post("/submit", api.Handlers.Application.SubmitApplication)
				r.Post("/save", api.Handlers.Application.SaveApplication)
				r.Get("/download-resume", api.Handlers.Application.DownloadResume)
				r.With(mw.Event.AttachEventRoleToContext()).Get("/{applicationId}/resume", api.Handlers.Application.GetResumePresignedUrl)

				// Getting a resume (Staff Only)
				r.With(ensureEventStaff).Get("/{applicationId}", api.Handlers.Application.GetApplication)

				// For statistics (Staff ONLY)
				r.With(ensureEventStaff).Get("/stats", api.Handlers.Application.GetApplicationStatistics)

				// For application review (Staff ONLY)
				r.With(ensureEventStaff).Get("/assigned", api.Handlers.Application.GetAssignedApplications)
				r.With(ensureEventStaff).Post("/{applicationId}/review", api.Handlers.Application.SubmitApplicationReview)

				// Review admin routes (For Event Admins only)
				r.With(ensureEventAdmin).Post("/reset-reviews", api.Handlers.Application.ResetApplicationReviews)
				r.With(ensureEventAdmin).Post("/assign-reviewers", api.Handlers.Application.AssignApplicationReviewers)

				//withdraw Acceptance
				r.Patch("/withdraw-acceptance", api.Handlers.Application.WithdrawAcceptance)

				//Accept acceptance
				r.Patch("/accept-acceptance", api.Handlers.Application.AcceptApplicationAcceptance)

				//Withdraw attendance
				r.Patch("/withdraw-attendance", api.Handlers.Application.WithdrawAttendance)

				//Waitlist application
				r.Patch("/join-waitlist", api.Handlers.Application.JoinWaitlist)
			})

			r.Route("/redeemables", func(r chi.Router) {
				r.Use(ensureEventStaff)
				// Get all redeemables and create new redeemable
				// eventId is available from parent route context
				r.Get("/", api.Handlers.Redeemables.GetRedeemables)
				r.Post("/", api.Handlers.Redeemables.CreateRedeemable)

				// Update and delete specific redeemable
				r.Route("/{redeemableId}", func(r chi.Router) {
					r.Patch("/", api.Handlers.Redeemables.UpdateRedeemable)
					r.Delete("/", api.Handlers.Redeemables.DeleteRedeemable)

					r.Route("/users/{userId}", func(r chi.Router) {
						r.Post("/", api.Handlers.Redeemables.RedeemRedeemable)
						r.Patch("/", api.Handlers.Redeemables.UpdateRedemption)
					})
				})
			})

			// Team routes
			r.Route("/teams", func(r chi.Router) {
				r.Post("/", api.Handlers.Teams.CreateTeam)
				r.Get("/", api.Handlers.Teams.GetEventTeams)
				r.Get("/me", api.Handlers.Teams.GetMyTeam)
				r.Get("/me/pending-joins", api.Handlers.Teams.GetMyPendingRequests)

				// Specific team routes within events
				r.Route("/{teamId}", func(r chi.Router) {
					r.Post("/join", api.Handlers.Teams.RequestToJoinTeam)
				})
			})
		})
	})

	// Protected test routes
	api.Router.Route("/protected", func(r chi.Router) {
		r.Use(mw.Auth.RequireAuth)

		r.Get("/basic", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Welcome, arbitrarily roled user!\n")); err != nil {
				log.Err(err)
			}
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequirePlatformRole([]sqlc.AuthUserRole{sqlc.AuthUserRoleUser}))
			r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Welcome, user!\n")); err != nil {
					log.Err(err)
				}
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequirePlatformRole([]sqlc.AuthUserRole{sqlc.AuthUserRoleSuperuser}))
			r.Get("/superuser", func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Welcome, superuser!\n")); err != nil {
					log.Err(err)
				}
			})
		})
	})
}
