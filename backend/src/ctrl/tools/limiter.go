package tools

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
)

const (
	defaultAuthRateLimitWindow = 60 * time.Second
	defaultAuthRateLimitMax    = 10
	authRateLimitPrefix        = "limit:auth:ip:"
	rateLimitedMessage         = "请求过于频繁，请稍后再试"
)

type rateLimitIncrementer func(context.Context, string, time.Duration) (int64, error)

var incrementFixedWindowScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return count
`)

func AuthRateLimitMiddleware() gin.HandlerFunc {
	var rateLimitConfig config.RateLimitConfig
	if config.Conf != nil {
		rateLimitConfig = config.Conf.AuthRateLimit
	}
	return newAuthRateLimitMiddleware(rateLimitConfig, incrementRateLimit)
}

func newAuthRateLimitMiddleware(
	rateLimitConfig config.RateLimitConfig,
	increment rateLimitIncrementer,
) gin.HandlerFunc {
	window, maxRequests := normalizedAuthRateLimit(rateLimitConfig)

	return func(ctx *gin.Context) {
		if !rateLimitConfig.Enable {
			ctx.Next()
			return
		}
		if increment == nil {
			log.Printf("auth rate limit: incrementer is not initialized")
			abortAuthRequestWithStatus(ctx, http.StatusInternalServerError, 500, internalErrorMessage)
			return
		}

		key := authRateLimitKey(ctx.ClientIP())
		count, err := increment(ctx.Request.Context(), key, window)
		if err != nil {
			log.Printf("auth rate limit: %v", err)
			abortAuthRequestWithStatus(ctx, http.StatusInternalServerError, 500, internalErrorMessage)
			return
		}
		if count > int64(maxRequests) {
			abortAuthRequestWithStatus(ctx, http.StatusTooManyRequests, 429, rateLimitedMessage)
			return
		}

		ctx.Next()
	}
}

func normalizedAuthRateLimit(rateLimitConfig config.RateLimitConfig) (time.Duration, int) {
	window := time.Duration(rateLimitConfig.WindowSeconds) * time.Second
	if window <= 0 {
		window = defaultAuthRateLimitWindow
	}
	maxRequests := rateLimitConfig.MaxRequests
	if maxRequests <= 0 {
		maxRequests = defaultAuthRateLimitMax
	}
	return window, maxRequests
}

func incrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	client := data.GetRedis()
	if client == nil {
		return 0, data.ErrRedisNotInitialized
	}

	count, err := incrementFixedWindowScript.Run(
		ctx,
		client,
		[]string{key},
		int64(window/time.Second),
	).Int64()
	if err != nil {
		return 0, fmt.Errorf("increment fixed-window counter: %w", err)
	}
	return count, nil
}

func authRateLimitKey(clientIP string) string {
	return authRateLimitPrefix + clientIP
}
