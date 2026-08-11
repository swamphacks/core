package announcements

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/swamphacks/core/apps/api/internal/config"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type handler struct {
	announcementsService *AnnouncementsService
	logger               zerolog.Logger
}

func NewHandler(announcementsService *AnnouncementsService, logger zerolog.Logger) *handler {
	return &handler{
		announcementsService: announcementsService,
		logger: logger.With().Str("handler", "AnnouncementsHandler").Str("domain", "announcements").Logger(),
	}
}

func (h *handler) getActiveAnnouncements(
	ctx context.Context,
	input *struct {
		HackathonID string `query:"hackathonId" required:"true"`
	},
	send sse.Sender,
) {
}

func (h *handler) getAllAnnouncements(
	ctx context.Context,
	input *struct {
		HackathonID string `query:"hackathonId" required:"true"`
	}) (*ListAnnouncementsOutput, error) {
	return nil, nil
}

func (h *handler) createAnnouncement(
	ctx context.Context,
	input *struct {
		Body *CreateAnnouncementRequest
	}) (*AnnouncementOutput, error) {
	return nil, nil
}

func (h *handler) updateAnnouncement(
	ctx context.Context,
	input *struct {
		AnnouncementID string                     `path:"announcementId"`
		HackathonID    string                     `query:"hackathonId" required:"true"`
		Body           *UpdateAnnouncementRequest
	}) (*AnnouncementOutput, error) {
	return nil, nil
}

func (h *handler) deleteAnnouncement(
	ctx context.Context,
	input *struct {
		AnnouncementID string `path:"announcementId"`
		HackathonID    string `query:"hackathonId" required:"true"`
	}) (*DeleteAnnouncementOutput, error) {
	return nil, nil
}

func (h *handler) dismissAnnouncement(
	ctx context.Context,
	input *struct {
		AnnouncementID string `path:"announcementId"`
	}) (*DismissAnnouncementOutput, error) {
	return nil, nil
}

func (h *handler) getDismissedAnnouncements(
	ctx context.Context,
	input *struct {
}) (*ListDismissedAnnouncementsOutput, error) {
	return nil, nil
}

var apiConfig = config.LoadConfig()
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")

		for _, allowedOrigin := range apiConfig.AllowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}

		return false
    },
}

var (
    clients = make(map[*websocket.Conn]bool)
    clientsMu sync.Mutex
)

func registerClient(conn *websocket.Conn) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    clients[conn] = true
}

func unregisterClient(conn *websocket.Conn) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    delete(clients, conn)
    conn.Close()
}

func (h *handler) announcementsWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
        return
    }
	registerClient(conn)
	defer unregisterClient(conn)

	// Listen for incoming messages
    for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)
	}
}