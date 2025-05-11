package main

import (
	"flag"
	"log"
	"net/http"
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
		//TODO : Add a function to start health checks
	}

	http.HandleFunc("GET /api/health/checks", listHealthchecksHandler)
	http.HandleFunc("GET /api/health/checks/{server_id}", getHealthcheckHandler)
	http.HandleFunc("POST /api/health/checks", createHealthcheckHandler)
	http.HandleFunc("POST /api/health/checks/{server_id}/try", executeHealthcheckHandler)
	http.HandleFunc("DELETE /api/health/checks/{server_id}", deleteHealthcheckHandler)

	if *ssl {
		if *sslCert != "" && *sslKey != "" {
			http.ListenAndServeTLS(*bind, *sslCert, *sslKey, nil)
		}
	}
	http.ListenAndServe(*bind, nil)
}
