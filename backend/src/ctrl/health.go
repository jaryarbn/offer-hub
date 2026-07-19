package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/service"
)

type HealthController struct {
	service *service.HealthService
}

func NewHealthController(service *service.HealthService) *HealthController {
	return &HealthController{service: service}
}

func (controller *HealthController) Check(ctx *gin.Context) {
	response := controller.service.Check(ctx.Request.Context())
	ctx.JSON(http.StatusOK, response)
}
