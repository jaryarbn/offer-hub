/** 评论目标类型；使用对象常量以兼容 erasableSyntaxOnly。 */
export const TargetType = {
  QUESTION: 1,
  INTERVIEW_EXPERIENCE: 2,
  COMMENT: 3,
} as const

/** GET /api/v1/open/list_comments 查询参数。 */
export interface ListCommentsParams {
  /** 1 题目，2 面经。 */
  target_type: 1 | 2
  target_id: string
  parent_id?: string
  sort_by?: 'create_time' | 'thumbs_up'
  sort_order?: 'asc' | 'desc'
  page?: number
  page_size?: number
  sub_comment_page?: number
  sub_comment_size?: number
}

/** 评论列表项；顶层评论通过 sub_comments 携带当前页回复。 */
export interface CommentInfo {
  comment_id: string
  user_id: string
  user_name: string
  user_avatar: string
  content: string
  parent_id: string
  reply_to: string
  reply_to_name: string
  /** 1 审核中，2 正常，3 拒绝，4 隐藏，5 删除。 */
  status: 1 | 2 | 3 | 4 | 5
  thumbs_up: number
  sub_comment_total: number
  user_liked: boolean
  sub_comments: CommentInfo[]
  create_time: string
  update_time: string
}

export interface ListCommentsResponseData {
  total: number
  list: CommentInfo[]
}

export interface ListCommentsResponse {
  code: number
  msg: string
  data: ListCommentsResponseData | null
}

/** POST /api/v1/comment/add 请求体，可用于顶层评论或回复。 */
export interface AddCommentParams {
  /** 1 题目，2 面经。 */
  target_type: 1 | 2
  target_id: string
  parent_id?: string
  reply_to?: string
  content: string
}

/** 回复评论复用 add 接口，但 parent_id 必填。 */
export interface ReplyCommentParams extends AddCommentParams {
  parent_id: string
}

export interface AddCommentResponseData {
  comment_id: string
  comment: CommentInfo
}

export interface AddCommentResponse {
  code: number
  msg: string
  data: AddCommentResponseData | null
}

/** POST /api/v1/comment/delete 请求体。 */
export interface DeleteCommentParams {
  comment_id: string
}

/** 删除成功时后端不返回 data；错误响应可能携带 data: null。 */
export interface DeleteCommentResponse {
  code: number
  msg: string
  data?: null
}

/** POST /api/v1/comment/update 请求体。 */
export interface UpdateCommentParams {
  comment_id: string
  content: string
}

export interface UpdateCommentResponseData {
  comment_id: string
}

export interface UpdateCommentResponse {
  code: number
  msg: string
  data: UpdateCommentResponseData | null
}

/** 评论点赞业务参数，service 会映射为 interaction 接口的 target_type=3。 */
export interface ThumbsUpCommentParams {
  comment_id: string
}

export interface ThumbsUpCommentResponseData {
  liked: boolean
  count: number
}

export interface ThumbsUpCommentResponse {
  code: number
  msg: string
  data: ThumbsUpCommentResponseData | null
}

/** 评论模块的前端业务方法；发表评论与回复共用同一后端路径。 */
export interface CommentApiService {
  getCommentList(params: ListCommentsParams): Promise<ListCommentsResponseData>
  postCommentAdd(params: AddCommentParams): Promise<AddCommentResponseData>
  postCommentReply(params: ReplyCommentParams): Promise<AddCommentResponseData>
  postCommentDelete(params: DeleteCommentParams): Promise<DeleteCommentResponse>
  postCommentUpdate(params: UpdateCommentParams): Promise<UpdateCommentResponseData>
  thumbsUpComment(params: ThumbsUpCommentParams): Promise<ThumbsUpCommentResponseData>
}
