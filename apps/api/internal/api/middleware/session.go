package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/swamphacks/core/apps/api/internal/db"
)

type contextKey string

const (
	sessionCookieName = "sh_session"
)

func SessionMiddleware(rdb *redis.Client, db *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || strings.TrimSpace(cookie.Value) == "" {
				http.Error(w, "unauthorized: missing session", http.StatusUnauthorized)
				return
			}
			sessionID := cookie.Value

			// ----------------------------
			// STEP 1: Try Redis
			// ----------------------------
			var userID string
			userID, err = rdb.Get(ctx, sessionID).Result()
			if err == redis.Nil {
				// ----------------------------
				// STEP 2: Not in Redis → check DB
				// ----------------------------
				session, err := db.Query.GetSessionByToken(ctx, sessionID)
				if err != nil {
					// Clean up cookie
					http.SetCookie(w, &http.Cookie{
						Name:     sessionCookieName,
						Value:    "",
						Path:     "/",
						MaxAge:   -1,
						HttpOnly: true,
						Secure:   true,
					})
					http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
					return
				}

				userID = session.UserID.String()

				// ----------------------------
				// STEP 3: Cache in Redis
				// ----------------------------
				_ = rdb.Set(ctx, sessionID, userID, 0).Err() // optional expiration
			} else if err != nil {
				// Redis connection error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// ----------------------------
			// STEP 4: Inject userID into context
			// ----------------------------
			const contextKeyUserID contextKey = "userID"
			ctx = context.WithValue(ctx, contextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
