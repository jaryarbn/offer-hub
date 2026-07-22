import apiClient from '@/lib/axios'
import type {
  AddCommentParams,
  AddCommentResponse,
  AddCommentResponseData,
  CommentApiService,
  DeleteCommentParams,
  DeleteCommentResponse,
  ListCommentsParams,
  ListCommentsResponse,
  ListCommentsResponseData,
  ReplyCommentParams,
  ThumbsUpCommentParams,
  ThumbsUpCommentResponse,
  ThumbsUpCommentResponseData,
  UpdateCommentParams,
  UpdateCommentResponse,
  UpdateCommentResponseData,
} from '@/types/comment'
import { compactParams } from '@/utils/query'

const openCommentPath = '/api/v1/open/list_comments'
const commentBasePath = '/api/v1/comment'

/** 评论接口业务错误，code 对应后端统一响应中的业务码。 */
export class CommentApiError extends Error {
  readonly code: number

  constructor(code: number, message: string) {
    super(message)
    this.name = 'CommentApiError'
    this.code = code
  }
}

function assertSuccess(response: { code: number; msg: string }, fallbackMessage: string): void {
  if (response.code !== 0) {
    throw new CommentApiError(response.code, response.msg.trim() || fallbackMessage)
  }
}

function requireData<T>(
  response: { code: number; msg: string; data: T | null },
  fallbackMessage: string
): T {
  assertSuccess(response, fallbackMessage)

  if (response.data === null || response.data === undefined) {
    throw new CommentApiError(-1, `${fallbackMessage}：响应数据为空`)
  }

  return response.data
}

/** 获取评论列表；parent_id 为空时返回顶层评论及其分页回复。 */
export async function getCommentList(
  params: ListCommentsParams
): Promise<ListCommentsResponseData> {
  const response = await apiClient.get<ListCommentsResponse>(openCommentPath, {
    params: compactParams(params),
  })

  return requireData(response.data, '获取评论列表失败')
}

/** 发表评论；传入 parent_id 时也可直接用于回复。 */
export async function postCommentAdd(params: AddCommentParams): Promise<AddCommentResponseData> {
  const response = await apiClient.post<AddCommentResponse>(`${commentBasePath}/add`, params)

  return requireData(response.data, '发表评论失败')
}

/** 回复评论，复用 POST /api/v1/comment/add，并在类型层要求 parent_id。 */
export async function postCommentReply(
  params: ReplyCommentParams
): Promise<AddCommentResponseData> {
  return postCommentAdd(params)
}

/** 删除当前用户自己的评论。 */
export async function postCommentDelete(
  params: DeleteCommentParams
): Promise<DeleteCommentResponse> {
  const response = await apiClient.post<DeleteCommentResponse>(`${commentBasePath}/delete`, params)

  assertSuccess(response.data, '删除评论失败')
  return response.data
}

/** 修改当前用户自己的评论。 */
export async function postCommentUpdate(
  params: UpdateCommentParams
): Promise<UpdateCommentResponseData> {
  const response = await apiClient.post<UpdateCommentResponse>(`${commentBasePath}/update`, params)

  return requireData(response.data, '修改评论失败')
}

/** 点赞评论；交互模块使用 target_type=3 表示评论。 */
export async function thumbsUpComment(
  params: ThumbsUpCommentParams
): Promise<ThumbsUpCommentResponseData> {
  const response = await apiClient.post<ThumbsUpCommentResponse>('/api/v1/interaction/like', {
    target_type: 3,
    target_id: params.comment_id,
  })

  return requireData(response.data, '点赞评论失败')
}

/** 评论模块 API 聚合对象。 */
export const commentApi: CommentApiService = {
  getCommentList,
  postCommentAdd,
  postCommentReply,
  postCommentDelete,
  postCommentUpdate,
  thumbsUpComment,
}
