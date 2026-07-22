import { LoaderCircle, MessageSquareText, RotateCcw } from 'lucide-react'

import { CommentItem } from '@/components/comment/CommentItem'
import type { CommentTargetType, CommentThreadStateProps } from '@/components/comment/types'
import { Button } from '@/components/ui/button'
import { useComments } from '@/hooks/useCommentQueries'

export interface CommentListProps extends CommentThreadStateProps {
  targetType: CommentTargetType
  targetId: string
}

export function CommentList({ targetType, targetId, ...threadState }: CommentListProps) {
  const commentsQuery = useComments({
    target_type: targetType,
    target_id: targetId,
    sort_by: 'create_time',
    sort_order: 'desc',
    page_size: 20,
    sub_comment_size: 5,
  })

  if (commentsQuery.isPending) {
    return <CommentListSkeleton />
  }

  if (commentsQuery.isError) {
    return (
      <div className="rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-6 text-center">
        <p className="text-sm text-destructive" role="alert">
          {commentsQuery.error.message || '评论加载失败'}
        </p>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="mt-3"
          onClick={() => void commentsQuery.refetch()}
        >
          <RotateCcw aria-hidden="true" />
          重新加载
        </Button>
      </div>
    )
  }

  const comments = commentsQuery.data.pages.flatMap(page => page.list)

  if (comments.length === 0) {
    return (
      <div className="rounded-lg border border-dashed border-border py-10 text-center">
        <MessageSquareText className="mx-auto size-8 text-muted-foreground" aria-hidden="true" />
        <p className="mt-3 text-sm font-medium">还没有评论</p>
        <p className="mt-1 text-xs text-muted-foreground">来分享第一条观点吧</p>
      </div>
    )
  }

  return (
    <div>
      <ul className="divide-y divide-border" role="list">
        {comments.map(comment => (
          <li key={comment.comment_id}>
            <CommentItem
              comment={comment}
              targetType={targetType}
              targetId={targetId}
              {...threadState}
            />
          </li>
        ))}
      </ul>

      {commentsQuery.hasNextPage ? (
        <div className="mt-5 flex justify-center">
          <Button
            type="button"
            variant="outline"
            disabled={commentsQuery.isFetchingNextPage}
            onClick={() => void commentsQuery.fetchNextPage()}
          >
            {commentsQuery.isFetchingNextPage ? (
              <LoaderCircle className="animate-spin" aria-hidden="true" />
            ) : null}
            {commentsQuery.isFetchingNextPage ? '加载中…' : '加载更多'}
          </Button>
        </div>
      ) : null}
    </div>
  )
}

function CommentListSkeleton() {
  return (
    <div className="space-y-5" aria-busy="true" aria-label="正在加载评论">
      {Array.from({ length: 3 }, (_, index) => (
        <div key={index} className="flex gap-3 py-3">
          <div className="size-9 shrink-0 animate-pulse rounded-full bg-muted" />
          <div className="flex-1 space-y-3">
            <div className="h-4 w-28 animate-pulse rounded bg-muted" />
            <div className="h-4 w-full animate-pulse rounded bg-muted" />
            <div className="h-4 w-2/3 animate-pulse rounded bg-muted" />
          </div>
        </div>
      ))}
    </div>
  )
}
