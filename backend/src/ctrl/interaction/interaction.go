package interaction

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

type Service interface {
	Like(context.Context, model.InteractionLikeReq, string) (model.InteractionLikeData, error)
	Unlike(context.Context, model.InteractionUnlikeReq, string) (model.InteractionUnlikeData, error)
	TagQuestion(context.Context, model.TagQuestionReq, string) error
}

type Controller struct {
	service Service
}

func NewController(interactionService Service) *Controller {
	return &Controller{service: interactionService}
}

func (controller *Controller) Like(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		respondError(ctx, 401, "未认证")
		return
	}

	var req model.InteractionLikeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, 400, "invalid parameters")
		return
	}

	responseData, err := controller.service.Like(ctx.Request.Context(), req, userID)
	if respondServiceError(ctx, "like interaction", err) {
		return
	}
	ctx.JSON(http.StatusOK, model.InteractionLikeResp{
		Code: 0,
		Msg:  "success",
		Data: responseData,
	})
}

func (controller *Controller) Unlike(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		respondError(ctx, 401, "未认证")
		return
	}

	var req model.InteractionUnlikeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, 400, "invalid parameters")
		return
	}

	responseData, err := controller.service.Unlike(ctx.Request.Context(), req, userID)
	if respondServiceError(ctx, "unlike interaction", err) {
		return
	}
	ctx.JSON(http.StatusOK, model.InteractionUnlikeResp{
		Code: 0,
		Msg:  "success",
		Data: responseData,
	})
}

func (controller *Controller) TagQuestion(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		respondError(ctx, 401, "未认证")
		return
	}

	var req model.TagQuestionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, 400, "invalid parameters")
		return
	}

	err := controller.service.TagQuestion(ctx.Request.Context(), req, userID)
	if respondServiceError(ctx, "tag question", err) {
		return
	}
	ctx.JSON(http.StatusOK, model.TagQuestionResp{Code: 0, Msg: "success", Data: nil})
}

func respondServiceError(ctx *gin.Context, operation string, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, service.ErrInvalidInteractionUserID):
		respondError(ctx, 401, "未认证")
	case errors.Is(err, service.ErrInvalidInteractionTargetType),
		errors.Is(err, service.ErrInvalidInteractionTargetID),
		errors.Is(err, service.ErrInvalidQuestionID),
		errors.Is(err, service.ErrInvalidQuestionTag):
		respondError(ctx, 400, "invalid parameters")
	case errors.Is(err, service.ErrInteractionTargetNotFound):
		respondError(ctx, 404, "target not found")
	default:
		log.Printf("%s: %v", operation, err)
		respondError(ctx, 500, "internal server error")
	}
	return true
}

func respondError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, gin.H{"code": code, "msg": message, "data": nil})
}
