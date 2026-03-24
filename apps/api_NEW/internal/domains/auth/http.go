package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
)

func RegisterRoutes(authHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-me",
		Method:      http.MethodGet,
		Summary:     "Get Me",
		Description: "Returns the authenticated user's profile",
		Tags:        []string{"Auth"},
		Path:        "/me",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized},
		Parameters: []*huma.Param{
			{
				Name:        authHandler.config.Cookie.SessionName,
				In:          "cookie",
				Required:    true,
				Schema:      &huma.Schema{Type: "string"},
				Description: "Session cookie used to authenticate the user",
			},
		},
	}, authHandler.handleGetMe)

	huma.Register(group, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Summary:     "Logout",
		Description: "Logs out the authenticated user by invalidating their session",
		Tags:        []string{"Auth"},
		Path:        "/logout",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
	}, authHandler.handleLogout)

	huma.Register(group, huma.Operation{
		OperationID: "oauth-callback",
		Method:      http.MethodGet,
		Summary:     "OAuth Callback",
		Description: "Handles the OAuth provider callback, validates state and nonce, and sets the session cookie.",
		Tags:        []string{"Auth"},
		Path:        "/callback",
		Middlewares: huma.Middlewares{mw.Auth.RawHTTPMiddlewareHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusNotImplemented, http.StatusBadRequest, http.StatusUnauthorized},
	}, authHandler.handleOAuthCallback)
}

type handler struct {
	authService *AuthService
	config      *config.Config
	logger      zerolog.Logger
}

func NewHandler(authService *AuthService, config *config.Config, logger zerolog.Logger) *handler {
	return &handler{
		authService: authService,
		config:      config,
		logger:      logger.With().Str("handler", "AuthHandler").Str("domain", "auth").Logger(),
	}
}

type GetMeOutput struct {
	Body *middleware.UserContext
}

func (h *handler) handleGetMe(ctx context.Context, input *struct{}) (*GetMeOutput, error) {
	user, err := h.authService.GetMe(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Your profile could not be loaded.")
	}

	return &GetMeOutput{Body: user}, nil
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (h *handler) handleLogout(ctx context.Context, input *struct{}) (*LogoutOutput, error) {
	err := h.authService.Logout(ctx)

	if err != nil {
		if errors.Is(err, ErrFetchSessionContextFailed) {
			return nil, huma.Error401Unauthorized("Not authorized.")
		} else {
			return nil, huma.Error500InternalServerError("Something went wrong while logging out.")
		}
	}

	res := &LogoutOutput{
		SetCookie: http.Cookie{
			Name:     h.config.Cookie.SessionName,
			Value:    "",
			Domain:   h.config.Cookie.Domain,
			Path:     "/",
			HttpOnly: true,
			Secure:   h.config.Cookie.Secure,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		},
	}

	return res, nil
}

type OAuthState struct {
	Nonce    string `json:"nonce"`
	Provider string `json:"provider"`
	Redirect string `json:"redirect"`
}

type OAuthCallbackOutput struct {
	SetCookie   []http.Cookie `header:"Set-Cookie"`
	RedirectUrl string        `header:"Location"`
	Status      int
}

func (h *handler) handleOAuthCallback(ctx context.Context, input *struct {
	Code      string `query:"code" required:"true" doc:"OAuth authorization code"`
	State     string `query:"state" required:"true" doc:"Base64 encoded OAuth state"`
	Nonce     string `cookie:"sh_auth_nonce" required:"true" doc:"Auth nonce cookie for CSRF protection"`
	UserAgent string `header:"User-Agent" doc:"Client user agent"`
}) (*OAuthCallbackOutput, error) {
	r := ctx.Value(middleware.RawRequestKey{}).(*http.Request)

	var ipAddress *string
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && ip != "" {
		ipAddress = &ip
	}

	if input.Code == "" || input.State == "" {
		return nil, huma.Error400BadRequest("Invalid callback. Please try again.")
	}

	decodedStateBytes, err := base64.URLEncoding.DecodeString(input.State)
	if err != nil {
		return nil, huma.Error400BadRequest("This callback was invalid. Please try again.")
	}

	var state OAuthState
	if err := json.Unmarshal(decodedStateBytes, &state); err != nil {
		return nil, huma.Error400BadRequest("This callback was invalid. Please try again.")
	}

	if input.Nonce != state.Nonce {
		return nil, huma.Error401Unauthorized("Failed to authenticate. Please try again.")
	}

	session, err := h.authService.AuthenticateWithOAuth(ctx, input.Code, state.Provider, ipAddress, &input.UserAgent)
	if err != nil {
		switch err {
		case ErrProviderUnsupported:
			return nil, huma.Error501NotImplemented("This provider is not supported.")
		case ErrAuthenticationFailed:
			return nil, huma.Error401Unauthorized("Failed to authenticate the user.")
		default:
			h.logger.Err(err).Msg("Something unexpected happened.")
			return nil, huma.Error500InternalServerError("Something went wrong")
		}
	}

	if isURL(state.Redirect) {
		return nil, huma.Error400BadRequest("invalid redirect path")
	}

	redirectPath := ensureLeadingSlash(state.Redirect)

	res := &OAuthCallbackOutput{
		SetCookie: []http.Cookie{
			{
				Name:     h.config.Cookie.SessionName,
				Value:    session.ID.String(),
				Domain:   h.config.Cookie.Domain,
				Path:     "/",
				HttpOnly: true,
				Secure:   h.config.Cookie.Secure,
				SameSite: http.SameSiteLaxMode,
				Expires:  session.ExpiresAt,
			},

			{
				Name:     "sh_auth_nonce",
				Value:    "",
				Domain:   h.config.Cookie.Domain,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
			},
		},

		RedirectUrl: h.config.ClientUrl + redirectPath,
		Status:      http.StatusSeeOther,
	}

	return res, nil
}

func ensureLeadingSlash(s string) string {
	if len(s) == 0 || s[0] != '/' {
		return "/" + s
	}
	return s
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}
