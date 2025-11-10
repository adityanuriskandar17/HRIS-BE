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