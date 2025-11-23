package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a successful response
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// Error sends an error response
func Error(w http.ResponseWriter, statusCode int, message string) {
	response := Response{
		Success: false,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// ReadJSON reads and parses JSON from request body
func ReadJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// Meta represents metadata in API responses
type Meta struct {
	TraceID string `json:"trace_id,omitempty"`
}

// ErrorDetail represents error information in API responses
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// EnvelopeResponse is the standard envelope for all API responses
type EnvelopeResponse struct {
	Success bool         `json:"success"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Meta    Meta         `json:"meta"`
}

// TokenResponseEnvelope wraps TokenResponse in the standard envelope
type TokenResponseEnvelope struct {
	Success bool `json:"success"`
	Data    struct {
		AccessToken      string `json:"accessToken"`
		RefreshToken     string `json:"refreshToken"`
		ExpiresIn        int64  `json:"expiresIn"`
		RefreshExpiresIn int64  `json:"refreshExpiresIn"`
		Role             string `json:"role"`
	} `json:"data"`
	Meta Meta `json:"meta"`
}
