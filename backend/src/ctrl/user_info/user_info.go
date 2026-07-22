package user_info

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	"offer-hub/backend/src/service"
)

const (
	userInfoSuccessMessage = "success"
	unauthorizedMessage    = "未认证"
	userNotFoundMessage    = "用户不存在"
	internalErrorMessage   = "服务器内部错误"
)

type Service interface {
	GetUserInfo(context.Context, string) (model.UserInfo, error)
}

type Controller struct {
	service Service
}

func NewController(userInfoService Service) *Controller {
	return &Controller{service: userInfoService}
}

func (controller *Controller) GetUserInfo(ctx *gin.Context) {
	// TODO: read user_id from Gin context after the JWT middleware is introduced.
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		respond(ctx, 401, unauthorizedMessage, nil)
		return
	}

	userInfo, err := controller.service.GetUserInfo(ctx.Request.Context(), userID)
	if errors.Is(err, service.ErrInvalidUserID) {
		respond(ctx, 401, unauthorizedMessage, nil)
		return
	}
	if errors.Is(err, service.ErrUserInfoNotFound) {
		respond(ctx, 404, userNotFoundMessage, nil)
		return
	}
	if err != nil {
		log.Printf("get user info: %v", err)
		respond(ctx, 500, internalErrorMessage, nil)
		return
	}

	respond(ctx, 0, userInfoSuccessMessage, &userInfo)
}

func respond(ctx *gin.Context, code int, message string, data *model.UserInfo) {
	ctx.JSON(http.StatusOK, model.GetUserInfoResp{
		Code: code,
		Msg:  message,
		Data: data,
	})
}
