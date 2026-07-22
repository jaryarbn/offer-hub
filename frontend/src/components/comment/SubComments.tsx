import { LoaderCircle, RotateCcw } from 'lucide-react'

import { CommentItem } from '@/components/comment/CommentItem'
import type { CommentTargetType, CommentThreadStateProps } from '@/components/comment/types'
import { Button } from '@/components/ui/button'
import { useComments } from '@/hooks/useCommentQueries'

export interface SubCommentsProps extends CommentThreadStateProps {
  targetType: CommentTargetType
  targetId: string
  parentId: string
}

/** 展开后按每页 5 条加载指定父评论的回复。 */
export function SubComments({ targetType, targetId, parentId, ...threadState }: SubCommentsProps) {
  const commentsQuery = useComments({
    target_type: targetType,
    target_id: targetId,
    parent_id: parentId,
    sort_by: 'create_time',
    sort_order: 'asc',
    page_size: 5,
  })

  if (commentsQuery.isPending) {
    return (
      <div className="mt-3 flex items-center gap-2 text-xs text-muted-foreground" role="status">
        <LoaderCircle className="size-3.5 animate-spin" aria-hidden="true" />
        正在加载回复…
      </div>
    )
  }

  if (commentsQuery.isError) {
    return (
      <div className="mt-3 flex items-center gap-2 text-sm text-destructive" role="alert">
        <span>回复加载失败</span>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => void commentsQuery.refetch()}
        >
          <RotateCcw aria-hidden="true" />
          重试
        </Button>
      </div>
    )
  }

  const comments = commentsQuery.data.pages.flatMap(page => page.list)

  return (
    <div className="mt-3 border-l-2 border-border pl-4">
      <ul className="divide-y divide-border" role="list">
        {comments.map(comment => (
          <li key={comment.comment_id}>
            <CommentItem
              comment={comment}
              targetType={targetType}
              targetId={targetId}
              isSub
              {...threadState}
            />
          </li>
        ))}
      </ul>

      {commentsQuery.hasNextPage ? (
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="mt-2"
          disabled={commentsQuery.isFetchingNextPage}
          onClick={() => void commentsQuery.fetchNextPage()}
        >
          {commentsQuery.isFetchingNextPage ? (
            <LoaderCircle className="animate-spin" aria-hidden="true" />
          ) : null}
          {commentsQuery.isFetchingNextPage ? '加载中…' : '查看更多回复'}
        </Button>
      ) : null}
    </div>
  )
}
