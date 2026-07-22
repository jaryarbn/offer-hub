import { useEffect, useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { Bookmark, Check, ChevronDown, LoaderCircle, ThumbsUp } from 'lucide-react'

import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { queryKeys } from '@/hooks/useQuestionQueries'
import { cn } from '@/lib/utils'
import { toggleLike } from '@/services/comment'
import { tagQuestion } from '@/services/question'
import { TargetType } from '@/types/comment'
import type { Question, QuestionTag } from '@/types/question'

interface DetailHeaderProps {
  question: Question
}

type ActiveQuestionTag = Exclude<QuestionTag, 0>

const tagOptions: Array<{
  value: ActiveQuestionTag
  label: string
  badgeClassName: string
}> = [
  {
    value: 1,
    label: '已掌握',
    badgeClassName: 'border-emerald-200 bg-emerald-50 text-emerald-700',
  },
  {
    value: 2,
    label: '晚点再刷',
    badgeClassName: 'border-amber-200 bg-amber-50 text-amber-700',
  },
  {
    value: 3,
    label: '未掌握',
    badgeClassName: 'border-red-200 bg-red-50 text-red-700',
  },
]

export function DetailHeader({ question }: DetailHeaderProps) {
  const queryClient = useQueryClient()
  const { isLoggedIn, setShowLoginDialog } = useLogin()
  const [isLiking, setIsLiking] = useState(false)
  const [isTagging, setIsTagging] = useState(false)
  const [menuOpen, setMenuOpen] = useState(false)
  const [interactionError, setInteractionError] = useState<string | null>(null)
  const closeTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const activeTag = tagOptions.find(option => option.value === question.user_tag)

  const cancelMenuClose = () => {
    if (closeTimerRef.current) {
      clearTimeout(closeTimerRef.current)
      closeTimerRef.current = null
    }
  }

  const openMenu = () => {
    cancelMenuClose()
    setMenuOpen(true)
  }

  const scheduleMenuClose = () => {
    cancelMenuClose()
    closeTimerRef.current = setTimeout(() => setMenuOpen(false), 120)
  }

  useEffect(() => () => {
    cancelMenuClose()
  })

  const handleLike = async () => {
    if (!isLoggedIn) {
      setShowLoginDialog(true)
      return
    }
    if (isLiking) {
      return
    }

    setIsLiking(true)
    setInteractionError(null)
    try {
      const result = await toggleLike(
        TargetType.QUESTION,
        question.question_id,
        !question.user_liked
      )
      queryClient.setQueryData<Question>(queryKeys.detail(question.question_id), current =>
        current
          ? {
              ...current,
              user_liked: result.liked,
              thumbs_up_count: result.count,
            }
          : current
      )
    } catch (error) {
      setInteractionError(error instanceof Error ? error.message : '更新点赞状态失败')
    } finally {
      setIsLiking(false)
    }
  }

  const handleTag = async (tag: ActiveQuestionTag) => {
    if (isTagging || tag === question.user_tag) {
      setMenuOpen(false)
      return
    }

    setIsTagging(true)
    setInteractionError(null)
    try {
      await tagQuestion(question.question_id, tag)
      queryClient.setQueryData<Question>(queryKeys.detail(question.question_id), current =>
        current ? { ...current, user_tag: tag } : current
      )
      setMenuOpen(false)
    } catch (error) {
      setInteractionError(error instanceof Error ? error.message : '更新掌握状态失败')
    } finally {
      setIsTagging(false)
    }
  }

  return (
    <header>
      <p className="font-mono text-xs text-muted-foreground">{question.question_id}</p>

      <div className="mt-3 flex flex-wrap items-center gap-3">
        <h1 className="min-w-0 text-2xl font-semibold leading-tight sm:text-3xl">
          {question.title}
        </h1>
        {activeTag ? (
          <span
            className={cn(
              'inline-flex shrink-0 items-center rounded-sm border px-2 py-1 text-xs font-medium',
              activeTag.badgeClassName
            )}
          >
            {activeTag.label}
          </span>
        ) : null}
      </div>

      <div className="mt-5 flex flex-wrap items-center justify-between gap-4">
        <div className="flex flex-wrap items-center gap-2">
          <span className="rounded-sm bg-primary px-2 py-1 text-xs font-medium text-primary-foreground">
            {getDifficultyLabel(question.difficulty)}
          </span>
          {question.tags.map(tag => (
            <span
              key={tag}
              className="rounded-sm bg-secondary px-2 py-1 text-xs text-secondary-foreground"
            >
              {tag}
            </span>
          ))}
        </div>

        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className={question.user_liked ? 'text-primary' : undefined}
            disabled={isLiking}
            aria-pressed={question.user_liked}
            onClick={() => void handleLike()}
          >
            {isLiking ? (
              <LoaderCircle className="animate-spin" aria-hidden="true" />
            ) : (
              <ThumbsUp aria-hidden="true" />
            )}
            {question.thumbs_up_count}
          </Button>

          {isLoggedIn ? (
            <DropdownMenu open={menuOpen} onOpenChange={setMenuOpen}>
              <DropdownMenuTrigger asChild>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={isTagging}
                  onPointerEnter={openMenu}
                  onPointerLeave={scheduleMenuClose}
                >
                  {isTagging ? (
                    <LoaderCircle className="animate-spin" aria-hidden="true" />
                  ) : (
                    <Bookmark aria-hidden="true" />
                  )}
                  标记
                  <ChevronDown aria-hidden="true" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align="end"
                onPointerEnter={openMenu}
                onPointerLeave={scheduleMenuClose}
              >
                {tagOptions.map(option => (
                  <DropdownMenuItem
                    key={option.value}
                    disabled={isTagging}
                    onSelect={() => void handleTag(option.value)}
                  >
                    <span
                      className={cn('size-2 rounded-full border', option.badgeClassName)}
                      aria-hidden="true"
                    />
                    <span className="flex-1">{option.label}</span>
                    {question.user_tag === option.value ? (
                      <Check className="text-primary" aria-hidden="true" />
                    ) : null}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          ) : null}
        </div>
      </div>

      {interactionError ? (
        <p className="mt-3 text-sm text-destructive" role="alert">
          {interactionError}
        </p>
      ) : null}
    </header>
  )
}

function getDifficultyLabel(difficulty: number): string {
  return ['', '简单', '中等', '困难'][difficulty] ?? `难度 ${difficulty}`
}
