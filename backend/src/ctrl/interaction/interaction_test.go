package interaction

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	"offer-hub/backend/src/service"
)

type interactionServiceStub struct {
	likeData    model.InteractionLikeData
	likeErr     error
	unlikeData  model.InteractionUnlikeData
	unlikeErr   error
	tagErr      error
	likeCalls   int
	unlikeCalls int
	tagCalls    int
	userID      string
	tagReq      model.TagQuestionReq
}

func (stub *interactionServiceStub) Like(
	_ context.Context,
	_ model.InteractionLikeReq,
	userID string,
) (model.InteractionLikeData, error) {
	stub.likeCalls++
	stub.userID = userID
	return stub.likeData, stub.likeErr
}

func (stub *interactionServiceStub) Unlike(
	_ context.Context,
	_ model.InteractionUnlikeReq,
	userID string,
) (model.InteractionUnlikeData, error) {
	stub.unlikeCalls++
	stub.userID = userID
	return stub.unlikeData, stub.unlikeErr
}

func (stub *interactionServiceStub) TagQuestion(
	_ context.Context,
	req model.TagQuestionReq,
	userID string,
) error {
	stub.tagCalls++
	stub.userID = userID
	stub.tagReq = req
	return stub.tagErr
}

func TestControllerLike(t *testing.T) {
	stub := &interactionServiceStub{likeData: model.InteractionLikeData{Liked: true, Count: 257}}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/like", `{"target_type":1,"target_id":"q001"}`, "user-1", controller.Like)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var response model.InteractionLikeResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != 0 || response.Msg != "success" || !response.Data.Liked || response.Data.Count != 257 {
		t.Fatalf("response = %#v", response)
	}
	if stub.likeCalls != 1 || stub.userID != "user-1" {
		t.Fatalf("like calls = %d, user_id = %q", stub.likeCalls, stub.userID)
	}
}

func TestControllerUnlike(t *testing.T) {
	stub := &interactionServiceStub{unlikeData: model.InteractionUnlikeData{Count: 256}}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/unlike", `{"target_type":3,"target_id":"c001"}`, "user-1", controller.Unlike)

	var response model.InteractionUnlikeResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != 0 || response.Msg != "success" || response.Data.Count != 256 {
		t.Fatalf("response = %#v", response)
	}
	if stub.unlikeCalls != 1 {
		t.Fatalf("unlike calls = %d, want 1", stub.unlikeCalls)
	}
}

func TestControllerTagQuestionAcceptsZero(t *testing.T) {
	stub := &interactionServiceStub{}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/tag", `{"question_id":"q001","tag":0}`, "user-1", controller.TagQuestion)

	var response model.TagQuestionResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != 0 || response.Msg != "success" || response.Data != nil {
		t.Fatalf("response = %#v", response)
	}
	if stub.tagCalls != 1 || stub.tagReq.Tag == nil || *stub.tagReq.Tag != 0 {
		t.Fatalf("tag calls = %d, request = %#v", stub.tagCalls, stub.tagReq)
	}
}

func TestControllerRejectsUnsupportedInteractionTarget(t *testing.T) {
	stub := &interactionServiceStub{}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/like", `{"target_type":2,"target_id":"e001"}`, "user-1", controller.Like)

	assertInteractionError(t, recorder, 400, "invalid parameters")
	if stub.likeCalls != 0 {
		t.Fatalf("like calls = %d, want 0", stub.likeCalls)
	}
}

func TestControllerRejectsMissingTagField(t *testing.T) {
	stub := &interactionServiceStub{}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/tag", `{"question_id":"q001"}`, "user-1", controller.TagQuestion)

	assertInteractionError(t, recorder, 400, "invalid parameters")
	if stub.tagCalls != 0 {
		t.Fatalf("tag calls = %d, want 0", stub.tagCalls)
	}
}

func TestControllerRequiresAuthenticatedUser(t *testing.T) {
	stub := &interactionServiceStub{}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/like", `{"target_type":1,"target_id":"q001"}`, "", controller.Like)

	assertInteractionError(t, recorder, 401, "未认证")
	if stub.likeCalls != 0 {
		t.Fatalf("like calls = %d, want 0", stub.likeCalls)
	}
}

func TestControllerMapsTargetNotFound(t *testing.T) {
	stub := &interactionServiceStub{likeErr: service.ErrInteractionTargetNotFound}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/like", `{"target_type":1,"target_id":"missing"}`, "user-1", controller.Like)

	assertInteractionError(t, recorder, 404, "target not found")
}

func TestControllerMapsTagValidationError(t *testing.T) {
	stub := &interactionServiceStub{tagErr: service.ErrInvalidQuestionTag}
	controller := NewController(stub)
	recorder := performInteractionRequest(t, http.MethodPost, "/tag", `{"question_id":"q001","tag":4}`, "user-1", controller.TagQuestion)

	assertInteractionError(t, recorder, 400, "invalid parameters")
}

func performInteractionRequest(
	t *testing.T,
	method string,
	path string,
	body string,
	userID string,
	handler gin.HandlerFunc,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Handle(method, path, handler)

	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	if userID != "" {
		request.Header.Set("user_id", userID)
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func assertInteractionError(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	wantCode int,
	wantMessage string,
) {
	t.Helper()
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Code != wantCode || response.Msg != wantMessage || response.Data != nil {
		t.Fatalf("response = %#v, want code=%d msg=%q data=nil", response, wantCode, wantMessage)
	}
}
