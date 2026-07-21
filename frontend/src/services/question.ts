import apiClient from "@/lib/axios"
import type { ApiResponse } from "@/types/api"
import type {
  ListQuestionParams,
  ListQuestionResponseData,
} from "@/types/question"
import { compactParams } from "@/utils/query"

export async function listQuestions(params: ListQuestionParams) {
  const response = await apiClient.get<ApiResponse<ListQuestionResponseData>>(
    "/api/v1/question/list",
    { params: compactParams(params) },
  )

  if (response.data.code !== 0) {
    throw new Error(response.data.msg)
  }

  return response.data.data
}
