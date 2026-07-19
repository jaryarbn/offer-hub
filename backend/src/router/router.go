package router

import (
	"errors"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/ctrl"
	"offer-hub/backend/src/data"
	"offer-hub/backend/src/service"
)

func RegisterRouter(engine *gin.Engine) error {
	if engine == nil {
		return errors.New("gin engine is nil")
	}

	initializedData := data.GetData()
	if initializedData == nil {
		return errors.New("data is not initialized")
	}

	healthService := service.NewHealthService(initializedData)
	healthController := ctrl.NewHealthController(healthService)
	engine.GET("/health", healthController.Check)
	return nil
}
