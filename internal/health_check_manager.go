package internal

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"time"

	"github.com/Tsm012/BlizServe/internal/models"
)

type HealthCheckManager struct {
	dataChannels   map[string]chan models.HealthCheckModel
	healthChecks   map[string]models.HealthCheckModel
	endpointLookup map[string]string
}

func NewHealthCheckManager() HealthCheckManager {
	manager := HealthCheckManager{
		healthChecks:   map[string]models.HealthCheckModel{},
		endpointLookup: map[string]string{},
		dataChannels:   map[string]chan models.HealthCheckModel{},
	}
	manager.LoadFromFile()
	return manager
}

func (hcm *HealthCheckManager) AddHealthCheck(healthCheckModel models.HealthCheckModel) {
	hcm.healthChecks[healthCheckModel.ID] = healthCheckModel
	hcm.endpointLookup[healthCheckModel.Endpoint] = healthCheckModel.ID
	dataChannel := make(chan models.HealthCheckModel)
	hcm.dataChannels[healthCheckModel.ID] = dataChannel
	go monitor(healthCheckModel, dataChannel)
	go listener(dataChannel)
	hcm.SaveToFile()
}

func (hcm *HealthCheckManager) ListHealthChecks() []models.HealthCheckModel {
	var healthChecks []models.HealthCheckModel
	for healthCheck := range maps.Values(hcm.healthChecks) {
		healthChecks = append(healthChecks, healthCheck)
	}
	return healthChecks
}

func (hcm *HealthCheckManager) GetHealthCheck(serverId string) (models.HealthCheckModel, error) {
	healthCheck := <-hcm.dataChannels[serverId]
	return healthCheck, nil
}

func (hcm *HealthCheckManager) DeleteHealthCheck(serverId string) {
	healthCheck, ok := hcm.healthChecks[serverId]
	if !ok {
		return
	}

	close(hcm.dataChannels[healthCheck.ID])
	delete(hcm.dataChannels, healthCheck.ID)
	delete(hcm.healthChecks, healthCheck.ID)
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
	for _, healthCheck := range hcm.healthChecks {
		hcm.endpointLookup[healthCheck.Endpoint] = healthCheck.ID
		dataChannel := make(chan models.HealthCheckModel)
		hcm.dataChannels[healthCheck.ID] = dataChannel
		go monitor(healthCheck, dataChannel)
		go listener(dataChannel)
	}
}

func listener(dataChannel chan models.HealthCheckModel) {
	for healthCheck := range dataChannel {
		// Print the ID and the received HealthCheckModel
		healthCheckData, err := json.Marshal(healthCheck)
		if err != nil {
			println("Error marshaling HealthCheckModel:", err.Error())
			continue
		}
		println("Received update for ", string(healthCheckData))
	}
}

func monitor(healthCheckModel models.HealthCheckModel, dataChannel chan models.HealthCheckModel) {
	for {
		select {
		case <-dataChannel:
			healthCheckModel = <-dataChannel
		default:
			if dataChannel == nil {
				break
			}
			performHealthCheck(&healthCheckModel)
			dataChannel <- healthCheckModel
		}
		time.Sleep(2 * time.Second)
	}
}

func performHealthCheck(healthCheckModel *models.HealthCheckModel) {
	start := time.Now()

	// Perform an HTTP GET request to the health check endpoint
	resp, err := http.Get(healthCheckModel.Endpoint)
	if err != nil {
		println("Error performing health check for endpoint:", healthCheckModel.Endpoint, "Error:", err.Error())
		return
	}
	defer resp.Body.Close()

	// Update the health check model with the response status code
	healthCheckModel.Code = int32(resp.StatusCode)
	healthCheckModel.Status = resp.Status
	healthCheckModel.Checked = time.Now().Unix()
	healthCheckModel.Duration = fmt.Sprintf("%d%s", time.Since(start).Milliseconds(), "ms")
	healthCheckModel.Error = ""
}
