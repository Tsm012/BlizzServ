package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Tsm012/BlizServe/internal/models"
)

type Handler struct {
	healthCheckManager HealthCheckManager
}

func NewHandler() *Handler {
	return &Handler{
		healthCheckManager: NewHealthCheckManager(),
	}
}

func (h *Handler) ListHealthChecksHandler(w http.ResponseWriter, r *http.Request) {
	healthChecks := h.healthCheckManager.ListHealthChecks()

	//page := r.URL.Query().Get("page")
	//fmt.Println(len(healthChecks))
	//	fmt.Fprintf(w, )

	healthChecksResponse := models.ListHealthCheckResponseModel{
		Page:         1,
		Total:        int32(len(healthChecks)),
		Size:         int32(len(healthChecks)),
		HealthChecks: healthChecks,
	}

	jsonData, err := json.Marshal(healthChecksResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonData))
}

func (h *Handler) GetHealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	healthCheck, err := h.healthCheckManager.GetHealthCheck(r.PathValue("server_id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonData, err := json.Marshal(healthCheck)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonData))
}

func (h *Handler) CreateHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload models.CreateHealthCheckRequestModel
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	if h.healthCheckManager.endpointLookup[payload.Endpoint] != "" {
		fmt.Fprintf(w, "Health check already exists for endpoint: %s", payload.Endpoint)
		return
	}

	healthCheckId := generateUUID()

	h.healthCheckManager.AddHealthCheck(models.HealthCheckModel{
		ID:       healthCheckId,
		Endpoint: payload.Endpoint,
	})

	response := models.CreateHealthCheckResponseModel{
		ID:       healthCheckId,
		Endpoint: payload.Endpoint,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "%s", string(jsonData))
}

func (h *Handler) ExecuteHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	fmt.Fprintf(w, "timeout is: %s", timeout)
	response := models.HealthCheckModel{
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

func (h *Handler) DeleteHealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	h.healthCheckManager.DeleteHealthCheck(r.PathValue("server_id"))
}
