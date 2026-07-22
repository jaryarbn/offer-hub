package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/config"
)

func TestAuthRateLimitMiddlewareLimitsByClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	counts := make(map[string]int64)
	var gotWindow time.Duration
	middleware := newAuthRateLimitMiddleware(
		config.RateLimitConfig{Enable: true, WindowSeconds: 60, MaxRequests: 1},
		func(_ context.Context, key string, window time.Duration) (int64, error) {
			gotWindow = window
			counts[key]++
			return counts[key], nil
		},
	)

	nextCalls := 0
	engine := gin.New()
	engine.POST("/auth/login", middleware, func(ctx *gin.Context) {
		nextCalls++
		ctx.Status(http.StatusNoContent)
	})

	first := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	first.RemoteAddr = "192.0.2.10:1234"
	firstResponse := httptest.NewRecorder()
	engine.ServeHTTP(firstResponse, first)
	if firstResponse.Code != http.StatusNoContent {
		t.Fatalf("first response status = %d, want 204", firstResponse.Code)
	}

	second := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	second.RemoteAddr = "192.0.2.10:5678"
	secondResponse := httptest.NewRecorder()
	engine.ServeHTTP(secondResponse, second)
	if secondResponse.Code != http.StatusTooManyRequests {
		t.Fatalf("second response status = %d, want 429", secondResponse.Code)
	}
	if secondResponse.Body.String() != `{"code":429,"msg":"请求过于频繁，请稍后再试","data":null}` {
		t.Fatalf("second response body = %s", secondResponse.Body.String())
	}
	if nextCalls != 1 {
		t.Fatalf("handler calls = %d, want 1", nextCalls)
	}
	if counts["limit:auth:ip:192.0.2.10"] != 2 {
		t.Fatalf("rate-limit counts = %#v", counts)
	}
	if gotWindow != time.Minute {
		t.Fatalf("window = %s, want 1m", gotWindow)
	}
}

func TestAuthRateLimitMiddlewareUsesDefaultsAndCanBeDisabled(t *testing.T) {
	window, maxRequests := normalizedAuthRateLimit(config.RateLimitConfig{})
	if window != time.Minute || maxRequests != 10 {
		t.Fatalf("defaults = %s/%d, want 1m/10", window, maxRequests)
	}

	incrementCalls := 0
	middleware := newAuthRateLimitMiddleware(
		config.RateLimitConfig{},
		func(context.Context, string, time.Duration) (int64, error) {
			incrementCalls++
			return 1, nil
		},
	)
	engine := gin.New()
	engine.Use(middleware)
	engine.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusNoContent || incrementCalls != 0 {
		t.Fatalf("disabled middleware = HTTP %d, increment calls %d", response.Code, incrementCalls)
	}
}
