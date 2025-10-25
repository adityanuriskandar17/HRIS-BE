package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/adityanuriskandar17/HRIS-BE/internal/config"
	"github.com/adityanuriskandar17/HRIS-BE/internal/db"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/handler"
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
	log.Fatal(http.ListenAndServe(addr, r))
}
