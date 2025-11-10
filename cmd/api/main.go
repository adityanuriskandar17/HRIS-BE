// @title HRIS API
// @description Human Resource Information System API
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	_ "github.com/adityanuriskandar17/HRIS-BE/docs"

	"github.com/adityanuriskandar17/HRIS-BE/internal/config"
	"github.com/adityanuriskandar17/HRIS-BE/internal/db"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
	"github.com/adityanuriskandar17/HRIS-BE/internal/repository"
	domainRepository "github.com/adityanuriskandar17/HRIS-BE/internal/domain/repository"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	"github.com/adityanuriskandar17/HRIS-BE/internal/telemetry"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	gdb, err := db.Open(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	// Enable UUID extension before running migrations
	if err := gdb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(gdb); err != nil {
		log.Fatal(err)
	}
	db.Seed(gdb)

	shutdownTelemetry, err := telemetry.Setup(context.Background(), cfg.Telemetry.ServiceName, cfg.Telemetry.JaegerEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTelemetry(context.Background()); err != nil {
			log.Printf("telemetry shutdown error: %v", err)
		}
	}()

	// Initialize repositories
	userRepo := repository.NewUserAccountRepository(gdb)
	unitRepo := repository.NewUnitRepository(gdb)
	positionRepo := repository.NewPositionRepository(gdb)
	employeeRepo := repository.NewEmployeeRepository(gdb)
	tenantRepo := domainRepository.NewTenantRepository(gdb)
	subscriptionRepo := repository.NewSubscriptionRepository(gdb)
	planRepo := repository.NewPlanRepository(gdb)
	invoiceRepo := repository.NewInvoiceRepository(gdb)
	companyRepo := domainRepository.NewCompanyRepository(gdb)

	// Initialize services
	tenantService := services.NewTenantService(tenantRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, planRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo, subscriptionRepo)
	companyService := services.NewCompanyService(companyRepo)

	// Initialize handlers
	tenantHandler := handler.NewTenantHandler(tenantService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	companyHandler := handler.NewCompanyHandler(companyService)

	r := httpx.NewRouter(func(api chi.Router) {
		authH := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
		api.Post("/auth/login", authH.Login)
		// TODO: add employee/attendance/leave handlers & middlewares

		masterH := handler.NewMasterDataHandler(unitRepo, positionRepo, employeeRepo)
		api.Route("/master", func(m chi.Router) {
			m.Get("/units", masterH.ListUnits)
			m.Post("/units", masterH.CreateUnit)

			m.Get("/positions", masterH.ListPositions)
			m.Post("/positions", masterH.CreatePosition)

			m.Get("/employees", masterH.ListEmployees)
			m.Post("/employees", masterH.CreateEmployee)
		})

		// Subscription routes
		api.Route("/subscriptions", func(m chi.Router) {
			m.Get("/", subscriptionHandler.GetAll)
			m.Post("/", subscriptionHandler.Create)
			m.Get("/{id}", subscriptionHandler.GetByID)
			m.Put("/{id}", subscriptionHandler.Update)
			m.Delete("/{id}", subscriptionHandler.Cancel)
			m.Post("/{id}/renew", subscriptionHandler.Renew)
			m.Get("/tenant/{tenantId}", subscriptionHandler.GetByTenantID)
		})

		// Invoice routes
		api.Route("/invoices", func(m chi.Router) {
			m.Get("/", invoiceHandler.GetAll)
			m.Post("/", invoiceHandler.Create)
			m.Get("/{id}", invoiceHandler.GetByID)
			m.Put("/{id}", invoiceHandler.Update)
			m.Post("/{id}/send", invoiceHandler.Send)
			m.Post("/{id}/pay", invoiceHandler.Pay)
			m.Get("/subscription/{subscriptionId}", invoiceHandler.GetBySubscriptionID)
		})
	}, tenantHandler, subscriptionHandler, invoiceHandler, companyHandler, cfg.Port)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", cfg.Port)
	handler := otelhttp.NewHandler(
		r,
		cfg.Telemetry.ServiceName,
		otelhttp.WithSpanNameFormatter(func(operation string, req *http.Request) string {
			if routeCtx := chi.RouteContext(req.Context()); routeCtx != nil {
				if route := routeCtx.RoutePattern(); route != "" {
					return req.Method + " " + route
				}
			}
			return req.Method + " " + req.URL.Path
		}),
	)
	log.Fatal(http.ListenAndServe(addr, handler))
}
