import { useState } from 'react'
import { AlertCircle, Eye, Flame, RotateCcw } from 'lucide-react'
import { Link } from 'react-router-dom'

import { Header } from '@/components/Header'
import { Button } from '@/components/ui/button'
import { useHotQuestions } from '@/hooks/useQuestionQueries'
import { cn } from '@/lib/utils'

const HOT_QUESTION_LIMIT = 20

const jobOptions = ['后端开发', '前端开发', 'AI开发'] as const

type JobName = (typeof jobOptions)[number]

const podiumClasses = [
  'border-amber-200 bg-amber-50 text-amber-700',
  'border-zinc-200 bg-zinc-100 text-zinc-600',
  'border-orange-200 bg-orange-50 text-orange-700',
] as const

export function HotContent() {
  const [jobName, setJobName] = useState<JobName>(jobOptions[0])
  const { data, isPending, isError, error, refetch } = useHotQuestions(HOT_QUESTION_LIMIT, jobName)

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Header
        sectionLabel="热门题目"
        action={
          <Link
            to="/questions"
            className="text-sm text-muted-foreground hover:text-foreground hover:underline"
          >
            返回题目列表
          </Link>
        }
      />

      <main className="mx-auto w-full max-w-4xl px-4 py-8 sm:px-6">
        <div className="flex items-center gap-3">
          <span className="flex size-9 items-center justify-center rounded-md bg-secondary text-secondary-foreground">
            <Flame className="size-4" aria-hidden="true" />
          </span>
          <div>
            <h1 className="text-2xl font-semibold">热门题目</h1>
            <p className="mt-1 text-sm text-muted-foreground">{jobName}</p>
          </div>
        </div>

        <div className="mt-7 overflow-x-auto border-b border-border">
          <div className="flex min-w-max gap-1" role="tablist" aria-label="职位方向">
            {jobOptions.map(option => {
              const isActive = option === jobName

              return (
                <button
                  key={option}
                  type="button"
                  role="tab"
                  className={cn(
                    'relative h-11 px-4 text-sm font-medium text-muted-foreground outline-none transition-colors',
                    'hover:text-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring',
                    isActive &&
                      'text-foreground after:absolute after:inset-x-4 after:bottom-0 after:h-0.5 after:bg-foreground'
                  )}
                  aria-selected={isActive}
                  onClick={() => setJobName(option)}
                >
                  {option}
                </button>
              )
            })}
          </div>
        </div>

        <section aria-label={`${jobName}热门题目`} className="mt-4">
          {isPending ? <HotQuestionListSkeleton /> : null}

          {isError ? (
            <div
              className="flex min-h-72 flex-col items-center justify-center border-y border-border px-4 text-center"
              role="alert"
            >
              <AlertCircle className="size-6 text-destructive" aria-hidden="true" />
              <h2 className="mt-3 text-sm font-semibold">热门题目加载失败</h2>
              <p className="mt-1 max-w-md text-sm text-muted-foreground">
                {error.message || '暂时无法获取热门题目，请稍后重试。'}
              </p>
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="mt-4"
                onClick={() => void refetch()}
              >
                <RotateCcw aria-hidden="true" />
                重新加载
              </Button>
            </div>
          ) : null}

          {!isPending && !isError && data?.list.length === 0 ? (
            <div className="flex min-h-72 flex-col items-center justify-center border-y border-border px-4 text-center">
              <Flame className="size-6 text-muted-foreground" aria-hidden="true" />
              <h2 className="mt-3 text-sm font-semibold">暂无热门题目</h2>
            </div>
          ) : null}

          {!isPending && !isError && data?.list.length ? (
            <ol className="divide-y divide-border">
              {data.list.map((question, index) => {
                const rank = index + 1

                return (
                  <li key={question.question_id}>
                    <Link
                      to={`/questions/${encodeURIComponent(question.question_id)}`}
                      className="group grid min-h-20 grid-cols-[2.5rem_minmax(0,1fr)_auto] items-center gap-3 px-2 py-3 outline-none transition-colors hover:bg-muted/50 focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring sm:gap-4 sm:px-3"
                    >
                      <span
                        className={cn(
                          'flex size-8 items-center justify-center rounded-md border font-mono text-sm font-semibold',
                          rank <= 3
                            ? podiumClasses[rank - 1]
                            : 'border-transparent text-muted-foreground'
                        )}
                        aria-label={`第 ${rank} 名`}
                      >
                        {String(rank).padStart(2, '0')}
                      </span>

                      <span className="min-w-0">
                        <span className="line-clamp-2 text-sm font-medium leading-5 group-hover:underline sm:text-base sm:leading-6">
                          {question.title}
                        </span>
                      </span>

                      <span className="flex shrink-0 items-center gap-1.5 text-xs text-muted-foreground sm:text-sm">
                        <Eye className="size-3.5" aria-hidden="true" />
                        <span>{formatViewCount(question.view_count)}</span>
                        <span className="sr-only">次浏览</span>
                      </span>
                    </Link>
                  </li>
                )
              })}
            </ol>
          ) : null}
        </section>
      </main>
    </div>
  )
}

function HotQuestionListSkeleton() {
  return (
    <div className="divide-y divide-border" aria-busy="true" aria-label="热门题目加载中">
      {Array.from({ length: 8 }, (_, index) => (
        <div
          key={index}
          className="grid min-h-20 grid-cols-[2.5rem_minmax(0,1fr)_3.5rem] items-center gap-3 px-2 py-3 sm:gap-4 sm:px-3"
        >
          <div className="size-8 animate-pulse rounded-md bg-muted" />
          <div className="h-4 w-4/5 animate-pulse rounded-sm bg-muted" />
          <div className="h-3 w-12 animate-pulse rounded-sm bg-muted" />
        </div>
      ))}
    </div>
  )
}

export function formatViewCount(viewCount: number): string {
  const normalizedCount = Math.max(0, Number.isFinite(viewCount) ? viewCount : 0)

  if (normalizedCount >= 10_000) {
    return `${(normalizedCount / 10_000).toFixed(1)}w`
  }
  if (normalizedCount >= 1_000) {
    return `${(normalizedCount / 1_000).toFixed(1)}k`
  }
  return String(Math.trunc(normalizedCount))
}
