import type { ApiResponse } from '@/types/api'

export type QuestionSortBy = 'create_time' | 'view_count' | 'thumbs_up_count' | 'dislike_count'

export type SortOrder = 'asc' | 'desc'

/** GET /api/v1/question/all/list 查询参数。 */
export interface GetQuestionBankSeriesParams {
  job_name?: string
}

export interface ListQuestionParams {
  bank_id?: string
  keyword?: string
  difficulty?: number
  tags?: string[]
  job_name?: string
  user_tag?: number
  sort_by?: QuestionSortBy
  sort_order?: SortOrder
  page?: number
  page_size?: number
}

/** 题目列表 Hook 使用的筛选条件，与列表接口查询参数保持一致。 */
export type QuestionFilters = ListQuestionParams

export type ListQuestionMetaParams = ListQuestionParams

/** GET /api/v1/question/detail 必填参数。 */
export interface GetQuestionDetailParams {
  question_id: string
}

export interface GetHotQuestionListParams {
  limit?: number
  job_name?: string
}

/** 题目列表项与详情共用的数据结构。 */
export interface Question {
  question_id: string
  bank_list: string[]
  title: string
  content: string
  difficulty: number
  tags: string[]
  status: number
  vip: boolean
  hot_degree: number
  view_count: number
  thumbs_up_count: number
  dislike_count: number
  order: number
  user_tag: number
  user_liked: boolean
  create_time: string
  update_time: string
}

/** 题目详情在列表字段基础上额外返回解析内容。 */
export interface QuestionDetail extends Question {
  analysis_content: string
}

/** 题库全量列表中的题库信息。 */
export interface QuestionBank {
  bank_id: string
  bank_name: string
  bank_logo: string
  desc: string
  count: number
  order: number
}

/** 一个题库系列及其所属题库。 */
export interface QuestionBankSeries {
  series_id: string
  series_name: string
  order: number
  bank_list: QuestionBank[]
}

/** 按职位方向聚合的题库系列。 */
export interface QuestionBankGroup {
  job_name: string
  series_list: QuestionBankSeries[]
}

/** 题目列表分页数据。 */
export interface ListQuestionResponseData {
  total: number
  list: Question[]
}

/** 题目导航使用的轻量元数据。 */
export interface QuestionMeta {
  question_id: string
  title: string
}

export interface ListQuestionMetaResponseData {
  total: number
  list: QuestionMeta[]
}

/** 热门题目接口只返回文档约定的四个字段。 */
export interface HotQuestion {
  question_id: string
  bank_list: string[]
  title: string
  view_count: number
}

export interface GetHotQuestionListResponseData {
  list: HotQuestion[]
}

export type GetQuestionBankSeriesResponse = ApiResponse<QuestionBankGroup[]>
export type ListQuestionResponse = ApiResponse<ListQuestionResponseData>
export type ListQuestionMetaResponse = ApiResponse<ListQuestionMetaResponseData>
export type GetQuestionDetailResponse = ApiResponse<QuestionDetail>
export type GetHotQuestionListResponse = ApiResponse<GetHotQuestionListResponseData>
