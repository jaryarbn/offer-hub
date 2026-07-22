import { useEffect, useState } from 'react'
import { LoaderCircle, Pencil, Reply, ThumbsUp, Trash2 } from 'lucide-react'

import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'
import { useDeleteComment } from '@/hooks/useCommentQueries'
import { thumbsUpComment } from '@/services/comment'
import type { CommentInfo } from '@/types/comment'

export interface CommentActionsProps {
  comment: CommentInfo
  isOwnComment: boolean
  onReply: () => void
  onEdit: () => void
}

/** 单条评论操作栏，点赞成功后再更新本地计数。 */
export function CommentActions({ comment, isOwnComment, onReply, onEdit }: CommentActionsProps) {
  const { isLoggedIn, setShowLoginDialog } = useLogin()
  const deleteComment = useDeleteComment()
  const [confirmingDelete, setConfirmingDelete] = useState(false)
  const [likeCount, setLikeCount] = useState(comment.thumbs_up)
  const [isLiked, setIsLiked] = useState(comment.user_liked)
  const [isLiking, setIsLiking] = useState(false)
  const [likeError, setLikeError] = useState<string | null>(null)

  useEffect(() => {
    setLikeCount(comment.thumbs_up)
    setIsLiked(comment.user_liked)
  }, [comment.thumbs_up, comment.user_liked])

  const handleReply = () => {
    if (!isLoggedIn) {
      setShowLoginDialog(true)
      return
    }
    onReply()
  }

  const handleLike = async () => {
    if (!isLoggedIn) {
      setShowLoginDialog(true)
      return
    }
    if (isLiking) {
      return
    }

    setIsLiking(true)
    setLikeError(null)
    try {
      const result = await thumbsUpComment({ comment_id: comment.comment_id })
      setLikeCount(result.count)
      setIsLiked(result.liked)
    } catch (error) {
      setLikeError(error instanceof Error ? error.message : '点赞失败')
    } finally {
      setIsLiking(false)
    }
  }

  const handleDelete = async () => {
    if (!confirmingDelete) {
      setConfirmingDelete(true)
      return
    }
    try {
      await deleteComment.mutateAsync({ comment_id: comment.comment_id })
      setConfirmingDelete(false)
    } catch {
      // React Query 已将错误写入 mutation.error，由操作栏下方统一展示。
    }
  }

  return (
    <>
      <div className="mt-3 flex flex-wrap items-center gap-1">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className={isLiked ? 'text-primary' : 'text-muted-foreground'}
          disabled={isLiking}
          aria-pressed={isLiked}
          onClick={() => void handleLike()}
        >
          {isLiking ? (
            <LoaderCircle className="animate-spin" aria-hidden="true" />
          ) : (
            <ThumbsUp aria-hidden="true" />
          )}
          {likeCount}
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="text-muted-foreground"
          onClick={handleReply}
        >
          <Reply aria-hidden="true" />
          回复
        </Button>

        {isOwnComment ? (
          <>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="text-muted-foreground"
              onClick={onEdit}
            >
              <Pencil aria-hidden="true" />
              编辑
            </Button>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="text-destructive"
              disabled={deleteComment.isPending}
              onClick={() => void handleDelete()}
            >
              {deleteComment.isPending ? (
                <LoaderCircle className="animate-spin" aria-hidden="true" />
              ) : (
                <Trash2 aria-hidden="true" />
              )}
              {confirmingDelete ? '确认删除' : '删除'}
            </Button>
            {confirmingDelete ? (
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => setConfirmingDelete(false)}
              >
                取消
              </Button>
            ) : null}
          </>
        ) : null}
      </div>

      {likeError ? (
        <p className="mt-2 text-xs text-destructive" role="alert">
          {likeError}
        </p>
      ) : null}
      {deleteComment.error ? (
        <p className="mt-2 text-xs text-destructive" role="alert">
          {deleteComment.error.message}
        </p>
      ) : null}
    </>
  )
}
