import { useState } from 'react'
import { MessageCircle } from 'lucide-react'

import { CommentInput } from '@/components/comment/CommentInput'
import { CommentList } from '@/components/comment/CommentList'
import type { CommentTargetType } from '@/components/comment/types'
import type { CommentInfo } from '@/types/comment'

export interface CommentsProps {
  targetType: CommentTargetType
  targetId: string
}

/** 评论模块入口，统一管理当前展开的回复框和编辑框。 */
export function Comments({ targetType, targetId }: CommentsProps) {
  const [replyingTo, setReplyingTo] = useState<string | null>(null)
  const [editingComment, setEditingComment] = useState<CommentInfo | null>(null)

  const handleReply = (commentId: string) => {
    setEditingComment(null)
    setReplyingTo(current => (current === commentId ? null : commentId))
  }

  const handleEdit = (comment: CommentInfo) => {
    setReplyingTo(null)
    setEditingComment(current => (current?.comment_id === comment.comment_id ? null : comment))
  }

  return (
    <section className="border-t border-border pt-8" aria-labelledby="comments-title">
      <div className="mb-5 flex items-center gap-2">
        <MessageCircle className="size-5 text-muted-foreground" aria-hidden="true" />
        <h2 id="comments-title" className="text-lg font-semibold">
          评论
        </h2>
      </div>

      <CommentInput targetType={targetType} targetId={targetId} />
      <CommentList
        targetType={targetType}
        targetId={targetId}
        replyingTo={replyingTo}
        editingComment={editingComment}
        onReply={handleReply}
        onEdit={handleEdit}
        onCloseReply={() => setReplyingTo(null)}
        onCloseEdit={() => setEditingComment(null)}
      />
    </section>
  )
}

export default Comments
