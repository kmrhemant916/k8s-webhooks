package controllers

import "github.com/kmrhemant916/k8s-webhooks/helpers"

type App struct {
	ReadConfig func(string) (*helpers.Config, error)
}

func NewApp() *App {
    app := &App{
        ReadConfig: helpers.ReadConfig,
    }
    return app
}