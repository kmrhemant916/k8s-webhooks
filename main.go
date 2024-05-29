package main

import (
	"net/http"

	"github.com/kmrhemant916/k8s-webhooks/helpers"
	"github.com/kmrhemant916/k8s-webhooks/routes"
)

const (
	Config = "config/config.yaml"
)

func main() {	
	c, err:= helpers.ReadConfig(Config)
    if err != nil {
        panic(err)
	}
	r := routes.SetupRoutes()
	http.ListenAndServe(":"+c.Service.Port, r)
}