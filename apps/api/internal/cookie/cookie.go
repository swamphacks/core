package cookie

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/config"
)

func SetSessionCookie(w http.ResponseWriter, sessionID uuid.UUID, expiresAt time.Time, cfg config.CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sh_session_id",
		Value:    sessionID.String(),
		Domain:   cfg.Domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
}

func ClearSessionCookie(w http.ResponseWriter, cfg config.CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sh_session_id",
		Value:    "",
		Domain:   cfg.Domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
