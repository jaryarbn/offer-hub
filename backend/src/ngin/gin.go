package ngin

import (
	"errors"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/config"
)

func CreateGin() (*gin.Engine, error) {
	if config.Conf == nil {
		return nil, errors.New("config is not initialized")
	}

	return gin.Default(), nil
}
