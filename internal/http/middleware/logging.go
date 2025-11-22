package middleware

import (
	"log"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		span := trace.SpanFromContext(r.Context())
		traceID := span.SpanContext().TraceID().String()

		log.Printf("method=%s path=%s status=%d duration=%s trace_id=%s", r.Method, r.URL.Path, ww.Status(), time.Since(start).String(), traceID)
	})
}
