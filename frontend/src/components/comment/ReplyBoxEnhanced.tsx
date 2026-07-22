import { useState, type FormEvent } from 'react'
import { LoaderCircle, Send } from 'lucide-react'

import type { CommentTargetType } from '@/components/comment/types'
import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'
import { useAddComment } from '@/hooks/useCommentQueries'

export interface ReplyBoxEnhancedProps {
  targetType: CommentTargetType
  targetId: string
  parentId: string
  replyTo: string
  replyToName: string
  onSubmitted: () => void
  onCancel: () => void
}

/** 回复输入框，提交时携带父评论 ID 和被回复用户 ID。 */
export function ReplyBoxEnhanced({
  targetType,
  targetId,
  parentId,
  replyTo,
  replyToName,
  onSubmitted,
  onCancel,
}: ReplyBoxEnhancedProps) {
  const { isLoggedIn, setShowLoginDialog } = useLogin()
  const addComment = useAddComment()
  const [content, setContent] = useState('')

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const normalizedContent = content.trim()

    if (!isLoggedIn) {
      setShowLoginDialog(true)
      return
    }
    if (!normalizedContent || addComment.isPending) {
      return
    }

    try {
      await addComment.mutateAsync({
        target_type: targetType,
        target_id: targetId,
        parent_id: parentId,
        reply_to: replyTo,
        content: normalizedContent,
      })
      setContent('')
      onSubmitted()
    } catch {
      // React Query 已将错误写入 mutation.error，由回复框下方统一展示。
    }
  }

  return (
    <form
      className="mt-3 rounded-md border border-input bg-muted/40 p-3"
      onSubmit={event => void handleSubmit(event)}
    >
      <p className="mb-2 text-xs text-muted-foreground">
        回复 <span className="font-medium text-foreground">@{replyToName}</span>
      </p>
      <label htmlFor={`reply-${replyTo}-${parentId}`} className="sr-only">
        回复 @{replyToName}
      </label>
      <textarea
        id={`reply-${replyTo}-${parentId}`}
        autoFocus
        rows={2}
        value={content}
        className="min-h-16 w-full resize-y rounded-md border border-input bg-background px-3 py-2 text-sm leading-6 outline-none focus:border-ring"
        placeholder="写下你的回复…"
        onChange={event => setContent(event.target.value)}
      />
      <div className="mt-2 flex items-center justify-end gap-2">
        <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
          取消
        </Button>
        <Button type="submit" size="sm" disabled={!content.trim() || addComment.isPending}>
          {addComment.isPending ? (
            <LoaderCircle className="animate-spin" aria-hidden="true" />
          ) : (
            <Send aria-hidden="true" />
          )}
          {addComment.isPending ? '提交中…' : '回复'}
        </Button>
      </div>
      {addComment.error ? (
        <p className="mt-2 text-sm text-destructive" role="alert">
          {addComment.error.message}
        </p>
      ) : null}
    </form>
  )
}
