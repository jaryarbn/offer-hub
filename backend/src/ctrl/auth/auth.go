package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	"offer-hub/backend/src/service"
)

const (
	registerSuccessMessage       = "注册成功"
	invalidRegisterParamsMessage = "用户名或密码不符合要求"
	usernameExistsMessage        = "用户名已存在"
	loginSuccessMessage          = "登录成功"
	invalidCredentialsMessage    = "用户名或密码错误"
	accountUnavailableMessage    = "账号不可用"
	logoutSuccessMessage         = "success"
	logoutConfirmationMessage    = "登出成功"
	invalidTokenMessage          = "token 无效或已过期"
	internalServerErrorMessage   = "服务器内部错误"
)

type Service interface {
	Register(context.Context, model.PasswordRegisterReq) error
	Login(context.Context, model.PasswordLoginReq) (model.PasswordLoginData, error)
	Logout(context.Context, string) error
}

type Controller struct {
	service Service
}

func NewController(authService Service) *Controller {
	return &Controller{service: authService}
}

func (controller *Controller) Login(ctx *gin.Context) {
	var req model.PasswordLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondLoginError(ctx, 400, invalidRegisterParamsMessage)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || utf8.RuneCountInString(req.Username) > 50 || utf8.RuneCountInString(req.Password) < 6 {
		respondLoginError(ctx, 400, invalidRegisterParamsMessage)
		return
	}

	loginData, err := controller.service.Login(ctx.Request.Context(), req)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, model.PasswordLoginResp{
			Code: 0,
			Msg:  loginSuccessMessage,
			Data: &loginData,
		})
	case errors.Is(err, service.ErrInvalidLoginParams):
		respondLoginError(ctx, 400, invalidRegisterParamsMessage)
	case errors.Is(err, service.ErrInvalidCredentials):
		respondLoginError(ctx, 500, invalidCredentialsMessage)
	case errors.Is(err, service.ErrAccountUnavailable):
		respondLoginError(ctx, 500, accountUnavailableMessage)
	default:
		log.Printf("login user: %v", err)
		respondLoginError(ctx, 500, internalServerErrorMessage)
	}
}

func (controller *Controller) Register(ctx *gin.Context) {
	var req model.PasswordRegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respond(ctx, 400, invalidRegisterParamsMessage)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || utf8.RuneCountInString(req.Username) > 50 || utf8.RuneCountInString(req.Password) < 6 {
		respond(ctx, 400, invalidRegisterParamsMessage)
		return
	}

	err := controller.service.Register(ctx.Request.Context(), req)
	switch {
	case err == nil:
		respond(ctx, 0, registerSuccessMessage)
	case errors.Is(err, service.ErrInvalidRegisterParams):
		respond(ctx, 400, invalidRegisterParamsMessage)
	case errors.Is(err, service.ErrUsernameAlreadyExists):
		respond(ctx, 500, usernameExistsMessage)
	default:
		log.Printf("register user: %v", err)
		respond(ctx, 500, internalServerErrorMessage)
	}
}

func (controller *Controller) Logout(ctx *gin.Context) {
	parts := strings.Fields(ctx.GetHeader("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		respondLogoutError(ctx, 401, invalidTokenMessage)
		return
	}

	if err := controller.service.Logout(ctx.Request.Context(), parts[1]); err != nil {
		if errors.Is(err, service.ErrInvalidLogoutToken) {
			respondLogoutError(ctx, 401, invalidTokenMessage)
			return
		}
		log.Printf("logout user: %v", err)
		respondLogoutError(ctx, 500, internalServerErrorMessage)
		return
	}

	ctx.JSON(http.StatusOK, model.PasswordLogoutResp{
		Code: 0,
		Msg:  logoutSuccessMessage,
		Data: &model.PasswordLogoutData{Message: logoutConfirmationMessage},
	})
}

func respond(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, model.PasswordRegisterResp{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}

func respondLoginError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, model.PasswordLoginResp{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}

func respondLogoutError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, model.PasswordLogoutResp{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}
