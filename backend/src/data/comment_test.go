package data

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCommentUserRecordTableName(t *testing.T) {
	if got := (CommentUserRecord{}).TableName(); got != userInfoTable {
		t.Fatalf("CommentUserRecord.TableName() = %q, want %q", got, userInfoTable)
	}
}

func TestBuildCommentFilter(t *testing.T) {
	tests := []struct {
		name     string
		parentID string
	}{
		{name: "top-level comments", parentID: ""},
		{name: "child comments", parentID: "comment-parent"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildCommentFilter(CommentFilter{
				TargetType: 1,
				TargetID:   "question-1",
				ParentID:   test.parentID,
			})
			want := bson.D{
				{Key: "target_type", Value: 1},
				{Key: "target_id", Value: "question-1"},
				{Key: "parent_id", Value: test.parentID},
				{Key: "status", Value: commentNormalStatus},
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("buildCommentFilter() = %#v, want %#v", got, want)
			}
		})
	}
}

func TestBuildParentCommentFilter(t *testing.T) {
	got := buildParentCommentFilter("comment-parent")
	want := bson.D{
		{Key: "comment_id", Value: "comment-parent"},
		{Key: "parent_id", Value: ""},
		{Key: "status", Value: commentNormalStatus},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildParentCommentFilter() = %#v, want %#v", got, want)
	}
}

func TestBuildCommentDeleteQueries(t *testing.T) {
	updatedAt := time.Date(2026, time.July, 22, 15, 30, 0, 0, time.UTC)

	wantFilter := bson.D{{Key: "comment_id", Value: "comment-1"}}
	if got := buildCommentIDFilter("comment-1"); !reflect.DeepEqual(got, wantFilter) {
		t.Fatalf("buildCommentIDFilter() = %#v, want %#v", got, wantFilter)
	}

	wantUpdate := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: commentDeletedStatus},
		{Key: "update_time", Value: updatedAt},
	}}}
	if got := buildSoftDeleteCommentUpdate(updatedAt); !reflect.DeepEqual(got, wantUpdate) {
		t.Fatalf("buildSoftDeleteCommentUpdate() = %#v, want %#v", got, wantUpdate)
	}

	wantContentUpdate := bson.D{{Key: "$set", Value: bson.D{
		{Key: "content", Value: "filtered content"},
		{Key: "update_time", Value: updatedAt},
	}}}
	if got := buildUpdateCommentContentUpdate("filtered content", updatedAt); !reflect.DeepEqual(got, wantContentUpdate) {
		t.Fatalf("buildUpdateCommentContentUpdate() = %#v, want %#v", got, wantContentUpdate)
	}
}

func TestBuildCommentSort(t *testing.T) {
	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
		want      bson.D
	}{
		{name: "likes ascending", sortBy: "thumbs_up", sortOrder: "asc", want: bson.D{{Key: "thumbs_up", Value: 1}}},
		{name: "time descending", sortBy: "create_time", sortOrder: "desc", want: bson.D{{Key: "create_time", Value: -1}}},
		{name: "unknown field falls back to newest", sortBy: "$where", sortOrder: "invalid", want: bson.D{{Key: "create_time", Value: -1}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := buildCommentSort(test.sortBy, test.sortOrder); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("buildCommentSort(%q, %q) = %#v, want %#v", test.sortBy, test.sortOrder, got, test.want)
			}
		})
	}
}

func TestBuildSubCommentQueries(t *testing.T) {
	parentIDs := []string{"comment-1", "comment-2"}
	wantFilter := bson.D{
		{Key: "parent_id", Value: bson.D{{Key: "$in", Value: parentIDs}}},
		{Key: "status", Value: commentNormalStatus},
	}
	if got := buildSubCommentsFilter(parentIDs); !reflect.DeepEqual(got, wantFilter) {
		t.Fatalf("buildSubCommentsFilter() = %#v, want %#v", got, wantFilter)
	}

	wantPipeline := mongo.Pipeline{
		{{Key: "$match", Value: wantFilter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$parent_id"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	if got := buildSubCommentCountPipeline(parentIDs); !reflect.DeepEqual(got, wantPipeline) {
		t.Fatalf("buildSubCommentCountPipeline() = %#v, want %#v", got, wantPipeline)
	}
}

func TestBuildLikedCommentsFilter(t *testing.T) {
	commentIDs := []string{"comment-1", "comment-2"}
	got := buildLikedCommentsFilter("user-1", commentIDs)
	want := bson.D{
		{Key: "user_id", Value: "user-1"},
		{Key: "target_type", Value: commentInteractionTargetType},
		{Key: "target_id", Value: bson.D{{Key: "$in", Value: commentIDs}}},
		{Key: "interaction_type", Value: commentLikeInteractionType},
		{Key: "status", Value: commentActiveInteractionStatus},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildLikedCommentsFilter() = %#v, want %#v", got, want)
	}
}

func TestNormalizeCommentPagination(t *testing.T) {
	if page, size := normalizeCommentPagination(0, 0); page != 1 || size != 20 {
		t.Fatalf("normalizeCommentPagination(0, 0) = (%d, %d), want (1, 20)", page, size)
	}
	if page, size := normalizeCommentPagination(3, 7); page != 3 || size != 7 {
		t.Fatalf("normalizeCommentPagination(3, 7) = (%d, %d), want (3, 7)", page, size)
	}
}
