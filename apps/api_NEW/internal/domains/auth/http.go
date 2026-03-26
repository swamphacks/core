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

func RegisterRoutes(authHandler *handler, group huma.API, mw *middleware.Middleware, config *config.Config) {
	// if config.AppEnv == "dev" {
	// 	huma.Register(group, huma.Operation{
	// 		OperationID: "login-with-discord",
	// 		Method:      http.MethodGet,
	// 		Summary:     "Login With Discord",
	// 		Description: "Redirects to discord oauth to login",
	// 		Tags:        []string{"Auth"},
	// 		Path:        "/login",
	// 		Middlewares: huma.Middlewares{mw.Auth.RawHTTPMiddlewareHuma},
	// 		Errors:      []int{http.StatusInternalServerError, http.StatusNotImplemented, http.StatusBadRequest, http.StatusUnauthorized},
	// 	}, authHandler.handleLogin)
	// }

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

var stateStore = map[string]bool{}

type LoginOutput struct {
	Location  string      `header:"Location"`
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (h *handler) handleLogin(ctx context.Context, input *struct{}) (*LoginOutput, error) {
	// Generate state
	state := OAuthState{
		Nonce:    "none",
		Provider: "discord",
		Redirect: "/docs",
	}

	stateJSON, _ := json.Marshal(state)

	params := url.Values{
		"client_id":     {h.config.Auth.Discord.ClientID},
		"redirect_uri":  {h.config.Auth.Discord.RedirectURI},
		"response_type": {"code"},
		"scope":         {"identify email"},
		"state":         {base64.StdEncoding.EncodeToString(stateJSON)},
	}

	resp := &LoginOutput{Location: "https://discord.com/oauth2/authorize?" + params.Encode()}
	resp.SetCookie = http.Cookie{
		Name:     "sh_auth_nonce",
		Value:    "test",
		SameSite: http.SameSiteLaxMode,
		// Path:     "/",
	}
	return resp, nil
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

	// if h.config.AppEnv == "prod" {
	// 	if input.Nonce != state.Nonce {
	// 		return nil, huma.Error401Unauthorized("Failed to authenticate. Please try again.")
	// 	}
	// }
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
