import { useQuery } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'

import { QuestionApiError, questionApi } from '@/services/question'
import type {
  GetHotQuestionListResponseData,
  ListQuestionResponseData,
  Question,
  QuestionBankGroup,
  QuestionFilters,
} from '@/types/question'

/**
 * 题目模块的查询键工厂。
 *
 * 页面刷新、手动失效缓存或预取数据时都应复用这里的 key，避免同一份
 * 服务端数据因为 key 写法不同而生成多份缓存。
 */
export const queryKeys = {
  all: ['questions'] as const,
  banks: (jobName?: string) => [...queryKeys.all, 'banks', jobName ?? 'all'] as const,
  lists: ['questions', 'list'] as const,
  list: (filters: QuestionFilters) => [...queryKeys.lists, filters] as const,
  details: ['questions', 'detail'] as const,
  detail: (questionId: string) => [...queryKeys.details, questionId] as const,
  hot: (limit: number, jobName?: string) =>
    [...queryKeys.all, 'hot', limit, jobName ?? 'all'] as const,
}

/** 获取按职位方向和系列分组的题库列表。 */
export function useQuestionBanks(
  jobName?: string
): UseQueryResult<QuestionBankGroup[], QuestionApiError> {
  const normalizedJobName = jobName?.trim() || undefined

  return useQuery<QuestionBankGroup[], QuestionApiError>({
    queryKey: queryKeys.banks(normalizedJobName),
    queryFn: () =>
      questionApi.getAllQuestionBanks(
        normalizedJobName ? { job_name: normalizedJobName } : undefined
      ),
  })
}

/**
 * 按筛选条件获取分页题目列表。
 * enabled 可用于等待页面依赖的路由参数或表单条件准备完成。
 */
export function useQuestionList(
  filters: QuestionFilters,
  enabled = true
): UseQueryResult<ListQuestionResponseData, QuestionApiError> {
  return useQuery<ListQuestionResponseData, QuestionApiError>({
    queryKey: queryKeys.list(filters),
    queryFn: () => questionApi.getQuestionList(filters),
    enabled,
  })
}

/**
 * 获取单道题目的完整详情。
 * questionId 为空时即使 enabled=true 也不会发送无效请求。
 */
export function useQuestionDetail(
  questionId: string,
  enabled = true
): UseQueryResult<Question, QuestionApiError> {
  const normalizedQuestionId = questionId.trim()

  return useQuery<Question, QuestionApiError>({
    queryKey: queryKeys.detail(normalizedQuestionId),
    queryFn: () => questionApi.getQuestionDetail({ question_id: normalizedQuestionId }),
    enabled: enabled && normalizedQuestionId.length > 0,
  })
}

/**
 * 获取指定题库的题目列表，可同时按关键词过滤。
 * 复用列表接口和 list 查询键，确保与相同条件的列表页共享缓存。
 */
export function useBankQuestionList(
  bankId: string,
  keyword?: string,
  enabled = true
): UseQueryResult<ListQuestionResponseData, QuestionApiError> {
  const normalizedBankId = bankId.trim()
  const normalizedKeyword = keyword?.trim() || undefined
  const filters: QuestionFilters = {
    bank_id: normalizedBankId,
    keyword: normalizedKeyword,
  }

  return useQuery<ListQuestionResponseData, QuestionApiError>({
    queryKey: queryKeys.list(filters),
    queryFn: () => questionApi.getQuestionList(filters),
    enabled: enabled && normalizedBankId.length > 0,
  })
}

/** 获取后端按 hot_degree 降序返回的热门题目。 */
export function useHotQuestions(
  limit = 10,
  jobName?: string
): UseQueryResult<GetHotQuestionListResponseData, QuestionApiError> {
  const normalizedLimit = Number.isInteger(limit) && limit > 0 ? limit : 10
  const normalizedJobName = jobName?.trim() || undefined

  return useQuery<GetHotQuestionListResponseData, QuestionApiError>({
    queryKey: queryKeys.hot(normalizedLimit, normalizedJobName),
    queryFn: () =>
      questionApi.getHotQuestionList({
        limit: normalizedLimit,
        job_name: normalizedJobName,
      }),
  })
}
