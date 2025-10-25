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

	// TODO: call migrations & seed admin

	r := httpx.NewRouter(func(api chi.Router) {
		authH := &handler.AuthHandler{DB: gdb, JWTSecret: cfg.JWTSecret}
		api.Post("/auth/login", authH.Login)
		// TODO: add employee/attendance/leave handlers & middlewares
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
