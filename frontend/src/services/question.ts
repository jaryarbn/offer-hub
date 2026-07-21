import apiClient from "@/lib/axios"
import type { ApiResponse } from "@/types/api"
import type {
  GetHotQuestionListParams,
  GetHotQuestionListResponse,
  GetHotQuestionListResponseData,
  GetQuestionBankSeriesParams,
  GetQuestionBankSeriesResponse,
  GetQuestionDetailParams,
  GetQuestionDetailResponse,
  ListQuestionMetaParams,
  ListQuestionMetaResponse,
  ListQuestionMetaResponseData,
  ListQuestionParams,
  ListQuestionResponse,
  ListQuestionResponseData,
  Question,
  QuestionBankGroup,
} from "@/types/question"
import { compactParams } from "@/utils/query"

const questionBasePath = "/api/v1/question"

/** 题目接口抛出的统一错误，code 对应后端业务响应 code。 */
export class QuestionApiError extends Error {
  readonly code: number

  constructor(code: number, message: string, options?: ErrorOptions) {
    super(message, options)
    this.name = "QuestionApiError"
    this.code = code
  }
}

export interface QuestionApiService {
  getAllQuestionBanks(
    params?: GetQuestionBankSeriesParams,
  ): Promise<QuestionBankGroup[]>
  getQuestionList(
    params?: ListQuestionParams,
  ): Promise<ListQuestionResponseData>
  getQuestionMetaList(
    params?: ListQuestionMetaParams,
  ): Promise<ListQuestionMetaResponseData>
  getQuestionDetail(params: GetQuestionDetailParams): Promise<Question>
  getHotQuestionList(
    params?: GetHotQuestionListParams,
  ): Promise<GetHotQuestionListResponseData>
}

function unwrapResponse<T>(
  response: ApiResponse<T> | null | undefined,
  fallbackMessage: string,
): T {
  if (
    !response ||
    typeof response.code !== "number" ||
    typeof response.msg !== "string" ||
    !("data" in response)
  ) {
    throw new QuestionApiError(-1, `${fallbackMessage}：响应格式错误`)
  }

  if (response.code !== 0) {
    throw new QuestionApiError(
      response.code,
      response.msg.trim() || fallbackMessage,
    )
  }

  if (response.data === null || response.data === undefined) {
    throw new QuestionApiError(-1, `${fallbackMessage}：响应数据为空`)
  }

  return response.data
}

function normalizeRequestError(error: unknown, fallbackMessage: string): never {
  if (error instanceof QuestionApiError) {
    throw error
  }

  // axios 响应拦截器会优先抛出后端的 { code, msg, data } 错误体。
  if (typeof error === "object" && error !== null) {
    const candidate = error as { code?: unknown; msg?: unknown }
    if (typeof candidate.msg === "string" && candidate.msg.trim()) {
      throw new QuestionApiError(
        typeof candidate.code === "number" ? candidate.code : -1,
        candidate.msg,
      )
    }
  }

  if (error instanceof Error) {
    throw new QuestionApiError(
      -1,
      `${fallbackMessage}：${error.message}`,
      { cause: error },
    )
  }

  throw new QuestionApiError(-1, fallbackMessage)
}

async function requestData<T>(
  request: () => Promise<ApiResponse<T>>,
  fallbackMessage: string,
): Promise<T> {
  try {
    return unwrapResponse(await request(), fallbackMessage)
  } catch (error) {
    normalizeRequestError(error, fallbackMessage)
  }
}

/** 题目模块 API 服务。所有方法成功时返回 data，失败时抛出 QuestionApiError。 */
export const questionApi: QuestionApiService = {
  /**
   * 获取按职位方向、系列分组的完整题库列表。
   * 题目数量只统计 status=1 的正常题目。
   *
   * @param params 可选的职位方向过滤条件
   * @returns 按 job_name 分组的题库系列列表
   * @throws {QuestionApiError} 请求失败、业务 code 非 0 或响应数据异常
   */
  async getAllQuestionBanks(params = {}) {
    return requestData<QuestionBankGroup[]>(async () => {
      const response = await apiClient.get<GetQuestionBankSeriesResponse>(
        `${questionBasePath}/all/list`,
        { params: compactParams(params) },
      )
      return response.data
    }, "获取题库列表失败")
  },

  /**
   * 获取题目列表，支持题库、关键词、难度、标签、排序和分页过滤。
   *
   * @param params 题目过滤、排序和分页参数
   * @returns 题目总数和当前页题目列表
   * @throws {QuestionApiError} 请求失败、业务 code 非 0 或响应数据异常
   */
  async getQuestionList(params = {}) {
    return requestData<ListQuestionResponseData>(async () => {
      const response = await apiClient.get<ListQuestionResponse>(
        `${questionBasePath}/list`,
        {
          params: compactParams(params),
          // Gin 使用 form:"tags"，数组必须编码为 tags=a&tags=b。
          paramsSerializer: { indexes: null },
        },
      )
      return response.data
    }, "获取题目列表失败")
  },

  /**
   * 获取题目导航元数据，列表项仅包含 question_id 和 title。
   *
   * @param params 与题目列表相同的过滤、排序和分页参数
   * @returns 题目总数和当前页元数据列表
   * @throws {QuestionApiError} 请求失败、业务 code 非 0 或响应数据异常
   */
  async getQuestionMetaList(params = {}) {
    return requestData<ListQuestionMetaResponseData>(async () => {
      const response = await apiClient.get<ListQuestionMetaResponse>(
        `${questionBasePath}/meta/list`,
        {
          params: compactParams(params),
          paramsSerializer: { indexes: null },
        },
      )
      return response.data
    }, "获取题目元数据失败")
  },

  /**
   * 获取一道正常题目的完整详情。
   *
   * @param params 包含必填 question_id 的查询参数
   * @returns 完整题目信息
   * @throws {QuestionApiError} question_id 为空、题目不存在或响应异常
   */
  async getQuestionDetail(params) {
    const questionId = params.question_id.trim()
    if (!questionId) {
      throw new QuestionApiError(400, "question_id 不能为空")
    }

    return requestData<Question>(async () => {
      const response = await apiClient.get<GetQuestionDetailResponse>(
        `${questionBasePath}/detail`,
        { params: { question_id: questionId } },
      )
      return response.data
    }, "获取题目详情失败")
  },

  /**
   * 获取热门题目列表，后端按 hot_degree 降序排列。
   *
   * @param params 可选的职位方向和数量限制，limit 默认值为 10
   * @returns 热门题目列表
   * @throws {QuestionApiError} 请求失败、业务 code 非 0 或响应数据异常
   */
  async getHotQuestionList(params = {}) {
    return requestData<GetHotQuestionListResponseData>(async () => {
      const response = await apiClient.get<GetHotQuestionListResponse>(
        `${questionBasePath}/hot/list`,
        { params: compactParams(params) },
      )
      return response.data
    }, "获取热门题目失败")
  },
}
