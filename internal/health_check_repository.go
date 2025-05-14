package internal

import (
	"github.com/Tsm012/BlizServe/internal/models"
)

// Simple implementation of an ordered map data structure
type HealthCheckRepository struct {
	Keys []string
	Data map[string]models.HealthCheckModel
}

func (om *HealthCheckRepository) Add(key string, value models.HealthCheckModel) {
	if _, exists := om.Data[key]; !exists {
		om.Keys = append(om.Keys, key)
	}
	om.Data[key] = value
}

func (om *HealthCheckRepository) Get(key string) (models.HealthCheckModel, bool) {
	val, exists := om.Data[key]
	return val, exists
}

func (om *HealthCheckRepository) Delete(key string) {
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

func (om *HealthCheckRepository) Set(key string, value models.HealthCheckModel) {
	if _, exists := om.Data[key]; !exists {
		om.Keys = append(om.Keys, key)
	}
	om.Data[key] = value
}

func (om *HealthCheckRepository) ToSortedList() []models.HealthCheckModel {
	var sortedValues []models.HealthCheckModel
	for _, key := range om.Keys {
		if value, exists := om.Data[key]; exists {
			sortedValues = append(sortedValues, value)
		}
	}
	return sortedValues
}
