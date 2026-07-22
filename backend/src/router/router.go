package router

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/ctrl"
	authctrl "offer-hub/backend/src/ctrl/auth"
	questionctrl "offer-hub/backend/src/ctrl/question"
	userinfoctrl "offer-hub/backend/src/ctrl/user_info"
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
	if config.Conf == nil {
		return errors.New("config is not initialized")
	}

	healthService := service.NewHealthService(initializedData)
	healthController := ctrl.NewHealthController(healthService)
	engine.GET("/health", healthController.Check)

	authService, err := service.NewAuthService(initializedData, config.Conf.JWT)
	if err != nil {
		return fmt.Errorf("initialize auth service: %w", err)
	}
	authController := authctrl.NewController(authService)
	engine.POST("/auth/register", authController.Register)
	engine.POST("/auth/login", authController.Login)
	engine.POST("/auth/logout", authController.Logout)

	userInfoService := service.NewUserInfoService(initializedData)
	userInfoController := userinfoctrl.NewController(userInfoService)
	engine.GET("/api/v1/user_info/get_user_info", userInfoController.GetUserInfo)

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
