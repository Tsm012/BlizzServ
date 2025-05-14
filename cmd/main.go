package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Tsm012/BlizServe/internal"
)

func main() {
	bind := flag.String("bind", "", "Address to bind the server")
	ssl := flag.Bool("ssl", false, "Enable SSL")
	sslCert := flag.String("sslcert", "", "Path to SSL certificate file")
	sslKey := flag.String("sslkey", "", "Path to SSL key file")
	checkFrequency := flag.Duration("checkfrequency", 0, "Frequency for health checks")
	flag.Parse()

	if *bind == "" {
		log.Fatal("The --bind flag is required")
	}

	if *checkFrequency == 0 {
		// 30 second frequency by default
		*checkFrequency = 30 * time.Second
	}

	healthCheckManager := internal.NewHealthCheckManager(*checkFrequency)
	handler := internal.NewHandler(healthCheckManager)

	http.HandleFunc("GET /api/health/checks", handler.ListHealthChecksHandler)
	http.HandleFunc("GET /api/health/checks/{server_id}", handler.GetHealthCheckHandler)
	http.HandleFunc("POST /api/health/checks", handler.CreateHealthCheckHandler)
	http.HandleFunc("POST /api/health/checks/{server_id}/try", handler.ExecuteHealthcheckHandler)
	http.HandleFunc("DELETE /api/health/checks/{server_id}", handler.DeleteHealthcheckHandler)

	if *ssl {
		if *sslCert != "" && *sslKey != "" {
			err := http.ListenAndServeTLS(*bind, *sslCert, *sslKey, nil)
			if err != nil {
				log.Fatalf("Failed to start server with SSL: %v", err)
			}
		}
	} else {
		err := http.ListenAndServe(*bind, nil)
		if err != nil {
			log.Fatalf("Failed to start server with SSL: %v", err)
		}
	}
}
