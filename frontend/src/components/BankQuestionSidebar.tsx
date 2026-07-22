import { Search } from "lucide-react"
import { useDeferredValue, useState } from "react"
import { Link } from "react-router-dom"

import { Input } from "@/components/ui/input"
import { useBankQuestionList } from "@/hooks/useQuestionQueries"
import { cn } from "@/lib/utils"

interface BankQuestionSidebarProps {
  bankId: string
  activeQuestionId: string
}

export function BankQuestionSidebar({
  bankId,
  activeQuestionId,
}: BankQuestionSidebarProps) {
  const [keyword, setKeyword] = useState("")
  const deferredKeyword = useDeferredValue(keyword)
  const { data, isPending, isError, error, refetch } = useBankQuestionList(
    bankId,
    deferredKeyword,
  )

  return (
    <aside aria-labelledby="bank-question-list-title" className="min-w-0">
      <div className="border-b border-border pb-4">
        <h2 id="bank-question-list-title" className="text-sm font-semibold">
          题目列表
        </h2>
        <p className="mt-1 text-xs text-muted-foreground">
          {data ? `当前题库共 ${data.total} 道题` : "读取当前题库"}
        </p>
        <label className="relative mt-3 block" htmlFor="bank-question-search">
          <span className="sr-only">搜索当前题库</span>
          <Search
            className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
            aria-hidden="true"
          />
          <Input
            id="bank-question-search"
            type="search"
            value={keyword}
            onChange={(event) => setKeyword(event.target.value)}
            placeholder="搜索当前题库"
            className="pl-9"
          />
        </label>
      </div>

      {isPending ? <BankQuestionListSkeleton /> : null}

      {isError ? (
        <div className="py-8 text-center" role="alert">
          <p className="text-sm font-medium">题目列表加载失败</p>
          <p className="mt-1 text-xs text-muted-foreground">
            {error.message || "请稍后重试"}
          </p>
          <button
            type="button"
            className="mt-3 text-xs font-medium underline underline-offset-4"
            onClick={() => void refetch()}
          >
            重新加载
          </button>
        </div>
      ) : null}

      {!isPending && !isError && !data?.list.length ? (
        <p className="py-8 text-center text-sm text-muted-foreground">
          没有匹配的题目
        </p>
      ) : null}

      {!isPending && !isError && data?.list.length ? (
        <ol className="max-h-[calc(100vh-13rem)] divide-y divide-border overflow-y-auto">
          {data.list.map((question, index) => {
            const isActive = question.question_id === activeQuestionId

            return (
              <li key={question.question_id}>
                <Link
                  to={`/questions/${encodeURIComponent(question.question_id)}`}
                  aria-current={isActive ? "page" : undefined}
                  className={cn(
                    "grid grid-cols-[2rem_minmax(0,1fr)] gap-2 px-2 py-3 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring",
                    isActive
                      ? "bg-secondary text-secondary-foreground"
                      : "hover:bg-muted/60",
                  )}
                >
                  <span className="pt-0.5 font-mono text-xs text-muted-foreground">
                    {String(index + 1).padStart(2, "0")}
                  </span>
                  <span className="line-clamp-2 leading-5">{question.title}</span>
                </Link>
              </li>
            )
          })}
        </ol>
      ) : null}
    </aside>
  )
}

function BankQuestionListSkeleton() {
  return (
    <div
      className="divide-y divide-border"
      aria-busy="true"
      aria-label="题目列表加载中"
    >
      {[0, 1, 2, 3, 4, 5].map((item) => (
        <div key={item} className="grid grid-cols-[2rem_1fr] gap-2 px-2 py-4">
          <div className="h-3 w-4 animate-pulse rounded-sm bg-muted" />
          <div className="h-4 w-4/5 animate-pulse rounded-sm bg-muted" />
        </div>
      ))}
    </div>
  )
}
