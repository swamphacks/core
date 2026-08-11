package announcements

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func RegisterHTTPRoutes(handler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "get-all-announcements",
		Method:        http.MethodGet,
		Summary:       "Get All Announcements",
		Description:   "Returns all announcements for the specified hackathon.",
		Tags:          []string{"Announcements"},
		Path:          "/",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.getAllAnnouncements)

	sse.Register(group, huma.Operation{
		OperationID:   "get-active-announcements",
		Method:        http.MethodGet,
		Summary:       "Get Active Announcements",
		Description:   "Returns active (non-expired) announcements for the specified hackathon.",
		Tags:          []string{"Announcements"},
		Path:          "/active",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	},
	map[string]any{
		"message": sqlc.Announcement{},
	},
	handler.getActiveAnnouncements)

	huma.Register(group, huma.Operation{
		OperationID:   "create-announcement",
		Method:        http.MethodPost,
		Summary:       "Create Announcement",
		Description:   "Create announcement for the specified hackathon.",
		Tags:          []string{"Announcements"},
		Path:          "/",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusCreated,
	}, handler.createAnnouncement)

	huma.Register(group, huma.Operation{
		OperationID:   "update-announcement",
		Method:        http.MethodPatch,
		Summary:       "Update Announcement",
		Description:   "Update announcement for the specified hackathon.",
		Tags:          []string{"Announcements"},
		Path:          "/{announcementId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.updateAnnouncement)

	huma.Register(group, huma.Operation{
		OperationID:   "delete-announcement",
		Method:        http.MethodDelete,
		Summary:       "Delete Announcement",
		Description:   "Delete announcement for the specified hackathon.",
		Tags:          []string{"Announcements"},
		Path:          "/{announcementId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.deleteAnnouncement)

	huma.Register(group, huma.Operation{
		OperationID:   "dismiss-announcement",
		Method:        http.MethodPost,
		Summary:       "Dismiss Announcement",
		Description:   "Dismiss announcement for a user.",
		Tags:          []string{"Announcements"},
		Path:          "/{announcementId}/dismiss",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusNoContent,
	}, handler.dismissAnnouncement)

	huma.Register(group, huma.Operation{
		OperationID:   "get-dismissed-announcements",
		Method:        http.MethodGet,
		Summary:       "Get Dismissed Announcements",
		Description:   "Returns all dismissed announcements for a user.",
		Tags:          []string{"Announcements"},
		Path:          "/dismissed",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, handler.getDismissedAnnouncements)
}