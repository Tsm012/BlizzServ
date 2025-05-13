package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Tsm012/BlizServe/internal/models"
)

type Handler struct {
	healthCheckManager HealthCheckManager
	pageSize           int
}

func NewHandler(healthCheckManager HealthCheckManager) *Handler {
	return &Handler{
		healthCheckManager: healthCheckManager,
		pageSize:           10,
	}
}

func (h *Handler) ListHealthChecksHandler(w http.ResponseWriter, r *http.Request) {
	pageParameter := r.URL.Query().Get("page")

	// string to int
	page, err := strconv.Atoi(pageParameter)
	if err != nil {
		page = 1
	}

	healthChecks := h.healthCheckManager.ListHealthChecks()

	start := h.pageSize * (page - 1)
	end := start + h.pageSize

	if start > len(healthChecks) {
		start = len(healthChecks)
	}

	if end > len(healthChecks) {
		end = len(healthChecks)
	}

	healthChecksSlice := healthChecks[start:end]

	healthChecksResponse := models.ListHealthCheckResponseModel{
		Page:         page,
		Total:        len(healthChecks),
		Size:         len(healthChecksSlice),
		HealthChecks: healthChecksSlice,
	}

	jsonData, err := json.Marshal(healthChecksResponse)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(jsonData))
}

func (h *Handler) GetHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	healthCheck, err := h.healthCheckManager.GetHealthCheck(r.PathValue("server_id"))
	if err != nil {
		http.Error(w, "Health check not found", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(healthCheck)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(jsonData))
}

func (h *Handler) CreateHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload models.CreateHealthCheckRequestModel
	err := decoder.Decode(&payload)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if payload.Endpoint == "" {
		http.Error(w, fmt.Sprintf("Endpoint cannot be blank: %s", payload.Endpoint), http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(payload.Endpoint)

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid Endpoint URL: %s", payload.Endpoint), http.StatusBadRequest)
		return
	}

	if h.healthCheckManager.endpointLookup[payload.Endpoint] != "" {
		http.Error(w, fmt.Sprintf("Health check already exists for endpoint: %s", payload.Endpoint), http.StatusConflict)
		return
	}

	healthCheckId := generateUUID()

	err = h.healthCheckManager.AddHealthCheck(models.HealthCheckModel{
		ID:       healthCheckId,
		Endpoint: payload.Endpoint,
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := models.CreateHealthCheckResponseModel{
		ID:       healthCheckId,
		Endpoint: payload.Endpoint,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(jsonData))
}

func (h *Handler) ExecuteHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Println("Error parsing duration:", err)
		return
	}

	healthCheck, err := h.healthCheckManager.GetHealthCheck(r.PathValue("server_id"))
	if err != nil {
		http.Error(w, "Health check not found", http.StatusNotFound)
		return
	}
	response := h.healthCheckManager.performHealthCheck(*healthCheck, duration)

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonData))
}

func (h *Handler) DeleteHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	h.healthCheckManager.DeleteHealthCheck(r.PathValue("server_id"))
}
