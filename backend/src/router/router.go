package router

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/ctrl"
	authctrl "offer-hub/backend/src/ctrl/auth"
	commentctrl "offer-hub/backend/src/ctrl/comment"
	interactionctrl "offer-hub/backend/src/ctrl/interaction"
	questionctrl "offer-hub/backend/src/ctrl/question"
	ctrltools "offer-hub/backend/src/ctrl/tools"
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
	authRateLimit := ctrltools.AuthRateLimitMiddleware()
	engine.POST("/auth/register", authRateLimit, authController.Register)
	engine.POST("/auth/login", authRateLimit, authController.Login)
	engine.POST("/auth/logout", ctrltools.JWTAuthMiddleware(), authController.Logout)

	apiV1Router := engine.Group("/api/v1")
	apiV1Router.Use(
		ctrltools.SoftJWTAuthMiddleware(),
		ctrltools.RateLimitMiddleware(&config.Conf.RateLimit),
	)

	userInfoService := service.NewUserInfoService(initializedData)
	userInfoController := userinfoctrl.NewController(userInfoService)
	userInfoRouter := apiV1Router.Group("/user_info")
	userInfoRouter.Use(ctrltools.JWTAuthMiddleware())
	userInfoRouter.GET("/get_user_info", userInfoController.GetUserInfo)

	questionService := service.NewQuestionService(initializedData)
	questionController := questionctrl.NewController(questionService)
	questionRouter := apiV1Router.Group("/question")
	questionRouter.GET("/all/list", questionController.GetQuestionBankSeries)
	questionRouter.GET("/list", questionController.ListQuestions)
	questionRouter.GET("/meta/list", questionController.ListQuestionMeta)
	questionRouter.GET("/detail", questionController.GetQuestionDetail)
	questionRouter.GET("/hot/list", questionController.GetHotQuestionList)

	commentService := service.NewCommentService(initializedData)
	commentController := commentctrl.NewController(commentService)
	commentRouter := apiV1Router.Group("/comment")
	commentRouter.Use(ctrltools.JWTAuthMiddleware())
	commentRouter.POST("/add", commentController.AddComment)
	commentRouter.POST("/delete", commentController.DeleteComment)
	commentRouter.POST("/update", commentController.UpdateComment)

	openRouter := apiV1Router.Group("/open")
	openRouter.GET("/list_comments", commentController.ListComments)

	interactionService := service.NewInteractionService(initializedData)
	interactionController := interactionctrl.NewController(interactionService)
	interactionRouter := apiV1Router.Group("/interaction")
	interactionRouter.Use(ctrltools.JWTAuthMiddleware())
	interactionRouter.POST("/like", interactionController.Like)
	interactionRouter.POST("/unlike", interactionController.Unlike)

	safeRouter := apiV1Router.Group("/safe")
	safeRouter.Use(ctrltools.JWTAuthMiddleware())
	safeRouter.POST("/tag_question", interactionController.TagQuestion)
	return nil
}
