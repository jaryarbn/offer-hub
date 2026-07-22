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
	interactionTargetQuestion = 1
	interactionTargetComment  = 3
	interactionTypeLike       = 1
	interactionStatusInactive = 0
	interactionStatusActive   = 1
)

var ErrInteractionTargetNotFound = errors.New("interaction target not found")

type interactionTargetSpec struct {
	collection string
	idField    string
	countField string
	status     int
}

type interactionTargetCountRecord struct {
	QuestionCount int64 `bson:"thumbs_up_count"`
	CommentCount  int64 `bson:"thumbs_up"`
}

func (record interactionTargetCountRecord) count(targetType int) int64 {
	if targetType == interactionTargetComment {
		return record.CommentCount
	}
	return record.QuestionCount
}

// GetInteractionTargetLikeCount verifies the target exists and returns its current like count.
func (data *Data) GetInteractionTargetLikeCount(
	ctx context.Context,
	targetType int,
	targetID string,
) (int64, error) {
	spec, ok := interactionTargetSpecFor(targetType)
	if !ok {
		return 0, fmt.Errorf("unsupported interaction target_type: %d", targetType)
	}

	var record interactionTargetCountRecord
	err := data.MongoDB.Collection(spec.collection).FindOne(
		ctx,
		buildInteractionTargetFilter(spec, targetID),
		options.FindOne().SetProjection(bson.D{
			{Key: "_id", Value: 0},
			{Key: spec.countField, Value: 1},
		}),
	).Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, fmt.Errorf("%w: target_type=%d target_id=%s", ErrInteractionTargetNotFound, targetType, targetID)
	}
	if err != nil {
		return 0, fmt.Errorf("query interaction target: %w", err)
	}
	return record.count(targetType), nil
}

// SetInteractionLikeStatus persists a like-state transition and reports whether the state changed.
func (data *Data) SetInteractionLikeStatus(
	ctx context.Context,
	userID string,
	targetType int,
	targetID string,
	status int,
	updatedAt time.Time,
) (bool, error) {
	if status == interactionStatusActive {
		var previous struct {
			Status int `bson:"status"`
		}
		err := data.MongoDB.Collection(userInteractionsCollection).FindOneAndUpdate(
			ctx,
			buildInteractionIdentityFilter(userID, targetType, targetID),
			buildActivateInteractionUpdate(userID, targetType, targetID, updatedAt),
			options.FindOneAndUpdate().
				SetUpsert(true).
				SetReturnDocument(options.Before).
				SetProjection(bson.D{{Key: "_id", Value: 0}, {Key: "status", Value: 1}}),
		).Decode(&previous)
		if errors.Is(err, mongo.ErrNoDocuments) {
			// ReturnDocument(Before) has no document when the upsert inserted a new record.
			return true, nil
		}
		if err != nil {
			return false, fmt.Errorf("activate user interaction: %w", err)
		}
		return previous.Status != interactionStatusActive, nil
	}

	filter := append(
		buildInteractionIdentityFilter(userID, targetType, targetID),
		bson.E{Key: "status", Value: interactionStatusActive},
	)
	result, err := data.MongoDB.Collection(userInteractionsCollection).UpdateOne(
		ctx,
		filter,
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "status", Value: interactionStatusInactive},
			{Key: "update_time", Value: updatedAt},
		}}},
	)
	if err != nil {
		return false, fmt.Errorf("deactivate user interaction: %w", err)
	}
	return result.ModifiedCount > 0, nil
}

// AdjustInteractionTargetLikeCount applies a transition delta and returns the non-negative count.
func (data *Data) AdjustInteractionTargetLikeCount(
	ctx context.Context,
	targetType int,
	targetID string,
	delta int64,
) (int64, error) {
	if delta == 0 {
		return data.GetInteractionTargetLikeCount(ctx, targetType, targetID)
	}

	spec, ok := interactionTargetSpecFor(targetType)
	if !ok {
		return 0, fmt.Errorf("unsupported interaction target_type: %d", targetType)
	}
	filter := buildInteractionTargetFilter(spec, targetID)
	if delta < 0 {
		filter = append(filter, bson.E{
			Key:   spec.countField,
			Value: bson.D{{Key: "$gt", Value: 0}},
		})
	}

	var record interactionTargetCountRecord
	err := data.MongoDB.Collection(spec.collection).FindOneAndUpdate(
		ctx,
		filter,
		bson.D{{Key: "$inc", Value: bson.D{{Key: spec.countField, Value: delta}}}},
		options.FindOneAndUpdate().
			SetReturnDocument(options.After).
			SetProjection(bson.D{
				{Key: "_id", Value: 0},
				{Key: spec.countField, Value: 1},
			}),
	).Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) && delta < 0 {
		// A zero counter intentionally does not match the decrement filter.
		return data.GetInteractionTargetLikeCount(ctx, targetType, targetID)
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, fmt.Errorf("%w: target_type=%d target_id=%s", ErrInteractionTargetNotFound, targetType, targetID)
	}
	if err != nil {
		return 0, fmt.Errorf("adjust interaction target like count: %w", err)
	}
	return record.count(targetType), nil
}

// UpsertUserQuestionTag creates or updates one user's tag for a question.
func (data *Data) UpsertUserQuestionTag(
	ctx context.Context,
	userID string,
	questionID string,
	tag int,
	updatedAt time.Time,
) error {
	_, err := data.MongoDB.Collection(userQuestionTagCollection).UpdateOne(
		ctx,
		buildUserQuestionTagFilter(userID, questionID),
		buildUserQuestionTagUpdate(userID, questionID, tag, updatedAt),
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("upsert user question tag: %w", err)
	}
	return nil
}

func interactionTargetSpecFor(targetType int) (interactionTargetSpec, bool) {
	switch targetType {
	case interactionTargetQuestion:
		return interactionTargetSpec{
			collection: questionCollection,
			idField:    "question_id",
			countField: "thumbs_up_count",
			status:     questionNormalStatus,
		}, true
	case interactionTargetComment:
		return interactionTargetSpec{
			collection: commentsCollection,
			idField:    "comment_id",
			countField: "thumbs_up",
			status:     commentNormalStatus,
		}, true
	default:
		return interactionTargetSpec{}, false
	}
}

func buildInteractionIdentityFilter(userID string, targetType int, targetID string) bson.D {
	return bson.D{
		{Key: "user_id", Value: userID},
		{Key: "target_type", Value: targetType},
		{Key: "target_id", Value: targetID},
		{Key: "interaction_type", Value: interactionTypeLike},
	}
}

func buildInteractionTargetFilter(spec interactionTargetSpec, targetID string) bson.D {
	return bson.D{
		{Key: spec.idField, Value: targetID},
		{Key: "status", Value: spec.status},
	}
}

func buildActivateInteractionUpdate(
	userID string,
	targetType int,
	targetID string,
	updatedAt time.Time,
) bson.D {
	return bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "target_type", Value: targetType},
			{Key: "target_id", Value: targetID},
			{Key: "interaction_type", Value: interactionTypeLike},
			{Key: "status", Value: interactionStatusActive},
			{Key: "update_time", Value: updatedAt},
		}},
		{Key: "$setOnInsert", Value: bson.D{{Key: "create_time", Value: updatedAt}}},
	}
}

func buildUserQuestionTagUpdate(
	userID string,
	questionID string,
	tag int,
	updatedAt time.Time,
) bson.D {
	return bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "question_id", Value: questionID},
			{Key: "tag", Value: tag},
			{Key: "update_time", Value: updatedAt},
		}},
		{Key: "$setOnInsert", Value: bson.D{{Key: "create_time", Value: updatedAt}}},
	}
}
