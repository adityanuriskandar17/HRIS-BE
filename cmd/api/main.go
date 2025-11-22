// @title HRIS API
// @description Human Resource Information System API
// @version 1.0
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/adityanuriskandar17/HRIS-BE/docs"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/adityanuriskandar17/HRIS-BE/internal/auth"
	"github.com/adityanuriskandar17/HRIS-BE/internal/config"
	"github.com/adityanuriskandar17/HRIS-BE/internal/db"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	domainRepository "github.com/adityanuriskandar17/HRIS-BE/internal/domain/repository"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
	httputil "github.com/adityanuriskandar17/HRIS-BE/internal/http/middleware"
	"github.com/adityanuriskandar17/HRIS-BE/internal/repository"
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
	if err := db.SeedReferenceData(gdb, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatal(err)
	}

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
	tokenSvc := auth.NewService(gdb, cfg.JWTSecret, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)
	authH := handler.NewAuthHandler(userRepo, tenantService, tokenSvc)
	authMw := &httputil.Authenticator{Secret: cfg.JWTSecret, DB: gdb}

	tenantHandler := handler.NewTenantHandler(tenantService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	companyHandler := handler.NewCompanyHandler(companyService)
	masterH := handler.NewMasterDataHandler(unitRepo, positionRepo, employeeRepo)

	r := httpx.NewRouter(cfg.Auth.AllowedOrigins, func(api chi.Router) {
		api.Post("/auth/login", authH.Login)
		api.Post("/auth/register", authH.Register)
		api.Post("/auth/refresh", authH.Refresh)

		api.Group(func(protected chi.Router) {
			protected.Use(authMw.Middleware)
			protected.Get("/auth/profile", authH.ProfileSelf)
			protected.Get("/auth/profile/{id}", authH.ProfileByID)

			protected.Group(func(admin chi.Router) {
				admin.Use(httputil.RequireRoles(model.RoleAdmin))
				admin.Post("/auth/employee", authH.CreateEmployeeAccount)
			})

			protected.Route("/master", func(m chi.Router) {
				m.Group(func(sec chi.Router) {
					sec.Use(httputil.RequireRoles(model.RoleAdmin, model.RoleHR))
					sec.Get("/units", masterH.ListUnits)
					sec.Post("/units", masterH.CreateUnit)

					sec.Get("/positions", masterH.ListPositions)
					sec.Post("/positions", masterH.CreatePosition)
				})

				m.Group(func(sec chi.Router) {
					sec.Use(httputil.RequireRoles(model.RoleAdmin, model.RoleHR, model.RoleManager))
					sec.Get("/employees", masterH.ListEmployees)
					sec.Post("/employees", masterH.CreateEmployee)
					sec.Get("/employees/{id}", masterH.GetEmployee)
					sec.Put("/employees/{id}", masterH.UpdateEmployee)
					sec.Delete("/employees/{id}", masterH.DeleteEmployee)
				})
			})
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
