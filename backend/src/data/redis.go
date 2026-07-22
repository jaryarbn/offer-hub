package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"offer-hub/backend/src/config"
)

var redisClient *redis.Client

var (
	ErrRedisNotInitialized = errors.New("Redis is not initialized")
	ErrInvalidTokenTTL     = errors.New("token TTL must be positive")
)

const latestTokenPrefix = "jwt:latest_token:"

func NewRedis(conf *config.TomlConfig) (*redis.Client, error) {
	if conf == nil {
		return nil, errors.New("config is not initialized")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.URL,
		Password: conf.Redis.Pwd,
		DB:       conf.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect to Redis: %w", err)
	}

	redisClient = client
	return client, nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() error {
	if redisClient == nil {
		return nil
	}

	err := redisClient.Close()
	redisClient = nil
	return err
}

func (data *Data) SaveLatestToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("user ID is empty")
	}
	if strings.TrimSpace(token) == "" {
		return errors.New("token is empty")
	}
	if ttl <= 0 {
		return ErrInvalidTokenTTL
	}
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	if err := redisClient.Set(ctx, latestTokenKey(userID), token, ttl).Err(); err != nil {
		return fmt.Errorf("write latest token: %w", err)
	}
	return nil
}

func (data *Data) DeleteLatestToken(ctx context.Context, userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("user ID is empty")
	}
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	if err := redisClient.Del(ctx, latestTokenKey(userID)).Err(); err != nil {
		return fmt.Errorf("delete latest token: %w", err)
	}
	return nil
}

func CheckLatestToken(ctx context.Context, userID, token string) (bool, error) {
	if strings.TrimSpace(userID) == "" {
		return false, errors.New("user ID is empty")
	}
	if strings.TrimSpace(token) == "" {
		return false, errors.New("token is empty")
	}
	if redisClient == nil {
		return false, ErrRedisNotInitialized
	}

	latestToken, err := redisClient.Get(ctx, latestTokenKey(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read latest token: %w", err)
	}
	return latestToken == token, nil
}

func latestTokenKey(userID string) string {
	return latestTokenPrefix + userID
}
