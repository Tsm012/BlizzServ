package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func listHealthchecksHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	fmt.Fprintf(w, "page is: %s", page)
}

func getHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Ideally not magic string
	fmt.Fprintf(w, "path is: %s", r.PathValue("server_id"))
}

func createHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload CreateHealthcheckPayload
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}
	response := CreateHealthcheckResponse{
		ID:       "12345",
		Endpoint: payload.Endpoint,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonData))
}

func executeHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	fmt.Fprintf(w, "timeout is: %s", timeout)
	response := HealthcheckPayload{
		ID:       "12345",
		Status:   "OK",
		Code:     200,
		Endpoint: "http://example.com",
		Checked:  1622547800,
		Duration: "100ms",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonData))
}

func deleteHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "path is: %s", r.PathValue("server_id"))
}
