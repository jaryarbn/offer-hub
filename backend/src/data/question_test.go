package data

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestBuildQuestionCountPipeline(t *testing.T) {
	bankIDs := []string{"bank-1", "bank-2"}
	bankFilter := bson.D{{Key: "$in", Value: bankIDs}}
	want := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "status", Value: 1},
			{Key: "bank_list", Value: bankFilter},
		}}},
		{{Key: "$unwind", Value: "$bank_list"}},
		{{Key: "$match", Value: bson.D{
			{Key: "bank_list", Value: bankFilter},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$bank_list"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	got := buildQuestionCountPipeline(bankIDs)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildQuestionCountPipeline() = %#v, want %#v", got, want)
	}
}

func TestBuildQuestionFilter(t *testing.T) {
	got := buildQuestionFilter(QuestionFilter{
		BankID:     "bank-1",
		Keyword:    "Go (基础)",
		Difficulty: 2,
		Tags:       []string{"Go", "并发"},
		JobName:    "后端开发",
	})
	want := bson.D{
		{Key: "status", Value: 1},
		{Key: "bank_list", Value: "bank-1"},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "title", Value: bson.D{
				{Key: "$regex", Value: `Go \(基础\)`},
				{Key: "$options", Value: "i"},
			}}},
			bson.D{{Key: "content", Value: bson.D{
				{Key: "$regex", Value: `Go \(基础\)`},
				{Key: "$options", Value: "i"},
			}}},
		}},
		{Key: "difficulty", Value: 2},
		{Key: "tags", Value: bson.D{{Key: "$all", Value: []string{"Go", "并发"}}}},
		{Key: "job_name", Value: "后端开发"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildQuestionFilter() = %#v, want %#v", got, want)
	}
}

func TestBuildQuestionFilterDefaultsToNormalStatus(t *testing.T) {
	got := buildQuestionFilter(QuestionFilter{})
	want := bson.D{{Key: "status", Value: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildQuestionFilter() = %#v, want %#v", got, want)
	}
}

func TestBuildQuestionDetailFilter(t *testing.T) {
	got := buildQuestionDetailFilter("question-1")
	want := bson.D{
		{Key: "question_id", Value: "question-1"},
		{Key: "status", Value: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildQuestionDetailFilter() = %#v, want %#v", got, want)
	}
}

func TestBuildHotQuestionQuery(t *testing.T) {
	filter := buildHotQuestionFilter("后端开发")
	wantFilter := bson.D{
		{Key: "status", Value: 1},
		{Key: "job_name", Value: "后端开发"},
	}
	if !reflect.DeepEqual(filter, wantFilter) {
		t.Fatalf("buildHotQuestionFilter() = %#v, want %#v", filter, wantFilter)
	}

	findOptions := buildHotQuestionFindOptions(7)
	wantSort := bson.D{{Key: "hot_degree", Value: -1}}
	if !reflect.DeepEqual(findOptions.Sort, wantSort) {
		t.Fatalf("hot question sort = %#v, want %#v", findOptions.Sort, wantSort)
	}
	if findOptions.Limit == nil || *findOptions.Limit != 7 {
		t.Fatalf("hot question limit = %#v, want 7", findOptions.Limit)
	}
	wantProjection := bson.D{
		{Key: "_id", Value: 0},
		{Key: "question_id", Value: 1},
		{Key: "bank_list", Value: 1},
		{Key: "title", Value: 1},
		{Key: "view_count", Value: 1},
	}
	if !reflect.DeepEqual(findOptions.Projection, wantProjection) {
		t.Fatalf("hot question projection = %#v, want %#v", findOptions.Projection, wantProjection)
	}
}

func TestBuildHotQuestionFilterWithoutJobName(t *testing.T) {
	got := buildHotQuestionFilter("")
	want := bson.D{{Key: "status", Value: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildHotQuestionFilter() = %#v, want %#v", got, want)
	}
}

func TestBuildQuestionSort(t *testing.T) {
	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
		want      bson.D
	}{
		{name: "allowed descending field", sortBy: "view_count", sortOrder: "desc", want: bson.D{{Key: "view_count", Value: -1}}},
		{name: "allowed ascending field", sortBy: "create_time", sortOrder: "asc", want: bson.D{{Key: "create_time", Value: 1}}},
		{name: "allowed dislike count field", sortBy: "dislike_count", sortOrder: "desc", want: bson.D{{Key: "dislike_count", Value: -1}}},
		{name: "unknown field falls back to order", sortBy: "$where", sortOrder: "desc", want: bson.D{{Key: "order", Value: -1}}},
		{name: "unknown direction falls back to ascending", sortBy: "view_count", sortOrder: "invalid", want: bson.D{{Key: "view_count", Value: 1}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildQuestionSort(test.sortBy, test.sortOrder)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("buildQuestionSort(%q, %q) = %#v, want %#v", test.sortBy, test.sortOrder, got, test.want)
			}
		})
	}
}
