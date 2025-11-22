package httputils

import (
	"net/http"

	"github.com/adityanuriskandar17/HRIS-BE/internal/http/envelope"
)

func OK(w http.ResponseWriter, r *http.Request, data any) {
	envelope.Respond(w, r, http.StatusOK, data, nil)
}

func Created(w http.ResponseWriter, r *http.Request, data any) {
	envelope.Respond(w, r, http.StatusCreated, data, nil)
}

func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, fields map[string]string) {
	envelope.Error(w, r, status, code, message, fields)
}
