package data

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestInteractionTargetSpecFor(t *testing.T) {
	tests := []struct {
		name       string
		targetType int
		want       interactionTargetSpec
		wantOK     bool
	}{
		{
			name:       "question",
			targetType: interactionTargetQuestion,
			want: interactionTargetSpec{
				collection: questionCollection,
				idField:    "question_id",
				countField: "thumbs_up_count",
				status:     questionNormalStatus,
			},
			wantOK: true,
		},
		{
			name:       "comment",
			targetType: interactionTargetComment,
			want: interactionTargetSpec{
				collection: commentsCollection,
				idField:    "comment_id",
				countField: "thumbs_up",
				status:     commentNormalStatus,
			},
			wantOK: true,
		},
		{name: "unsupported", targetType: 2, wantOK: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := interactionTargetSpecFor(test.targetType)
			if ok != test.wantOK || !reflect.DeepEqual(got, test.want) {
				t.Fatalf("interactionTargetSpecFor(%d) = %#v, %t; want %#v, %t", test.targetType, got, ok, test.want, test.wantOK)
			}
		})
	}
}

func TestBuildInteractionIdentityFilter(t *testing.T) {
	got := buildInteractionIdentityFilter("user-1", interactionTargetComment, "comment-1")
	want := bson.D{
		{Key: "user_id", Value: "user-1"},
		{Key: "target_type", Value: interactionTargetComment},
		{Key: "target_id", Value: "comment-1"},
		{Key: "interaction_type", Value: interactionTypeLike},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildInteractionIdentityFilter() = %#v, want %#v", got, want)
	}
}

func TestBuildInteractionTargetFilterUsesDictionaryFields(t *testing.T) {
	questionSpec, _ := interactionTargetSpecFor(interactionTargetQuestion)
	commentSpec, _ := interactionTargetSpecFor(interactionTargetComment)

	if got, want := buildInteractionTargetFilter(questionSpec, "q1"), (bson.D{
		{Key: "question_id", Value: "q1"},
		{Key: "status", Value: questionNormalStatus},
	}); !reflect.DeepEqual(got, want) {
		t.Fatalf("question target filter = %#v, want %#v", got, want)
	}
	if got, want := buildInteractionTargetFilter(commentSpec, "c1"), (bson.D{
		{Key: "comment_id", Value: "c1"},
		{Key: "status", Value: commentNormalStatus},
	}); !reflect.DeepEqual(got, want) {
		t.Fatalf("comment target filter = %#v, want %#v", got, want)
	}
}

func TestBuildActivateInteractionUpdate(t *testing.T) {
	now := time.Date(2026, 7, 22, 6, 0, 0, 0, time.UTC)
	got := buildActivateInteractionUpdate("user-1", 1, "question-1", now)
	want := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: "user-1"},
			{Key: "target_type", Value: 1},
			{Key: "target_id", Value: "question-1"},
			{Key: "interaction_type", Value: interactionTypeLike},
			{Key: "status", Value: interactionStatusActive},
			{Key: "update_time", Value: now},
		}},
		{Key: "$setOnInsert", Value: bson.D{{Key: "create_time", Value: now}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildActivateInteractionUpdate() = %#v, want %#v", got, want)
	}
}

func TestBuildUserQuestionTagUpdate(t *testing.T) {
	now := time.Date(2026, 7, 22, 6, 30, 0, 0, time.UTC)
	got := buildUserQuestionTagUpdate("user-1", "question-1", 2, now)
	want := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: "user-1"},
			{Key: "question_id", Value: "question-1"},
			{Key: "tag", Value: 2},
			{Key: "update_time", Value: now},
		}},
		{Key: "$setOnInsert", Value: bson.D{{Key: "create_time", Value: now}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildUserQuestionTagUpdate() = %#v, want %#v", got, want)
	}
}
