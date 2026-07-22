package tools

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"offer-hub/backend/src/config"
)

const middlewareTestSecret = "middleware-test-secret"

type testClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func TestJWTAuthMiddlewareAllowsValidTokenAndOverwritesUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token := signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
		UserID: "verified-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, middlewareTestSecret)

	checkerCalls := 0
	engine := gin.New()
	engine.Use(newJWTAuthMiddleware(
		config.JWTConfig{Secret: middlewareTestSecret, Enable: true},
		func(_ context.Context, userID string) (int, error) {
			if userID != "verified-user" {
				t.Fatalf("status user id = %q, want verified-user", userID)
			}
			return 1, nil
		},
		func(_ context.Context, userID, gotToken string) (bool, error) {
			checkerCalls++
			if userID != "verified-user" {
				t.Fatalf("SSO user id = %q, want verified-user", userID)
			}
			if gotToken != token {
				t.Fatalf("SSO token = %q, want signed token", gotToken)
			}
			return true, nil
		},
	))
	engine.GET("/protected", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user_id": ctx.GetHeader("user_id")})
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("user_id", "attacker-controlled")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want 200", response.Code)
	}
	if got := response.Body.String(); got != `{"user_id":"verified-user"}` {
		t.Fatalf("response body = %s", got)
	}
	if checkerCalls != 1 {
		t.Fatalf("SSO checker calls = %d, want 1", checkerCalls)
	}
}

func TestJWTAuthMiddlewareRejectsUnauthenticatedRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	testCases := []struct {
		name          string
		authorization string
		token         string
		userStatus    int
		statusErr     error
		checker       latestTokenChecker
		wantCode      int
		wantMessage   string
	}{
		{
			name:        "missing authorization",
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name:          "malformed authorization",
			authorization: "Token abc",
			wantCode:      401,
			wantMessage:   unauthorizedMessage,
		},
		{
			name: "expired token",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "expired-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute)),
				},
			}, middlewareTestSecret),
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name: "missing expiration",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "no-exp-user",
			}, middlewareTestSecret),
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name: "missing user id",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name: "wrong signing method",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS512, testClaims{
				UserID: "wrong-algorithm-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name: "invalid signature",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "wrong-signature-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, "different-secret"),
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name: "SSO token mismatch",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "superseded-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			userStatus: 1,
			checker: func(context.Context, string, string) (bool, error) {
				return false, nil
			},
			wantCode:    401,
			wantMessage: invalidatedTokenMessage,
		},
		{
			name: "SSO store unavailable",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "valid-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			userStatus: 1,
			checker: func(context.Context, string, string) (bool, error) {
				return false, errors.New("redis unavailable")
			},
			wantCode:    500,
			wantMessage: internalErrorMessage,
		},
		{
			name: "inactive user",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "disabled-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			userStatus:  -1,
			wantCode:    401,
			wantMessage: accountUnavailableMessage,
		},
		{
			name: "user store unavailable",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "valid-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			statusErr:   errors.New("mysql unavailable"),
			wantCode:    500,
			wantMessage: internalErrorMessage,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			nextCalled := false
			checker := testCase.checker
			if checker == nil {
				checker = func(context.Context, string, string) (bool, error) {
					return true, nil
				}
			}
			userStatus := testCase.userStatus
			if userStatus == 0 {
				userStatus = 1
			}

			engine := gin.New()
			engine.Use(newJWTAuthMiddleware(
				config.JWTConfig{Secret: middlewareTestSecret, Enable: true},
				func(context.Context, string) (int, error) {
					return userStatus, testCase.statusErr
				},
				checker,
			))
			engine.GET("/protected", func(ctx *gin.Context) {
				nextCalled = true
				ctx.Status(http.StatusNoContent)
			})

			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			authorization := testCase.authorization
			if testCase.token != "" {
				authorization = "Bearer " + testCase.token
			}
			if authorization != "" {
				request.Header.Set("Authorization", authorization)
			}
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			if nextCalled {
				t.Fatal("protected handler was called")
			}
			if response.Code != http.StatusOK {
				t.Fatalf("HTTP status = %d, want 200", response.Code)
			}
			wantBody := `{"code":` + strconv.Itoa(testCase.wantCode) + `,"msg":"` + testCase.wantMessage + `","data":null}`
			if got := response.Body.String(); got != wantBody {
				t.Fatalf("response body = %s, want %s", got, wantBody)
			}
		})
	}
}

func TestSoftJWTAuthMiddlewareAddsVerifiedUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token := signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
		UserID: "verified-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, middlewareTestSecret)

	engine := gin.New()
	engine.Use(newSoftJWTAuthMiddleware(config.JWTConfig{Secret: middlewareTestSecret}, nil))
	engine.GET("/optional", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user_id": ctx.GetHeader("user_id")})
	})

	request := httptest.NewRequest(http.MethodGet, "/optional", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("user_id", "attacker-controlled")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want 200", response.Code)
	}
	if got := response.Body.String(); got != `{"user_id":"verified-user"}` {
		t.Fatalf("response body = %s", got)
	}
}

func TestSoftJWTAuthMiddlewareAlwaysAllowsAnonymousRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	testCases := []struct {
		name          string
		authorization string
	}{
		{name: "missing authorization"},
		{name: "malformed authorization", authorization: "Token abc"},
		{name: "invalid token", authorization: "Bearer not-a-jwt"},
		{
			name: "expired token",
			authorization: "Bearer " + signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "expired-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute)),
				},
			}, middlewareTestSecret),
		},
		{
			name: "wrong signing method",
			authorization: "Bearer " + signMiddlewareTestToken(t, jwt.SigningMethodHS512, testClaims{
				UserID: "wrong-algorithm-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			nextCalled := false
			engine := gin.New()
			engine.Use(newSoftJWTAuthMiddleware(config.JWTConfig{Secret: middlewareTestSecret}, nil))
			engine.GET("/optional", func(ctx *gin.Context) {
				nextCalled = true
				ctx.JSON(http.StatusOK, gin.H{"user_id": ctx.GetHeader("user_id")})
			})

			request := httptest.NewRequest(http.MethodGet, "/optional", nil)
			request.Header.Set("user_id", "attacker-controlled")
			if testCase.authorization != "" {
				request.Header.Set("Authorization", testCase.authorization)
			}
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			if !nextCalled {
				t.Fatal("optional handler was not called")
			}
			if response.Code != http.StatusOK {
				t.Fatalf("HTTP status = %d, want 200", response.Code)
			}
			if got := response.Body.String(); got != `{"user_id":""}` {
				t.Fatalf("response body = %s, want empty user_id", got)
			}
		})
	}
}

func TestSoftJWTAuthMiddlewareRejectsInvalidatedSSOToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token := signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
		UserID: "verified-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, middlewareTestSecret)

	for _, testCase := range []struct {
		name    string
		checker latestTokenChecker
	}{
		{
			name: "token mismatch",
			checker: func(context.Context, string, string) (bool, error) {
				return false, nil
			},
		},
		{
			name: "Redis has no latest token",
			checker: func(context.Context, string, string) (bool, error) {
				return false, nil
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			nextCalled := false
			engine := gin.New()
			engine.Use(newSoftJWTAuthMiddleware(
				config.JWTConfig{Secret: middlewareTestSecret, Enable: true},
				testCase.checker,
			))
			engine.GET("/optional", func(ctx *gin.Context) {
				nextCalled = true
			})

			request := httptest.NewRequest(http.MethodGet, "/optional", nil)
			request.Header.Set("Authorization", "Bearer "+token)
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			if nextCalled {
				t.Fatal("optional handler was called")
			}
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("HTTP status = %d, want 401", response.Code)
			}
			if got := response.Body.String(); got != `{"code":401,"msg":"Token已失效，请重新登录","data":null}` {
				t.Fatalf("response body = %s", got)
			}
		})
	}
}

func TestSoftJWTAuthMiddlewareAllowsCurrentSSOToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token := signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
		UserID: "verified-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, middlewareTestSecret)
	checkerCalls := 0
	engine := gin.New()
	engine.Use(newSoftJWTAuthMiddleware(
		config.JWTConfig{Secret: middlewareTestSecret, Enable: true},
		func(_ context.Context, userID, gotToken string) (bool, error) {
			checkerCalls++
			return userID == "verified-user" && gotToken == token, nil
		},
	))
	engine.GET("/optional", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user_id": ctx.GetHeader("user_id")})
	})

	request := httptest.NewRequest(http.MethodGet, "/optional", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.String() != `{"user_id":"verified-user"}` {
		t.Fatalf("response = HTTP %d %s", response.Code, response.Body.String())
	}
	if checkerCalls != 1 {
		t.Fatalf("SSO checker calls = %d, want 1", checkerCalls)
	}
}

func TestSoftJWTAuthMiddlewareReturnsInternalErrorWhenSSOCheckFails(t *testing.T) {
	token := signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
		UserID: "verified-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, middlewareTestSecret)
	engine := gin.New()
	engine.Use(newSoftJWTAuthMiddleware(
		config.JWTConfig{Secret: middlewareTestSecret, Enable: true},
		func(context.Context, string, string) (bool, error) {
			return false, errors.New("redis unavailable")
		},
	))
	engine.GET("/optional", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	request := httptest.NewRequest(http.MethodGet, "/optional", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("HTTP status = %d, want 500", response.Code)
	}
	if response.Body.String() != `{"code":500,"msg":"服务器内部错误","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
}

func TestSoftJWTAuthMiddlewareKeepsMissingAuthorizationAnonymousWithSSO(t *testing.T) {
	checkerCalls := 0
	engine := gin.New()
	engine.Use(newSoftJWTAuthMiddleware(
		config.JWTConfig{Secret: middlewareTestSecret, Enable: true},
		func(context.Context, string, string) (bool, error) {
			checkerCalls++
			return false, nil
		},
	))
	engine.GET("/optional", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user_id": ctx.GetHeader("user_id")})
	})

	request := httptest.NewRequest(http.MethodGet, "/optional", nil)
	request.Header.Set("user_id", "attacker-controlled")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.String() != `{"user_id":""}` {
		t.Fatalf("anonymous response = HTTP %d %s", response.Code, response.Body.String())
	}
	if checkerCalls != 0 {
		t.Fatalf("SSO checker calls = %d, want 0", checkerCalls)
	}
}

func TestBearerTokenRequiresExactBearerPrefix(t *testing.T) {
	validToken := "header.payload.signature"
	for _, authorization := range []string{
		"bearer " + validToken,
		"BEARER " + validToken,
		" Bearer " + validToken,
		"Bearer  " + validToken,
		"Token " + validToken,
	} {
		if token, ok := bearerToken(authorization); ok || token != "" {
			t.Fatalf("bearerToken(%q) = %q, %t; want rejected", authorization, token, ok)
		}
	}

	if token, ok := bearerToken("Bearer " + validToken); !ok || token != validToken {
		t.Fatalf("valid bearer token = %q, %t", token, ok)
	}
}

func signMiddlewareTestToken(
	t *testing.T,
	method jwt.SigningMethod,
	claims jwt.Claims,
	secret string,
) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return token
}
