package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/kmrhemant916/k8s-webhooks/helpers"
	"github.com/kmrhemant916/k8s-webhooks/routes"
)

const (
	Config = "config/config.yaml"
	CertFile = "/etc/webhook/certs/tls.crt"
	KeyFile = "/etc/webhook/certs/tls.key"
)

func mustLoadCertificate(certFile string, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %v", err)
	}
	return cert
}

func main() {	
	c, err:= helpers.ReadConfig(Config)
    if err != nil {
        panic(err)
	}
	r := routes.SetupRoutes()
	server := &http.Server{
		Addr:    ":443",
		Handler: r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				mustLoadCertificate(CertFile, KeyFile),
			},
		},
	}
	log.Fatal(server.ListenAndServeTLS(CertFile, KeyFile))
	http.ListenAndServe(":"+c.Service.Port, r)
}