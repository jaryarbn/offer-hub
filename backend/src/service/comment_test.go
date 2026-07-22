package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

type commentDataStub struct {
	insertedRecord    data.CommentRecord
	insertErr         error
	insertCalls       int
	incrementParentID string
	incrementErr      error
	incrementCalls    int
	getCommentID      string
	getCommentRecord  data.CommentRecord
	getCommentErr     error
	getCommentCalls   int
	deletedCommentID  string
	deletedAt         time.Time
	deleteErr         error
	deleteCalls       int
	updatedCommentID  string
	updatedContent    string
	updatedAt         time.Time
	updateErr         error
	updateCalls       int
	filter            data.CommentFilter
	filterRecords     []data.CommentRecord
	filterTotal       int64
	filterErr         error
	childParentIDs    []string
	childRecords      []data.CommentRecord
	childErr          error
	childCalls        int
	countParentIDs    []string
	childCounts       map[string]int64
	countErr          error
	countCalls        int
	userIDs           []string
	users             []data.CommentUserRecord
	usersErr          error
	userCalls         int
	likedUserID       string
	likedCommentIDs   []string
	liked             map[string]bool
	likedErr          error
	likedCalls        int
}

func (stub *commentDataStub) InsertComment(
	_ context.Context,
	record data.CommentRecord,
) error {
	stub.insertCalls++
	stub.insertedRecord = record
	return stub.insertErr
}

func (stub *commentDataStub) IncrementCommentChildCount(
	_ context.Context,
	parentID string,
) error {
	stub.incrementCalls++
	stub.incrementParentID = parentID
	return stub.incrementErr
}

func (stub *commentDataStub) GetCommentByID(
	_ context.Context,
	commentID string,
) (data.CommentRecord, error) {
	stub.getCommentCalls++
	stub.getCommentID = commentID
	return stub.getCommentRecord, stub.getCommentErr
}

func (stub *commentDataStub) SoftDeleteComment(
	_ context.Context,
	commentID string,
	updatedAt time.Time,
) error {
	stub.deleteCalls++
	stub.deletedCommentID = commentID
	stub.deletedAt = updatedAt
	return stub.deleteErr
}

func (stub *commentDataStub) UpdateCommentContent(
	_ context.Context,
	commentID string,
	content string,
	updatedAt time.Time,
) error {
	stub.updateCalls++
	stub.updatedCommentID = commentID
	stub.updatedContent = content
	stub.updatedAt = updatedAt
	return stub.updateErr
}

func (stub *commentDataStub) FilterComments(
	_ context.Context,
	filter data.CommentFilter,
) ([]data.CommentRecord, int64, error) {
	stub.filter = filter
	return stub.filterRecords, stub.filterTotal, stub.filterErr
}

func (stub *commentDataStub) ListSubCommentsByParentIDs(
	_ context.Context,
	parentIDs []string,
) ([]data.CommentRecord, error) {
	stub.childCalls++
	stub.childParentIDs = append([]string(nil), parentIDs...)
	return stub.childRecords, stub.childErr
}

func (stub *commentDataStub) CountSubCommentsByParentIDs(
	_ context.Context,
	parentIDs []string,
) (map[string]int64, error) {
	stub.countCalls++
	stub.countParentIDs = append([]string(nil), parentIDs...)
	return stub.childCounts, stub.countErr
}

func (stub *commentDataStub) ListCommentUsers(
	_ context.Context,
	userIDs []string,
) ([]data.CommentUserRecord, error) {
	stub.userCalls++
	stub.userIDs = append([]string(nil), userIDs...)
	return stub.users, stub.usersErr
}

func (stub *commentDataStub) ListLikedCommentIDs(
	_ context.Context,
	userID string,
	commentIDs []string,
) (map[string]bool, error) {
	stub.likedCalls++
	stub.likedUserID = userID
	stub.likedCommentIDs = append([]string(nil), commentIDs...)
	return stub.liked, stub.likedErr
}

func TestCommentServiceAddCommentFiltersAndReturnsCompleteComment(t *testing.T) {
	fixedTime := time.Date(2026, time.July, 22, 21, 15, 0, 0, time.FixedZone("CST", 8*60*60))
	stub := &commentDataStub{users: []data.CommentUserRecord{
		{UserID: "user-1", Username: "alice-account", NickName: "Alice", Avatar: "alice.png"},
		{UserID: "user-2", Username: "bob", NickName: "Bob", Avatar: "bob.png"},
	}}
	commentService := NewCommentService(stub)
	commentService.newCommentID = func() string { return "comment-new" }
	commentService.now = func() time.Time { return fixedTime }
	var filteredInput string
	commentService.filterSensitiveWords = func(content string) string {
		filteredInput = content
		return "filtered content"
	}

	got, err := commentService.AddComment(context.Background(), model.AddCommentReq{
		TargetType: 1,
		TargetID:   " question-1 ",
		ParentID:   " parent-1 ",
		ReplyTo:    " user-2 ",
		Content:    " original content ",
	}, " user-1 ")
	if err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}
	if filteredInput != "original content" {
		t.Fatalf("sensitive-word filter input = %q, want original content", filteredInput)
	}
	wantRecord := data.CommentRecord{
		ID:         "comment-new",
		CommentID:  "comment-new",
		TargetType: 1,
		TargetID:   "question-1",
		QuestionID: "question-1",
		UserID:     "user-1",
		Content:    "filtered content",
		ParentID:   "parent-1",
		ReplyTo:    "user-2",
		ChildCount: 0,
		ThumbsUp:   0,
		ViewCount:  0,
		Status:     2,
		CreateTime: fixedTime.UTC(),
		UpdateTime: fixedTime.UTC(),
	}
	if !reflect.DeepEqual(stub.insertedRecord, wantRecord) {
		t.Fatalf("inserted comment = %#v, want %#v", stub.insertedRecord, wantRecord)
	}
	if stub.insertCalls != 1 || stub.incrementCalls != 1 || stub.incrementParentID != "parent-1" {
		t.Fatalf("write calls = insert %d, increment %d parent %q", stub.insertCalls, stub.incrementCalls, stub.incrementParentID)
	}
	if !reflect.DeepEqual(stub.userIDs, []string{"user-1", "user-2"}) {
		t.Fatalf("comment user IDs = %#v", stub.userIDs)
	}
	if got.CommentID != "comment-new" || got.Comment.CommentID != "comment-new" {
		t.Fatalf("response IDs = %#v", got)
	}
	if got.Comment.UserName != "Alice" || got.Comment.UserAvatar != "alice.png" || got.Comment.ReplyToName != "Bob" {
		t.Fatalf("response user metadata = %#v", got.Comment)
	}
	if got.Comment.Content != "filtered content" || got.Comment.Status != 2 || got.Comment.UserLiked {
		t.Fatalf("response comment = %#v", got.Comment)
	}
	if got.Comment.SubComments == nil || len(got.Comment.SubComments) != 0 {
		t.Fatalf("response sub_comments = %#v, want []", got.Comment.SubComments)
	}
	if got.Comment.CreateTime != "2026-07-22 13:15:00" || got.Comment.UpdateTime != "2026-07-22 13:15:00" {
		t.Fatalf("response times = (%q, %q)", got.Comment.CreateTime, got.Comment.UpdateTime)
	}
}

func TestCommentServiceAddCommentValidatesInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.AddCommentReq
		userID  string
		wantErr error
	}{
		{name: "missing user", req: model.AddCommentReq{TargetType: 1, TargetID: "question-1", Content: "content"}, wantErr: ErrInvalidCommentUserID},
		{name: "invalid target type", req: model.AddCommentReq{TargetType: 3, TargetID: "question-1", Content: "content"}, userID: "user-1", wantErr: ErrInvalidCommentTargetType},
		{name: "blank target ID", req: model.AddCommentReq{TargetType: 1, TargetID: "  ", Content: "content"}, userID: "user-1", wantErr: ErrInvalidCommentTargetID},
		{name: "blank content", req: model.AddCommentReq{TargetType: 1, TargetID: "question-1", Content: "  "}, userID: "user-1", wantErr: ErrInvalidCommentContent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := &commentDataStub{}
			commentService := NewCommentService(stub)
			_, err := commentService.AddComment(context.Background(), test.req, test.userID)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("AddComment() error = %v, want %v", err, test.wantErr)
			}
			if stub.insertCalls != 0 || stub.incrementCalls != 0 || stub.userCalls != 0 {
				t.Fatalf("unexpected data calls: %#v", stub)
			}
		})
	}
}

func TestCommentServiceDeletesOwnedComment(t *testing.T) {
	fixedTime := time.Date(2026, time.July, 22, 23, 0, 0, 0, time.FixedZone("CST", 8*60*60))
	stub := &commentDataStub{getCommentRecord: data.CommentRecord{
		CommentID: "comment-1",
		UserID:    "user-1",
	}}
	commentService := NewCommentService(stub)
	commentService.now = func() time.Time { return fixedTime }

	err := commentService.DeleteComment(
		context.Background(),
		model.DeleteCommentReq{CommentID: " comment-1 "},
		" user-1 ",
	)
	if err != nil {
		t.Fatalf("DeleteComment() error = %v", err)
	}
	if stub.getCommentCalls != 1 || stub.getCommentID != "comment-1" {
		t.Fatalf("comment lookup = calls %d, ID %q", stub.getCommentCalls, stub.getCommentID)
	}
	if stub.deleteCalls != 1 || stub.deletedCommentID != "comment-1" {
		t.Fatalf("soft delete = calls %d, ID %q", stub.deleteCalls, stub.deletedCommentID)
	}
	if !stub.deletedAt.Equal(fixedTime.UTC()) {
		t.Fatalf("soft delete time = %v, want %v", stub.deletedAt, fixedTime.UTC())
	}
}

func TestCommentServiceRejectsDeletingAnotherUsersComment(t *testing.T) {
	stub := &commentDataStub{getCommentRecord: data.CommentRecord{
		CommentID: "comment-1",
		UserID:    "user-owner",
	}}
	commentService := NewCommentService(stub)

	err := commentService.DeleteComment(
		context.Background(),
		model.DeleteCommentReq{CommentID: "comment-1"},
		"user-other",
	)
	if !errors.Is(err, ErrCommentForbidden) {
		t.Fatalf("DeleteComment() error = %v, want %v", err, ErrCommentForbidden)
	}
	if stub.deleteCalls != 0 {
		t.Fatalf("soft delete calls = %d, want 0", stub.deleteCalls)
	}
}

func TestCommentServicePreservesCommentNotFoundError(t *testing.T) {
	stub := &commentDataStub{getCommentErr: data.ErrCommentNotFound}
	commentService := NewCommentService(stub)

	err := commentService.DeleteComment(
		context.Background(),
		model.DeleteCommentReq{CommentID: "comment-missing"},
		"user-1",
	)
	if !errors.Is(err, data.ErrCommentNotFound) {
		t.Fatalf("DeleteComment() error = %v, want wrapped %v", err, data.ErrCommentNotFound)
	}
	if stub.deleteCalls != 0 {
		t.Fatalf("soft delete calls = %d, want 0", stub.deleteCalls)
	}
}

func TestCommentServiceValidatesDeleteInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.DeleteCommentReq
		userID  string
		wantErr error
	}{
		{name: "missing user", req: model.DeleteCommentReq{CommentID: "comment-1"}, wantErr: ErrInvalidCommentUserID},
		{name: "blank comment ID", req: model.DeleteCommentReq{CommentID: "  "}, userID: "user-1", wantErr: ErrInvalidCommentID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := &commentDataStub{}
			commentService := NewCommentService(stub)
			err := commentService.DeleteComment(context.Background(), test.req, test.userID)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("DeleteComment() error = %v, want %v", err, test.wantErr)
			}
			if stub.getCommentCalls != 0 || stub.deleteCalls != 0 {
				t.Fatalf("unexpected data calls: %#v", stub)
			}
		})
	}
}

func TestCommentServiceUpdatesOwnedCommentWithFilteredContent(t *testing.T) {
	fixedTime := time.Date(2026, time.July, 22, 23, 30, 0, 0, time.FixedZone("CST", 8*60*60))
	stub := &commentDataStub{getCommentRecord: data.CommentRecord{
		CommentID: "comment-1",
		UserID:    "user-1",
	}}
	commentService := NewCommentService(stub)
	commentService.now = func() time.Time { return fixedTime }
	var filterInput string
	commentService.filterSensitiveWords = func(content string) string {
		filterInput = content
		return "filtered content"
	}

	got, err := commentService.UpdateComment(
		context.Background(),
		model.UpdateCommentReq{CommentID: " comment-1 ", Content: " new content "},
		" user-1 ",
	)
	if err != nil {
		t.Fatalf("UpdateComment() error = %v", err)
	}
	if got.CommentID != "comment-1" {
		t.Fatalf("UpdateComment() data = %#v, want comment-1", got)
	}
	if stub.getCommentCalls != 1 || stub.getCommentID != "comment-1" {
		t.Fatalf("comment lookup = calls %d, ID %q", stub.getCommentCalls, stub.getCommentID)
	}
	if filterInput != "new content" {
		t.Fatalf("sensitive-word filter input = %q, want new content", filterInput)
	}
	if stub.updateCalls != 1 || stub.updatedCommentID != "comment-1" || stub.updatedContent != "filtered content" {
		t.Fatalf(
			"comment update = calls %d, ID %q, content %q",
			stub.updateCalls,
			stub.updatedCommentID,
			stub.updatedContent,
		)
	}
	if !stub.updatedAt.Equal(fixedTime.UTC()) {
		t.Fatalf("comment update time = %v, want %v", stub.updatedAt, fixedTime.UTC())
	}
}

func TestCommentServiceRejectsUpdatingAnotherUsersComment(t *testing.T) {
	stub := &commentDataStub{getCommentRecord: data.CommentRecord{
		CommentID: "comment-1",
		UserID:    "user-owner",
	}}
	commentService := NewCommentService(stub)
	filterCalls := 0
	commentService.filterSensitiveWords = func(content string) string {
		filterCalls++
		return content
	}

	_, err := commentService.UpdateComment(
		context.Background(),
		model.UpdateCommentReq{CommentID: "comment-1", Content: "new content"},
		"user-other",
	)
	if !errors.Is(err, ErrCommentForbidden) {
		t.Fatalf("UpdateComment() error = %v, want %v", err, ErrCommentForbidden)
	}
	if filterCalls != 0 || stub.updateCalls != 0 {
		t.Fatalf("filter calls = %d, update calls = %d; want 0, 0", filterCalls, stub.updateCalls)
	}
}

func TestCommentServiceValidatesUpdateInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.UpdateCommentReq
		userID  string
		wantErr error
	}{
		{
			name:    "missing user",
			req:     model.UpdateCommentReq{CommentID: "comment-1", Content: "content"},
			wantErr: ErrInvalidCommentUserID,
		},
		{
			name:    "blank comment ID",
			req:     model.UpdateCommentReq{CommentID: "  ", Content: "content"},
			userID:  "user-1",
			wantErr: ErrInvalidCommentID,
		},
		{
			name:    "blank content",
			req:     model.UpdateCommentReq{CommentID: "comment-1", Content: "  "},
			userID:  "user-1",
			wantErr: ErrInvalidCommentContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := &commentDataStub{}
			commentService := NewCommentService(stub)
			_, err := commentService.UpdateComment(context.Background(), test.req, test.userID)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("UpdateComment() error = %v, want %v", err, test.wantErr)
			}
			if stub.getCommentCalls != 0 || stub.updateCalls != 0 {
				t.Fatalf("unexpected data calls: %#v", stub)
			}
		})
	}
}

func TestCommentServicePreservesUpdateNotFoundError(t *testing.T) {
	stub := &commentDataStub{getCommentErr: data.ErrCommentNotFound}
	commentService := NewCommentService(stub)

	_, err := commentService.UpdateComment(
		context.Background(),
		model.UpdateCommentReq{CommentID: "comment-missing", Content: "new content"},
		"user-1",
	)
	if !errors.Is(err, data.ErrCommentNotFound) {
		t.Fatalf("UpdateComment() error = %v, want wrapped %v", err, data.ErrCommentNotFound)
	}
	if stub.updateCalls != 0 {
		t.Fatalf("update calls = %d, want 0", stub.updateCalls)
	}
}

func TestCommentServiceListsTopCommentsWithPagedChildren(t *testing.T) {
	createdAt := time.Date(2026, time.July, 22, 10, 30, 0, 0, time.UTC)
	stub := &commentDataStub{
		filterRecords: []data.CommentRecord{
			{CommentID: "comment-1", UserID: "user-1", Content: "top one", Status: 2, ThumbsUp: 3, CreateTime: createdAt, UpdateTime: createdAt},
			{CommentID: "comment-2", UserID: "user-2", Content: "top two", Status: 2, CreateTime: createdAt, UpdateTime: createdAt},
		},
		filterTotal: 2,
		childRecords: []data.CommentRecord{
			{CommentID: "child-1-1", ParentID: "comment-1", UserID: "user-3", ReplyTo: "user-1", Content: "first child", Status: 2, CreateTime: createdAt, UpdateTime: createdAt},
			{CommentID: "child-1-2", ParentID: "comment-1", UserID: "user-4", Content: "second child", Status: 2, CreateTime: createdAt, UpdateTime: createdAt},
			{CommentID: "child-2-1", ParentID: "comment-2", UserID: "user-5", ReplyTo: "user-2", Content: "other child", Status: 2, CreateTime: createdAt, UpdateTime: createdAt},
		},
		childCounts: map[string]int64{"comment-1": 2, "comment-2": 1},
		users: []data.CommentUserRecord{
			{UserID: "user-1", Username: "alice-account", NickName: "Alice", Avatar: "alice.png"},
			{UserID: "user-2", Username: "bob", Avatar: "bob.png"},
			{UserID: "user-3", Username: "carol", NickName: "Carol", Avatar: "carol.png"},
			{UserID: "user-5", Username: "eve", NickName: "Eve", Avatar: "eve.png"},
		},
		liked: map[string]bool{"comment-1": true, "child-1-1": true},
	}
	commentService := NewCommentService(stub)
	req := model.ListCommentsReq{
		TargetType:     1,
		TargetID:       "question-1",
		SortBy:         "thumbs_up",
		SortOrder:      "desc",
		Page:           2,
		PageSize:       10,
		SubCommentPage: 1,
		SubCommentSize: 1,
	}

	got, err := commentService.ListComments(context.Background(), req, " user-current ")
	if err != nil {
		t.Fatalf("ListComments() error = %v", err)
	}
	wantFilter := data.CommentFilter{
		TargetType: 1,
		TargetID:   "question-1",
		SortBy:     "thumbs_up",
		SortOrder:  "desc",
		Page:       2,
		PageSize:   10,
	}
	if !reflect.DeepEqual(stub.filter, wantFilter) {
		t.Fatalf("comment filter = %#v, want %#v", stub.filter, wantFilter)
	}
	if !reflect.DeepEqual(stub.childParentIDs, []string{"comment-1", "comment-2"}) {
		t.Fatalf("child parent IDs = %#v", stub.childParentIDs)
	}
	if !reflect.DeepEqual(stub.countParentIDs, []string{"comment-1", "comment-2"}) {
		t.Fatalf("count parent IDs = %#v", stub.countParentIDs)
	}
	if !reflect.DeepEqual(stub.userIDs, []string{"user-1", "user-2", "user-3", "user-5"}) {
		t.Fatalf("user IDs = %#v", stub.userIDs)
	}
	if stub.likedUserID != "user-current" {
		t.Fatalf("liked user ID = %q, want user-current", stub.likedUserID)
	}
	if !reflect.DeepEqual(stub.likedCommentIDs, []string{"comment-1", "comment-2", "child-1-1", "child-2-1"}) {
		t.Fatalf("liked comment IDs = %#v", stub.likedCommentIDs)
	}

	if got.Total != 2 || len(got.List) != 2 {
		t.Fatalf("response data = %#v", got)
	}
	first := got.List[0]
	if first.UserName != "Alice" || first.UserAvatar != "alice.png" || !first.UserLiked {
		t.Fatalf("first top comment user metadata = %#v", first)
	}
	if first.SubCommentTotal != 2 || len(first.SubComments) != 1 {
		t.Fatalf("first top comment children = %#v", first)
	}
	child := first.SubComments[0]
	if child.CommentID != "child-1-1" || child.UserName != "Carol" || child.ReplyToName != "Alice" || !child.UserLiked {
		t.Fatalf("first child = %#v", child)
	}
	if child.SubComments == nil || len(child.SubComments) != 0 || child.SubCommentTotal != 0 {
		t.Fatalf("child nesting fields = %#v", child)
	}
	if got.List[1].UserName != "bob" || got.List[1].SubCommentTotal != 1 || len(got.List[1].SubComments) != 1 {
		t.Fatalf("second top comment = %#v", got.List[1])
	}
}

func TestCommentServiceListsChildrenWithoutNestedLookup(t *testing.T) {
	stub := &commentDataStub{
		filterRecords: []data.CommentRecord{
			{CommentID: "child-1", ParentID: "comment-1", UserID: "user-1", Status: 2},
		},
		filterTotal: 1,
		users: []data.CommentUserRecord{
			{UserID: "user-1", Username: "alice", Avatar: "alice.png"},
		},
	}
	commentService := NewCommentService(stub)

	got, err := commentService.ListComments(context.Background(), model.ListCommentsReq{
		TargetType: 1,
		TargetID:   "question-1",
		ParentID:   "comment-1",
		Page:       1,
		PageSize:   20,
	}, "")
	if err != nil {
		t.Fatalf("ListComments() error = %v", err)
	}
	if stub.childCalls != 0 || stub.countCalls != 0 {
		t.Fatalf("nested lookup calls = child %d, count %d; want 0", stub.childCalls, stub.countCalls)
	}
	if stub.likedCalls != 0 {
		t.Fatalf("liked lookup calls = %d, want 0 for guest", stub.likedCalls)
	}
	if got.Total != 1 || len(got.List) != 1 || got.List[0].UserName != "alice" {
		t.Fatalf("response data = %#v", got)
	}
	if got.List[0].SubComments == nil || len(got.List[0].SubComments) != 0 {
		t.Fatalf("child sub_comments = %#v, want []", got.List[0].SubComments)
	}
}

func TestCommentServiceReturnsEmptyListWithoutMetadataQueries(t *testing.T) {
	stub := &commentDataStub{filterRecords: []data.CommentRecord{}, filterTotal: 0}
	commentService := NewCommentService(stub)

	got, err := commentService.ListComments(context.Background(), model.ListCommentsReq{}, "user-1")
	if err != nil {
		t.Fatalf("ListComments() error = %v", err)
	}
	if got.List == nil || len(got.List) != 0 {
		t.Fatalf("ListComments() list = %#v, want []", got.List)
	}
	if stub.childCalls != 0 || stub.countCalls != 0 || stub.userCalls != 0 || stub.likedCalls != 0 {
		t.Fatalf("unexpected metadata calls: %#v", stub)
	}
}

func TestPaginateSubCommentsSupportsLaterPages(t *testing.T) {
	grouped := map[string][]data.CommentRecord{
		"comment-1": {
			{CommentID: "child-1"},
			{CommentID: "child-2"},
			{CommentID: "child-3"},
		},
	}
	got := paginateSubComments(grouped, 2, 2)
	if len(got["comment-1"]) != 1 || got["comment-1"][0].CommentID != "child-3" {
		t.Fatalf("paginateSubComments() = %#v", got)
	}
}
