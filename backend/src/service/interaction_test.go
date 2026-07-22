package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"offer-hub/backend/src/model"
)

type interactionDataStub struct {
	targetCount       int64
	targetCountErr    error
	statusChanged     bool
	statusErr         error
	adjustedCount     int64
	adjustErr         error
	tagErr            error
	statusCalls       []int
	adjustCalls       []int64
	tagCalls          int
	tagUserID         string
	tagQuestionID     string
	tagValue          int
	tagUpdatedAt      time.Time
	interactionUserID string
}

func (stub *interactionDataStub) GetInteractionTargetLikeCount(
	context.Context,
	int,
	string,
) (int64, error) {
	return stub.targetCount, stub.targetCountErr
}

func (stub *interactionDataStub) SetInteractionLikeStatus(
	_ context.Context,
	userID string,
	_ int,
	_ string,
	status int,
	_ time.Time,
) (bool, error) {
	stub.interactionUserID = userID
	stub.statusCalls = append(stub.statusCalls, status)
	return stub.statusChanged, stub.statusErr
}

func (stub *interactionDataStub) AdjustInteractionTargetLikeCount(
	_ context.Context,
	_ int,
	_ string,
	delta int64,
) (int64, error) {
	stub.adjustCalls = append(stub.adjustCalls, delta)
	return stub.adjustedCount, stub.adjustErr
}

func (stub *interactionDataStub) UpsertUserQuestionTag(
	_ context.Context,
	userID string,
	questionID string,
	tag int,
	updatedAt time.Time,
) error {
	stub.tagCalls++
	stub.tagUserID = userID
	stub.tagQuestionID = questionID
	stub.tagValue = tag
	stub.tagUpdatedAt = updatedAt
	return stub.tagErr
}

func TestInteractionServiceLikeIncrementsOnlyOnStateTransition(t *testing.T) {
	stub := &interactionDataStub{targetCount: 10, statusChanged: true, adjustedCount: 11}
	service := NewInteractionService(stub)
	service.now = func() time.Time { return time.Date(2026, 7, 22, 14, 0, 0, 0, time.FixedZone("CST", 8*60*60)) }

	got, err := service.Like(context.Background(), model.InteractionLikeReq{
		TargetType: 1,
		TargetID:   " question-1 ",
	}, " user-1 ")
	if err != nil {
		t.Fatalf("Like() error = %v", err)
	}
	if !got.Liked || got.Count != 11 {
		t.Fatalf("Like() = %#v, want liked=true count=11", got)
	}
	if len(stub.statusCalls) != 1 || stub.statusCalls[0] != activeInteractionStatus {
		t.Fatalf("status calls = %#v, want [1]", stub.statusCalls)
	}
	if len(stub.adjustCalls) != 1 || stub.adjustCalls[0] != 1 {
		t.Fatalf("adjust calls = %#v, want [1]", stub.adjustCalls)
	}
	if stub.interactionUserID != "user-1" {
		t.Fatalf("interaction user_id = %q, want user-1", stub.interactionUserID)
	}
}

func TestInteractionServiceLikeIsIdempotent(t *testing.T) {
	stub := &interactionDataStub{targetCount: 10, statusChanged: false}
	service := NewInteractionService(stub)

	got, err := service.Like(context.Background(), model.InteractionLikeReq{
		TargetType: 3,
		TargetID:   "comment-1",
	}, "user-1")
	if err != nil {
		t.Fatalf("Like() error = %v", err)
	}
	if !got.Liked || got.Count != 10 {
		t.Fatalf("Like() = %#v, want liked=true count=10", got)
	}
	if len(stub.adjustCalls) != 0 {
		t.Fatalf("adjust calls = %#v, want none", stub.adjustCalls)
	}
}

func TestInteractionServiceUnlikeDecrementsOnlyOnStateTransition(t *testing.T) {
	stub := &interactionDataStub{targetCount: 10, statusChanged: true, adjustedCount: 9}
	service := NewInteractionService(stub)

	got, err := service.Unlike(context.Background(), model.InteractionUnlikeReq{
		TargetType: 3,
		TargetID:   "comment-1",
	}, "user-1")
	if err != nil {
		t.Fatalf("Unlike() error = %v", err)
	}
	if got.Count != 9 {
		t.Fatalf("Unlike() = %#v, want count=9", got)
	}
	if len(stub.statusCalls) != 1 || stub.statusCalls[0] != inactiveInteractionStatus {
		t.Fatalf("status calls = %#v, want [0]", stub.statusCalls)
	}
	if len(stub.adjustCalls) != 1 || stub.adjustCalls[0] != -1 {
		t.Fatalf("adjust calls = %#v, want [-1]", stub.adjustCalls)
	}
}

func TestInteractionServiceUnlikeIsIdempotent(t *testing.T) {
	stub := &interactionDataStub{targetCount: 0, statusChanged: false}
	service := NewInteractionService(stub)

	got, err := service.Unlike(context.Background(), model.InteractionUnlikeReq{
		TargetType: 1,
		TargetID:   "question-1",
	}, "user-1")
	if err != nil {
		t.Fatalf("Unlike() error = %v", err)
	}
	if got.Count != 0 {
		t.Fatalf("Unlike() = %#v, want count=0", got)
	}
	if len(stub.adjustCalls) != 0 {
		t.Fatalf("adjust calls = %#v, want none", stub.adjustCalls)
	}
}

func TestInteractionServiceRejectsInvalidLikeRequests(t *testing.T) {
	tests := []struct {
		name    string
		req     model.InteractionLikeReq
		userID  string
		wantErr error
	}{
		{name: "missing user", req: model.InteractionLikeReq{TargetType: 1, TargetID: "q1"}, wantErr: ErrInvalidInteractionUserID},
		{name: "unsupported target type", req: model.InteractionLikeReq{TargetType: 2, TargetID: "q1"}, userID: "u1", wantErr: ErrInvalidInteractionTargetType},
		{name: "blank target ID", req: model.InteractionLikeReq{TargetType: 1, TargetID: "  "}, userID: "u1", wantErr: ErrInvalidInteractionTargetID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := &interactionDataStub{}
			service := NewInteractionService(stub)
			_, err := service.Like(context.Background(), test.req, test.userID)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Like() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestInteractionServiceTagQuestionAcceptsZeroAndUpserts(t *testing.T) {
	stub := &interactionDataStub{}
	service := NewInteractionService(stub)
	fixedNow := time.Date(2026, 7, 22, 14, 30, 0, 0, time.FixedZone("CST", 8*60*60))
	service.now = func() time.Time { return fixedNow }
	tag := 0

	err := service.TagQuestion(context.Background(), model.TagQuestionReq{
		QuestionID: " question-1 ",
		Tag:        &tag,
	}, " user-1 ")
	if err != nil {
		t.Fatalf("TagQuestion() error = %v", err)
	}
	if stub.tagCalls != 1 || stub.tagUserID != "user-1" || stub.tagQuestionID != "question-1" || stub.tagValue != 0 {
		t.Fatalf("tag call = count:%d user:%q question:%q tag:%d", stub.tagCalls, stub.tagUserID, stub.tagQuestionID, stub.tagValue)
	}
	if !stub.tagUpdatedAt.Equal(fixedNow.UTC()) {
		t.Fatalf("tag update_time = %v, want %v", stub.tagUpdatedAt, fixedNow.UTC())
	}
}

func TestInteractionServiceRejectsInvalidQuestionTags(t *testing.T) {
	invalidLow := -1
	invalidHigh := 4
	tests := []struct {
		name    string
		req     model.TagQuestionReq
		userID  string
		wantErr error
	}{
		{name: "missing user", req: model.TagQuestionReq{QuestionID: "q1", Tag: &invalidLow}, wantErr: ErrInvalidInteractionUserID},
		{name: "blank question ID", req: model.TagQuestionReq{QuestionID: " ", Tag: &invalidLow}, userID: "u1", wantErr: ErrInvalidQuestionID},
		{name: "missing tag", req: model.TagQuestionReq{QuestionID: "q1"}, userID: "u1", wantErr: ErrInvalidQuestionTag},
		{name: "tag below range", req: model.TagQuestionReq{QuestionID: "q1", Tag: &invalidLow}, userID: "u1", wantErr: ErrInvalidQuestionTag},
		{name: "tag above range", req: model.TagQuestionReq{QuestionID: "q1", Tag: &invalidHigh}, userID: "u1", wantErr: ErrInvalidQuestionTag},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := &interactionDataStub{}
			service := NewInteractionService(stub)
			err := service.TagQuestion(context.Background(), test.req, test.userID)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("TagQuestion() error = %v, want %v", err, test.wantErr)
			}
			if stub.tagCalls != 0 {
				t.Fatalf("tag calls = %d, want 0", stub.tagCalls)
			}
		})
	}
}
