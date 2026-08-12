package apikeys

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
)

func RegisterRoutes(handler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "get-all-apikeys",
		Method:        http.MethodGet,
		Summary:       "Get All API Keys",
		Description:   "Returns all API keys.",
		Tags:          []string{"API Keys"},
		Path:          "/",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.getAllApiKeys)

	// TODO: create middleware for either admin or api key itself
	huma.Register(group, huma.Operation{
		OperationID:   "get-apikey",
		Method:        http.MethodGet,
		Summary:       "Get API Key",
		Description:   "Returns specified API key.",
		Tags:          []string{"API Keys"},
		Path:          "/{apiKeyId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.getApiKey)

	huma.Register(group, huma.Operation{
		OperationID:   "create-apikey",
		Method:        http.MethodPost,
		Summary:       "Create API Key",
		Description:   "Creates an API key.",
		Tags:          []string{"API Keys"},
		Path:          "/",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusCreated,
	}, handler.createApiKey)

	// TODO: create middleware for either admin or api key itself
	huma.Register(group, huma.Operation{
		OperationID:   "delete-apikey",
		Method:        http.MethodDelete,
		Summary:       "Delete API Key",
		Description:   "Deletes an API key.",
		Tags:          []string{"API Keys"},
		Path:          "/{apiKeyId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.deleteApiKey)
}