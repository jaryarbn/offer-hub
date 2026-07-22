import type { CommentInfo, ListCommentsParams } from '@/types/comment'

export type CommentTargetType = ListCommentsParams['target_type']

export interface CommentThreadStateProps {
  replyingTo: string | null
  editingComment: CommentInfo | null
  onReply: (commentId: string) => void
  onEdit: (comment: CommentInfo) => void
  onCloseReply: () => void
  onCloseEdit: () => void
}
