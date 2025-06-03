package response

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewError(errorCode, msg string) ErrorResponse {
	return ErrorResponse{
		Error:   errorCode,
		Message: msg,
	}
}

func SendError(w http.ResponseWriter, status int, errorResponse ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Err(err).Str("function", "sendError").Msg("Failed to encode and send JSON object")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
