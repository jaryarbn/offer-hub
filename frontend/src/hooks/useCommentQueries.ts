import { useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { InfiniteData, UseInfiniteQueryResult, UseMutationResult } from '@tanstack/react-query'

import { CommentApiError, commentApi } from '@/services/comment'
import type {
  AddCommentParams,
  AddCommentResponseData,
  DeleteCommentParams,
  DeleteCommentResponse,
  ListCommentsParams,
  ListCommentsResponseData,
  UpdateCommentParams,
  UpdateCommentResponseData,
} from '@/types/comment'

/** 评论模块查询键工厂，失效缓存、预取和手动刷新时应统一复用。 */
export const queryKeys = {
  all: ['comments'] as const,
  lists: ['comments', 'list'] as const,
  list: (params: ListCommentsParams) => [...queryKeys.lists, params] as const,
}

/** 获取指定题目或面经的评论列表。 */
export function useComments(
  params: ListCommentsParams,
  enabled = true
): UseInfiniteQueryResult<InfiniteData<ListCommentsResponseData, number>, CommentApiError> {
  const normalizedParams: ListCommentsParams = {
    ...params,
    target_id: params.target_id.trim(),
    parent_id: params.parent_id?.trim() || undefined,
  }

  const initialPage = normalizedParams.page ?? 1
  const pageSize = normalizedParams.page_size ?? 20

  return useInfiniteQuery<
    ListCommentsResponseData,
    CommentApiError,
    InfiniteData<ListCommentsResponseData, number>,
    ReturnType<typeof queryKeys.list>,
    number
  >({
    queryKey: queryKeys.list(normalizedParams),
    queryFn: ({ pageParam }) =>
      commentApi.getCommentList({ ...normalizedParams, page: pageParam, page_size: pageSize }),
    initialPageParam: initialPage,
    getNextPageParam: (lastPage, pages, lastPageParam) => {
      const loadedCount = pages.reduce((total, page) => total + page.list.length, 0)
      return loadedCount < lastPage.total ? lastPageParam + 1 : undefined
    },
    enabled: enabled && normalizedParams.target_id.length > 0,
  })
}

/** 添加顶层评论或回复，成功后刷新所有评论列表缓存。 */
export function useAddComment(): UseMutationResult<
  AddCommentResponseData,
  CommentApiError,
  AddCommentParams
> {
  const queryClient = useQueryClient()

  return useMutation<AddCommentResponseData, CommentApiError, AddCommentParams>({
    mutationFn: params => commentApi.postCommentAdd(params),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.lists })
    },
  })
}

/** 编辑当前用户自己的评论，成功后刷新所有评论列表缓存。 */
export function useUpdateComment(): UseMutationResult<
  UpdateCommentResponseData,
  CommentApiError,
  UpdateCommentParams
> {
  const queryClient = useQueryClient()

  return useMutation<UpdateCommentResponseData, CommentApiError, UpdateCommentParams>({
    mutationFn: params => commentApi.postCommentUpdate(params),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.lists })
    },
  })
}

/** 删除当前用户自己的评论，成功后刷新所有评论列表缓存。 */
export function useDeleteComment(): UseMutationResult<
  DeleteCommentResponse,
  CommentApiError,
  DeleteCommentParams
> {
  const queryClient = useQueryClient()

  return useMutation<DeleteCommentResponse, CommentApiError, DeleteCommentParams>({
    mutationFn: params => commentApi.postCommentDelete(params),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.lists })
    },
  })
}
