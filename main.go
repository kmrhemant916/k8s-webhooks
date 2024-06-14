package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/kmrhemant916/k8s-webhooks/routes"
)

const (
	Config   = "config/config.yaml"
	CertFile = "/etc/webhook/certs/tls.crt"
	KeyFile  = "/etc/webhook/certs/tls.key"
	CAFile   = "/etc/webhook/certs/ca.crt"
)

func main() {
	r := routes.SetupRoutes()

	// Load server certificate
	cert, err := tls.LoadX509KeyPair(CertFile, KeyFile)
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile(CAFile)
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("failed to append CA certificate to pool")
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
			ClientAuth:   tls.VerifyClientCertIfGiven, // Optional client certificate
		},
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}
