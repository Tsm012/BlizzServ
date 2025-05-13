package internal

import (
	"github.com/Tsm012/BlizServe/internal/models"
)

// OrderedHealthCheckMap structure
type OrderedHealthCheckMap struct {
	Keys []string
	Data map[string]models.HealthCheckModel
}

func (om *OrderedHealthCheckMap) Add(key string, value models.HealthCheckModel) {
	if _, exists := om.Data[key]; !exists {
		om.Keys = append(om.Keys, key) // Preserve order
	}
	om.Data[key] = value
}

func (om *OrderedHealthCheckMap) Get(key string) (models.HealthCheckModel, bool) {
	val, exists := om.Data[key]
	return val, exists
}

func (om *OrderedHealthCheckMap) Delete(key string) {
	if _, exists := om.Data[key]; exists {
		delete(om.Data, key)
		for i, k := range om.Keys {
			if k == key {
				om.Keys = append(om.Keys[:i], om.Keys[i+1:]...)
				break
			}
		}
	}
}

func (om *OrderedHealthCheckMap) Set(key string, value models.HealthCheckModel) {
	if _, exists := om.Data[key]; !exists {
		om.Keys = append(om.Keys, key) // Add key if it doesn't exist
	}
	om.Data[key] = value // Update or set the value
}

func (om *OrderedHealthCheckMap) ToSortedList() []models.HealthCheckModel {
	var sortedValues []models.HealthCheckModel
	for _, key := range om.Keys {
		if value, exists := om.Data[key]; exists {
			sortedValues = append(sortedValues, value)
		}
	}
	return sortedValues
}
