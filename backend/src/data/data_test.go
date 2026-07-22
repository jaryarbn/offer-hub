package data

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewDataRejectsNilConfig(t *testing.T) {
	initializedData, err := NewData(nil)
	if err == nil {
		t.Fatal("NewData(nil) error = nil, want error")
	}
	if initializedData != nil {
		t.Fatalf("NewData(nil) data = %#v, want nil", initializedData)
	}
}

func TestNewRedisRejectsNilConfig(t *testing.T) {
	client, err := NewRedis(nil)
	if err == nil {
		t.Fatal("NewRedis(nil) error = nil, want error")
	}
	if client != nil {
		t.Fatalf("NewRedis(nil) client = %#v, want nil", client)
	}
}

func TestLatestTokenKey(t *testing.T) {
	if got := latestTokenKey("user-123"); got != "jwt:latest_token:user-123" {
		t.Fatalf("latestTokenKey() = %q, want jwt:latest_token:user-123", got)
	}
}

func TestSaveLatestTokenValidatesBeforeRedisAccess(t *testing.T) {
	data := &Data{}
	if err := data.SaveLatestToken(context.Background(), "", "token", time.Minute); err == nil {
		t.Fatal("empty user ID error = nil")
	}
	if err := data.SaveLatestToken(context.Background(), "user-123", "", time.Minute); err == nil {
		t.Fatal("empty token error = nil")
	}
	if err := data.SaveLatestToken(context.Background(), "user-123", "token", 0); !errors.Is(err, ErrInvalidTokenTTL) {
		t.Fatalf("zero TTL error = %v, want ErrInvalidTokenTTL", err)
	}
	if err := data.SaveLatestToken(context.Background(), "user-123", "token", time.Minute); !errors.Is(err, ErrRedisNotInitialized) {
		t.Fatalf("uninitialized Redis error = %v, want ErrRedisNotInitialized", err)
	}
}

func TestCheckLatestTokenValidatesBeforeRedisAccess(t *testing.T) {
	if _, err := CheckLatestToken(context.Background(), "", "token"); err == nil {
		t.Fatal("empty user ID error = nil")
	}
	if _, err := CheckLatestToken(context.Background(), "user-123", ""); err == nil {
		t.Fatal("empty token error = nil")
	}
	if _, err := CheckLatestToken(context.Background(), "user-123", "token"); !errors.Is(err, ErrRedisNotInitialized) {
		t.Fatalf("uninitialized Redis error = %v, want ErrRedisNotInitialized", err)
	}
}
