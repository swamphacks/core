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
)

/*

Will remove comments when fully done


	register for a workshop
	get a workshop
	get all workshops
	update a workshop (staff only)
	delete a workshop (staff only)
	create a workshop (staff only)
	unregister for a workshop



	ask questions about the path and what pattern they would like me to follow
*/

func RegisterRoutes(workshopHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "register-for-workshop",
		Method: http.MethodPost,
		Summary: "Register for a workshop",
		Description: "Lets a user register for a workshop.", //will workshops have a capacity?
		Tags: []string{"Workshops"},
		Path: "/{workshopID}/register", //leave this empty
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma}, //for staff stuff give mw.Auth.RequireStaffHuma
		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleRegisterForWorkshop) 

	huma.Register(group, huma.Operation{
		OperationID: "get-workshop",
		Method: http.MethodGet,
		Summary: "Get a workshop",
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
		Summary: "Get all workshops",
		Description: "Returns a list of all workshops.",
		Tags: []string{"Workshops"},
		Path: "",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	}, workshopHandler.handleGetAllWorkshops)

	// huma.Register(group, huma.Operation{
	// 	OperationID: "update-workshop",
	// 	Method: http.MethodPatch,
	// 	Summary: "Update a workshop",
	// 	Description: "Updates a workshop based on the provided id and body. Only the fields provided in the body will be updated.",
	// 	Tags: []string{"Workshops"},
	// 	Path: "/{workshopId}",
	// 	Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
	// 	Errors: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest},
	// 	Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
	// }, workshopHandler.handleUpdateWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "delete-workshop",
		Method: http.MethodDelete,
		Summary: "Delete a workshop",
		Description: "Deletes a workshop based on the provided id.",
		Tags: []string{"Workshops"},
		Path: "/{workshopId}",
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},

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
		// Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},

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
}


func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
	Title      string `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Location  string `json:"location"`
	Description string `json:"description"`
	Presenter string `json:"presenter"`
	Attendees int `json:"attendees"`
}

type CreateWorkshopInput struct {
	Title        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Presenter   string    `json:"presenter"`
}

type WorkshopIdInput struct{
	WorkshopID uuid.UUID 
}

type GetWorkshopOutput struct {
	Body OpenWorkshop
}

type GetAllWorkshopsOutput struct {
	Body []sqlc.Workshop 
}



func (h *handler) handleRegisterForWorkshop(ctx context.Context, input *WorkshopIdInput)(*GetWorkshopOutput, error) {
	userID := ctx.Value(middleware.UserContextKey).(uuid.UUID)
	workshop, err := h.workshopService.RegisterWorkshop(ctx, userID, input.WorkshopID)

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
//I am trying to follow the structure of the other files
//but for this one I think I may need some assistance
//I am very new with go so I am not sure if here I should return the new user?
//or just the workshop? Idea i have for now is registerForWorkshop would update attendee number

func (h *handler) handleUnregisterForWorkshop(ctx context.Context, input *WorkshopIdInput) (*DeleteWorkshopOutput, error) {
	userID := ctx.Value(middleware.UserContextKey).(uuid.UUID)
	err := h.workshopService.UnregisterWorkshop(ctx, userID, input.WorkshopID)

	if err != nil {
		if errors.Is(err, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Couldn't unregister for workshop")
	}

	return &DeleteWorkshopOutput{Status: http.StatusOK}, nil
}

func (h *handler) handleGetWorkshop(ctx context.Context, input *WorkshopIdInput) (*GetWorkshopOutput, error) {
	workshop, err := h.workshopService.GetWorkshop(ctx, input.WorkshopID)

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

// func (h *handler) handleUpdateWorkshop(ctx context.Context, input *WorkshopIdInput) (*GetWorkshopOutput, error) {
// 	workshop, err := h.workshopService.UpdateWorkshop(ctx, input.WorkshopID)

// 	if err != nil {
// 		if errors.Is(err, database.ErrEntityNotFound) {
// 			return nil, huma.Error404NotFound("Workshop not found")
// 		}
// 		return nil, huma.Error500InternalServerError("Failed to update workshop")
// 	}

// 	return &GetWorkshopOutput{Body: OpenWorkshop{
// 		ID: workshop.ID.String(),
// 		Title: workshop.Title,
// 		StartTime: workshop.StartTime,
// 		EndTime: workshop.EndTime,
// 		Location: deref(workshop.Location),
// 		Description: deref(workshop.Description),
// 		Presenter: deref(workshop.Presenter),
// 		Attendees: int(workshop.CurrAttendees),
// 	}}, nil
// }

type DeleteWorkshopOutput struct {
	Status int
}

func (h *handler) handleDeleteWorkshop(ctx context.Context, input *WorkshopIdInput) (*DeleteWorkshopOutput, error) {
	err := h.workshopService.DeleteWorkshop(ctx, input.WorkshopID)

	if err != nil {
		if errors.Is(err, database.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Workshop not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete workshop")
	}

	return &DeleteWorkshopOutput{Status: http.StatusOK}, nil
}

func (h *handler) handleCreateWorkshop(ctx context.Context, input *CreateWorkshopInput) (*GetWorkshopOutput, error) {

	params := sqlc.CreateWorkshopParams{
		Title:       input.Title,
		Description: &input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Location:    &input.Location,
		Presenter:   &input.Presenter,
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

