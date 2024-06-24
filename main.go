package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kmrhemant916/k8s-webhooks/helpers"
	"github.com/kmrhemant916/k8s-webhooks/routes"
)

func main() {
    configPath := flag.String("config", "config/config.yaml", "config file for webhook")
    certFile := flag.String("certFile", "/etc/webhook/certs/tls.crt", "certificate file")
	keyFile := flag.String("keyFile", "/etc/webhook/certs/tls.key", "certificate key file")
	caFile := flag.String("caFile", "/etc/webhook/certs/ca.crt", "CA file")
    flag.Parse()
    config, err := helpers.ReadConfig(*configPath)
    if err != nil {
        log.Printf("Error reading config: %v", err)
        panic(err)
    }
	r := routes.SetupRoutes(config)
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}
	caCert, err := os.ReadFile(*caFile)
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
			ClientAuth:   tls.VerifyClientCertIfGiven,
		},
	}
	log.Fatal(server.ListenAndServeTLS("", ""))
}
