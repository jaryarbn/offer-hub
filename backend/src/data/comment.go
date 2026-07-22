package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	commentsCollection             = "comments"
	commentInteractionsCollection  = "user_interactions"
	commentNormalStatus            = 2
	commentDeletedStatus           = 5
	commentInteractionTargetType   = 3
	commentLikeInteractionType     = 1
	commentActiveInteractionStatus = 1
	defaultCommentPage             = 1
	defaultCommentPageSize         = 20
)

var commentSortFields = map[string]string{
	"create_time": "create_time",
	"thumbs_up":   "thumbs_up",
}

var ErrCommentNotFound = errors.New("comment not found")

// CommentRecord mirrors the persisted fields in the comments collection.
type CommentRecord struct {
	ID         string    `bson:"_id,omitempty"`
	CommentID  string    `bson:"comment_id"`
	TargetType int       `bson:"target_type"`
	TargetID   string    `bson:"target_id"`
	QuestionID string    `bson:"question_id"`
	UserID     string    `bson:"user_id"`
	Content    string    `bson:"content"`
	ParentID   string    `bson:"parent_id"`
	ReplyTo    string    `bson:"reply_to"`
	ChildCount int       `bson:"child_count"`
	ThumbsUp   int       `bson:"thumbs_up"`
	ViewCount  int       `bson:"view_count"`
	Status     int       `bson:"status"`
	CreateTime time.Time `bson:"create_time"`
	UpdateTime time.Time `bson:"update_time"`
}

// CommentUserRecord contains the user fields needed to render a comment.
type CommentUserRecord struct {
	UserID   string `gorm:"column:user_id"`
	Username string `gorm:"column:username"`
	NickName string `gorm:"column:nick_name"`
	Avatar   string `gorm:"column:avatar"`
}

func (CommentUserRecord) TableName() string {
	return userInfoTable
}

type CommentFilter struct {
	TargetType int
	TargetID   string
	ParentID   string
	SortBy     string
	SortOrder  string
	Page       int
	PageSize   int
}

func (data *Data) InsertComment(ctx context.Context, record CommentRecord) error {
	if _, err := data.MongoDB.Collection(commentsCollection).InsertOne(ctx, record); err != nil {
		return fmt.Errorf("insert comment: %w", err)
	}
	return nil
}

func (data *Data) IncrementCommentChildCount(ctx context.Context, parentID string) error {
	result, err := data.MongoDB.Collection(commentsCollection).UpdateOne(
		ctx,
		buildParentCommentFilter(parentID),
		bson.D{{Key: "$inc", Value: bson.D{{Key: "child_count", Value: 1}}}},
	)
	if err != nil {
		return fmt.Errorf("increment parent comment child_count: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: %s", ErrCommentNotFound, parentID)
	}
	return nil
}

func (data *Data) GetCommentByID(ctx context.Context, commentID string) (CommentRecord, error) {
	var record CommentRecord
	err := data.MongoDB.Collection(commentsCollection).
		FindOne(ctx, buildCommentIDFilter(commentID)).
		Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return CommentRecord{}, fmt.Errorf("%w: %s", ErrCommentNotFound, commentID)
	}
	if err != nil {
		return CommentRecord{}, fmt.Errorf("query comment: %w", err)
	}
	return record, nil
}

func (data *Data) SoftDeleteComment(
	ctx context.Context,
	commentID string,
	updatedAt time.Time,
) error {
	result, err := data.MongoDB.Collection(commentsCollection).UpdateOne(
		ctx,
		buildCommentIDFilter(commentID),
		buildSoftDeleteCommentUpdate(updatedAt),
	)
	if err != nil {
		return fmt.Errorf("soft-delete comment: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: %s", ErrCommentNotFound, commentID)
	}
	return nil
}

func (data *Data) UpdateCommentContent(
	ctx context.Context,
	commentID string,
	content string,
	updatedAt time.Time,
) error {
	result, err := data.MongoDB.Collection(commentsCollection).UpdateOne(
		ctx,
		buildCommentIDFilter(commentID),
		buildUpdateCommentContentUpdate(content, updatedAt),
	)
	if err != nil {
		return fmt.Errorf("update comment content: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: %s", ErrCommentNotFound, commentID)
	}
	return nil
}

func (data *Data) FilterComments(
	ctx context.Context,
	filter CommentFilter,
) ([]CommentRecord, int64, error) {
	query := buildCommentFilter(filter)
	collection := data.MongoDB.Collection(commentsCollection)

	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("count %s: %w", commentsCollection, err)
	}

	page, pageSize := normalizeCommentPagination(filter.Page, filter.PageSize)
	cursor, err := collection.Find(
		ctx,
		query,
		options.Find().
			SetSort(buildCommentSort(filter.SortBy, filter.SortOrder)).
			SetSkip(int64((page-1)*pageSize)).
			SetLimit(int64(pageSize)),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("query %s: %w", commentsCollection, err)
	}
	defer cursor.Close(ctx)

	records := make([]CommentRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, 0, fmt.Errorf("decode %s: %w", commentsCollection, err)
	}
	return records, total, nil
}

func (data *Data) ListSubCommentsByParentIDs(
	ctx context.Context,
	parentIDs []string,
) ([]CommentRecord, error) {
	if len(parentIDs) == 0 {
		return make([]CommentRecord, 0), nil
	}

	cursor, err := data.MongoDB.Collection(commentsCollection).Find(
		ctx,
		buildSubCommentsFilter(parentIDs),
		options.Find().SetSort(bson.D{{Key: "create_time", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("query child comments: %w", err)
	}
	defer cursor.Close(ctx)

	records := make([]CommentRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode child comments: %w", err)
	}
	return records, nil
}

func (data *Data) CountSubCommentsByParentIDs(
	ctx context.Context,
	parentIDs []string,
) (map[string]int64, error) {
	counts := make(map[string]int64, len(parentIDs))
	if len(parentIDs) == 0 {
		return counts, nil
	}

	cursor, err := data.MongoDB.Collection(commentsCollection).Aggregate(
		ctx,
		buildSubCommentCountPipeline(parentIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("count child comments: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		ParentID string `bson:"_id"`
		Count    int64  `bson:"count"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode child comment counts: %w", err)
	}
	for _, result := range results {
		counts[result.ParentID] = result.Count
	}
	return counts, nil
}

func (data *Data) ListCommentUsers(
	ctx context.Context,
	userIDs []string,
) ([]CommentUserRecord, error) {
	if len(userIDs) == 0 {
		return make([]CommentUserRecord, 0), nil
	}

	records := make([]CommentUserRecord, 0)
	err := data.MySQL.WithContext(ctx).
		Select("user_id", "username", "nick_name", "avatar").
		Where("user_id IN ?", userIDs).
		Find(&records).
		Error
	if err != nil {
		return nil, fmt.Errorf("query comment users: %w", err)
	}
	return records, nil
}

func (data *Data) ListLikedCommentIDs(
	ctx context.Context,
	userID string,
	commentIDs []string,
) (map[string]bool, error) {
	liked := make(map[string]bool, len(commentIDs))
	if userID == "" || len(commentIDs) == 0 {
		return liked, nil
	}

	cursor, err := data.MongoDB.Collection(commentInteractionsCollection).Find(
		ctx,
		buildLikedCommentsFilter(userID, commentIDs),
		options.Find().SetProjection(bson.D{
			{Key: "_id", Value: 0},
			{Key: "target_id", Value: 1},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("query liked comments: %w", err)
	}
	defer cursor.Close(ctx)

	var records []struct {
		CommentID string `bson:"target_id"`
	}
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode liked comments: %w", err)
	}
	for _, record := range records {
		liked[record.CommentID] = true
	}
	return liked, nil
}

func buildCommentFilter(filter CommentFilter) bson.D {
	return bson.D{
		{Key: "target_type", Value: filter.TargetType},
		{Key: "target_id", Value: filter.TargetID},
		{Key: "parent_id", Value: filter.ParentID},
		{Key: "status", Value: commentNormalStatus},
	}
}

func buildParentCommentFilter(parentID string) bson.D {
	return bson.D{
		{Key: "comment_id", Value: parentID},
		{Key: "parent_id", Value: ""},
		{Key: "status", Value: commentNormalStatus},
	}
}

func buildCommentIDFilter(commentID string) bson.D {
	return bson.D{{Key: "comment_id", Value: commentID}}
}

func buildSoftDeleteCommentUpdate(updatedAt time.Time) bson.D {
	return bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: commentDeletedStatus},
		{Key: "update_time", Value: updatedAt},
	}}}
}

func buildUpdateCommentContentUpdate(content string, updatedAt time.Time) bson.D {
	return bson.D{{Key: "$set", Value: bson.D{
		{Key: "content", Value: content},
		{Key: "update_time", Value: updatedAt},
	}}}
}

func buildCommentSort(sortBy, sortOrder string) bson.D {
	field, exists := commentSortFields[sortBy]
	if !exists {
		field = "create_time"
	}

	direction := -1
	if sortOrder == "asc" {
		direction = 1
	}
	return bson.D{{Key: field, Value: direction}}
}

func buildSubCommentsFilter(parentIDs []string) bson.D {
	return bson.D{
		{Key: "parent_id", Value: bson.D{{Key: "$in", Value: parentIDs}}},
		{Key: "status", Value: commentNormalStatus},
	}
}

func buildSubCommentCountPipeline(parentIDs []string) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: buildSubCommentsFilter(parentIDs)}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$parent_id"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
}

func buildLikedCommentsFilter(userID string, commentIDs []string) bson.D {
	return bson.D{
		{Key: "user_id", Value: userID},
		{Key: "target_type", Value: commentInteractionTargetType},
		{Key: "target_id", Value: bson.D{{Key: "$in", Value: commentIDs}}},
		{Key: "interaction_type", Value: commentLikeInteractionType},
		{Key: "status", Value: commentActiveInteractionStatus},
	}
}

func normalizeCommentPagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = defaultCommentPage
	}
	if pageSize <= 0 {
		pageSize = defaultCommentPageSize
	}
	return page, pageSize
}
