package tools

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
)

const (
	defaultRateLimitWindow     = 60 * time.Second
	defaultRateLimitMax        = 20
	rateLimitPrefix            = "rate_limit:"
	defaultAuthRateLimitWindow = 60 * time.Second
	defaultAuthRateLimitMax    = 10
	authRateLimitPrefix        = "limit:auth:ip:"
	rateLimitedMessage         = "请求过于频繁，请稍后再试"
)

type rateLimitRedisClient interface {
	Incr(context.Context, string) *redis.IntCmd
	Expire(context.Context, string, time.Duration) *redis.BoolCmd
}

type RateLimiter struct {
	redisClient rateLimitRedisClient
	limit       int
	window      time.Duration
}

func (limiter *RateLimiter) IsRequestAllowed(keyID string) bool {
	if limiter == nil || limiter.redisClient == nil || limiter.limit <= 0 {
		return false
	}

	key, ok := rateLimitKey(keyID, time.Now(), limiter.window)
	if !ok {
		return false
	}

	ctx := context.Background()
	count, err := limiter.redisClient.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("rate limiter: increment counter: %v", err)
		return false
	}

	if count == 1 {
		if err := limiter.redisClient.Expire(ctx, key, limiter.window+time.Second).Err(); err != nil {
			log.Printf("rate limiter: set counter TTL: %v", err)
			return false
		}
	}

	return count <= int64(limiter.limit)
}

// RateLimitMiddleware must run after JWT middleware so user_id is a verified
// server-provided identity. Anonymous requests fall back to the client IP.
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	rateLimitConfig := resolveRateLimitConfig(cfg)
	window, limit := normalizedRateLimit(rateLimitConfig)
	limiter := &RateLimiter{
		redisClient: data.GetRedis(),
		limit:       limit,
		window:      window,
	}

	return newRateLimitMiddleware(rateLimitConfig.Enable, limiter.IsRequestAllowed)
}

func newRateLimitMiddleware(enabled bool, isAllowed func(string) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !enabled {
			ctx.Next()
			return
		}

		keyID := rateLimitKeyID(ctx)
		if isAllowed == nil || !isAllowed(keyID) {
			abortAuthRequestWithStatus(
				ctx,
				http.StatusTooManyRequests,
				http.StatusTooManyRequests,
				rateLimitedMessage,
			)
			return
		}

		ctx.Next()
	}
}

func resolveRateLimitConfig(cfg *config.Config) config.RateLimitConfig {
	if cfg != nil {
		return *cfg
	}
	if config.Conf != nil {
		return config.Conf.RateLimit
	}
	return config.RateLimitConfig{}
}

func normalizedRateLimit(rateLimitConfig config.RateLimitConfig) (time.Duration, int) {
	window := time.Duration(rateLimitConfig.WindowSeconds) * time.Second
	if window <= 0 {
		window = defaultRateLimitWindow
	}

	limit := rateLimitConfig.MaxRequests
	if limit <= 0 {
		limit = defaultRateLimitMax
	}
	return window, limit
}

func rateLimitKeyID(ctx *gin.Context) string {
	if userID := strings.TrimSpace(ctx.GetHeader("user_id")); userID != "" {
		return "user:" + userID
	}
	return "ip:" + ctx.ClientIP()
}

func rateLimitKey(keyID string, now time.Time, window time.Duration) (string, bool) {
	windowSeconds := int64(window / time.Second)
	if strings.TrimSpace(keyID) == "" || windowSeconds <= 0 {
		return "", false
	}

	windowStart := now.Unix() / windowSeconds
	return fmt.Sprintf("%s%s:%d", rateLimitPrefix, keyID, windowStart), true
}

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
