package announcements

import (
	"github.com/go-chi/chi/v5"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
)

func RegisterWSRoutes(handler *handler, router chi.Router, mw *middleware.Middleware) {
	// TODO: new discord only middleware
	router.Get(
		"/announcements/ws",
		handler.announcementsWS,
	)
}