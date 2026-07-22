import { RotateCcw } from "lucide-react"
import { Link } from "react-router-dom"

import { QuestionPagination } from "@/components/QuestionPagination"
import { Button } from "@/components/ui/button"
import type { Question } from "@/types/question"

export const QUESTION_PAGE_SIZE = 20

interface QuestionListContentProps {
  questions?: Question[]
  page: number
  totalPages: number
  isPending: boolean
  isError: boolean
  errorMessage?: string
  onRetry: () => void
  onPageChange: (page: number) => void
}

export function QuestionListContent({
  questions,
  page,
  totalPages,
  isPending,
  isError,
  errorMessage,
  onRetry,
  onPageChange,
}: QuestionListContentProps) {
  if (isPending) return <QuestionListSkeleton />

  if (isError) {
    return (
      <div
        className="flex min-h-56 flex-col items-center justify-center border-y border-border px-4 text-center"
        role="alert"
      >
        <p className="text-sm font-medium">题目加载失败</p>
        <p className="mt-1 text-sm text-muted-foreground">
          {errorMessage || "请稍后重试"}
        </p>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="mt-3"
          onClick={onRetry}
        >
          <RotateCcw aria-hidden="true" />
          重新加载
        </Button>
      </div>
    )
  }

  if (!questions?.length) {
    return (
      <div className="flex min-h-56 items-center justify-center border-y border-border text-sm text-muted-foreground">
        没有匹配的题目
      </div>
    )
  }

  return (
    <>
      <ol className="divide-y divide-border border-y border-border">
        {questions.map((question, index) => (
          <li key={question.question_id}>
            <Link
              to={`/questions/${encodeURIComponent(question.question_id)}`}
              className="group grid min-h-24 grid-cols-[2.5rem_minmax(0,1fr)] gap-3 py-5 outline-none hover:bg-muted/40 focus-visible:ring-2 focus-visible:ring-ring sm:grid-cols-[3rem_minmax(0,1fr)_auto] sm:items-center"
            >
              <span className="px-1 pt-0.5 font-mono text-xs text-muted-foreground sm:pt-0">
                {String((page - 1) * QUESTION_PAGE_SIZE + index + 1).padStart(2, "0")}
              </span>
              <span className="min-w-0">
                <span className="block text-sm font-medium leading-6 group-hover:underline">
                  {question.title}
                </span>
                <span className="mt-2 flex flex-wrap items-center gap-2">
                  <span className="text-xs text-muted-foreground">
                    {getDifficultyLabel(question.difficulty)}
                  </span>
                  {question.tags.slice(0, 3).map((tag) => (
                    <span
                      key={tag}
                      className="rounded-sm bg-secondary px-1.5 py-0.5 text-xs text-secondary-foreground"
                    >
                      {tag}
                    </span>
                  ))}
                </span>
              </span>
              <span className="col-start-2 text-xs text-muted-foreground sm:col-start-auto sm:px-2">
                {question.view_count} 次浏览
              </span>
            </Link>
          </li>
        ))}
      </ol>

      {totalPages > 1 ? (
        <QuestionPagination page={page} totalPages={totalPages} onChange={onPageChange} />
      ) : null}
    </>
  )
}

function QuestionListSkeleton() {
  return (
    <div
      className="divide-y divide-border border-y border-border"
      aria-busy="true"
      aria-label="题目加载中"
    >
      {[0, 1, 2, 3, 4].map((item) => (
        <div
          key={item}
          className="grid min-h-24 grid-cols-[2.5rem_1fr] gap-3 py-5 sm:grid-cols-[3rem_1fr_auto]"
        >
          <div className="h-3 w-5 animate-pulse rounded-sm bg-muted" />
          <div>
            <div className="h-4 w-2/3 animate-pulse rounded-sm bg-muted" />
            <div className="mt-3 h-3 w-1/3 animate-pulse rounded-sm bg-muted" />
          </div>
        </div>
      ))}
    </div>
  )
}

function getDifficultyLabel(difficulty: number): string {
  return ["", "简单", "中等", "困难"][difficulty] ?? `难度 ${difficulty}`
}
