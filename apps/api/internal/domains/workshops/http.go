package workshops

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"

)


func RegisterRoutes(workshopHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "register-for-workshop",
		Method: http.MethodPost,
		Summary: "Register for a workshop",
		Description: "Lets a user register for a workshop.", 
		Tags: []string{"Workshops"},
		Path: "/{workshopId}/register", 
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma}, 
		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleRegisterForWorkshop) 

	huma.Register(group, huma.Operation{
		OperationID: "get-workshop",
		Method: http.MethodGet,
		Summary: "Get a workshop off of workshopID",
		Description: "Returns a workshop based on the provided id.",
		Tags: []string{"Workshops"},
		Path: "/{workshopId}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleGetWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "get-all-workshops",
		Method: http.MethodGet,
		Summary: "Get all UPCOMING workshops",
		Description: "Returns a list of all workshops.",
		Tags: []string{"Workshops"},
		Path: "",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleGetAllWorkshops)

	huma.Register(group, huma.Operation{
		OperationID: "update-workshop",
		Method: http.MethodPatch,
		Summary: "Update a workshop",
		Description: "Updates a workshop based on the provided id and body. Only the fields provided in the body will be updated.",
		Tags: []string{"Workshops"},
		Path: "/{workshopId}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},

		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleUpdateWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "delete-workshop",
		Method: http.MethodDelete,
		Summary: "Delete a workshop",
		Description: "Deletes a workshop based on the provided id.",
		Tags: []string{"Workshops"},
		Path: "/{workshopId}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},

		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleDeleteWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "create-workshop",
		Method: http.MethodPost,
		Summary: "Create a workshop",
		Description: "Creates a workshop based on the provided body.",
		Tags: []string{"Workshops"},
		Path: "",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},

		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleCreateWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "unregister-for-workshop",
		Method: http.MethodDelete,
		Summary: "Unregister for a workshop",
		Description: "Lets a user unregister for a workshop.",
		Tags: []string{"Workshops"},
		Path: "/{workshopId}/register",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleUnregisterForWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "view-all-workshops",
		Method: http.MethodGet,
		Summary: "View all workshops ever made",
		Description: "Returns a list of all workshops, including past workshops.",
		Tags: []string{"Workshops"},
		Path: "/view-all",
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},

		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleViewAllWorkshops)

	huma.Register(group, huma.Operation{
		OperationID: "delete-all-workshops",
		Method: http.MethodDelete,
		Summary: "Delete all workshops",
		Description: "Deletes all the workshops",
		Tags: []string{"Workshops"},
		Path: "/delete-all",
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},

		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleDeleteAllWorkshops)
}


func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func coalesceString(input, fallback string) string {
    if input == "" {
        return fallback
    }
    return input
}

func coalesceTime(input, fallback time.Time) time.Time {
    if input.IsZero() {
        return fallback
    }
    return input
}

type handler struct {
	workshopService *WorkshopService
	logger zerolog.Logger
}

func NewHandler(workshopService *WorkshopService, logger zerolog.Logger) *handler {
	
	return &handler{
		workshopService: workshopService,
		logger: logger,
	}
}

type OpenWorkshop struct {
	ID		  string `json:"id"`
	Title      string `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Location  string `json:"location"`
	Description string `json:"description"`
	Presenter string `json:"presenter"`
	Attendees int `json:"attendees"`
}

type UpdateWorkshopInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Location    *string    `json:"location"`
	Presenter   *string    `json:"presenter"`
}

type CreateWorkshopInput struct {
	Title        string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Presenter   string    `json:"presenter"`
}


type GetWorkshopOutput struct {
	Body OpenWorkshop
}

type GetAllWorkshopsOutput struct {
	Body []sqlc.Workshop 
}

type DeleteWorkshopOutput struct {
	Status int
}


func (h *handler) handleRegisterForWorkshop(ctx context.Context, input *struct {
	WorkshopID string `path:"workshopId"`
})(*GetWorkshopOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}
	
	workshopID, err := uuid.Parse(input.WorkshopID)


	if err != nil {
		return nil, huma.Error400BadRequest("Invalid workshop id")
	}


	workshop, err := h.workshopService.RegisterWorkshop(ctx, userCtx.UserID, workshopID)

	if err != nil {
		return nil, huma.Error400BadRequest("Couldn't register for workshop")
	}

	return &GetWorkshopOutput{Body: OpenWorkshop{
		ID: workshop.ID.String(),
		Title: workshop.Title,
		StartTime: workshop.StartTime,
		EndTime: workshop.EndTime,
		Location: deref(workshop.Location),
		Description: deref(workshop.Description),
		Presenter: deref(workshop.Presenter),
		Attendees: int(workshop.CurrAttendees),
	}}, nil
}


func (h *handler) handleUnregisterForWorkshop(ctx context.Context, input *struct {
	WorkshopID string `path:"workshopId"`
}) (*GetWorkshopOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	workshopID, err := uuid.Parse(input.WorkshopID)


	if err != nil {
		return nil, huma.Error400BadRequest("Invalid workshop id")
	}

	err = h.workshopService.UnregisterWorkshop(ctx, userCtx.UserID, workshopID)

	if err != nil {
		if errors.Is(err, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Couldn't unregister for workshop")
	}

	workshop, repoErr := h.workshopService.GetWorkshop(ctx, workshopID)

	if repoErr != nil {
		return nil, huma.Error500InternalServerError("Failed to get workshop")
	}

	return &GetWorkshopOutput{Body: OpenWorkshop{
		ID: workshop.ID.String(),
		Title: workshop.Title,
		StartTime: workshop.StartTime,
		EndTime: workshop.EndTime,
		Location: deref(workshop.Location),
		Description: deref(workshop.Description),
		Presenter: deref(workshop.Presenter),
		Attendees: int(workshop.CurrAttendees),
	}}, nil
}

func (h *handler) handleGetWorkshop(ctx context.Context, input *struct {
	WorkshopID string `path:"workshopId"`
}) (*GetWorkshopOutput, error) {
	workshopID, err := uuid.Parse(input.WorkshopID)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid workshop id")
	}

	workshop, err := h.workshopService.GetWorkshop(ctx, workshopID)

	if err != nil {
		if errors.Is(err, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get workshop")
	}

	return &GetWorkshopOutput{Body: OpenWorkshop{
		ID: workshop.ID.String(),
		Title: workshop.Title,
		StartTime: workshop.StartTime,
		EndTime: workshop.EndTime,
		Location: deref(workshop.Location),
		Description: deref(workshop.Description),
		Presenter: deref(workshop.Presenter),
		Attendees: int(workshop.CurrAttendees),
	}}, nil
}

func (h *handler) handleGetAllWorkshops(ctx context.Context, input *struct{}) (*GetAllWorkshopsOutput, error) {
	workshops, err := h.workshopService.GetAllWorkshops(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get workshops")
	}

	return &GetAllWorkshopsOutput{Body: workshops}, nil
}

func (h *handler) handleViewAllWorkshops(ctx context.Context, input *struct{}) (*GetAllWorkshopsOutput, error) {
	workshops, err := h.workshopService.ViewAllWorkshops(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get workshops")
	}

	return &GetAllWorkshopsOutput{Body: workshops}, nil
}

func (h *handler) handleUpdateWorkshop(ctx context.Context, input *struct{
	WorkshopID string `path:"workshopId"`
	Body UpdateWorkshopInput
}) (*GetWorkshopOutput, error) {
	
	workshopID, err := uuid.Parse(input.WorkshopID)


	if err != nil {
		return nil, huma.Error400BadRequest("Invalid workshop id")
	}

	currWorkshop, getErr := h.workshopService.GetWorkshop(ctx, workshopID)

	if getErr != nil {
		if errors.Is(getErr, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get workshop")
	}


	title := coalesceString(deref(input.Body.Title), currWorkshop.Title)
	desc := coalesceString(deref(input.Body.Description), deref(currWorkshop.Description))
	startTime := coalesceTime(derefTime(input.Body.StartTime), currWorkshop.StartTime)
	endTime := coalesceTime(derefTime(input.Body.EndTime), currWorkshop.EndTime)
	loc := coalesceString(deref(input.Body.Location), deref(currWorkshop.Location))
	pres := coalesceString(deref(input.Body.Presenter), deref(currWorkshop.Presenter))

	params := sqlc.UpdateWorkshopParams{
		Title:       title,
		Description: &desc,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    &loc,
		Presenter:   &pres,
		WorkshopID: workshopID,
	}

	workshop, repoErr := h.workshopService.UpdateWorkshop(ctx, params)

	if repoErr != nil {
		if errors.Is(repoErr, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Failed to update workshop")
	}

	return &GetWorkshopOutput{Body: OpenWorkshop{
		ID: workshop.ID.String(),
		Title: workshop.Title,
		StartTime: workshop.StartTime,
		EndTime: workshop.EndTime,
		Location: deref(workshop.Location),
		Description: deref(workshop.Description),
		Presenter: deref(workshop.Presenter),
		Attendees: int(workshop.CurrAttendees),
	}}, nil
}


func (h *handler) handleDeleteWorkshop(ctx context.Context, input *struct {
	WorkshopID string `path:"workshopId"`
}) (*DeleteWorkshopOutput, error) {
	workshopID, err := uuid.Parse(input.WorkshopID)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid workshop id")
	}

	err = h.workshopService.DeleteWorkshop(ctx, workshopID)

	if err != nil {
		if errors.Is(err, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete workshop")
	}

	return &DeleteWorkshopOutput{Status: http.StatusOK}, nil
}

func (h *handler) handleDeleteAllWorkshops(ctx context.Context, input *struct{}) (*DeleteWorkshopOutput, error) {

	err := h.workshopService.DeleteAllWorkshops(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete workshop")
	}

	return &DeleteWorkshopOutput{Status: http.StatusOK}, nil
}

func (h *handler) handleCreateWorkshop(ctx context.Context, input *struct{
	Body CreateWorkshopInput
}) (*GetWorkshopOutput, error) {

	params := sqlc.CreateWorkshopParams{
		Title:       input.Body.Title,
		Description: &input.Body.Description,
		StartTime:   input.Body.StartTime,
		EndTime:     input.Body.EndTime,
		Location:    &input.Body.Location,
		Presenter:   &input.Body.Presenter,
	}

	workshop, err := h.workshopService.CreateWorkshop(ctx, params)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create workshop")
	}

	return &GetWorkshopOutput{Body: OpenWorkshop{
		ID: workshop.ID.String(),
		Title: workshop.Title,
		StartTime: workshop.StartTime,
		EndTime: workshop.EndTime,
		Location: deref(workshop.Location),
		Description: deref(workshop.Description),
		Presenter: deref(workshop.Presenter),
		Attendees: int(workshop.CurrAttendees),
	}}, nil
}

