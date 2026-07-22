package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

type QuestionData interface {
	ListQuestionBankSeries(context.Context, string) ([]data.QuestionBankSeriesRecord, error)
	ListQuestionBanks(context.Context, string) ([]data.QuestionBankRecord, error)
	CountNormalQuestionsByBank(context.Context, []string) (map[string]int64, error)
	FilterQuestions(context.Context, data.QuestionFilter) ([]data.QuestionRecord, int64, error)
	GetQuestionByID(context.Context, string) (data.QuestionRecord, error)
	IsQuestionLiked(context.Context, string, string) (bool, error)
	GetUserQuestionTag(context.Context, string, string) (int, error)
	ListHotQuestions(context.Context, string, int) ([]data.HotQuestionRecord, error)
}

var ErrQuestionNotFound = data.ErrQuestionNotFound

type QuestionService struct {
	data QuestionData
}

func NewQuestionService(questionData QuestionData) *QuestionService {
	return &QuestionService{data: questionData}
}

const (
	visitorContentLength = 150
	timeFormat           = "2006-01-02 15:04:05"
)

func (service *QuestionService) ListQuestions(
	ctx context.Context,
	req model.ListQuestionReq,
) (model.ListQuestionResponseData, error) {
	records, total, err := service.data.FilterQuestions(ctx, toQuestionFilter(req))
	if err != nil {
		return model.ListQuestionResponseData{}, fmt.Errorf("filter questions: %w", err)
	}

	questions := make([]model.OneQuestion, 0, len(records))
	for _, record := range records {
		questions = append(questions, toOneQuestion(
			record,
			truncateRunes(record.Content, visitorContentLength),
		))
	}

	return model.ListQuestionResponseData{Total: total, List: questions}, nil
}

func (service *QuestionService) ListQuestionsMeta(
	ctx context.Context,
	req model.ListQuestionMetaReq,
) (model.ListQuestionMetaResp, error) {
	records, total, err := service.data.FilterQuestions(ctx, toQuestionFilter(req))
	if err != nil {
		return model.ListQuestionMetaResp{}, fmt.Errorf("filter question metadata: %w", err)
	}

	questions := make([]model.QuestionMetaInfo, 0, len(records))
	for _, record := range records {
		questions = append(questions, model.QuestionMetaInfo{
			QuestionID: record.QuestionID,
			Title:      record.Title,
		})
	}

	return model.ListQuestionMetaResp{
		Code: 0,
		Msg:  "success",
		Data: model.ListQuestionMetaResponseData{
			Total: total,
			List:  questions,
		},
	}, nil
}

func (service *QuestionService) GetQuestionDetail(
	ctx context.Context,
	req model.GetQuestionDetailReq,
	userID string,
) (model.QuestionDetail, error) {
	record, err := service.data.GetQuestionByID(ctx, req.QuestionID)
	if err != nil {
		return model.QuestionDetail{}, fmt.Errorf("get question detail: %w", err)
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return model.QuestionDetail{
			OneQuestion: toOneQuestion(record, truncateRunes(record.Content, visitorContentLength)),
		}, nil
	}

	question := model.QuestionDetail{
		OneQuestion:     toOneQuestion(record, record.Content),
		AnalysisContent: record.AnalysisContent,
	}

	question.UserLiked, err = service.data.IsQuestionLiked(ctx, userID, record.QuestionID)
	if err != nil {
		return model.QuestionDetail{}, fmt.Errorf("get question like state: %w", err)
	}
	question.UserTag, err = service.data.GetUserQuestionTag(ctx, userID, record.QuestionID)
	if err != nil {
		return model.QuestionDetail{}, fmt.Errorf("get question tag state: %w", err)
	}
	return question, nil
}

func (service *QuestionService) GetHotQuestions(
	ctx context.Context,
	req model.GetHotQuestionsReq,
) (model.GetHotQuestionsResp, error) {
	records, err := service.data.ListHotQuestions(ctx, req.JobName, req.Limit)
	if err != nil {
		return model.GetHotQuestionsResp{}, fmt.Errorf("get hot questions: %w", err)
	}

	questions := make([]model.HotQuestionInfo, 0, len(records))
	for _, record := range records {
		questions = append(questions, model.HotQuestionInfo{
			QuestionID: record.QuestionID,
			BankList:   nonNilStrings(record.BankList),
			Title:      record.Title,
			ViewCount:  record.ViewCount,
		})
	}

	return model.GetHotQuestionsResp{
		Code: 0,
		Msg:  "success",
		Data: model.HotQuestionListData{List: questions},
	}, nil
}

func toOneQuestion(record data.QuestionRecord, content string) model.OneQuestion {
	return model.OneQuestion{
		QuestionID:    record.QuestionID,
		BankList:      nonNilStrings(record.BankList),
		Title:         record.Title,
		Content:       content,
		Difficulty:    record.Difficulty,
		Tags:          nonNilStrings(record.Tags),
		Status:        record.Status,
		VIP:           record.VIP,
		HotDegree:     record.HotDegree,
		ViewCount:     record.ViewCount,
		ThumbsUpCount: record.ThumbsUpCount,
		DislikeCount:  record.DislikeCount,
		Order:         record.Order,
		UserTag:       0,
		UserLiked:     false,
		CreateTime:    formatTime(record.CreateTime),
		UpdateTime:    formatTime(record.UpdateTime),
	}
}

func toQuestionFilter(req model.ListQuestionReq) data.QuestionFilter {
	return data.QuestionFilter{
		BankID:     req.BankID,
		Keyword:    req.Keyword,
		Difficulty: req.Difficulty,
		Tags:       req.Tags,
		JobName:    req.JobName,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}
}

func truncateRunes(value string, maxLength int) string {
	runes := []rune(value)
	if len(runes) <= maxLength {
		return value
	}
	return string(runes[:maxLength])
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return make([]string, 0)
	}
	return values
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(timeFormat)
}

func (service *QuestionService) GetAllQuestionList(
	ctx context.Context,
	req model.GetAllQuestionListReq,
) ([]model.GetAllQuestionListData, error) {
	seriesRecords, err := service.data.ListQuestionBankSeries(ctx, req.JobName)
	if err != nil {
		return nil, fmt.Errorf("list question bank series: %w", err)
	}
	if len(seriesRecords) == 0 {
		return make([]model.GetAllQuestionListData, 0), nil
	}

	bankRecords, err := service.data.ListQuestionBanks(ctx, req.JobName)
	if err != nil {
		return nil, fmt.Errorf("list question banks: %w", err)
	}

	bankIDs := make([]string, 0, len(bankRecords))
	for _, bank := range bankRecords {
		bankIDs = append(bankIDs, bank.BankID)
	}
	counts, err := service.data.CountNormalQuestionsByBank(ctx, bankIDs)
	if err != nil {
		return nil, fmt.Errorf("count normal questions by bank: %w", err)
	}

	result := make([]model.GetAllQuestionListData, 0)
	jobIndexes := make(map[string]int)
	seriesIndexes := make(map[string]seriesPosition)
	for _, record := range seriesRecords {
		jobIndex, exists := jobIndexes[record.JobName]
		if !exists {
			jobIndex = len(result)
			jobIndexes[record.JobName] = jobIndex
			result = append(result, model.GetAllQuestionListData{
				JobName:    record.JobName,
				SeriesList: make([]model.QuestionBankSeries, 0),
			})
		}

		seriesIndex := len(result[jobIndex].SeriesList)
		result[jobIndex].SeriesList = append(result[jobIndex].SeriesList, model.QuestionBankSeries{
			SeriesID:   record.SeriesID,
			SeriesName: record.SeriesName,
			Order:      record.Order,
			BankList:   make([]model.QuestionBank, 0),
		})
		seriesIndexes[seriesKey(record.JobName, record.SeriesID)] = seriesPosition{
			job:    jobIndex,
			series: seriesIndex,
		}
	}

	for _, record := range bankRecords {
		position, exists := seriesIndexes[seriesKey(record.JobName, record.SeriesID)]
		if !exists {
			continue
		}
		series := &result[position.job].SeriesList[position.series]
		series.BankList = append(series.BankList, model.QuestionBank{
			BankID:   record.BankID,
			BankName: record.BankName,
			BankLogo: record.BankLogo,
			Desc:     record.Desc,
			Count:    counts[record.BankID],
			Order:    record.Order,
		})
	}

	return result, nil
}

type seriesPosition struct {
	job    int
	series int
}

func seriesKey(jobName, seriesID string) string {
	return jobName + "\x00" + seriesID
}
