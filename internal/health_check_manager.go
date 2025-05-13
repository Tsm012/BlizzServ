package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Tsm012/BlizServe/internal/models"
)

type HealthCheckManager struct {
	dataChannels   map[string]chan models.HealthCheckModel
	healthChecks   OrderedHealthCheckMap
	endpointLookup map[string]string
	checkFrequency time.Duration
	defaultTimeout time.Duration
}

func NewHealthCheckManager(checkFrequency time.Duration) HealthCheckManager {
	manager := HealthCheckManager{
		healthChecks: OrderedHealthCheckMap{
			Keys: []string{},
			Data: map[string]models.HealthCheckModel{},
		},
		endpointLookup: map[string]string{},
		dataChannels:   map[string]chan models.HealthCheckModel{},
		checkFrequency: checkFrequency,
		defaultTimeout: 15 * time.Second,
	}
	manager.LoadFromFile()
	return manager
}

func (hcm *HealthCheckManager) AddHealthCheck(healthCheckModel models.HealthCheckModel) error {
	hcm.healthChecks.Add(healthCheckModel.ID, healthCheckModel)
	hcm.endpointLookup[healthCheckModel.Endpoint] = healthCheckModel.ID
	dataChannel := make(chan models.HealthCheckModel)
	hcm.dataChannels[healthCheckModel.ID] = dataChannel
	go hcm.monitor(healthCheckModel, hcm.checkFrequency, dataChannel)
	go hcm.listener(dataChannel)
	hcm.SaveToFile()
	return nil
}

func (hcm *HealthCheckManager) ListHealthChecks() []models.HealthCheckModel {
	return hcm.healthChecks.ToSortedList()
}

func (hcm *HealthCheckManager) GetHealthCheck(serverId string) (*models.HealthCheckModel, error) {
	healthCheck, ok := hcm.healthChecks.Get(serverId)
	if !ok {
		return nil, errors.New("health check not found")
	}
	return &healthCheck, nil
}

func (hcm *HealthCheckManager) DeleteHealthCheck(serverId string) {
	healthCheck, ok := hcm.healthChecks.Get(serverId)
	if !ok {
		return
	}

	close(hcm.dataChannels[healthCheck.ID])
	delete(hcm.dataChannels, healthCheck.ID)
	hcm.healthChecks.Delete(healthCheck.ID)
	delete(hcm.endpointLookup, healthCheck.Endpoint)

	hcm.SaveToFile()
}

func (hcm *HealthCheckManager) SaveToFile() error {
	file, err := os.Create(".save")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(hcm.healthChecks); err != nil {
		return err
	}

	return nil
}

func (hcm *HealthCheckManager) LoadFromFile() {
	file, err := os.Open(".save")
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&hcm.healthChecks); err != nil {
		return
	}

	// Recreate the endpoint lookup map
	hcm.endpointLookup = make(map[string]string)
	for _, healthCheckKey := range hcm.healthChecks.Keys {
		healthCheck, _ := hcm.healthChecks.Get(healthCheckKey)
		hcm.endpointLookup[healthCheck.Endpoint] = healthCheck.ID
		dataChannel := make(chan models.HealthCheckModel)
		hcm.dataChannels[healthCheck.ID] = dataChannel
		go hcm.monitor(healthCheck, hcm.checkFrequency, dataChannel)
		go hcm.listener(dataChannel)
	}
}

func (hcm *HealthCheckManager) listener(dataChannel chan models.HealthCheckModel) {
	for healthCheck := range dataChannel {
		// Print the ID and the received HealthCheckModel
		healthCheckData, err := json.Marshal(healthCheck)
		if err != nil {
			log.Println("Error marshaling HealthCheckModel:", err.Error())
			continue
		}
		log.Println("Received update for ", string(healthCheckData))
		hcm.healthChecks.Set(healthCheck.ID, healthCheck)
	}
}

func (hcm *HealthCheckManager) monitor(healthCheckModel models.HealthCheckModel, frequency time.Duration, dataChannel chan models.HealthCheckModel) {
	for {
		select {
		case <-dataChannel:
			healthCheckModel = <-dataChannel
		default:
			if dataChannel == nil {
				break
			}

			dataChannel <- hcm.performHealthCheck(healthCheckModel, hcm.defaultTimeout)
		}
		time.Sleep(frequency)
	}
}

func (hcm *HealthCheckManager) performHealthCheck(healthCheckModel models.HealthCheckModel, timeout time.Duration) models.HealthCheckModel {
	start := time.Now()

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: timeout}

	// Perform an HTTP GET request to the health check endpoint
	resp, err := client.Get(healthCheckModel.Endpoint)
	if err != nil {
		log.Println("Error performing health check for endpoint:", healthCheckModel.Endpoint, "Error:", err.Error())
		healthCheckModel.Error = err.Error()
		return healthCheckModel
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		healthCheckModel.Error = fmt.Sprintf("Unsuccessful status code: %d", resp.StatusCode)
	}

	// Update the health check model with the response status code
	healthCheckModel.Code = int32(resp.StatusCode)
	healthCheckModel.Status = resp.Status
	healthCheckModel.Checked = time.Now().Unix()
	healthCheckModel.Duration = fmt.Sprintf("%d%s", time.Since(start).Milliseconds(), "ms")

	// Return the updated health check model
	return healthCheckModel
}
