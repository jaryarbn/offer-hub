package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
	"offer-hub/backend/src/tools"
)

const (
	defaultSubCommentPage = 1
	defaultSubCommentSize = 5
	normalCommentStatus   = 2
)

var (
	ErrInvalidCommentUserID     = errors.New("comment user_id is empty")
	ErrInvalidCommentID         = errors.New("comment_id is empty")
	ErrInvalidCommentTargetType = errors.New("invalid comment target_type")
	ErrInvalidCommentTargetID   = errors.New("comment target_id is empty")
	ErrInvalidCommentContent    = errors.New("comment content is empty")
	ErrCommentForbidden         = errors.New("comment does not belong to current user")
	ErrCommentNotFound          = data.ErrCommentNotFound
)

type CommentData interface {
	InsertComment(context.Context, data.CommentRecord) error
	IncrementCommentChildCount(context.Context, string) error
	GetCommentByID(context.Context, string) (data.CommentRecord, error)
	SoftDeleteComment(context.Context, string, time.Time) error
	UpdateCommentContent(context.Context, string, string, time.Time) error
	FilterComments(context.Context, data.CommentFilter) ([]data.CommentRecord, int64, error)
	ListSubCommentsByParentIDs(context.Context, []string) ([]data.CommentRecord, error)
	CountSubCommentsByParentIDs(context.Context, []string) (map[string]int64, error)
	ListCommentUsers(context.Context, []string) ([]data.CommentUserRecord, error)
	ListLikedCommentIDs(context.Context, string, []string) (map[string]bool, error)
}

func (service *CommentService) DeleteComment(
	ctx context.Context,
	req model.DeleteCommentReq,
	userID string,
) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidCommentUserID
	}

	req.CommentID = strings.TrimSpace(req.CommentID)
	if req.CommentID == "" {
		return ErrInvalidCommentID
	}

	record, err := service.data.GetCommentByID(ctx, req.CommentID)
	if err != nil {
		return fmt.Errorf("get comment: %w", err)
	}
	if strings.TrimSpace(record.UserID) != userID {
		return ErrCommentForbidden
	}

	if err := service.data.SoftDeleteComment(ctx, req.CommentID, service.now().UTC()); err != nil {
		return fmt.Errorf("soft-delete comment: %w", err)
	}
	return nil
}

func (service *CommentService) UpdateComment(
	ctx context.Context,
	req model.UpdateCommentReq,
	userID string,
) (model.UpdateCommentData, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return model.UpdateCommentData{}, ErrInvalidCommentUserID
	}

	req.CommentID = strings.TrimSpace(req.CommentID)
	if req.CommentID == "" {
		return model.UpdateCommentData{}, ErrInvalidCommentID
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return model.UpdateCommentData{}, ErrInvalidCommentContent
	}

	record, err := service.data.GetCommentByID(ctx, req.CommentID)
	if err != nil {
		return model.UpdateCommentData{}, fmt.Errorf("get comment: %w", err)
	}
	if strings.TrimSpace(record.UserID) != userID {
		return model.UpdateCommentData{}, ErrCommentForbidden
	}

	filteredContent := service.filterSensitiveWords(req.Content)
	if err := service.data.UpdateCommentContent(
		ctx,
		req.CommentID,
		filteredContent,
		service.now().UTC(),
	); err != nil {
		return model.UpdateCommentData{}, fmt.Errorf("update comment content: %w", err)
	}
	return model.UpdateCommentData{CommentID: req.CommentID}, nil
}

type CommentService struct {
	data                 CommentData
	newCommentID         func() string
	now                  func() time.Time
	filterSensitiveWords func(string) string
}

func NewCommentService(commentData CommentData) *CommentService {
	return &CommentService{
		data:                 commentData,
		newCommentID:         uuid.NewString,
		now:                  time.Now,
		filterSensitiveWords: tools.FilterSensitiveWords,
	}
}

func (service *CommentService) AddComment(
	ctx context.Context,
	req model.AddCommentReq,
	userID string,
) (model.AddCommentData, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return model.AddCommentData{}, ErrInvalidCommentUserID
	}
	if req.TargetType != 1 && req.TargetType != 2 {
		return model.AddCommentData{}, ErrInvalidCommentTargetType
	}

	req.TargetID = strings.TrimSpace(req.TargetID)
	if req.TargetID == "" {
		return model.AddCommentData{}, ErrInvalidCommentTargetID
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return model.AddCommentData{}, ErrInvalidCommentContent
	}

	req.ParentID = strings.TrimSpace(req.ParentID)
	req.ReplyTo = strings.TrimSpace(req.ReplyTo)
	createdAt := service.now().UTC()
	commentID := service.newCommentID()
	record := data.CommentRecord{
		ID:         commentID,
		CommentID:  commentID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		UserID:     userID,
		Content:    service.filterSensitiveWords(req.Content),
		ParentID:   req.ParentID,
		ReplyTo:    req.ReplyTo,
		ChildCount: 0,
		ThumbsUp:   0,
		ViewCount:  0,
		Status:     normalCommentStatus,
		CreateTime: createdAt,
		UpdateTime: createdAt,
	}
	if req.TargetType == 1 {
		record.QuestionID = req.TargetID
	}

	if err := service.data.InsertComment(ctx, record); err != nil {
		return model.AddCommentData{}, fmt.Errorf("insert comment: %w", err)
	}
	if req.ParentID != "" {
		if err := service.data.IncrementCommentChildCount(ctx, req.ParentID); err != nil {
			return model.AddCommentData{}, fmt.Errorf("increment parent comment child count: %w", err)
		}
	}

	users, err := service.loadCommentUsers(ctx, []data.CommentRecord{record})
	if err != nil {
		return model.AddCommentData{}, err
	}
	return model.AddCommentData{
		CommentID: commentID,
		Comment:   toCommentInfo(record, users, nil),
	}, nil
}

func (service *CommentService) ListComments(
	ctx context.Context,
	req model.ListCommentsReq,
	userID string,
) (model.ListCommentsData, error) {
	records, total, err := service.data.FilterComments(ctx, data.CommentFilter{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		ParentID:   req.ParentID,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return model.ListCommentsData{}, fmt.Errorf("filter comments: %w", err)
	}
	if len(records) == 0 {
		return model.ListCommentsData{Total: total, List: make([]model.CommentInfo, 0)}, nil
	}

	userID = strings.TrimSpace(userID)
	if req.ParentID != "" {
		comments, err := service.enrichComments(ctx, records, userID)
		if err != nil {
			return model.ListCommentsData{}, err
		}
		return model.ListCommentsData{Total: total, List: comments}, nil
	}

	parentIDs := collectCommentIDs(records)
	childRecords, err := service.data.ListSubCommentsByParentIDs(ctx, parentIDs)
	if err != nil {
		return model.ListCommentsData{}, fmt.Errorf("list child comments: %w", err)
	}
	childCounts, err := service.data.CountSubCommentsByParentIDs(ctx, parentIDs)
	if err != nil {
		return model.ListCommentsData{}, fmt.Errorf("count child comments: %w", err)
	}

	selectedChildren := paginateSubComments(
		groupCommentsByParent(childRecords),
		req.SubCommentPage,
		req.SubCommentSize,
	)
	returnedRecords := append(make([]data.CommentRecord, 0, len(records)+len(childRecords)), records...)
	for _, parentID := range parentIDs {
		returnedRecords = append(returnedRecords, selectedChildren[parentID]...)
	}

	users, liked, err := service.loadCommentMetadata(ctx, returnedRecords, userID)
	if err != nil {
		return model.ListCommentsData{}, err
	}

	comments := make([]model.CommentInfo, 0, len(records))
	for _, record := range records {
		comment := toCommentInfo(record, users, liked)
		comment.SubCommentTotal = childCounts[record.CommentID]
		children := selectedChildren[record.CommentID]
		comment.SubComments = make([]model.CommentInfo, 0, len(children))
		for _, child := range children {
			comment.SubComments = append(comment.SubComments, toCommentInfo(child, users, liked))
		}
		comments = append(comments, comment)
	}

	return model.ListCommentsData{Total: total, List: comments}, nil
}

func (service *CommentService) enrichComments(
	ctx context.Context,
	records []data.CommentRecord,
	userID string,
) ([]model.CommentInfo, error) {
	users, liked, err := service.loadCommentMetadata(ctx, records, userID)
	if err != nil {
		return nil, err
	}

	comments := make([]model.CommentInfo, 0, len(records))
	for _, record := range records {
		comments = append(comments, toCommentInfo(record, users, liked))
	}
	return comments, nil
}

func (service *CommentService) loadCommentMetadata(
	ctx context.Context,
	records []data.CommentRecord,
	userID string,
) (map[string]data.CommentUserRecord, map[string]bool, error) {
	users, err := service.loadCommentUsers(ctx, records)
	if err != nil {
		return nil, nil, err
	}

	liked := make(map[string]bool)
	if userID != "" {
		liked, err = service.data.ListLikedCommentIDs(ctx, userID, collectCommentIDs(records))
		if err != nil {
			return nil, nil, fmt.Errorf("list liked comments: %w", err)
		}
	}
	return users, liked, nil
}

func (service *CommentService) loadCommentUsers(
	ctx context.Context,
	records []data.CommentRecord,
) (map[string]data.CommentUserRecord, error) {
	userRecords, err := service.data.ListCommentUsers(ctx, collectCommentUserIDs(records))
	if err != nil {
		return nil, fmt.Errorf("list comment users: %w", err)
	}
	users := make(map[string]data.CommentUserRecord, len(userRecords))
	for _, record := range userRecords {
		users[record.UserID] = record
	}
	return users, nil
}

func toCommentInfo(
	record data.CommentRecord,
	users map[string]data.CommentUserRecord,
	liked map[string]bool,
) model.CommentInfo {
	user := users[record.UserID]
	replyToUser := users[record.ReplyTo]
	return model.CommentInfo{
		CommentID:       record.CommentID,
		UserID:          record.UserID,
		UserName:        commentUserName(user),
		UserAvatar:      user.Avatar,
		Content:         record.Content,
		ParentID:        record.ParentID,
		ReplyTo:         record.ReplyTo,
		ReplyToName:     commentUserName(replyToUser),
		Status:          record.Status,
		ThumbsUp:        record.ThumbsUp,
		SubCommentTotal: 0,
		UserLiked:       liked[record.CommentID],
		SubComments:     make([]model.CommentInfo, 0),
		CreateTime:      formatTime(record.CreateTime),
		UpdateTime:      formatTime(record.UpdateTime),
	}
}

func commentUserName(record data.CommentUserRecord) string {
	if strings.TrimSpace(record.NickName) != "" {
		return record.NickName
	}
	return record.Username
}

func collectCommentIDs(records []data.CommentRecord) []string {
	ids := make([]string, 0, len(records))
	seen := make(map[string]struct{}, len(records))
	for _, record := range records {
		if record.CommentID == "" {
			continue
		}
		if _, exists := seen[record.CommentID]; exists {
			continue
		}
		seen[record.CommentID] = struct{}{}
		ids = append(ids, record.CommentID)
	}
	return ids
}

func collectCommentUserIDs(records []data.CommentRecord) []string {
	ids := make([]string, 0, len(records)*2)
	seen := make(map[string]struct{}, len(records)*2)
	for _, record := range records {
		for _, userID := range []string{record.UserID, record.ReplyTo} {
			userID = strings.TrimSpace(userID)
			if userID == "" {
				continue
			}
			if _, exists := seen[userID]; exists {
				continue
			}
			seen[userID] = struct{}{}
			ids = append(ids, userID)
		}
	}
	return ids
}

func groupCommentsByParent(records []data.CommentRecord) map[string][]data.CommentRecord {
	grouped := make(map[string][]data.CommentRecord)
	for _, record := range records {
		grouped[record.ParentID] = append(grouped[record.ParentID], record)
	}
	return grouped
}

func paginateSubComments(
	grouped map[string][]data.CommentRecord,
	page int,
	pageSize int,
) map[string][]data.CommentRecord {
	if page <= 0 {
		page = defaultSubCommentPage
	}
	if pageSize <= 0 {
		pageSize = defaultSubCommentSize
	}

	result := make(map[string][]data.CommentRecord, len(grouped))
	start := (page - 1) * pageSize
	for parentID, records := range grouped {
		if start >= len(records) {
			result[parentID] = make([]data.CommentRecord, 0)
			continue
		}
		end := min(start+pageSize, len(records))
		result[parentID] = records[start:end]
	}
	return result
}
