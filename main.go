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
	Config = "config/config.yaml"
	CertFile = "/etc/webhook/certs/tls.crt"
	KeyFile = "/etc/webhook/certs/tls.key"
)

func main() {	
	r := routes.SetupRoutes()
    cert, err := os.ReadFile(CertFile)
    if err != nil {
        log.Fatalf("failed to read certificate: %v", err)
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(cert)
    server := &http.Server{
        Addr:    ":443",
        Handler: r,
        TLSConfig: &tls.Config{
            ClientCAs:  certPool,
            ClientAuth: tls.RequireAndVerifyClientCert,
        },
    }
    log.Fatal(server.ListenAndServeTLS(CertFile, KeyFile))
}