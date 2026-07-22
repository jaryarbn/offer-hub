package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"offer-hub/backend/src/data"
)

type userInfoDataStub struct {
	record data.UserInfoRecord
	err    error
	userID string
}

func (stub *userInfoDataStub) GetUserByID(_ context.Context, userID string) (data.UserInfoRecord, error) {
	stub.userID = userID
	return stub.record, stub.err
}

func TestUserInfoServiceReturnsSnakeCaseDataShape(t *testing.T) {
	createdAt := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	stub := &userInfoDataStub{record: data.UserInfoRecord{
		UserID:       "abc123",
		Username:     "testuser",
		NickName:     "测试用户",
		Avatar:       "https://example.com/avatar.jpg",
		Introduction: "个人简介",
		VIP:          false,
		Sex:          1,
		Phone:        "",
		Email:        "",
		UserStatus:   1,
		UserType:     1,
		CreateTime:   createdAt,
		UpdateTime:   createdAt,
	}}
	userInfoService := NewUserInfoService(stub)

	got, err := userInfoService.GetUserInfo(context.Background(), "  abc123  ")
	if err != nil {
		t.Fatalf("GetUserInfo() error = %v", err)
	}
	if stub.userID != "abc123" {
		t.Fatalf("queried user_id = %q, want abc123", stub.userID)
	}
	if got.UserID != "abc123" || got.Username != "testuser" || got.NickName != "测试用户" {
		t.Fatalf("basic user info = %#v", got)
	}
	if got.AvatarURL != got.Avatar {
		t.Fatalf("avatar_url = %q, want avatar %q", got.AvatarURL, got.Avatar)
	}
	if got.CreateTime != "2024-01-01 12:00:00" || got.UpdateTime != "2024-01-01 12:00:00" {
		t.Fatalf("formatted timestamps = (%q, %q)", got.CreateTime, got.UpdateTime)
	}
}

func TestUserInfoServiceRejectsEmptyUserID(t *testing.T) {
	stub := &userInfoDataStub{}
	userInfoService := NewUserInfoService(stub)

	_, err := userInfoService.GetUserInfo(context.Background(), "   ")
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("GetUserInfo() error = %v, want ErrInvalidUserID", err)
	}
	if stub.userID != "" {
		t.Fatalf("empty user_id reached data layer: %q", stub.userID)
	}
}

func TestUserInfoServicePreservesNotFoundAndStorageErrors(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
		want error
	}{
		{name: "not found", err: data.ErrUserNotFound, want: ErrUserInfoNotFound},
		{name: "storage error", err: errors.New("database unavailable"), want: nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := NewUserInfoService(&userInfoDataStub{err: test.err})
			_, err := service.GetUserInfo(context.Background(), "abc123")
			if test.want != nil {
				if !errors.Is(err, test.want) {
					t.Fatalf("GetUserInfo() error = %v, want %v", err, test.want)
				}
				return
			}
			if err == nil {
				t.Fatal("GetUserInfo() error = nil, want storage error")
			}
		})
	}
}
