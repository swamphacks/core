package handlers
import (
    "encoding/json"
    "net/http"
    
    "github.com/go-chi/chi/v5"
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

// GetUserEventRoleByDiscordID
//
//	@Summary		Get user event role by Discord ID
//	@Description	Get the event role for a user based on their Discord account ID
//	@Tags			Discord
//	@Param			discord_id	path		string	true	"Discord account ID"
//	@Success		200			{object}	map[string]interface{}	"event_id and role"
//	@Failure		404			{object}	response.ErrorResponse	"User or role not found"
//	@Failure		500			{object}	response.ErrorResponse	"Internal server error"
//	@Router			/discord/user/{discord_id}/role [get]

func (h *DiscordHandler) GetUserEventRoleByDiscordID(w http.ResponseWriter, r *http.Request) {
	discordID := chi.URLParam(r, "discord_id")
	if discordID == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("bad_request", "discord_id is required"))
		return
	}

	eventID, role, err := h.discordService.GetUserEventRoleByDiscordID(r.Context(), discordID)
	if err != nil {
		if err == services.ErrDiscordAccountNotFound || err == services.ErrNoEventRole {
            res.SendError(w, http.StatusNotFound, res.NewError("not_found", err.Error()))
            return
        }
        h.logger.Err(err).Msg("failed to get user event role")
        res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to get user role"))
        return
	}

	response := map[string]interface{}{
		"event_id": eventID,
		"role": role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}