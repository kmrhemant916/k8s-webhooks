package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kmrhemant916/k8s-webhooks/controllers"
	"github.com/kmrhemant916/k8s-webhooks/helpers"
)


func SetupRoutes(config *helpers.Config) (*chi.Mux){
	app := controllers.NewApp()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/mutate", func(w http.ResponseWriter, r *http.Request) {
        app.Mutate(w, r, config)
    })
	router.Get("/healthz", app.Healthz)
	return router
}