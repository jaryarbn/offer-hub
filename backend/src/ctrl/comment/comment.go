package comment

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
	defaultPage           = 1
	defaultPageSize       = 20
	defaultSortBy         = "create_time"
	defaultSortOrder      = "desc"
	defaultSubCommentPage = 1
	defaultSubCommentSize = 5
)

type Service interface {
	AddComment(context.Context, model.AddCommentReq, string) (model.AddCommentData, error)
	DeleteComment(context.Context, model.DeleteCommentReq, string) error
	UpdateComment(context.Context, model.UpdateCommentReq, string) (model.UpdateCommentData, error)
	ListComments(context.Context, model.ListCommentsReq, string) (model.ListCommentsData, error)
}

type Controller struct {
	service Service
}

func NewController(commentService Service) *Controller {
	return &Controller{service: commentService}
}

func (controller *Controller) AddComment(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未认证", "data": nil})
		return
	}

	var req model.AddCommentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}

	responseData, err := controller.service.AddComment(ctx.Request.Context(), req, userID)
	if errors.Is(err, service.ErrInvalidCommentUserID) {
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未认证", "data": nil})
		return
	}
	if errors.Is(err, service.ErrInvalidCommentTargetType) ||
		errors.Is(err, service.ErrInvalidCommentTargetID) ||
		errors.Is(err, service.ErrInvalidCommentContent) {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}
	if err != nil {
		log.Printf("add comment: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, model.AddCommentResp{
		Code: 0,
		Msg:  "success",
		Data: responseData,
	})
}

func (controller *Controller) DeleteComment(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未认证", "data": nil})
		return
	}

	var req model.DeleteCommentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}

	err := controller.service.DeleteComment(ctx.Request.Context(), req, userID)
	switch {
	case errors.Is(err, service.ErrInvalidCommentUserID):
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未认证", "data": nil})
	case errors.Is(err, service.ErrInvalidCommentID):
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
	case errors.Is(err, service.ErrCommentForbidden):
		ctx.JSON(http.StatusOK, gin.H{"code": 403, "msg": "forbidden", "data": nil})
	case errors.Is(err, service.ErrCommentNotFound):
		ctx.JSON(http.StatusOK, gin.H{"code": 404, "msg": "comment not found", "data": nil})
	case err != nil:
		log.Printf("delete comment: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
	default:
		ctx.JSON(http.StatusOK, model.DeleteCommentResp{Code: 0, Msg: "success"})
	}
}

func (controller *Controller) UpdateComment(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	if userID == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未认证", "data": nil})
		return
	}

	var req model.UpdateCommentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}

	responseData, err := controller.service.UpdateComment(ctx.Request.Context(), req, userID)
	switch {
	case errors.Is(err, service.ErrInvalidCommentUserID):
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未认证", "data": nil})
	case errors.Is(err, service.ErrInvalidCommentID),
		errors.Is(err, service.ErrInvalidCommentContent):
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
	case errors.Is(err, service.ErrCommentForbidden):
		ctx.JSON(http.StatusOK, gin.H{"code": 403, "msg": "forbidden", "data": nil})
	case errors.Is(err, service.ErrCommentNotFound):
		ctx.JSON(http.StatusOK, gin.H{"code": 404, "msg": "comment not found", "data": nil})
	case err != nil:
		log.Printf("update comment: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
	default:
		ctx.JSON(http.StatusOK, model.UpdateCommentResp{
			Code: 0,
			Msg:  "success",
			Data: responseData,
		})
	}
}

func (controller *Controller) ListComments(ctx *gin.Context) {
	var req model.ListCommentsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}

	req.TargetID = strings.TrimSpace(req.TargetID)
	req.ParentID = strings.TrimSpace(req.ParentID)
	if req.TargetID == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}
	applyListCommentDefaults(&req)

	responseData, err := controller.service.ListComments(
		ctx.Request.Context(),
		req,
		strings.TrimSpace(ctx.GetHeader("user_id")),
	)
	if err != nil {
		log.Printf("list comments: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, model.ListCommentsResp{
		Code: 0,
		Msg:  "success",
		Data: responseData,
	})
}

func applyListCommentDefaults(req *model.ListCommentsReq) {
	if req.Page <= 0 {
		req.Page = defaultPage
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	if req.SubCommentPage <= 0 {
		req.SubCommentPage = defaultSubCommentPage
	}
	if req.SubCommentSize <= 0 {
		req.SubCommentSize = defaultSubCommentSize
	}
	if req.SortBy == "" {
		req.SortBy = defaultSortBy
	}
	if req.SortOrder == "" {
		req.SortOrder = defaultSortOrder
	}
}
