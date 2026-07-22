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
		config.JWTConfig{Secret: middlewareTestSecret},
		func(_ context.Context, gotToken string) (bool, error) {
			checkerCalls++
			if gotToken != token {
				t.Fatalf("blacklist token = %q, want signed token", gotToken)
			}
			return false, nil
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
		t.Fatalf("blacklist checker calls = %d, want 1", checkerCalls)
	}
}

func TestJWTAuthMiddlewareRejectsUnauthenticatedRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	testCases := []struct {
		name          string
		authorization string
		token         string
		checker       tokenBlacklistChecker
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
			name: "blacklisted token",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "logged-out-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			checker: func(context.Context, string) (bool, error) {
				return true, nil
			},
			wantCode:    401,
			wantMessage: unauthorizedMessage,
		},
		{
			name: "blacklist unavailable",
			token: signMiddlewareTestToken(t, jwt.SigningMethodHS256, testClaims{
				UserID: "valid-user",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}, middlewareTestSecret),
			checker: func(context.Context, string) (bool, error) {
				return false, errors.New("redis unavailable")
			},
			wantCode:    500,
			wantMessage: internalErrorMessage,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			nextCalled := false
			checker := testCase.checker
			if checker == nil {
				checker = func(context.Context, string) (bool, error) {
					return false, nil
				}
			}

			engine := gin.New()
			engine.Use(newJWTAuthMiddleware(
				config.JWTConfig{Secret: middlewareTestSecret},
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
	engine.Use(newSoftJWTAuthMiddleware(config.JWTConfig{Secret: middlewareTestSecret}))
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
			engine.Use(newSoftJWTAuthMiddleware(config.JWTConfig{Secret: middlewareTestSecret}))
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
