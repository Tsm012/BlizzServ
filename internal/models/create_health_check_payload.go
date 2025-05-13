package models

type CreateHealthCheckRequestModel struct {
	Endpoint string `json:"endpoint"`
}

type CreateHealthCheckResponseModel struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}

type ListHealthCheckResponseModel struct {
	Page         int                `json:"page"`
	Total        int                `json:"total"`
	Size         int                `json:"size"`
	HealthChecks []HealthCheckModel `json:"items"`
}
