package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kmrhemant916/k8s-webhooks/controllers"
)


func SetupRoutes() (*chi.Mux){
	app := controllers.NewApp()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/mutate", app.Mutate)
	return router
}