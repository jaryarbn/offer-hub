import { useState } from 'react'
import { ChevronDown } from 'lucide-react'

import { CommentActions } from '@/components/comment/CommentActions'
import { CommentEditForm } from '@/components/comment/CommentEditForm'
import { ReplyBoxEnhanced } from '@/components/comment/ReplyBoxEnhanced'
import { SubComments } from '@/components/comment/SubComments'
import type { CommentTargetType, CommentThreadStateProps } from '@/components/comment/types'
import { useLogin } from '@/components/provider/LoginProvider'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import type { CommentInfo } from '@/types/comment'

export interface CommentItemProps extends CommentThreadStateProps {
  comment: CommentInfo
  targetType: CommentTargetType
  targetId: string
  isSub?: boolean
}

export function CommentItem({
  comment,
  targetType,
  targetId,
  isSub = false,
  replyingTo,
  editingComment,
  onReply,
  onEdit,
  onCloseReply,
  onCloseEdit,
}: CommentItemProps) {
  const { userInfo } = useLogin()
  const [showSubComments, setShowSubComments] = useState(false)
  const isOwnComment = userInfo?.user_id === comment.user_id
  const isReplying = replyingTo === comment.comment_id
  const isEditing = editingComment?.comment_id === comment.comment_id
  const parentId = isSub ? comment.parent_id : comment.comment_id
  const displayName = comment.user_name.trim() || '用户'

  return (
    <article className={isSub ? 'py-4' : 'py-5'}>
      <div className="flex gap-3">
        <Avatar className={isSub ? 'size-8' : 'size-9'}>
          <AvatarImage src={comment.user_avatar || undefined} alt={`${displayName}的头像`} />
          <AvatarFallback>{displayName.slice(0, 1).toUpperCase()}</AvatarFallback>
        </Avatar>

        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-baseline gap-x-2 gap-y-1">
            <span className="text-sm font-medium">{displayName}</span>
            {comment.reply_to_name ? (
              <span className="text-xs text-muted-foreground">回复 @{comment.reply_to_name}</span>
            ) : null}
            <time className="text-xs text-muted-foreground" dateTime={comment.create_time}>
              {formatCommentTime(comment.create_time)}
            </time>
          </div>

          {isEditing ? (
            <CommentEditForm comment={comment} onCancel={onCloseEdit} onSaved={onCloseEdit} />
          ) : (
            <p className="mt-2 whitespace-pre-wrap break-words text-sm leading-6 text-foreground/90">
              {comment.content}
            </p>
          )}

          <CommentActions
            comment={comment}
            isOwnComment={isOwnComment}
            onReply={() => onReply(comment.comment_id)}
            onEdit={() => onEdit(comment)}
          />

          {isReplying ? (
            <ReplyBoxEnhanced
              targetType={targetType}
              targetId={targetId}
              parentId={parentId}
              replyTo={comment.user_id}
              replyToName={displayName}
              onSubmitted={onCloseReply}
              onCancel={onCloseReply}
            />
          ) : null}

          {!isSub && comment.sub_comment_total > 0 ? (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="mt-2 text-muted-foreground"
              aria-expanded={showSubComments}
              onClick={() => setShowSubComments(value => !value)}
            >
              <ChevronDown
                className={
                  showSubComments ? 'rotate-180 transition-transform' : 'transition-transform'
                }
                aria-hidden="true"
              />
              {showSubComments ? '收起回复' : `查看 ${comment.sub_comment_total} 条回复`}
            </Button>
          ) : null}

          {showSubComments ? (
            <SubComments
              targetType={targetType}
              targetId={targetId}
              parentId={comment.comment_id}
              replyingTo={replyingTo}
              editingComment={editingComment}
              onReply={onReply}
              onEdit={onEdit}
              onCloseReply={onCloseReply}
              onCloseEdit={onCloseEdit}
            />
          ) : null}
        </div>
      </div>
    </article>
  )
}

function formatCommentTime(value: string): string {
  const normalized = value.includes('T') ? value : value.replace(' ', 'T')
  const date = new Date(normalized.endsWith('Z') ? normalized : `${normalized}Z`)

  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}
