package cookie

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/config"
)

var SessionCookieName = "sh_session_id"

var SessionCookieHumaParam *huma.Param = &huma.Param{
	Name:        SessionCookieName,
	In:          "cookie",
	Required:    true,
	Schema:      &huma.Schema{Type: "string"},
	Description: "Session cookie used to authenticate the user",
}

func SetSessionCookie(w http.ResponseWriter, sessionID uuid.UUID, expiresAt time.Time, cfg config.CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID.String(),
		Domain:   cfg.Domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
}
