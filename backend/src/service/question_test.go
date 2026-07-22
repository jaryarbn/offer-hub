package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

type questionDataStub struct {
	series       []data.QuestionBankSeriesRecord
	banks        []data.QuestionBankRecord
	counts       map[string]int64
	seriesFilter string
	bankFilter   string
	countBankIDs []string
	questions    []data.QuestionRecord
	total        int64
	filter       data.QuestionFilter
	question     data.QuestionRecord
	questionID   string
	questionErr  error
	liked        bool
	likedUserID  string
	likedQID     string
	likedErr     error
	userTag      int
	tagUserID    string
	tagQID       string
	userTagErr   error
	hotQuestions []data.HotQuestionRecord
	hotJobName   string
	hotLimit     int
}

func (stub *questionDataStub) ListQuestionBankSeries(
	_ context.Context,
	jobName string,
) ([]data.QuestionBankSeriesRecord, error) {
	stub.seriesFilter = jobName
	return stub.series, nil
}

func (stub *questionDataStub) ListQuestionBanks(
	_ context.Context,
	jobName string,
) ([]data.QuestionBankRecord, error) {
	stub.bankFilter = jobName
	return stub.banks, nil
}

func (stub *questionDataStub) CountNormalQuestionsByBank(
	_ context.Context,
	bankIDs []string,
) (map[string]int64, error) {
	stub.countBankIDs = bankIDs
	return stub.counts, nil
}

func (stub *questionDataStub) FilterQuestions(
	_ context.Context,
	filter data.QuestionFilter,
) ([]data.QuestionRecord, int64, error) {
	stub.filter = filter
	return stub.questions, stub.total, nil
}

func (stub *questionDataStub) GetQuestionByID(
	_ context.Context,
	questionID string,
) (data.QuestionRecord, error) {
	stub.questionID = questionID
	return stub.question, stub.questionErr
}

func (stub *questionDataStub) IsQuestionLiked(
	_ context.Context,
	userID string,
	questionID string,
) (bool, error) {
	stub.likedUserID = userID
	stub.likedQID = questionID
	return stub.liked, stub.likedErr
}

func (stub *questionDataStub) GetUserQuestionTag(
	_ context.Context,
	userID string,
	questionID string,
) (int, error) {
	stub.tagUserID = userID
	stub.tagQID = questionID
	return stub.userTag, stub.userTagErr
}

func (stub *questionDataStub) ListHotQuestions(
	_ context.Context,
	jobName string,
	limit int,
) ([]data.HotQuestionRecord, error) {
	stub.hotJobName = jobName
	stub.hotLimit = limit
	return stub.hotQuestions, nil
}

func TestQuestionServiceListQuestions(t *testing.T) {
	content := strings.Repeat("题", 201)
	createdAt := time.Date(2026, time.July, 21, 10, 30, 0, 0, time.Local)
	stub := &questionDataStub{
		questions: []data.QuestionRecord{{
			QuestionID:    "question-1",
			BankList:      nil,
			Title:         "Go 并发",
			Content:       content,
			Difficulty:    2,
			Tags:          nil,
			Status:        1,
			VIP:           true,
			HotDegree:     9,
			ViewCount:     10,
			ThumbsUpCount: 8,
			DislikeCount:  1,
			Order:         7,
			CreateTime:    createdAt,
			UpdateTime:    createdAt.Add(time.Hour),
		}},
		total: 1,
	}
	questionService := NewQuestionService(stub)
	req := model.ListQuestionReq{
		BankID: "bank-1", Keyword: "Go", Difficulty: 2,
		Tags: []string{"Go", "并发"}, JobName: "后端开发", UserTag: 2,
		SortBy: "view_count", SortOrder: "desc", Page: 3, PageSize: 10,
	}

	got, err := questionService.ListQuestions(context.Background(), req, "user-1")
	if err != nil {
		t.Fatalf("ListQuestions() error = %v", err)
	}
	wantFilter := data.QuestionFilter{
		BankID: "bank-1", Keyword: "Go", Difficulty: 2,
		Tags: []string{"Go", "并发"}, JobName: "后端开发", UserID: "user-1", UserTag: 2,
		SortBy: "view_count", SortOrder: "desc", Page: 3, PageSize: 10,
	}
	if !reflect.DeepEqual(stub.filter, wantFilter) {
		t.Fatalf("FilterQuestions filter = %#v, want %#v", stub.filter, wantFilter)
	}
	if got.Total != 1 || len(got.List) != 1 {
		t.Fatalf("ListQuestions() = %#v, want one result", got)
	}
	question := got.List[0]
	if len([]rune(question.Content)) != visitorContentLength {
		t.Fatalf("content rune length = %d, want %d", len([]rune(question.Content)), visitorContentLength)
	}
	if question.BankList == nil || question.Tags == nil {
		t.Fatalf("list fields must be non-nil: bank_list=%#v tags=%#v", question.BankList, question.Tags)
	}
	if question.UserTag != 0 || question.UserLiked {
		t.Fatalf("visitor fields = (%d, %t), want (0, false)", question.UserTag, question.UserLiked)
	}
	if question.CreateTime != "2026-07-21 10:30:00" || question.UpdateTime != "2026-07-21 11:30:00" {
		t.Fatalf("formatted times = (%q, %q)", question.CreateTime, question.UpdateTime)
	}
}

func TestQuestionServiceListQuestionsReturnsEmptyArray(t *testing.T) {
	questionService := NewQuestionService(&questionDataStub{})

	got, err := questionService.ListQuestions(context.Background(), model.ListQuestionReq{}, "")
	if err != nil {
		t.Fatalf("ListQuestions() error = %v", err)
	}
	if got.List == nil || len(got.List) != 0 {
		t.Fatalf("ListQuestions().List = %#v, want non-nil empty slice", got.List)
	}
}

func TestQuestionServiceListQuestionsMetaUsesSharedFilter(t *testing.T) {
	stub := &questionDataStub{
		questions: []data.QuestionRecord{
			{QuestionID: "question-1", Title: "Go 并发", Content: "must not be exposed"},
			{QuestionID: "question-2", Title: "MongoDB 索引", Content: "must not be exposed"},
		},
		total: 8,
	}
	questionService := NewQuestionService(stub)
	req := model.ListQuestionMetaReq{
		BankID: "bank-1", Keyword: "Go", Difficulty: 2,
		Tags: []string{"Go", "并发"}, JobName: "后端开发", UserTag: 1,
		SortBy: "dislike_count", SortOrder: "desc", Page: 2, PageSize: 5,
	}

	got, err := questionService.ListQuestionsMeta(context.Background(), req, "user-2")
	if err != nil {
		t.Fatalf("ListQuestionsMeta() error = %v", err)
	}
	wantFilter := data.QuestionFilter{
		BankID: "bank-1", Keyword: "Go", Difficulty: 2,
		Tags: []string{"Go", "并发"}, JobName: "后端开发", UserID: "user-2", UserTag: 1,
		SortBy: "dislike_count", SortOrder: "desc", Page: 2, PageSize: 5,
	}
	if !reflect.DeepEqual(stub.filter, wantFilter) {
		t.Fatalf("FilterQuestions filter = %#v, want %#v", stub.filter, wantFilter)
	}
	want := model.ListQuestionMetaResp{
		Code: 0,
		Msg:  "success",
		Data: model.ListQuestionMetaResponseData{
			Total: 8,
			List: []model.QuestionMetaInfo{
				{QuestionID: "question-1", Title: "Go 并发"},
				{QuestionID: "question-2", Title: "MongoDB 索引"},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListQuestionsMeta() = %#v, want %#v", got, want)
	}
}

func TestQuestionServiceListQuestionsMetaReturnsEmptyArray(t *testing.T) {
	questionService := NewQuestionService(&questionDataStub{})

	got, err := questionService.ListQuestionsMeta(context.Background(), model.ListQuestionMetaReq{}, "")
	if err != nil {
		t.Fatalf("ListQuestionsMeta() error = %v", err)
	}
	if got.Data.List == nil || len(got.Data.List) != 0 {
		t.Fatalf("ListQuestionsMeta().Data.List = %#v, want non-nil empty slice", got.Data.List)
	}
}

func TestQuestionServiceGetQuestionDetailTruncatesVisitorContent(t *testing.T) {
	content := strings.Repeat("完整内容", 60)
	createdAt := time.Date(2026, time.July, 21, 10, 30, 0, 0, time.Local)
	stub := &questionDataStub{question: data.QuestionRecord{
		QuestionID:      "question-1",
		BankList:        nil,
		Title:           "Go 并发",
		Content:         content,
		AnalysisContent: "仅登录用户可见的解析",
		Difficulty:      2,
		Tags:            nil,
		Status:          1,
		VIP:             true,
		HotDegree:       9,
		ViewCount:       10,
		ThumbsUpCount:   8,
		DislikeCount:    1,
		Order:           7,
		CreateTime:      createdAt,
		UpdateTime:      createdAt.Add(time.Hour),
	}}
	questionService := NewQuestionService(stub)

	got, err := questionService.GetQuestionDetail(
		context.Background(),
		model.GetQuestionDetailReq{QuestionID: "question-1"},
		"",
	)
	if err != nil {
		t.Fatalf("GetQuestionDetail() error = %v", err)
	}
	if stub.questionID != "question-1" {
		t.Fatalf("GetQuestionByID questionID = %q", stub.questionID)
	}
	if len([]rune(got.Content)) != visitorContentLength {
		t.Fatalf("content rune length = %d, want %d", len([]rune(got.Content)), visitorContentLength)
	}
	if got.AnalysisContent != "" {
		t.Fatalf("analysis_content = %q, want empty", got.AnalysisContent)
	}
	if got.BankList == nil || got.Tags == nil {
		t.Fatalf("list fields must be non-nil: bank_list=%#v tags=%#v", got.BankList, got.Tags)
	}
	if got.UserTag != 0 || got.UserLiked {
		t.Fatalf("visitor fields = (%d, %t), want (0, false)", got.UserTag, got.UserLiked)
	}
	if stub.likedUserID != "" || stub.tagUserID != "" {
		t.Fatalf("visitor interaction lookups = (%q, %q), want none", stub.likedUserID, stub.tagUserID)
	}
	if got.CreateTime != "2026-07-21 10:30:00" || got.UpdateTime != "2026-07-21 11:30:00" {
		t.Fatalf("formatted times = (%q, %q)", got.CreateTime, got.UpdateTime)
	}
}

func TestQuestionServiceGetQuestionDetailReturnsAuthenticatedUserState(t *testing.T) {
	content := strings.Repeat("完整内容", 60)
	stub := &questionDataStub{
		question: data.QuestionRecord{
			QuestionID:      "question-1",
			Content:         content,
			AnalysisContent: "完整解析",
		},
		liked:   true,
		userTag: 3,
	}
	questionService := NewQuestionService(stub)

	got, err := questionService.GetQuestionDetail(
		context.Background(),
		model.GetQuestionDetailReq{QuestionID: "question-1"},
		"user-1",
	)
	if err != nil {
		t.Fatalf("GetQuestionDetail() error = %v", err)
	}
	if got.Content != content || got.AnalysisContent != "完整解析" {
		t.Fatalf("authenticated content = (%q, %q)", got.Content, got.AnalysisContent)
	}
	if !got.UserLiked || got.UserTag != 3 {
		t.Fatalf("authenticated user state = (%t, %d), want (true, 3)", got.UserLiked, got.UserTag)
	}
	if stub.likedUserID != "user-1" || stub.likedQID != "question-1" {
		t.Fatalf("like lookup args = (%q, %q)", stub.likedUserID, stub.likedQID)
	}
	if stub.tagUserID != "user-1" || stub.tagQID != "question-1" {
		t.Fatalf("tag lookup args = (%q, %q)", stub.tagUserID, stub.tagQID)
	}
}

func TestQuestionServiceGetQuestionDetailReturnsInteractionLookupError(t *testing.T) {
	stub := &questionDataStub{
		question: data.QuestionRecord{QuestionID: "question-1"},
		likedErr: errors.New("mongo unavailable"),
	}
	questionService := NewQuestionService(stub)

	_, err := questionService.GetQuestionDetail(
		context.Background(),
		model.GetQuestionDetailReq{QuestionID: "question-1"},
		"user-1",
	)
	if err == nil || !strings.Contains(err.Error(), "get question like state") {
		t.Fatalf("GetQuestionDetail() error = %v", err)
	}
}

func TestQuestionServiceGetQuestionDetailReturnsTagLookupError(t *testing.T) {
	stub := &questionDataStub{
		question:   data.QuestionRecord{QuestionID: "question-1"},
		userTagErr: errors.New("mongo unavailable"),
	}
	questionService := NewQuestionService(stub)

	_, err := questionService.GetQuestionDetail(
		context.Background(),
		model.GetQuestionDetailReq{QuestionID: "question-1"},
		"user-1",
	)
	if err == nil || !strings.Contains(err.Error(), "get question tag state") {
		t.Fatalf("GetQuestionDetail() error = %v", err)
	}
}

func TestQuestionServiceGetQuestionDetailPreservesNotFoundError(t *testing.T) {
	questionService := NewQuestionService(&questionDataStub{questionErr: data.ErrQuestionNotFound})

	_, err := questionService.GetQuestionDetail(
		context.Background(),
		model.GetQuestionDetailReq{QuestionID: "missing"},
		"",
	)
	if !errors.Is(err, ErrQuestionNotFound) {
		t.Fatalf("GetQuestionDetail() error = %v, want ErrQuestionNotFound", err)
	}
}

func TestQuestionServiceGetHotQuestions(t *testing.T) {
	stub := &questionDataStub{hotQuestions: []data.HotQuestionRecord{
		{QuestionID: "question-1", BankList: nil, Title: "Go 并发", ViewCount: 100},
		{QuestionID: "question-2", BankList: []string{"bank-2"}, Title: "MongoDB 索引", ViewCount: 80},
	}}
	questionService := NewQuestionService(stub)

	got, err := questionService.GetHotQuestions(context.Background(), model.GetHotQuestionsReq{
		Limit:   2,
		JobName: "后端开发",
	})
	if err != nil {
		t.Fatalf("GetHotQuestions() error = %v", err)
	}
	if stub.hotJobName != "后端开发" || stub.hotLimit != 2 {
		t.Fatalf("ListHotQuestions args = (%q, %d)", stub.hotJobName, stub.hotLimit)
	}
	want := model.GetHotQuestionsResp{
		Code: 0,
		Msg:  "success",
		Data: model.HotQuestionListData{List: []model.HotQuestionInfo{
			{QuestionID: "question-1", BankList: []string{}, Title: "Go 并发", ViewCount: 100},
			{QuestionID: "question-2", BankList: []string{"bank-2"}, Title: "MongoDB 索引", ViewCount: 80},
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetHotQuestions() = %#v, want %#v", got, want)
	}
}

func TestQuestionServiceGetHotQuestionsReturnsEmptyArray(t *testing.T) {
	questionService := NewQuestionService(&questionDataStub{})

	got, err := questionService.GetHotQuestions(context.Background(), model.GetHotQuestionsReq{Limit: 10})
	if err != nil {
		t.Fatalf("GetHotQuestions() error = %v", err)
	}
	if got.Data.List == nil || len(got.Data.List) != 0 {
		t.Fatalf("GetHotQuestions().Data.List = %#v, want non-nil empty slice", got.Data.List)
	}
}

func TestQuestionServiceGetAllQuestionList(t *testing.T) {
	stub := &questionDataStub{
		series: []data.QuestionBankSeriesRecord{
			{SeriesID: "backend-basic", SeriesName: "基础", JobName: "后端开发", Order: 1},
			{SeriesID: "frontend-basic", SeriesName: "基础", JobName: "前端开发", Order: 1},
			{SeriesID: "backend-db", SeriesName: "数据库", JobName: "后端开发", Order: 2},
		},
		banks: []data.QuestionBankRecord{
			{BankID: "go", SeriesID: "backend-basic", BankName: "Go", JobName: "后端开发", Order: 1},
			{BankID: "js", SeriesID: "frontend-basic", BankName: "JavaScript", JobName: "前端开发", Order: 1},
			{BankID: "mysql", SeriesID: "backend-db", BankName: "MySQL", JobName: "后端开发", Order: 2},
		},
		counts: map[string]int64{"go": 12, "mysql": 8},
	}
	questionService := NewQuestionService(stub)

	got, err := questionService.GetAllQuestionList(context.Background(), model.GetAllQuestionListReq{
		JobName: "后端开发",
	})
	if err != nil {
		t.Fatalf("GetAllQuestionList() error = %v", err)
	}

	want := []model.GetAllQuestionListData{
		{
			JobName: "后端开发",
			SeriesList: []model.QuestionBankSeries{
				{SeriesID: "backend-basic", SeriesName: "基础", Order: 1, BankList: []model.QuestionBank{
					{BankID: "go", BankName: "Go", Count: 12, Order: 1},
				}},
				{SeriesID: "backend-db", SeriesName: "数据库", Order: 2, BankList: []model.QuestionBank{
					{BankID: "mysql", BankName: "MySQL", Count: 8, Order: 2},
				}},
			},
		},
		{
			JobName: "前端开发",
			SeriesList: []model.QuestionBankSeries{
				{SeriesID: "frontend-basic", SeriesName: "基础", Order: 1, BankList: []model.QuestionBank{
					{BankID: "js", BankName: "JavaScript", Count: 0, Order: 1},
				}},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetAllQuestionList() = %#v, want %#v", got, want)
	}
	if stub.seriesFilter != "后端开发" || stub.bankFilter != "后端开发" {
		t.Fatalf("job_name filters = (%q, %q), want both %q", stub.seriesFilter, stub.bankFilter, "后端开发")
	}
	if !reflect.DeepEqual(stub.countBankIDs, []string{"go", "js", "mysql"}) {
		t.Fatalf("count bank IDs = %#v", stub.countBankIDs)
	}
}

func TestQuestionServiceGetAllQuestionListReturnsEmptyArray(t *testing.T) {
	stub := &questionDataStub{}
	questionService := NewQuestionService(stub)

	got, err := questionService.GetAllQuestionList(context.Background(), model.GetAllQuestionListReq{})
	if err != nil {
		t.Fatalf("GetAllQuestionList() error = %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("GetAllQuestionList() = %#v, want non-nil empty slice", got)
	}
}
