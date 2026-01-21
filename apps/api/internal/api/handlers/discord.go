package handlers
import (
    "encoding/json"
    "net/http"
    
    "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
    "github.com/rs/zerolog"
    res "github.com/swamphacks/core/apps/api/internal/api/response"
    "github.com/swamphacks/core/apps/api/internal/services"
)

type DiscordHandler struct {
	discordService *services.DiscordService
	logger         zerolog.Logger
}

func NewDiscordHandler(discordService *services.DiscordService, logger zerolog.Logger) *DiscordHandler {
	return &DiscordHandler{
		discordService: discordService,
		logger: logger.With().Str("handler", "DiscordHandler").Str("component", "discord").Logger(),
	}
}

// GetEventAttendeesWithDiscord
//
//	@Summary		Get event attendees with Discord IDs
//	@Description	Get all attendees for an event who have Discord accounts linked
//	@Tags			Discord
//	@Param			event_id	path		string	true	"Event ID (UUID)"
//	@Success		200			{array}		sqlc.GetEventAttendeesWithDiscordRow	"List of attendees with Discord IDs"
//	@Failure		400			{object}	response.ErrorResponse	"Invalid event ID"
//	@Failure		500			{object}	response.ErrorResponse	"Internal server error"
//	@Router			/discord/event/{event_id}/attendees [get]
func (h *DiscordHandler) GetEventAttendeesWithDiscord(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "event_id")
	if eventIDStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_request", "event_id is required"))
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_request", "invalid event ID format"))
		return
	}

	attendees, err := h.discordService.GetEventAttendeesWithDiscord(r.Context(), eventID)
	if err != nil {
		h.logger.Err(err).Msg("failed to get event attendees with discord")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to get attendees"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attendees); err != nil {
		h.logger.Err(err).Msg("failed to encode response")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to encode response"))
		return
	}
}

// GetUserEventRoleByDiscordIDAndEventId
//
//	@Summary		Get user event role by Discord ID and Event ID
//	@Description	Get the event role for a user based on their Discord account ID and a specific event ID
//	@Tags			Discord
//	@Param			eventId		path		string	true	"Event ID (UUID)"
//	@Param			discordId	path		string	true	"Discord account ID"
//	@Success		200			{object}	map[string]interface{}	"role"
//	@Failure		400			{object}	response.ErrorResponse	"Invalid event ID or discord ID"
//	@Failure		404			{object}	response.ErrorResponse	"User or role not found"
//	@Failure		500			{object}	response.ErrorResponse	"Internal server error"
//	@Router			/events/{eventId}/discord/{discordId} [get]
func (h *DiscordHandler) GetUserEventRoleByDiscordIDAndEventId(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	if eventIDStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_request", "eventId is required"))
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_request", "invalid event ID format"))
		return
	}

	discordID := chi.URLParam(r, "discordId")
	if discordID == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_request", "discordId is required"))
		return
	}

	role, err := h.discordService.GetUserEventRoleByDiscordIDAndEventId(r.Context(), discordID, eventID)
	if err != nil {
		if err == services.ErrNoEventRole {
			res.SendError(w, http.StatusNotFound, res.NewError("not_found", err.Error()))
			return
		}
		h.logger.Err(err).Msg("failed to get user event role")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to get user role"))
		return
	}

	response := map[string]interface{}{
		"role": role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}