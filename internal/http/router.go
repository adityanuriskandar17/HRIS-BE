package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type HandlerRegistrar func(r chi.Router)

func NewRouter(register HandlerRegistrar) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		OK(w, map[string]any{"status": "ok"})
	})

	r.Route("/api/v1", func(api chi.Router) { register(api) })
	return r
}
