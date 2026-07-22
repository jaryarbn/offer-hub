package model

type ListQuestionReq struct {
	BankID     string   `json:"bank_id" form:"bank_id"`
	Keyword    string   `json:"keyword" form:"keyword"`
	Difficulty int      `json:"difficulty" form:"difficulty"`
	Tags       []string `json:"tags" form:"tags"`
	JobName    string   `json:"job_name" form:"job_name"`
	UserTag    int      `json:"user_tag" form:"user_tag"`
	SortBy     string   `json:"sort_by" form:"sort_by"`
	SortOrder  string   `json:"sort_order" form:"sort_order"`
	Page       int      `json:"page" form:"page"`
	PageSize   int      `json:"page_size" form:"page_size"`
}

type OneQuestion struct {
	QuestionID    string   `json:"question_id"`
	BankList      []string `json:"bank_list"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Difficulty    int      `json:"difficulty"`
	Tags          []string `json:"tags"`
	Status        int      `json:"status"`
	VIP           bool     `json:"vip"`
	HotDegree     int      `json:"hot_degree"`
	ViewCount     int      `json:"view_count"`
	ThumbsUpCount int      `json:"thumbs_up_count"`
	DislikeCount  int      `json:"dislike_count"`
	Order         int64    `json:"order"`
	UserTag       int      `json:"user_tag"`
	UserLiked     bool     `json:"user_liked"`
	CreateTime    string   `json:"create_time"`
	UpdateTime    string   `json:"update_time"`
}

type ListQuestionResp struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data ListQuestionResponseData `json:"data"`
}

type ListQuestionResponseData struct {
	Total int64         `json:"total"`
	List  []OneQuestion `json:"list"`
}

type ListQuestionMetaReq = ListQuestionReq

type ListQuestionMetaResp struct {
	Code int                          `json:"code"`
	Msg  string                       `json:"msg"`
	Data ListQuestionMetaResponseData `json:"data"`
}

type ListQuestionMetaResponseData struct {
	Total int64              `json:"total"`
	List  []QuestionMetaInfo `json:"list"`
}

type QuestionMetaInfo struct {
	QuestionID string `json:"question_id"`
	Title      string `json:"title"`
}

type GetQuestionDetailReq struct {
	QuestionID string `json:"question_id" form:"question_id"`
}

type GetQuestionDetailResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data QuestionDetail `json:"data"`
}

type QuestionDetail struct {
	OneQuestion
	AnalysisContent string `json:"analysis_content"`
}

type HotQuestionInfo struct {
	QuestionID string   `json:"question_id"`
	BankList   []string `json:"bank_list"`
	Title      string   `json:"title"`
	ViewCount  int      `json:"view_count"`
}

type GetHotQuestionsReq struct {
	Limit   int    `json:"limit" form:"limit"`
	JobName string `json:"job_name" form:"job_name"`
}

type GetHotQuestionsResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data HotQuestionListData `json:"data"`
}

type GetHotQuestionListResp = GetHotQuestionsResp

type HotQuestionListData struct {
	List []HotQuestionInfo `json:"list"`
}

type GetAllQuestionListReq struct {
	JobName string `json:"job_name" form:"job_name"`
}

type GetAllQuestionListResp struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data []GetAllQuestionListData `json:"data"`
}

type GetAllQuestionListData struct {
	JobName    string               `json:"job_name"`
	SeriesList []QuestionBankSeries `json:"series_list"`
}

type GetQuestionBankSeriesData = GetAllQuestionListData

type QuestionBankSeries struct {
	SeriesID   string         `json:"series_id"`
	SeriesName string         `json:"series_name"`
	Order      int64          `json:"order"`
	BankList   []QuestionBank `json:"bank_list"`
}

type QuestionBank struct {
	BankID   string `json:"bank_id"`
	BankName string `json:"bank_name"`
	BankLogo string `json:"bank_logo"`
	Desc     string `json:"desc"`
	Count    int64  `json:"count"`
	Order    int64  `json:"order"`
}
