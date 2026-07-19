package service

import (
	"context"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

type HealthService struct {
	data *data.Data
}

func NewHealthService(initializedData *data.Data) *HealthService {
	return &HealthService{data: initializedData}
}

func (service *HealthService) Check(ctx context.Context) model.HealthResponse {
	storage := service.data.Ping(ctx)
	status := "ok"
	for _, storageStatus := range storage {
		if storageStatus != "up" {
			status = "degraded"
			break
		}
	}

	return model.HealthResponse{
		Status:  status,
		Storage: storage,
	}
}
