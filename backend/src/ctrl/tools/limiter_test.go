package tools

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"offer-hub/backend/src/config"
)

func TestRateLimiterUsesFixedWindowKeyAndTTL(t *testing.T) {
	client := newFakeRateLimitRedisClient()
	limiter := &RateLimiter{
		redisClient: client,
		limit:       2,
		window:      time.Minute,
	}

	if !limiter.IsRequestAllowed("user:user-123") {
		t.Fatal("first request should be allowed")
	}
	if !limiter.IsRequestAllowed("user:user-123") {
		t.Fatal("second request should be allowed")
	}
	if limiter.IsRequestAllowed("user:user-123") {
		t.Fatal("third request should be rejected")
	}

	if len(client.counts) != 1 {
		t.Fatalf("counter keys = %#v, want one fixed-window key", client.counts)
	}
	for key, count := range client.counts {
		if !strings.HasPrefix(key, "rate_limit:user:user-123:") {
			t.Fatalf("counter key = %q", key)
		}
		if count != 3 {
			t.Fatalf("counter value = %d, want 3", count)
		}
		if client.expirations[key] != time.Minute+time.Second {
			t.Fatalf("counter TTL = %s, want 1m1s", client.expirations[key])
		}
	}
	if client.expireCalls != 1 {
		t.Fatalf("Expire calls = %d, want 1", client.expireCalls)
	}
}

func TestRateLimiterRejectsRedisFailuresAndInvalidState(t *testing.T) {
	if (&RateLimiter{}).IsRequestAllowed("ip:192.0.2.1") {
		t.Fatal("limiter without Redis should reject requests")
	}

	client := newFakeRateLimitRedisClient()
	client.incrErr = errors.New("redis unavailable")
	limiter := &RateLimiter{redisClient: client, limit: 1, window: time.Minute}
	if limiter.IsRequestAllowed("ip:192.0.2.1") {
		t.Fatal("Redis INCR failure should reject requests")
	}

	client.incrErr = nil
	client.expireErr = errors.New("expire failed")
	if limiter.IsRequestAllowed("ip:192.0.2.2") {
		t.Fatal("Redis EXPIRE failure should reject requests")
	}
}

func TestRateLimitKeyUsesUnixWindow(t *testing.T) {
	key, ok := rateLimitKey("ip:192.0.2.10", time.Unix(125, 0), time.Minute)
	if !ok {
		t.Fatal("rateLimitKey() rejected valid input")
	}
	if key != "rate_limit:ip:192.0.2.10:2" {
		t.Fatalf("rateLimitKey() = %q", key)
	}

	if _, ok := rateLimitKey("", time.Now(), time.Minute); ok {
		t.Fatal("rateLimitKey() accepted an empty identity")
	}
	if _, ok := rateLimitKey("ip:192.0.2.10", time.Now(), 0); ok {
		t.Fatal("rateLimitKey() accepted a zero window")
	}
}

func TestRateLimitMiddlewareUsesUserIDThenClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var keyIDs []string
	middleware := newRateLimitMiddleware(true, func(keyID string) bool {
		keyIDs = append(keyIDs, keyID)
		return true
	})

	engine := gin.New()
	engine.Use(middleware)
	engine.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	authenticated := httptest.NewRequest(http.MethodGet, "/", nil)
	authenticated.Header.Set("user_id", "user-123")
	authenticated.RemoteAddr = "192.0.2.10:1234"
	authenticatedResponse := httptest.NewRecorder()
	engine.ServeHTTP(authenticatedResponse, authenticated)

	anonymous := httptest.NewRequest(http.MethodGet, "/", nil)
	anonymous.RemoteAddr = "192.0.2.20:5678"
	anonymousResponse := httptest.NewRecorder()
	engine.ServeHTTP(anonymousResponse, anonymous)

	if authenticatedResponse.Code != http.StatusNoContent || anonymousResponse.Code != http.StatusNoContent {
		t.Fatalf(
			"response status = authenticated %d, anonymous %d",
			authenticatedResponse.Code,
			anonymousResponse.Code,
		)
	}
	if len(keyIDs) != 2 || keyIDs[0] != "user:user-123" || keyIDs[1] != "ip:192.0.2.20" {
		t.Fatalf("rate-limit identities = %#v", keyIDs)
	}
}

func TestRateLimitMiddlewareRejectsWithDocumentedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	nextCalls := 0
	engine := gin.New()
	engine.Use(newRateLimitMiddleware(true, func(string) bool { return false }))
	engine.GET("/", func(ctx *gin.Context) {
		nextCalls++
		ctx.Status(http.StatusNoContent)
	})

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("response status = %d, want 429", response.Code)
	}
	if response.Body.String() != `{"code":429,"msg":"请求过于频繁，请稍后再试","data":null}` {
		t.Fatalf("response body = %s", response.Body.String())
	}
	if nextCalls != 0 {
		t.Fatalf("handler calls = %d, want 0", nextCalls)
	}
}

func TestRateLimitMiddlewareCanBeDisabledAndReadsGlobalConfig(t *testing.T) {
	previous := config.Conf
	config.Conf = &config.TomlConfig{
		RateLimit: config.RateLimitConfig{Enable: true, WindowSeconds: 30, MaxRequests: 7},
	}
	t.Cleanup(func() { config.Conf = previous })

	resolved := resolveRateLimitConfig(nil)
	window, limit := normalizedRateLimit(resolved)
	if !resolved.Enable || window != 30*time.Second || limit != 7 {
		t.Fatalf("resolved rate-limit config = %+v, %s/%d", resolved, window, limit)
	}

	nextCalls := 0
	engine := gin.New()
	engine.Use(RateLimitMiddleware(&config.Config{Enable: false}))
	engine.GET("/", func(ctx *gin.Context) {
		nextCalls++
		ctx.Status(http.StatusNoContent)
	})
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusNoContent || nextCalls != 1 {
		t.Fatalf("disabled middleware = HTTP %d, handler calls %d", response.Code, nextCalls)
	}
}

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

type fakeRateLimitRedisClient struct {
	counts      map[string]int64
	expirations map[string]time.Duration
	expireCalls int
	incrErr     error
	expireErr   error
}

func newFakeRateLimitRedisClient() *fakeRateLimitRedisClient {
	return &fakeRateLimitRedisClient{
		counts:      make(map[string]int64),
		expirations: make(map[string]time.Duration),
	}
}

func (client *fakeRateLimitRedisClient) Incr(_ context.Context, key string) *redis.IntCmd {
	if client.incrErr != nil {
		return redis.NewIntResult(0, client.incrErr)
	}
	client.counts[key]++
	return redis.NewIntResult(client.counts[key], nil)
}

func (client *fakeRateLimitRedisClient) Expire(
	_ context.Context,
	key string,
	expiration time.Duration,
) *redis.BoolCmd {
	client.expireCalls++
	if client.expireErr != nil {
		return redis.NewBoolResult(false, client.expireErr)
	}
	client.expirations[key] = expiration
	return redis.NewBoolResult(true, nil)
}
