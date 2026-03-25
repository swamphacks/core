package redeemables

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func RegisterRoutes(redeemablesHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "get-redeemables",
		Method:        http.MethodGet,
		Summary:       "Get Redeemables",
		Description:   "Returns a list of all redeemable items",
		Tags:          []string{"Redeemables"},
		Path:          "",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, redeemablesHandler.handleGetRedeemables)

	huma.Register(group, huma.Operation{
		OperationID:   "create-redeemable",
		Method:        http.MethodPost,
		Summary:       "Create Redeemable",
		Description:   "Creates a new redeemable item",
		Tags:          []string{"Redeemables"},
		Path:          "",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, redeemablesHandler.handleCreateRedeemable)

	huma.Register(group, huma.Operation{
		OperationID:   "update-redeemable",
		Method:        http.MethodPatch,
		Summary:       "Update Redeemable",
		Description:   "Update specific fields (name, stock, max per user) of a redeemable",
		Tags:          []string{"Redeemables"},
		Path:          "/{redeemableId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, redeemablesHandler.handleUpdateRedeemable)

	huma.Register(group, huma.Operation{
		OperationID:   "delete-redeemable",
		Method:        http.MethodDelete,
		Summary:       "Delete Redeemable",
		Description:   "Deletes a redeemable by id",
		Tags:          []string{"Redeemables"},
		Path:          "/{redeemableId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, redeemablesHandler.handleDeleteRedeemable)

	huma.Register(group, huma.Operation{
		OperationID:   "redeem-redeemable",
		Method:        http.MethodPost,
		Summary:       "Redeem Redeemable",
		Description:   "Redeems a redeemable by id. Creates a redemption record linking a specific user to a redeemable item",
		Tags:          []string{"Redeemables"},
		Path:          "/{redeemableId}/users/{userId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, redeemablesHandler.handleRedeemRedeemable)

	huma.Register(group, huma.Operation{
		OperationID:   "update-redemption",
		Method:        http.MethodPatch,
		Summary:       "Update Redemption",
		Description:   "Updates a redemption created by the user.",
		Tags:          []string{"Redeemables"},
		Path:          "/{redeemableId}/users/{userId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusBadRequest},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, redeemablesHandler.handleUpdateRedemption)
}

type handler struct {
	redeemablesService *RedeemablesService
	config             *config.Config
	logger             zerolog.Logger
}

func NewHandler(redeemablesService *RedeemablesService, config *config.Config, logger zerolog.Logger) *handler {
	return &handler{
		redeemablesService: redeemablesService,
		config:             config,
		logger:             logger,
	}
}

type GetRedeemablesOutput struct {
	Body *[]sqlc.GetRedeemablesRow
}

func (h *handler) handleGetRedeemables(ctx context.Context, input *struct{}) (*GetRedeemablesOutput, error) {
	redeemables, err := h.redeemablesService.GetRedeemables(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get redeemables")
	}

	return &GetRedeemablesOutput{Body: redeemables}, nil
}

type CreateRedeemableRequest struct {
	Name          string `json:"name" minLength:"1"`
	Amount        int    `json:"amount" minimum:"1"`
	MaxUserAmount int    `json:"max_user_amount"`
}

type CreateRedeemableOutput struct {
	Body *sqlc.Redeemable
}

func (h *handler) handleCreateRedeemable(ctx context.Context, input *struct {
	Body CreateRedeemableRequest
}) (*CreateRedeemableOutput, error) {
	redeemable, err := h.redeemablesService.CreateRedeemable(ctx, input.Body.Name, input.Body.Amount, input.Body.MaxUserAmount)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create redeemable")
	}

	return &CreateRedeemableOutput{Body: redeemable}, nil
}

type UpdateRedeemableRequest struct {
	Name          *string `json:"name,omitempty"`
	Amount        *int    `json:"total_stock,omitempty"`
	MaxUserAmount *int    `json:"max_user_amount,omitempty"`
}

type UpdateRedeemableOutput struct {
	Status int
}

func (h *handler) handleUpdateRedeemable(ctx context.Context, input *struct {
	RedeemableId string `path:"redeemableId"`
	Body         UpdateRedeemableRequest
}) (*UpdateRedeemableOutput, error) {
	redeemableId, err := uuid.Parse(input.RedeemableId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid redeemable id")
	}

	_, err = h.redeemablesService.UpdateRedeemable(ctx, redeemableId, input.Body.Name, input.Body.Amount, input.Body.MaxUserAmount)

	if err != nil {
		return nil, huma.Error500InternalServerError("Fail to update redeemable")
	}

	return &UpdateRedeemableOutput{Status: http.StatusOK}, nil
}

type DeleteRedeemableOutput struct {
	Status int
}

func (h *handler) handleDeleteRedeemable(ctx context.Context, input *struct {
	RedeemableId string `path:"redeemableId"`
}) (*DeleteRedeemableOutput, error) {
	redeemableId, err := uuid.Parse(input.RedeemableId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid redeemable id")
	}

	err = h.redeemablesService.DeleteRedeemable(ctx, redeemableId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete redeemable")
	}

	return &DeleteRedeemableOutput{Status: http.StatusOK}, nil
}

type RedeemRedeemableOutput struct {
	Status int
}

func (h *handler) handleRedeemRedeemable(ctx context.Context, input *struct {
	RedeemableId string `path:"redeemableId"`
}) (*RedeemRedeemableOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	redeemableId, err := uuid.Parse(input.RedeemableId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid redeemable id")
	}

	err = h.redeemablesService.RedeemRedeemable(ctx, redeemableId, *userId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to redeem redeemable")
	}

	return &RedeemRedeemableOutput{Status: http.StatusOK}, nil
}

type UpdateRedemptionRequest struct {
	Amount int `json:"new_amount,omitempty"`
}

type UpdateRedemptionOutput struct {
	Status int
}

func (h *handler) handleUpdateRedemption(ctx context.Context, input *struct {
	RedeemableId string `path:"redeemableId"`
	Body         UpdateRedemptionRequest
}) (*UpdateRedeemableOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	redeemableId, err := uuid.Parse(input.RedeemableId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid redeemable id")
	}

	err = h.redeemablesService.UpdateRedemption(ctx, redeemableId, *userId, input.Body.Amount)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update redemption")
	}

	return &UpdateRedeemableOutput{Status: http.StatusOK}, nil
}
