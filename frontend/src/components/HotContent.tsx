import { Flame, RotateCcw } from 'lucide-react'
import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { useHotQuestions } from '@/hooks/useQuestionQueries'
import { cn } from '@/lib/utils'

export interface HotQuestionsProps {
  limit?: number
  showHeader?: boolean
  variant?: 'sidebar' | 'section'
}

export function HotQuestions({
  limit = 10,
  showHeader = true,
  variant = 'sidebar',
}: HotQuestionsProps) {
  const { data, isPending, isError, refetch } = useHotQuestions(limit)

  const content = (
    <>
      {showHeader ? (
        <div className="flex items-center gap-2 border-b border-border pb-3">
          <Flame className="size-4 text-muted-foreground" aria-hidden="true" />
          <h2 id="hot-questions-title" className="text-sm font-semibold">
            热门题目
          </h2>
        </div>
      ) : null}

      {isPending ? <HotContentSkeleton /> : null}

      {isError ? (
        <div className="py-8 text-center" role="alert">
          <p className="text-sm text-muted-foreground">热门题目加载失败</p>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="mt-3"
            onClick={() => void refetch()}
          >
            <RotateCcw aria-hidden="true" />
            重试
          </Button>
        </div>
      ) : null}

      {!isPending && !isError && data?.list.length === 0 ? (
        <p className="py-8 text-center text-sm text-muted-foreground">暂无热门题目</p>
      ) : null}

      {!isPending && !isError && data?.list.length ? (
        <ol className="divide-y divide-border">
          {data.list.map((question, index) => (
            <li key={question.question_id}>
              <Link
                to={`/questions/${encodeURIComponent(question.question_id)}`}
                className={cn(
                  'group grid gap-3 py-3 outline-none focus-visible:ring-2 focus-visible:ring-ring',
                  variant === 'section'
                    ? 'grid-cols-[2rem_minmax(0,1fr)_auto] items-center sm:py-4'
                    : 'grid-cols-[1.5rem_1fr]'
                )}
              >
                <span
                  className={cn(
                    'pt-0.5 font-mono text-xs text-muted-foreground',
                    variant === 'section' && 'text-center text-sm font-semibold'
                  )}
                >
                  {String(index + 1).padStart(2, '0')}
                </span>
                <span className="min-w-0">
                  <span
                    className={cn(
                      'line-clamp-2 font-medium group-hover:underline',
                      variant === 'section' ? 'text-base leading-6' : 'text-sm leading-5'
                    )}
                  >
                    {question.title}
                  </span>
                  {variant === 'sidebar' ? (
                    <span className="mt-1 block text-xs text-muted-foreground">
                      {question.view_count} 次浏览
                    </span>
                  ) : null}
                </span>
                {variant === 'section' ? (
                  <span className="shrink-0 text-xs text-muted-foreground sm:text-sm">
                    {question.view_count} 次浏览
                  </span>
                ) : null}
              </Link>
            </li>
          ))}
        </ol>
      ) : null}
    </>
  )

  if (variant === 'section') {
    return <div aria-label="热门题目">{content}</div>
  }

  return (
    <aside aria-labelledby="hot-questions-title" className="lg:border-l lg:border-border lg:pl-6">
      {content}
    </aside>
  )
}

export function HotContent() {
  return <HotQuestions />
}

function HotContentSkeleton() {
  return (
    <div className="divide-y divide-border" aria-busy="true" aria-label="热门题目加载中">
      {[0, 1, 2, 3, 4].map(item => (
        <div key={item} className="grid grid-cols-[1.5rem_1fr] gap-2 py-4">
          <div className="h-3 w-4 animate-pulse rounded-sm bg-muted" />
          <div>
            <div className="h-4 w-full animate-pulse rounded-sm bg-muted" />
            <div className="mt-2 h-3 w-20 animate-pulse rounded-sm bg-muted" />
          </div>
        </div>
      ))}
    </div>
  )
}
