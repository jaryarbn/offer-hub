package model

type HealthResponse struct {
	Status  string            `json:"status"`
	Storage map[string]string `json:"storage"`
}
