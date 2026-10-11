package nfc

import (
	"context"
	"net/http"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5"
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
		Tags: []string{"NFC"},
		Path: "/redeemables/meals",
		Errors: []int{http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetMeals)

	huma.Register(group, huma.Operation{
		OperationID: "GetTshirts",
		Method:      http.MethodGet,
		Summary:     "Get all Tshirts from redeemables table",
		Description: "Get all Tshirts from redeemables table",
		Tags: []string{"NFC"},
		Path: "/redeemables/Tshirt",
		Errors: []int{http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetTshirts)

	huma.Register(group, huma.Operation{
		OperationID: "CheckinUser",
		Method:      http.MethodPost,
		Summary:     "Checkin user with NFC tag",
		Description: "Checkin user with NFC tag",
		Tags: []string{"NFC"},
		Path: "/checkin/nfc-links",
		Errors: []int{http.StatusUnprocessableEntity, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleCheckinUser)

	huma.Register(group, huma.Operation{
		OperationID: "TagToWorkshop",
		Method:      http.MethodPost,
		Summary:     "Tag user to workshop with NFC tag",
		Description: "Tag user to workshop with NFC tag",
		Tags: []string{"NFC"},
		Path: "/workshops/tag",
		Errors: []int{http.StatusUnprocessableEntity, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleTagToWorkshop)

	huma.Register(group, huma.Operation{
		OperationID: "TagToRedeemable",
		Method:      http.MethodPost,
		Summary:     "Tag user to redeemable with NFC tag",
		Description: "Tag user to redeemable with NFC tag",
		Tags: []string{"NFC"},
		Path: "/redeemables/tag",
		Errors: []int{http.StatusUnprocessableEntity, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleTagToRedeemable)

	huma.Register(group, huma.Operation{
		OperationID: "GetWorkshops",
		Method:      http.MethodGet,
		Summary:     "Get all workshops",
		Description: "Get all workshops",
		Tags: []string{"NFC"},
		Path: "/workshops",
		Errors: []int{http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.handleGetWorkshops)

	huma.Register(group, huma.Operation{
		OperationID: "GetSocials",
		Method:      http.MethodGet,
		Summary:     "Get all socials",
		Description: "Get all socials",
		Tags: []string{"NFC"},
		Path: "/workshops/socials",
		Errors: []int{ http.StatusInternalServerError, http.StatusNotFound},
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
	Body struct {
		TagID  string    `json:"nfc_id"`
		UserID uuid.UUID `json:"event_id"`
	}
}

type CheckinUserOutput struct {
	Body sqlc.NfcTagsUser `json:"body"`
}

type ExpectedBody struct {
	Res bool   `json:"res"`
	Msg string `json:"msg"`
	Evidence *string `json:"evidence"`
}

type ExpectedOutput struct {
	Body ExpectedBody
}

type TagToWorkshopInput struct {
	Body struct {
		TagID      string    `json:"nfc_id"`
		WorkshopID uuid.UUID `json:"event_id"`
	}
}

type TagToWorkshopOutput struct {
	Body int64 `json:"body"`
}

type TagToRedeemableInput struct {
	Body struct {
		TagID        string    `json:"nfc_id"`
		RedeemableID uuid.UUID `json:"event_id"`
	}
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

func (h *handler) handleTagToRedeemable(ctx context.Context, input *TagToRedeemableInput) (*ExpectedOutput, error) {
	err := h.nfcService.TagToRedeemable(ctx, input.Body.TagID, input.Body.RedeemableID)
	if err != nil {
		h.logger.Err(err).Msg("failed to tag redeemable")
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "P0001":
				return &ExpectedOutput{Body: ExpectedBody{Res: false, Msg: "Hacker reached the limit for this item"}}, nil
			case "23503":
				return &ExpectedOutput{Body: ExpectedBody{Res: false, Msg: "Tag has not been checked in"}}, nil
			}
		}
		return nil, huma.Error500InternalServerError("failed to tag redeemable")
	}
	return &ExpectedOutput{Body: ExpectedBody{Res: true, Msg: "Successfully tagged user to redeemable"}}, nil
}

func (h *handler) handleTagToWorkshop(ctx context.Context, input *TagToWorkshopInput) (*ExpectedOutput, error) {
	rows, err := h.nfcService.TagToWorkshop(ctx, input.Body.TagID, input.Body.WorkshopID)
	if err != nil {
		h.logger.Err(err).Msg("failed to tag workshop")
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return &ExpectedOutput{Body: ExpectedBody{Res: false, Msg: "Tag has not been checked in"}}, nil
		}
		return nil, huma.Error500InternalServerError("failed to tag workshop")
	}
	if rows == 0 {
		return &ExpectedOutput{Body: ExpectedBody{Res: false, Msg: "Already registered for this workshop"}}, nil
	}
	return &ExpectedOutput{Body: ExpectedBody{Res: true, Msg: "Successfully tagged user to workshop"}}, nil
}

func (h *handler) handleCheckinUser(ctx context.Context, input *CheckinUserInput) (*ExpectedOutput, error) {
	user, err := h.nfcService.CheckinUser(ctx, input.Body.TagID, input.Body.UserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return &ExpectedOutput{Body: ExpectedBody{Res: false, Msg: "Tag or user already checked in"}}, nil
	}
	if err != nil {
		h.logger.Err(err).Msg("failed to check in user")
		return nil, huma.Error500InternalServerError("failed to check in user")
	}
	id := user.ID.String()
	return &ExpectedOutput{Body: ExpectedBody{Res: true, Msg: "Successfully checked in user", Evidence: &id}}, nil
}