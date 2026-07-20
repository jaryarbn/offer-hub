package question

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	backendservice "offer-hub/backend/src/service"
)

type questionServiceStub struct {
	req         model.GetAllQuestionListReq
	listReq     model.ListQuestionReq
	metaReq     model.ListQuestionMetaReq
	metaResp    model.ListQuestionMetaResp
	detailReq   model.GetQuestionDetailReq
	detail      model.OneQuestion
	detailErr   error
	detailCalls int
	hotReq      model.GetHotQuestionsReq
	hotResp     model.GetHotQuestionsResp
}

func (stub *questionServiceStub) GetAllQuestionList(
	_ context.Context,
	req model.GetAllQuestionListReq,
) ([]model.GetAllQuestionListData, error) {
	stub.req = req
	return make([]model.GetAllQuestionListData, 0), nil
}

func (stub *questionServiceStub) ListQuestionsMeta(
	_ context.Context,
	req model.ListQuestionMetaReq,
) (model.ListQuestionMetaResp, error) {
	stub.metaReq = req
	return stub.metaResp, nil
}

func (stub *questionServiceStub) GetQuestionDetail(
	_ context.Context,
	req model.GetQuestionDetailReq,
) (model.OneQuestion, error) {
	stub.detailCalls++
	stub.detailReq = req
	return stub.detail, stub.detailErr
}

func (stub *questionServiceStub) GetHotQuestions(
	_ context.Context,
	req model.GetHotQuestionsReq,
) (model.GetHotQuestionsResp, error) {
	stub.hotReq = req
	return stub.hotResp, nil
}

func (stub *questionServiceStub) ListQuestions(
	_ context.Context,
	req model.ListQuestionReq,
) (model.ListQuestionResponseData, error) {
	stub.listReq = req
	return model.ListQuestionResponseData{
		List: make([]model.OneQuestion, 0),
	}, nil
}

func TestListQuestionsBindsQueryAndAppliesDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/list", controller.ListQuestions)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/question/list?bank_id=bank-1&keyword=Go&difficulty=2&tags=Go&tags=%E5%B9%B6%E5%8F%91&job_name=%E5%90%8E%E7%AB%AF%E5%BC%80%E5%8F%91&user_tag=1",
		nil,
	)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	wantReq := model.ListQuestionReq{
		BankID: "bank-1", Keyword: "Go", Difficulty: 2,
		Tags: []string{"Go", "并发"}, JobName: "后端开发", UserTag: 1,
		SortBy: "order", SortOrder: "asc", Page: 1, PageSize: 20,
	}
	if !reflect.DeepEqual(stub.listReq, wantReq) {
		t.Fatalf("bound request = %#v, want %#v", stub.listReq, wantReq)
	}

	var body model.ListQuestionResp
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != 0 || body.Msg != "success" || body.Data.List == nil || len(body.Data.List) != 0 {
		t.Fatalf("response body = %#v", body)
	}
}

func TestListQuestionsMetaBindsQueryAndReturnsMinimalItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{
		metaResp: model.ListQuestionMetaResp{
			Code: 0,
			Msg:  "success",
			Data: model.ListQuestionMetaResponseData{
				Total: 1,
				List:  []model.QuestionMetaInfo{{QuestionID: "question-1", Title: "Go 并发"}},
			},
		},
	}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/meta/list", controller.ListQuestionsMeta)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/question/meta/list?bank_id=bank-1&keyword=Go&difficulty=2&tags=Go&tags=%E5%B9%B6%E5%8F%91&job_name=%E5%90%8E%E7%AB%AF%E5%BC%80%E5%8F%91&user_tag=1&sort_by=dislike_count&sort_order=desc&page=2&page_size=5",
		nil,
	)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	wantReq := model.ListQuestionMetaReq{
		BankID: "bank-1", Keyword: "Go", Difficulty: 2,
		Tags: []string{"Go", "并发"}, JobName: "后端开发", UserTag: 1,
		SortBy: "dislike_count", SortOrder: "desc", Page: 2, PageSize: 5,
	}
	if !reflect.DeepEqual(stub.metaReq, wantReq) {
		t.Fatalf("bound request = %#v, want %#v", stub.metaReq, wantReq)
	}
	wantBody := `{"code":0,"msg":"success","data":{"total":1,"list":[{"question_id":"question-1","title":"Go 并发"}]}}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestGetQuestionDetailValidatesRequiredQuestionID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/detail", controller.GetQuestionDetail)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/question/detail?question_id=++", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	wantBody := `{"code":400,"data":null,"msg":"question_id is required"}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
	if stub.detailCalls != 0 {
		t.Fatalf("service calls = %d, want 0", stub.detailCalls)
	}
}

func TestGetQuestionDetailReturnsCompleteQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{detail: model.OneQuestion{
		QuestionID: "question-1",
		Title:      "Go 并发",
		Content:    "完整题目内容",
		BankList:   []string{},
		Tags:       []string{},
		Status:     1,
	}}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/detail", controller.GetQuestionDetail)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/question/detail?question_id=question-1", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	var body model.GetQuestionDetailResp
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != 0 || body.Msg != "success" || body.Data.Content != "完整题目内容" {
		t.Fatalf("response body = %#v", body)
	}
	if stub.detailReq.QuestionID != "question-1" {
		t.Fatalf("bound question_id = %q", stub.detailReq.QuestionID)
	}
}

func TestGetQuestionDetailReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{detailErr: backendservice.ErrQuestionNotFound}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/detail", controller.GetQuestionDetail)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/question/detail?question_id=missing", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	wantBody := `{"code":404,"data":null,"msg":"question not found"}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestGetHotQuestionsAppliesDefaultLimitAndReturnsMinimalItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{hotResp: model.GetHotQuestionsResp{
		Code: 0,
		Msg:  "success",
		Data: model.HotQuestionListData{List: []model.HotQuestionInfo{{
			QuestionID: "question-1",
			BankList:   []string{"bank-1"},
			Title:      "Go 并发",
			ViewCount:  100,
		}}},
	}}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/hot/list", controller.GetHotQuestions)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/question/hot/list?job_name=%E5%90%8E%E7%AB%AF%E5%BC%80%E5%8F%91",
		nil,
	)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if stub.hotReq.Limit != 10 || stub.hotReq.JobName != "后端开发" {
		t.Fatalf("bound request = %#v", stub.hotReq)
	}
	wantBody := `{"code":0,"msg":"success","data":{"list":[{"question_id":"question-1","bank_list":["bank-1"],"title":"Go 并发","view_count":100}]}}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestGetAllQuestionListBindsQueryAndReturnsEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &questionServiceStub{}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/question/all/list", controller.GetAllQuestionList)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/question/all/list?job_name=%E5%90%8E%E7%AB%AF%E5%BC%80%E5%8F%91",
		nil,
	)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	if stub.req.JobName != "后端开发" {
		t.Fatalf("bound job_name = %q, want %q", stub.req.JobName, "后端开发")
	}
	wantBody := `{"code":0,"msg":"success","data":[]}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}
