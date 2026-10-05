package nfc

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
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/parse"
)

func RegisterRoutes(nfcHandler *handler, group huma.API, mw *middleware.Middleware) {

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "GetMeals",
		Method: http.MethodGet,
		Summary: "Get all Meals from redeemables table",
		Description: "Get all Meals from redeemables table",
		Tags: []string{"NFC"},
		Path: "/redeemables/meals",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.GetMeals)

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "GetTshirts",
		Method: http.MethodGet,
		Summary: "Get all Tshirts from redeemables table",
		Description: "Get all Tshirts from redeemables table",
		Tags: []string{"NFC"},
		Path: "/redeemables/Tshirt",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.GetTshirts)

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "CheckinUser",
		Method: http.MethodPost,
		Summary: "Checkin user with NFC tag",
		Description: "Checkin user with NFC tag",
		Tags: []string{"NFC"},
		Path: "/checkin/nfc-links",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.CheckinUser)

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "TagToWorkshop",
		Method: http.MethodPost,
		Summary: "Tag user to workshop with NFC tag",
		Description: "Tag user to workshop with NFC tag",
		Tags: []string{"NFC"},
		Path: "/workshops/tag",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.TagToWorkshop)

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "TagToRedeemable",
		Method: http.MethodPost,
		Summary: "Tag user to redeemable with NFC tag",
		Description: "Tag user to redeemable with NFC tag",
		Tags: []string{"NFC"},
		Path: "/redeemables/tag",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.TagToRedeemable)

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "GetWorkshops",
		Method: http.MethodGet,
		Summary: "Get all workshops",
		Description: "Get all workshops",
		Tags: []string{"NFC"},
		Path: "/workshops/",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.GetWorkshops)

	huma.RegisterRoute(group, huma.Operation{
		OperationID: "GetSocials",
		Method: http.MethodGet,
		Summary: "Get all socials",
		Description: "Get all socials",
		Tags: []string{"NFC"},
		Path: "/workshops/socials",
		Errors: []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNotFound},
	}, nfcHandler.GetSocials)


		
}

type handler struct {
	nfcService *NfcService
	logger     *zerolog.Logger
}

func newHandler(nfcService *NfcService, logger *zerolog.Logger) *handler {
	return &handler{
		nfcService: nfcService,
		logger:     logger,
	}
}

type GetMealsOutput struct {
	Body []sqlc.getMealsRow `json:"body"`
}

type GetTshirtsOutput struct {
	Body []sqlc.getTshirtsRow `json:"body"`
}

type GetWorkshopsOutput struct {
	Body []sqlc.getWorkshopsRow `json:"body"`
}

type GetSocialsOutput struct {
	Body []sqlc.getSocialsRow `json:"body"`
}

type CheckinUserInput struct {
	TagID string `json:"tag_id"`
	UserID uuid.UUID `json:"event_id"`
}

