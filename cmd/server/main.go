package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"https/github.com/proxy1301sl/auth/internal/handler"
	"https/github.com/proxy1301sl/auth/middlewareauth"
	"https/github.com/proxy1301sl/auth/repo"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	storage, err := repo.NewStorage("")
	if err != nil {
		log.Fatal(err)
	}
	r.Post("/api/register", handler.Register(storage))

	r.Post("/api/login", handler.Login(storage))

	r.Group(func(r chi.Router) {
		r.Use(middlewareauth.AuthMiddleware)
		r.Get("/api/profile", handler.Profile(storage))

	})
	log.Fatal(http.ListenAndServe(":8080", r))

}
