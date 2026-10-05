package nfc

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func RegisterRoutes(nfcHandler *handler, group huma.API, mw *middleware.Middleware) {

	huma.Register(group, huma.Operation{
		OperationID: "GetMeals",
		Method:      http.MethodGet,
		Summary:     "Get all Meals from redeemables table",
		Description: "Get all Meals from redeemables table",
		Tags:        []string{"NFC"},
		Path:        "/redeemables/meals",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetMeals)

	huma.Register(group, huma.Operation{
		OperationID: "GetTshirts",
		Method:      http.MethodGet,
		Summary:     "Get all Tshirts from redeemables table",
		Description: "Get all Tshirts from redeemables table",
		Tags:        []string{"NFC"},
		Path:        "/redeemables/Tshirt",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetTshirts)

	huma.Register(group, huma.Operation{
		OperationID: "CheckinUser",
		Method:      http.MethodPost,
		Summary:     "Checkin user with NFC tag",
		Description: "Checkin user with NFC tag",
		Tags:        []string{"NFC"},
		Path:        "/checkin/nfc-links",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleCheckinUser)

	huma.Register(group, huma.Operation{
		OperationID: "TagToWorkshop",
		Method:      http.MethodPost,
		Summary:     "Tag user to workshop with NFC tag",
		Description: "Tag user to workshop with NFC tag",
		Tags:        []string{"NFC"},
		Path:        "/workshops/tag",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleTagToWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "TagToRedeemable",
		Method:      http.MethodPost,
		Summary:     "Tag user to redeemable with NFC tag",
		Description: "Tag user to redeemable with NFC tag",
		Tags:        []string{"NFC"},
		Path:        "/redeemables/tag",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleTagToRedeemable)

	huma.Register(group, huma.Operation{
		OperationID: "GetWorkshops",
		Method:      http.MethodGet,
		Summary:     "Get all workshops",
		Description: "Get all workshops",
		Tags:        []string{"NFC"},
		Path:        "/workshops/",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetWorkshops)

	huma.Register(group, huma.Operation{
		OperationID: "GetSocials",
		Method:      http.MethodGet,
		Summary:     "Get all socials",
		Description: "Get all socials",
		Tags:        []string{"NFC"},
		Path:        "/workshops/socials",
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetSocials)

}

type handler struct {
	nfcService *NfcService
	logger     zerolog.Logger
}

func NewHandler(nfcService *NfcService, logger zerolog.Logger) *handler {
	return &handler{
		nfcService: nfcService,
		logger:     logger,
	}
}

type GetMealsOutput struct {
	Body []sqlc.GetMealsRow `json:"body"`
}

type GetTshirtsOutput struct {
	Body []sqlc.GetTshirtsRow `json:"body"`
}

type GetWorkshopsOutput struct {
	Body []sqlc.GetWorkshopsRow `json:"body"`
}

type GetSocialsOutput struct {
	Body []sqlc.GetSocialsRow `json:"body"`
}

type CheckinUserInput struct {
	TagID  string    `json:"tag_id"`
	UserID uuid.UUID `json:"event_id"`
}

type CheckinUserOutput struct {
	Body sqlc.NfcTagsUser `json:"body"`
}

type expectedOutput struct {
	Res bool   `json:"res"`
	Msg string `json:"msg"`
}

type TagToWorkshopInput struct {
	TagID      string    `json:"tag_id"`
	WorkshopID uuid.UUID `json:"workshop_id"`
}

type TagToWorkshopOutput struct {
	Body int64 `json:"body"`
}

type TagToRedeemableInput struct {
	TagID        string    `json:"tag_id"`
	RedeemableID uuid.UUID `json:"redeemable_id"`
}

type TagToRedeemableOutput struct {
	Body string `json:"body"`
}

func (h *handler) handleGetMeals(ctx context.Context, input *struct{}) (*GetMealsOutput, error) {
	meals, err := h.nfcService.GetMeals(ctx)
	if err != nil {
		return nil, err
	}

	return &GetMealsOutput{
		Body: meals,
	}, nil
}

func (h *handler) handleGetTshirts(ctx context.Context, input *struct{}) (*GetTshirtsOutput, error) {
	tshirts, err := h.nfcService.GetTshirts(ctx)
	if err != nil {
		return nil, err
	}

	return &GetTshirtsOutput{
		Body: tshirts,
	}, nil
}

func (h *handler) handleGetWorkshops(ctx context.Context, input *struct{}) (*GetWorkshopsOutput, error) {
	workshops, err := h.nfcService.GetWorkshops(ctx)
	if err != nil {
		return nil, err
	}

	return &GetWorkshopsOutput{
		Body: workshops,
	}, nil
}

func (h *handler) handleGetSocials(ctx context.Context, input *struct{}) (*GetSocialsOutput, error) {
	socials, err := h.nfcService.GetSocials(ctx)
	if err != nil {
		return nil, err
	}

	return &GetSocialsOutput{
		Body: socials,
	}, nil
}

func (h *handler) handleCheckinUser(ctx context.Context, input *CheckinUserInput) (*CheckinUserOutput, error) {
	user, err := h.nfcService.CheckinUser(ctx, sqlc.CheckinUserParams{
		TagID:  input.TagID,
		UserID: &input.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &CheckinUserOutput{
		Body: *user,
	}, nil
}

func (h *handler) handleTagToWorkshop(ctx context.Context, input *TagToWorkshopInput) (*TagToWorkshopOutput, error) {
	workshop, err := h.nfcService.TagToWorkshop(ctx, sqlc.TagToWorkshopParams{
		TagID:      input.TagID,
		WorkshopID: input.WorkshopID,
	})
	if err != nil {
		return nil, err
	}

	return &TagToWorkshopOutput{
		Body: workshop,
	}, nil
}

func (h *handler) handleTagToRedeemable(ctx context.Context, input *TagToRedeemableInput) (*TagToRedeemableOutput, error) {
	err := h.nfcService.TagToRedeemable(ctx, sqlc.TagToRedeemableParams{
		TagID:        input.TagID,
		RedeemableID: input.RedeemableID,
	})
	if err != nil {
		return nil, err
	}

	return &TagToRedeemableOutput{
		Body: "Successfully tagged user to redeemable",
	}, nil
}
