package http

import (
	"net/http"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HandlerRegistrar func(r chi.Router)

func NewRouter(register HandlerRegistrar, tenantHandler *handler.TenantHandler, subscriptionHandler *handler.SubscriptionHandler, invoiceHandler *handler.InvoiceHandler, port string) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		OK(w, map[string]any{"status": "ok"})
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

	// Swagger documentation route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
	))

	// Add explicit route for /swagger/ to redirect to /swagger/index.html
	r.Get("/swagger/", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)

	return r
}
