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

func TestTokenBlacklistKey(t *testing.T) {
	if got := blacklistKey("signed-token"); got != "blacklist:signed-token" {
		t.Fatalf("blacklistKey() = %q, want blacklist:signed-token", got)
	}
}

func TestAddTokenToBlacklistValidatesBeforeRedisAccess(t *testing.T) {
	data := &Data{}
	if err := data.AddTokenToBlacklist(context.Background(), "", time.Minute); err == nil {
		t.Fatal("empty token error = nil")
	}
	if err := data.AddTokenToBlacklist(context.Background(), "token", 0); !errors.Is(err, ErrInvalidBlacklistTTL) {
		t.Fatalf("zero TTL error = %v, want ErrInvalidBlacklistTTL", err)
	}
	if err := data.AddTokenToBlacklist(context.Background(), "token", time.Minute); !errors.Is(err, ErrRedisNotInitialized) {
		t.Fatalf("uninitialized Redis error = %v, want ErrRedisNotInitialized", err)
	}
}
