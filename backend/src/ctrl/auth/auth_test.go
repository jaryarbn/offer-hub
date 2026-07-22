package auth

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	backendservice "offer-hub/backend/src/service"
)

type authServiceStub struct {
	req         model.PasswordRegisterReq
	err         error
	calls       int
	loginReq    model.PasswordLoginReq
	loginData   model.PasswordLoginData
	loginErr    error
	loginCalls  int
	logoutErr   error
	logoutToken string
	logoutCalls int
}

func (stub *authServiceStub) Login(
	_ context.Context,
	req model.PasswordLoginReq,
) (model.PasswordLoginData, error) {
	stub.loginCalls++
	stub.loginReq = req
	return stub.loginData, stub.loginErr
}

func (stub *authServiceStub) Logout(_ context.Context, token string) error {
	stub.logoutCalls++
	stub.logoutToken = token
	return stub.logoutErr
}

func (stub *authServiceStub) Register(_ context.Context, req model.PasswordRegisterReq) error {
	stub.calls++
	stub.req = req
	return stub.err
}

func TestRegisterReturnsSuccess(t *testing.T) {
	stub := &authServiceStub{}
	response := performRegisterRequest(stub, `{"username":"  testuser  ","password":"123456"}`)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != `{"code":0,"msg":"注册成功","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
	if stub.req.Username != "testuser" || stub.req.Password != "123456" {
		t.Fatalf("service request = %#v", stub.req)
	}
}

func TestRegisterRejectsInvalidParameters(t *testing.T) {
	tests := []string{
		`{"username":"","password":"123456"}`,
		`{"username":"   ","password":"123456"}`,
		`{"username":"testuser","password":"12345"}`,
		`{"username":"` + strings.Repeat("a", 51) + `","password":"123456"}`,
		`{"username":`,
	}
	for _, body := range tests {
		stub := &authServiceStub{}
		response := performRegisterRequest(stub, body)

		if response.Body.String() != `{"code":400,"msg":"用户名或密码不符合要求","data":null}` {
			t.Fatalf("request body = %q, response = %s", body, response.Body.String())
		}
		if stub.calls != 0 {
			t.Fatalf("request body = %q, service calls = %d, want 0", body, stub.calls)
		}
	}
}

func TestRegisterReturnsUsernameExists(t *testing.T) {
	stub := &authServiceStub{err: backendservice.ErrUsernameAlreadyExists}
	response := performRegisterRequest(stub, `{"username":"testuser","password":"123456"}`)

	if response.Body.String() != `{"code":500,"msg":"用户名已存在","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestRegisterMapsServiceValidationError(t *testing.T) {
	stub := &authServiceStub{err: backendservice.ErrInvalidRegisterParams}
	response := performRegisterRequest(stub, `{"username":"testuser","password":"123456"}`)

	if response.Body.String() != `{"code":400,"msg":"用户名或密码不符合要求","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestRegisterDoesNotExposeInternalErrors(t *testing.T) {
	stub := &authServiceStub{err: errors.New("database connection contains internal details")}
	response := performRegisterRequest(stub, `{"username":"testuser","password":"123456"}`)

	if response.Body.String() != `{"code":500,"msg":"服务器内部错误","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestLoginReturnsTokenAndCamelCaseUserInfo(t *testing.T) {
	stub := &authServiceStub{loginData: model.PasswordLoginData{
		Token: "signed-token",
		UserInfo: model.PasswordLoginUserInfo{
			UserID: "abc123", Username: "testuser", NickName: "测试用户",
			Avatar: "https://example.com/avatar.jpg", Sex: 1, VIP: false,
			Phone: "", Email: "", UserStatus: 1, UserType: 1,
		},
	}}
	response := performLoginRequest(stub, `{"username":"  testuser  ","password":"123456"}`)

	wantBody := `{"code":0,"msg":"登录成功","data":{"token":"signed-token","userInfo":{"userId":"abc123","username":"testuser","nickName":"测试用户","avatar":"https://example.com/avatar.jpg","sex":1,"vip":false,"phone":"","email":"","userStatus":1,"userType":1}}}`
	if response.Code != http.StatusOK || response.Body.String() != wantBody {
		t.Fatalf("HTTP status = %d, response body = %s", response.Code, response.Body.String())
	}
	if stub.loginReq.Username != "testuser" || stub.loginReq.Password != "123456" {
		t.Fatalf("login service request = %#v", stub.loginReq)
	}
}

func TestLoginRejectsInvalidParameters(t *testing.T) {
	stub := &authServiceStub{}
	response := performLoginRequest(stub, `{"username":"testuser","password":"12345"}`)

	if response.Body.String() != `{"code":400,"msg":"用户名或密码不符合要求","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
	if stub.loginCalls != 0 {
		t.Fatalf("login service calls = %d, want 0", stub.loginCalls)
	}
}

func TestLoginReturnsInvalidCredentials(t *testing.T) {
	stub := &authServiceStub{loginErr: backendservice.ErrInvalidCredentials}
	response := performLoginRequest(stub, `{"username":"testuser","password":"123456"}`)

	if response.Body.String() != `{"code":500,"msg":"用户名或密码错误","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestLoginReturnsAccountUnavailable(t *testing.T) {
	stub := &authServiceStub{loginErr: backendservice.ErrAccountUnavailable}
	response := performLoginRequest(stub, `{"username":"testuser","password":"123456"}`)

	if response.Body.String() != `{"code":500,"msg":"账号不可用","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestLoginDoesNotExposeInternalErrors(t *testing.T) {
	stub := &authServiceStub{loginErr: errors.New("database connection contains internal details")}
	response := performLoginRequest(stub, `{"username":"testuser","password":"123456"}`)

	if response.Body.String() != `{"code":500,"msg":"服务器内部错误","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestLogoutReturnsSuccess(t *testing.T) {
	stub := &authServiceStub{}
	response := performLogoutRequest(stub, "Bearer signed-token")

	wantBody := `{"code":0,"msg":"success","data":{"message":"登出成功"}}`
	if response.Code != http.StatusOK || response.Body.String() != wantBody {
		t.Fatalf("HTTP status = %d, response body = %s", response.Code, response.Body.String())
	}
	if stub.logoutToken != "signed-token" {
		t.Fatalf("logout token = %q, want signed-token", stub.logoutToken)
	}
}

func TestLogoutRequiresBearerToken(t *testing.T) {
	for _, header := range []string{"", "Basic signed-token", "Bearer", "Bearer one two"} {
		stub := &authServiceStub{}
		response := performLogoutRequest(stub, header)

		if response.Code != http.StatusOK {
			t.Fatalf("Authorization %q HTTP status = %d, want %d", header, response.Code, http.StatusOK)
		}
		if response.Body.String() != `{"code":401,"msg":"token 无效或已过期","data":null}` {
			t.Fatalf("Authorization %q response body = %s", header, response.Body.String())
		}
		if stub.logoutCalls != 0 {
			t.Fatalf("Authorization %q service calls = %d, want 0", header, stub.logoutCalls)
		}
	}
}

func TestLogoutMapsInvalidTokenAndInternalErrors(t *testing.T) {
	invalidTokenStub := &authServiceStub{logoutErr: backendservice.ErrInvalidLogoutToken}
	invalidResponse := performLogoutRequest(invalidTokenStub, "Bearer invalid-token")
	if invalidResponse.Code != http.StatusOK || invalidResponse.Body.String() != `{"code":401,"msg":"token 无效或已过期","data":null}` {
		t.Fatalf("invalid token response = HTTP %d %s", invalidResponse.Code, invalidResponse.Body.String())
	}

	internalErrorStub := &authServiceStub{logoutErr: errors.New("redis unavailable")}
	internalResponse := performLogoutRequest(internalErrorStub, "Bearer signed-token")
	if internalResponse.Code != http.StatusOK || internalResponse.Body.String() != `{"code":500,"msg":"服务器内部错误","data":null}` {
		t.Fatalf("internal error response = HTTP %d %s", internalResponse.Code, internalResponse.Body.String())
	}
}

func performRegisterRequest(service Service, body string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	controller := NewController(service)
	engine := gin.New()
	engine.POST("/auth/register", controller.Register)

	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	return response
}

func performLoginRequest(service Service, body string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	controller := NewController(service)
	engine := gin.New()
	engine.POST("/auth/login", controller.Login)

	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	return response
}

func performLogoutRequest(service Service, authorization string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	controller := NewController(service)
	engine := gin.New()
	engine.POST("/auth/logout", controller.Logout)

	request := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	return response
}
