import { useState, type FormEvent } from 'react'
import { LoaderCircle, Send } from 'lucide-react'

import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'
import { useAddComment } from '@/hooks/useCommentQueries'
import type { CommentTargetType } from '@/components/comment/types'

export interface CommentInputProps {
  targetType: CommentTargetType
  targetId: string
}

/** 发表顶层评论；游客聚焦输入框时打开登录弹窗。 */
export function CommentInput({ targetType, targetId }: CommentInputProps) {
  const { isLoggedIn, setShowLoginDialog } = useLogin()
  const addComment = useAddComment()
  const [content, setContent] = useState('')

  const openLoginDialog = () => {
    if (!isLoggedIn) {
      setShowLoginDialog(true)
    }
  }

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
        content: normalizedContent,
      })
      setContent('')
    } catch {
      // React Query 已将错误写入 mutation.error，由表单下方统一展示。
    }
  }

  return (
    <form className="mb-7" onSubmit={event => void handleSubmit(event)}>
      <label htmlFor="new-comment" className="sr-only">
        发表评论
      </label>
      <div className="rounded-lg border border-input bg-card p-3 transition-colors focus-within:border-ring">
        <textarea
          id="new-comment"
          value={content}
          readOnly={!isLoggedIn}
          rows={3}
          className="min-h-20 w-full resize-y bg-transparent text-sm leading-6 outline-none placeholder:text-muted-foreground"
          placeholder={isLoggedIn ? '分享你的思考…' : '登录后参与讨论'}
          onClick={openLoginDialog}
          onFocus={openLoginDialog}
          onChange={event => setContent(event.target.value)}
          aria-describedby={addComment.error ? 'new-comment-error' : undefined}
        />
        <div className="mt-2 flex items-center justify-between gap-3 border-t border-border pt-3">
          <p className="text-xs text-muted-foreground">请友善交流，评论内容会经过敏感词过滤</p>
          <Button
            type="submit"
            size="sm"
            disabled={addComment.isPending || (isLoggedIn && !content.trim())}
          >
            {addComment.isPending ? (
              <LoaderCircle className="animate-spin" aria-hidden="true" />
            ) : (
              <Send aria-hidden="true" />
            )}
            {addComment.isPending ? '发表中…' : '发表评论'}
          </Button>
        </div>
      </div>
      {addComment.error ? (
        <p id="new-comment-error" className="mt-2 text-sm text-destructive" role="alert">
          {addComment.error.message}
        </p>
      ) : null}
    </form>
  )
}
