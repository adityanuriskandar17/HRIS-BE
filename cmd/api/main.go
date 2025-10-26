package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/adityanuriskandar17/HRIS-BE/internal/config"
	"github.com/adityanuriskandar17/HRIS-BE/internal/db"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
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

	if err := db.AutoMigrate(gdb); err != nil {
		log.Fatal(err)
	}
	if err := db.SeedReferenceData(gdb); err != nil {
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

	r := httpx.NewRouter(func(api chi.Router) {
		authH := &handler.AuthHandler{DB: gdb, JWTSecret: cfg.JWTSecret}
		api.Post("/auth/login", authH.Login)
		// TODO: add employee/attendance/leave handlers & middlewares

		masterH := &handler.MasterDataHandler{DB: gdb}
		api.Route("/master", func(m chi.Router) {
			m.Get("/units", masterH.ListUnits)
			m.Post("/units", masterH.CreateUnit)

			m.Get("/positions", masterH.ListPositions)
			m.Post("/positions", masterH.CreatePosition)

			m.Get("/employees", masterH.ListEmployees)
			m.Post("/employees", masterH.CreateEmployee)
		})
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
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
