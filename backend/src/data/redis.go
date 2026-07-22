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
	ErrInvalidBlacklistTTL = errors.New("blacklist TTL must be positive")
)

const tokenBlacklistPrefix = "blacklist:"

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

func (data *Data) AddTokenToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is empty")
	}
	if ttl <= 0 {
		return ErrInvalidBlacklistTTL
	}
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	if err := redisClient.Set(ctx, blacklistKey(token), "1", ttl).Err(); err != nil {
		return fmt.Errorf("write token blacklist: %w", err)
	}
	return nil
}

func IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	if strings.TrimSpace(token) == "" {
		return false, errors.New("token is empty")
	}
	if redisClient == nil {
		return false, ErrRedisNotInitialized
	}

	exists, err := redisClient.Exists(ctx, blacklistKey(token)).Result()
	if err != nil {
		return false, fmt.Errorf("read token blacklist: %w", err)
	}
	return exists > 0, nil
}

func blacklistKey(token string) string {
	return tokenBlacklistPrefix + token
}
