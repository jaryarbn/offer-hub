package data

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	questionBankSeriesCollection  = "question_bank_series"
	questionBankCollection        = "question_bank"
	questionCollection            = "question"
	userInteractionsCollection    = "user_interactions"
	userQuestionTagCollection     = "user_question_tag"
	questionNormalStatus          = 1
	questionInteractionTargetType = 1
	questionLikeInteractionType   = 1
	activeInteractionStatus       = 1
)

var ErrQuestionNotFound = errors.New("question not found")

type QuestionBankSeriesRecord struct {
	SeriesID   string `bson:"series_id"`
	SeriesName string `bson:"series_name"`
	JobName    string `bson:"job_name"`
	Order      int64  `bson:"order"`
}

type QuestionBankRecord struct {
	BankID   string `bson:"bank_id"`
	SeriesID string `bson:"series_id"`
	BankName string `bson:"bank_name"`
	BankLogo string `bson:"bank_logo"`
	Desc     string `bson:"desc"`
	JobName  string `bson:"job_name"`
	Order    int64  `bson:"order"`
}

type QuestionRecord struct {
	QuestionID      string    `bson:"question_id"`
	BankList        []string  `bson:"bank_list"`
	JobName         string    `bson:"job_name"`
	Title           string    `bson:"title"`
	Content         string    `bson:"content"`
	AnalysisContent string    `bson:"analysis_content"`
	Difficulty      int       `bson:"difficulty"`
	Tags            []string  `bson:"tags"`
	Status          int       `bson:"status"`
	VIP             bool      `bson:"vip"`
	HotDegree       int       `bson:"hot_degree"`
	ViewCount       int       `bson:"view_count"`
	ThumbsUpCount   int       `bson:"thumbs_up_count"`
	DislikeCount    int       `bson:"dislike_count"`
	Order           int64     `bson:"order"`
	CreateTime      time.Time `bson:"create_time"`
	UpdateTime      time.Time `bson:"update_time"`
}

type HotQuestionRecord struct {
	QuestionID string   `bson:"question_id"`
	BankList   []string `bson:"bank_list"`
	Title      string   `bson:"title"`
	ViewCount  int      `bson:"view_count"`
}

type QuestionFilter struct {
	BankID     string
	Keyword    string
	Difficulty int
	Tags       []string
	JobName    string
	SortBy     string
	SortOrder  string
	Page       int
	PageSize   int
}

var questionSortFields = map[string]string{
	"view_count":      "view_count",
	"thumbs_up_count": "thumbs_up_count",
	"dislike_count":   "dislike_count",
	"create_time":     "create_time",
}

func (data *Data) FilterQuestions(
	ctx context.Context,
	filter QuestionFilter,
) ([]QuestionRecord, int64, error) {
	query := buildQuestionFilter(filter)
	collection := data.MongoDB.Collection(questionCollection)

	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("count %s: %w", questionCollection, err)
	}

	findOptions := options.Find().
		SetSort(buildQuestionSort(filter.SortBy, filter.SortOrder)).
		SetSkip(int64((filter.Page - 1) * filter.PageSize)).
		SetLimit(int64(filter.PageSize))
	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("query %s: %w", questionCollection, err)
	}
	defer cursor.Close(ctx)

	records := make([]QuestionRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, 0, fmt.Errorf("decode %s: %w", questionCollection, err)
	}
	return records, total, nil
}

func (data *Data) GetQuestionByID(ctx context.Context, questionID string) (QuestionRecord, error) {
	var record QuestionRecord
	err := data.MongoDB.Collection(questionCollection).
		FindOne(ctx, buildQuestionDetailFilter(questionID)).
		Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return QuestionRecord{}, fmt.Errorf("%w: %s", ErrQuestionNotFound, questionID)
	}
	if err != nil {
		return QuestionRecord{}, fmt.Errorf("query %s by question_id: %w", questionCollection, err)
	}
	return record, nil
}

func (data *Data) IsQuestionLiked(ctx context.Context, userID, questionID string) (bool, error) {
	err := data.MongoDB.Collection(userInteractionsCollection).
		FindOne(
			ctx,
			buildQuestionLikeFilter(userID, questionID),
			options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}}),
		).
		Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("query %s for question like: %w", userInteractionsCollection, err)
	}
	return true, nil
}

func (data *Data) GetUserQuestionTag(ctx context.Context, userID, questionID string) (int, error) {
	var record struct {
		Tag int `bson:"tag"`
	}
	err := data.MongoDB.Collection(userQuestionTagCollection).
		FindOne(
			ctx,
			buildUserQuestionTagFilter(userID, questionID),
			options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 0}, {Key: "tag", Value: 1}}),
		).
		Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("query %s for question tag: %w", userQuestionTagCollection, err)
	}
	return record.Tag, nil
}

func buildQuestionLikeFilter(userID, questionID string) bson.D {
	return bson.D{
		{Key: "user_id", Value: userID},
		{Key: "target_type", Value: questionInteractionTargetType},
		{Key: "target_id", Value: questionID},
		{Key: "interaction_type", Value: questionLikeInteractionType},
		{Key: "status", Value: activeInteractionStatus},
	}
}

func buildUserQuestionTagFilter(userID, questionID string) bson.D {
	return bson.D{
		{Key: "user_id", Value: userID},
		{Key: "question_id", Value: questionID},
	}
}

func (data *Data) ListHotQuestions(
	ctx context.Context,
	jobName string,
	limit int,
) ([]HotQuestionRecord, error) {
	cursor, err := data.MongoDB.Collection(questionCollection).Find(
		ctx,
		buildHotQuestionFilter(jobName),
		buildHotQuestionFindOptions(limit),
	)
	if err != nil {
		return nil, fmt.Errorf("query hot questions: %w", err)
	}
	defer cursor.Close(ctx)

	records := make([]HotQuestionRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode hot questions: %w", err)
	}
	return records, nil
}

func buildHotQuestionFilter(jobName string) bson.D {
	filter := bson.D{{Key: "status", Value: questionNormalStatus}}
	if jobName != "" {
		filter = append(filter, bson.E{Key: "job_name", Value: jobName})
	}
	return filter
}

func buildHotQuestionFindOptions(limit int) *options.FindOptions {
	return options.Find().
		SetSort(bson.D{{Key: "hot_degree", Value: -1}}).
		SetLimit(int64(limit)).
		SetProjection(bson.D{
			{Key: "_id", Value: 0},
			{Key: "question_id", Value: 1},
			{Key: "bank_list", Value: 1},
			{Key: "title", Value: 1},
			{Key: "view_count", Value: 1},
		})
}

func buildQuestionDetailFilter(questionID string) bson.D {
	return bson.D{
		{Key: "question_id", Value: questionID},
		{Key: "status", Value: questionNormalStatus},
	}
}

func buildQuestionFilter(filter QuestionFilter) bson.D {
	query := bson.D{{Key: "status", Value: questionNormalStatus}}
	if filter.BankID != "" {
		query = append(query, bson.E{Key: "bank_list", Value: filter.BankID})
	}
	if filter.Keyword != "" {
		keyword := bson.D{
			{Key: "$regex", Value: regexp.QuoteMeta(filter.Keyword)},
			{Key: "$options", Value: "i"},
		}
		query = append(query, bson.E{Key: "$or", Value: bson.A{
			bson.D{{Key: "title", Value: keyword}},
			bson.D{{Key: "content", Value: keyword}},
		}})
	}
	if filter.Difficulty != 0 {
		query = append(query, bson.E{Key: "difficulty", Value: filter.Difficulty})
	}
	if len(filter.Tags) > 0 {
		query = append(query, bson.E{Key: "tags", Value: bson.D{{Key: "$all", Value: filter.Tags}}})
	}
	if filter.JobName != "" {
		query = append(query, bson.E{Key: "job_name", Value: filter.JobName})
	}
	return query
}

func buildQuestionSort(sortBy, sortOrder string) bson.D {
	field, exists := questionSortFields[sortBy]
	if !exists {
		field = "order"
	}
	direction := 1
	if sortOrder == "desc" {
		direction = -1
	}
	return bson.D{{Key: field, Value: direction}}
}

func (data *Data) ListQuestionBankSeries(ctx context.Context, jobName string) ([]QuestionBankSeriesRecord, error) {
	filter := bson.D{}
	if jobName != "" {
		filter = append(filter, bson.E{Key: "job_name", Value: jobName})
	}

	cursor, err := data.MongoDB.Collection(questionBankSeriesCollection).Find(
		ctx,
		filter,
		options.Find().SetSort(bson.D{{Key: "order", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", questionBankSeriesCollection, err)
	}
	defer cursor.Close(ctx)

	records := make([]QuestionBankSeriesRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode %s: %w", questionBankSeriesCollection, err)
	}
	return records, nil
}

func (data *Data) ListQuestionBanks(ctx context.Context, jobName string) ([]QuestionBankRecord, error) {
	filter := bson.D{}
	if jobName != "" {
		filter = append(filter, bson.E{Key: "job_name", Value: jobName})
	}

	cursor, err := data.MongoDB.Collection(questionBankCollection).Find(
		ctx,
		filter,
		options.Find().SetSort(bson.D{{Key: "order", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", questionBankCollection, err)
	}
	defer cursor.Close(ctx)

	records := make([]QuestionBankRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode %s: %w", questionBankCollection, err)
	}
	return records, nil
}

func (data *Data) CountNormalQuestionsByBank(ctx context.Context, bankIDs []string) (map[string]int64, error) {
	counts := make(map[string]int64, len(bankIDs))
	if len(bankIDs) == 0 {
		return counts, nil
	}

	cursor, err := data.MongoDB.Collection(questionCollection).Aggregate(
		ctx,
		buildQuestionCountPipeline(bankIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("aggregate %s counts: %w", questionCollection, err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		BankID string `bson:"_id"`
		Count  int64  `bson:"count"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode %s counts: %w", questionCollection, err)
	}
	for _, result := range results {
		counts[result.BankID] = result.Count
	}
	return counts, nil
}

func buildQuestionCountPipeline(bankIDs []string) mongo.Pipeline {
	bankFilter := bson.D{{Key: "$in", Value: bankIDs}}
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "status", Value: questionNormalStatus},
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
}
