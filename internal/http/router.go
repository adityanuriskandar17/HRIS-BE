package http

import (
	"net/http"

	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/httputils"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HandlerRegistrar func(r chi.Router)

func NewRouter(allowedOrigins []string, register HandlerRegistrar, tenantHandler *handler.TenantHandler, subscriptionHandler *handler.SubscriptionHandler, invoiceHandler *handler.InvoiceHandler, companyHandler *handler.CompanyHandler, port string) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP)

	if len(allowedOrigins) > 0 {
		corsOpts := cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}
		r.Use(cors.New(corsOpts).Handler)
	}
	r.Use(middleware.RequestLogger)
	r.Use(chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httputils.OK(w, r, map[string]any{"status": "ok"})
	})

	r.Route("/api/v1", func(api chi.Router) { register(api) })
	


	r.Route("/tenants", func(tenant chi.Router) {
		tenant.Use(chimw.RealIP, chimw.Logger, chimw.Recoverer)
		tenant.Get("/", tenantHandler.GetAll)
		tenant.Get("/{id}", tenantHandler.GetByID)
		tenant.Post("/", tenantHandler.Create)
		tenant.Put("/{id}", tenantHandler.Update)
	})

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

	r.Route("/companies", func(company chi.Router) {
		company.Use(chimw.RealIP, chimw.Logger, chimw.Recoverer)
		company.Get("/{id}", companyHandler.GetCompanyProfile)
		company.Patch("/{id}", companyHandler.UpdateCompanyProfile)
		company.Get("/{id}/settings", companyHandler.GetCompanySettings)
		company.Patch("/{id}/settings", companyHandler.UpdateCompanySettings)
		company.Get("/{id}/limits", companyHandler.GetCompanyLimits)
	})

	// Swagger documentation route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
	))

	// Add explicit route for /swagger/ to redirect to /swagger/index.html
	r.Get("/swagger/", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)

	return r
}
