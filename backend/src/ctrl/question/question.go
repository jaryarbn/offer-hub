package question

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
	GetAllQuestionList(context.Context, model.GetAllQuestionListReq) ([]model.GetAllQuestionListData, error)
	ListQuestions(context.Context, model.ListQuestionReq, string) (model.ListQuestionResponseData, error)
	ListQuestionsMeta(context.Context, model.ListQuestionMetaReq, string) (model.ListQuestionMetaResp, error)
	GetQuestionDetail(context.Context, model.GetQuestionDetailReq, string) (model.QuestionDetail, error)
	GetHotQuestions(context.Context, model.GetHotQuestionsReq) (model.GetHotQuestionsResp, error)
}

type Controller struct {
	service Service
}

func NewController(questionService Service) *Controller {
	return &Controller{service: questionService}
}

const (
	defaultPage      = 1
	defaultPageSize  = 20
	defaultSortBy    = "order"
	defaultSortOrder = "asc"
	defaultHotLimit  = 10
)

func (controller *Controller) ListQuestions(ctx *gin.Context) {
	var req model.ListQuestionReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}
	applyListQuestionDefaults(&req)

	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	data, err := controller.service.ListQuestions(ctx.Request.Context(), req, userID)
	if err != nil {
		log.Printf("list questions: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, model.ListQuestionResp{Code: 0, Msg: "success", Data: data})
}

func applyListQuestionDefaults(req *model.ListQuestionReq) {
	if req.Page <= 0 {
		req.Page = defaultPage
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	if req.SortBy == "" {
		req.SortBy = defaultSortBy
	}
	if req.SortOrder == "" {
		req.SortOrder = defaultSortOrder
	}
}

func (controller *Controller) ListQuestionsMeta(ctx *gin.Context) {
	var req model.ListQuestionMetaReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}
	applyListQuestionDefaults(&req)

	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	response, err := controller.service.ListQuestionsMeta(ctx.Request.Context(), req, userID)
	if err != nil {
		log.Printf("list question metadata: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListQuestionMeta is the route-facing name for the question metadata handler.
func (controller *Controller) ListQuestionMeta(ctx *gin.Context) {
	controller.ListQuestionsMeta(ctx)
}

func (controller *Controller) GetQuestionDetail(ctx *gin.Context) {
	var req model.GetQuestionDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}
	req.QuestionID = strings.TrimSpace(req.QuestionID)
	if req.QuestionID == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "question_id is required", "data": nil})
		return
	}

	userID := strings.TrimSpace(ctx.GetHeader("user_id"))
	question, err := controller.service.GetQuestionDetail(ctx.Request.Context(), req, userID)
	if errors.Is(err, service.ErrQuestionNotFound) {
		ctx.JSON(http.StatusOK, gin.H{"code": 404, "msg": "question not found", "data": nil})
		return
	}
	if err != nil {
		log.Printf("get question detail: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, model.GetQuestionDetailResp{Code: 0, Msg: "success", Data: question})
}

// GetQuestionBankSeries is the route-facing name for the question bank series handler.
func (controller *Controller) GetQuestionBankSeries(ctx *gin.Context) {
	controller.GetAllQuestionList(ctx)
}

// GetHotQuestionList is the route-facing name for the hot question handler.
func (controller *Controller) GetHotQuestionList(ctx *gin.Context) {
	controller.GetHotQuestions(ctx)
}

func (controller *Controller) GetHotQuestions(ctx *gin.Context) {
	var req model.GetHotQuestionsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "msg": "invalid parameters", "data": nil})
		return
	}
	if req.Limit <= 0 {
		req.Limit = defaultHotLimit
	}

	response, err := controller.service.GetHotQuestions(ctx.Request.Context(), req)
	if err != nil {
		log.Printf("get hot questions: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "internal server error", "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (controller *Controller) GetAllQuestionList(ctx *gin.Context) {
	var req model.GetAllQuestionListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, model.GetAllQuestionListResp{
			Code: 400,
			Msg:  "invalid parameters",
			Data: nil,
		})
		return
	}

	data, err := controller.service.GetAllQuestionList(ctx.Request.Context(), req)
	if err != nil {
		log.Printf("get all question list: %v", err)
		ctx.JSON(http.StatusOK, model.GetAllQuestionListResp{
			Code: 500,
			Msg:  "internal server error",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.GetAllQuestionListResp{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}
