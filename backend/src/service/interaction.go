package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

const (
	questionInteractionTarget = 1
	commentInteractionTarget  = 3
	inactiveInteractionStatus = 0
	activeInteractionStatus   = 1
	minimumQuestionTag        = 0
	maximumQuestionTag        = 3
)

var (
	ErrInvalidInteractionUserID     = errors.New("interaction user_id is empty")
	ErrInvalidInteractionTargetType = errors.New("invalid interaction target_type")
	ErrInvalidInteractionTargetID   = errors.New("interaction target_id is empty")
	ErrInteractionTargetNotFound    = data.ErrInteractionTargetNotFound
	ErrInvalidQuestionID            = errors.New("question_id is empty")
	ErrInvalidQuestionTag           = errors.New("invalid question tag")
)

type InteractionData interface {
	GetInteractionTargetLikeCount(context.Context, int, string) (int64, error)
	SetInteractionLikeStatus(context.Context, string, int, string, int, time.Time) (bool, error)
	AdjustInteractionTargetLikeCount(context.Context, int, string, int64) (int64, error)
	UpsertUserQuestionTag(context.Context, string, string, int, time.Time) error
}

type InteractionService struct {
	data InteractionData
	now  func() time.Time
}

func NewInteractionService(interactionData InteractionData) *InteractionService {
	return &InteractionService{data: interactionData, now: time.Now}
}

func (service *InteractionService) Like(
	ctx context.Context,
	req model.InteractionLikeReq,
	userID string,
) (model.InteractionLikeData, error) {
	userID, err := normalizeInteractionRequest(&req, userID)
	if err != nil {
		return model.InteractionLikeData{}, err
	}

	currentCount, err := service.data.GetInteractionTargetLikeCount(ctx, req.TargetType, req.TargetID)
	if err != nil {
		return model.InteractionLikeData{}, fmt.Errorf("get like target count: %w", err)
	}
	changed, err := service.data.SetInteractionLikeStatus(
		ctx,
		userID,
		req.TargetType,
		req.TargetID,
		activeInteractionStatus,
		service.now().UTC(),
	)
	if err != nil {
		return model.InteractionLikeData{}, fmt.Errorf("activate interaction: %w", err)
	}
	if !changed {
		return model.InteractionLikeData{Liked: true, Count: currentCount}, nil
	}

	count, err := service.data.AdjustInteractionTargetLikeCount(ctx, req.TargetType, req.TargetID, 1)
	if err != nil {
		rollbackErr := service.rollbackInteractionStatus(
			ctx,
			userID,
			req.TargetType,
			req.TargetID,
			inactiveInteractionStatus,
		)
		return model.InteractionLikeData{}, errors.Join(
			fmt.Errorf("increment like count: %w", err),
			rollbackErr,
		)
	}
	return model.InteractionLikeData{Liked: true, Count: count}, nil
}

func (service *InteractionService) Unlike(
	ctx context.Context,
	req model.InteractionUnlikeReq,
	userID string,
) (model.InteractionUnlikeData, error) {
	userID, err := normalizeInteractionRequest(&req, userID)
	if err != nil {
		return model.InteractionUnlikeData{}, err
	}

	currentCount, err := service.data.GetInteractionTargetLikeCount(ctx, req.TargetType, req.TargetID)
	if err != nil {
		return model.InteractionUnlikeData{}, fmt.Errorf("get unlike target count: %w", err)
	}
	changed, err := service.data.SetInteractionLikeStatus(
		ctx,
		userID,
		req.TargetType,
		req.TargetID,
		inactiveInteractionStatus,
		service.now().UTC(),
	)
	if err != nil {
		return model.InteractionUnlikeData{}, fmt.Errorf("deactivate interaction: %w", err)
	}
	if !changed {
		return model.InteractionUnlikeData{Count: currentCount}, nil
	}

	count, err := service.data.AdjustInteractionTargetLikeCount(ctx, req.TargetType, req.TargetID, -1)
	if err != nil {
		rollbackErr := service.rollbackInteractionStatus(
			ctx,
			userID,
			req.TargetType,
			req.TargetID,
			activeInteractionStatus,
		)
		return model.InteractionUnlikeData{}, errors.Join(
			fmt.Errorf("decrement like count: %w", err),
			rollbackErr,
		)
	}
	return model.InteractionUnlikeData{Count: count}, nil
}

func (service *InteractionService) TagQuestion(
	ctx context.Context,
	req model.TagQuestionReq,
	userID string,
) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidInteractionUserID
	}

	req.QuestionID = strings.TrimSpace(req.QuestionID)
	if req.QuestionID == "" {
		return ErrInvalidQuestionID
	}
	if req.Tag == nil || *req.Tag < minimumQuestionTag || *req.Tag > maximumQuestionTag {
		return ErrInvalidQuestionTag
	}

	if err := service.data.UpsertUserQuestionTag(
		ctx,
		userID,
		req.QuestionID,
		*req.Tag,
		service.now().UTC(),
	); err != nil {
		return fmt.Errorf("upsert question tag: %w", err)
	}
	return nil
}

func normalizeInteractionRequest(req *model.InteractionLikeReq, userID string) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", ErrInvalidInteractionUserID
	}
	if req.TargetType != questionInteractionTarget && req.TargetType != commentInteractionTarget {
		return "", ErrInvalidInteractionTargetType
	}

	req.TargetID = strings.TrimSpace(req.TargetID)
	if req.TargetID == "" {
		return "", ErrInvalidInteractionTargetID
	}
	return userID, nil
}

func (service *InteractionService) rollbackInteractionStatus(
	ctx context.Context,
	userID string,
	targetType int,
	targetID string,
	status int,
) error {
	_, err := service.data.SetInteractionLikeStatus(
		ctx,
		userID,
		targetType,
		targetID,
		status,
		service.now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("rollback interaction status: %w", err)
	}
	return nil
}
