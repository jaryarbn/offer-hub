package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

var (
	ErrInvalidUserID    = errors.New("user_id is empty")
	ErrUserInfoNotFound = data.ErrUserNotFound
)

type UserInfoData interface {
	GetUserByID(context.Context, string) (data.UserInfoRecord, error)
}

type UserInfoService struct {
	data UserInfoData
}

func NewUserInfoService(userInfoData UserInfoData) *UserInfoService {
	return &UserInfoService{data: userInfoData}
}

func (service *UserInfoService) GetUserInfo(ctx context.Context, userID string) (model.UserInfo, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return model.UserInfo{}, ErrInvalidUserID
	}

	record, err := service.data.GetUserByID(ctx, userID)
	if errors.Is(err, data.ErrUserNotFound) {
		return model.UserInfo{}, ErrUserInfoNotFound
	}
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("get user info: %w", err)
	}

	return model.UserInfo{
		UserID:       record.UserID,
		Username:     record.Username,
		NickName:     record.NickName,
		Avatar:       record.Avatar,
		VIP:          record.VIP,
		Sex:          record.Sex,
		Phone:        record.Phone,
		Email:        record.Email,
		Introduction: record.Introduction,
		AvatarURL:    record.Avatar,
		UserStatus:   record.UserStatus,
		UserType:     record.UserType,
		CreateTime:   formatTime(record.CreateTime),
		UpdateTime:   formatTime(record.UpdateTime),
	}, nil
}
