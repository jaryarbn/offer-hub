import { useState, type FormEvent } from 'react'
import { LoaderCircle } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { useUpdateComment } from '@/hooks/useCommentQueries'
import type { CommentInfo } from '@/types/comment'

export interface CommentEditFormProps {
  comment: CommentInfo
  onCancel: () => void
  onSaved: () => void
}

export function CommentEditForm({ comment, onCancel, onSaved }: CommentEditFormProps) {
  const updateComment = useUpdateComment()
  const [content, setContent] = useState(comment.content)

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const normalizedContent = content.trim()
    if (!normalizedContent || updateComment.isPending) {
      return
    }

    try {
      await updateComment.mutateAsync({
        comment_id: comment.comment_id,
        content: normalizedContent,
      })
      onSaved()
    } catch {
      // React Query 已将错误写入 mutation.error，由编辑框下方统一展示。
    }
  }

  return (
    <form className="mt-3" onSubmit={event => void handleSubmit(event)}>
      <label htmlFor={`edit-${comment.comment_id}`} className="sr-only">
        编辑评论
      </label>
      <textarea
        id={`edit-${comment.comment_id}`}
        autoFocus
        rows={3}
        value={content}
        className="w-full resize-y rounded-md border border-input bg-background px-3 py-2 text-sm leading-6 outline-none focus:border-ring"
        onChange={event => setContent(event.target.value)}
      />
      <div className="mt-2 flex justify-end gap-2">
        <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
          取消
        </Button>
        <Button type="submit" size="sm" disabled={!content.trim() || updateComment.isPending}>
          {updateComment.isPending ? (
            <LoaderCircle className="animate-spin" aria-hidden="true" />
          ) : null}
          {updateComment.isPending ? '保存中…' : '保存'}
        </Button>
      </div>
      {updateComment.error ? (
        <p className="mt-2 text-sm text-destructive" role="alert">
          {updateComment.error.message}
        </p>
      ) : null}
    </form>
  )
}
