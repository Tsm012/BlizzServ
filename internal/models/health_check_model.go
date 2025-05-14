package models

type HealthCheckModel struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Code     int    `json:"code"`
	Endpoint string `json:"endpoint"`
	Checked  int64  `json:"checked"`
	Duration string `json:"duration"`
	Error    string `json:"error"`
}
