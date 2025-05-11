package models

type CreateHealthCheckRequestModel struct {
	Endpoint string `json:"endpoint"`
}

type CreateHealthCheckResponseModel struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}

type ListHealthCheckResponseModel struct {
	Page         int32              `json:"page"`
	Total        int32              `json:"total"`
	Size         int32              `json:"size"`
	HealthChecks []HealthCheckModel `json:"items"`
}
