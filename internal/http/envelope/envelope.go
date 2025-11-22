package envelope

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type Meta struct {
	TraceID string `json:"trace_id,omitempty"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Response struct {
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Meta    Meta         `json:"meta"`
}

func traceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	sc := span.SpanContext()
	if !sc.HasTraceID() {
		return ""
	}
	return sc.TraceID().String()
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data any, errDetail *ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	envelope := Response{
		Success: errDetail == nil,
		Meta: Meta{
			TraceID: traceIDFromContext(r.Context()),
		},
	}

	if errDetail != nil {
		envelope.Error = errDetail
	} else {
		envelope.Data = data
	}

	_ = json.NewEncoder(w).Encode(envelope)
}

func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, fields map[string]string) {
	if fields == nil || len(fields) == 0 {
		fields = nil
	}

	Respond(w, r, status, nil, &ErrorDetail{
		Code:    code,
		Message: message,
		Fields:  fields,
	})
}
