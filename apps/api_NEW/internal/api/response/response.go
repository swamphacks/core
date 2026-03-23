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

// Send marshals any successful payload struct to JSON, sets the status code,
// and writes the response.
func Send(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// If encoding fails, log the error and fall back to a plain text error.
			// This is crucial because the header has already been written.
			log.Err(err).Str("function", "Send").Msg("Failed to encode and send JSON success object")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
