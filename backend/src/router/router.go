package router

import (
	"errors"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/ctrl"
	questionctrl "offer-hub/backend/src/ctrl/question"
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

	questionService := service.NewQuestionService(initializedData)
	questionController := questionctrl.NewController(questionService)
	questionRouter := engine.Group("/api/v1/question")
	questionRouter.GET("/all/list", questionController.GetQuestionBankSeries)
	questionRouter.GET("/list", questionController.ListQuestions)
	questionRouter.GET("/meta/list", questionController.ListQuestionMeta)
	questionRouter.GET("/detail", questionController.GetQuestionDetail)
	questionRouter.GET("/hot/list", questionController.GetHotQuestionList)
	return nil
}
