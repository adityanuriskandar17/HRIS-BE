package http

import (
	"net/http"

	"github.com/adityanuriskandar17/HRIS-BE/internal/config"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/httputils"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HandlerRegistrar func(r chi.Router)

func NewRouter(corsConfig config.CORSConfig, register HandlerRegistrar, tenantHandler *handler.TenantHandler, subscriptionHandler *handler.SubscriptionHandler, invoiceHandler *handler.InvoiceHandler, companyHandler *handler.CompanyHandler, port string) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP)

	// Configure CORS
	corsOpts := cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		ExposedHeaders:   corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	}
	r.Use(cors.New(corsOpts).Handler)
	r.Use(middleware.RequestLogger)
	r.Use(chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httputils.OK(w, r, map[string]any{"status": "ok"})
	})

	r.Route("/api/v1", func(api chi.Router) { register(api) })

	r.Route("/subscriptions", func(subscription chi.Router) {
		subscription.Use(chimw.RealIP, chimw.Logger, chimw.Recoverer)
		subscription.Get("/", subscriptionHandler.GetAll)
		subscription.Get("/{id}", subscriptionHandler.GetByID)
		subscription.Get("/tenant/{tenantId}", subscriptionHandler.GetByTenantID)
		subscription.Post("/", subscriptionHandler.Create)
		subscription.Put("/{id}", subscriptionHandler.Update)
		subscription.Post("/{id}/cancel", subscriptionHandler.Cancel)
		subscription.Post("/{id}/renew", subscriptionHandler.Renew)
	})

	r.Route("/invoices", func(invoice chi.Router) {
		invoice.Use(chimw.RealIP, chimw.Logger, chimw.Recoverer)
		invoice.Get("/", invoiceHandler.GetAll)
		invoice.Get("/{id}", invoiceHandler.GetByID)
		invoice.Get("/subscription/{subscriptionId}", invoiceHandler.GetBySubscriptionID)
		invoice.Post("/", invoiceHandler.Create)
		invoice.Put("/{id}", invoiceHandler.Update)
		invoice.Post("/{id}/send", invoiceHandler.Send)
		invoice.Post("/{id}/pay", invoiceHandler.Pay)
	})

	// Swagger documentation route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
	))

	// Add explicit route for /swagger/ to redirect to /swagger/index.html
	r.Get("/swagger/", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)

	return r
}
