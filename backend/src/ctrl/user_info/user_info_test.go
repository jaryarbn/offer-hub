package user_info

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	backendservice "offer-hub/backend/src/service"
)

type userInfoServiceStub struct {
	userID string
	data   model.UserInfo
	err    error
	calls  int
}

func (stub *userInfoServiceStub) GetUserInfo(_ context.Context, userID string) (model.UserInfo, error) {
	stub.calls++
	stub.userID = userID
	return stub.data, stub.err
}

func TestGetUserInfoReadsUserIDHeaderAndReturnsSnakeCaseFields(t *testing.T) {
	stub := &userInfoServiceStub{data: model.UserInfo{
		UserID:       "abc123",
		Username:     "testuser",
		NickName:     "测试用户",
		Avatar:       "https://example.com/avatar.jpg",
		VIP:          false,
		Sex:          1,
		Phone:        "",
		Email:        "",
		Introduction: "个人简介",
		AvatarURL:    "https://example.com/avatar.jpg",
		UserStatus:   1,
		UserType:     1,
		CreateTime:   "2024-01-01 12:00:00",
		UpdateTime:   "2024-01-01 12:00:00",
	}}
	response := performGetUserInfoRequest(stub, "  abc123  ")

	wantBody := `{"code":0,"msg":"success","data":{"user_id":"abc123","username":"testuser","nick_name":"测试用户","avatar":"https://example.com/avatar.jpg","vip":false,"sex":1,"phone":"","email":"","introduction":"个人简介","avatar_url":"https://example.com/avatar.jpg","user_status":1,"user_type":1,"create_time":"2024-01-01 12:00:00","update_time":"2024-01-01 12:00:00"}}`
	if response.Code != http.StatusOK || response.Body.String() != wantBody {
		t.Fatalf("HTTP status = %d, response body = %s", response.Code, response.Body.String())
	}
	if stub.userID != "abc123" {
		t.Fatalf("service user_id = %q, want abc123", stub.userID)
	}
}

func TestGetUserInfoRequiresUserIDHeader(t *testing.T) {
	stub := &userInfoServiceStub{}
	response := performGetUserInfoRequest(stub, "")

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != `{"code":401,"msg":"未认证","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
	if stub.calls != 0 {
		t.Fatalf("service calls = %d, want 0", stub.calls)
	}
}

func TestGetUserInfoMapsNotFoundAndStorageErrors(t *testing.T) {
	notFoundStub := &userInfoServiceStub{err: backendservice.ErrUserInfoNotFound}
	notFoundResponse := performGetUserInfoRequest(notFoundStub, "abc123")
	if notFoundResponse.Body.String() != `{"code":404,"msg":"用户不存在","data":null}` {
		t.Fatalf("not found response = %s", notFoundResponse.Body.String())
	}

	storageStub := &userInfoServiceStub{err: errors.New("database unavailable")}
	storageResponse := performGetUserInfoRequest(storageStub, "abc123")
	if storageResponse.Body.String() != `{"code":500,"msg":"服务器内部错误","data":null}` {
		t.Fatalf("storage error response = %s", storageResponse.Body.String())
	}
}

func performGetUserInfoRequest(service Service, userID string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	controller := NewController(service)
	engine := gin.New()
	engine.GET("/api/v1/user_info/get_user_info", controller.GetUserInfo)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/user_info/get_user_info", nil)
	if userID != "" {
		request.Header.Set("user_id", userID)
	}
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	return response
}
