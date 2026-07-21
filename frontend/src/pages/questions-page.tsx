import { useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { BookOpen, RotateCcw, Search } from "lucide-react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useQuestionList } from "@/hooks/useQuestionQueries"

const searchSchema = z.object({
  keyword: z.string().trim().max(80, "关键词不能超过 80 个字符"),
})

type SearchValues = z.infer<typeof searchSchema>

export function QuestionsPage() {
  const [keyword, setKeyword] = useState("")
  const { data, isPending, isFetching, isError, refetch } = useQuestionList({
    keyword,
    page: 1,
    page_size: 20,
  })
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SearchValues>({
    resolver: zodResolver(searchSchema),
    defaultValues: { keyword: "" },
  })

  const submitSearch = ({ keyword: nextKeyword }: SearchValues) => {
    setKeyword(nextKeyword)
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b border-border bg-background">
        <div className="mx-auto flex h-14 max-w-6xl items-center gap-3 px-4 sm:px-6">
          <div className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <BookOpen className="size-4" aria-hidden="true" />
          </div>
          <span className="text-sm font-semibold">Offer Hub</span>
          <span className="h-4 w-px bg-border" aria-hidden="true" />
          <span className="text-sm text-muted-foreground">题目练习</span>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6">
        <div className="mb-6 flex flex-col gap-4 border-b border-border pb-6 md:flex-row md:items-end md:justify-between">
          <div>
            <h1 className="text-2xl font-semibold">题目列表</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              {data ? `共 ${data.total} 道题目` : "正在读取题库"}
            </p>
          </div>

          <form
            className="flex w-full max-w-md items-start gap-2"
            onSubmit={handleSubmit(submitSearch)}
          >
            <div className="min-w-0 flex-1">
              <Input
                type="search"
                placeholder="搜索标题或内容"
                aria-label="搜索题目"
                aria-invalid={Boolean(errors.keyword)}
                {...register("keyword")}
              />
              {errors.keyword ? (
                <p className="mt-1 text-xs text-destructive">
                  {errors.keyword.message}
                </p>
              ) : null}
            </div>
            <Button type="submit" disabled={isFetching}>
              <Search aria-hidden="true" />
              搜索
            </Button>
          </form>
        </div>

        {isPending ? <QuestionListSkeleton /> : null}

        {isError ? (
          <div className="flex min-h-56 flex-col items-center justify-center border-y border-border text-center">
            <p className="text-sm font-medium">题目加载失败</p>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="mt-3"
              onClick={() => void refetch()}
            >
              <RotateCcw aria-hidden="true" />
              重新加载
            </Button>
          </div>
        ) : null}

        {!isPending && !isError && data?.list.length === 0 ? (
          <div className="flex min-h-56 items-center justify-center border-y border-border text-sm text-muted-foreground">
            没有匹配的题目
          </div>
        ) : null}

        {!isPending && !isError && data?.list.length ? (
          <ol className="divide-y divide-border border-y border-border">
            {data.list.map((question, index) => (
              <li
                key={question.question_id}
                className="grid min-h-24 grid-cols-[2.5rem_1fr] gap-3 py-5 sm:grid-cols-[3rem_1fr_auto] sm:items-center"
              >
                <span className="pt-0.5 font-mono text-xs text-muted-foreground sm:pt-0">
                  {String(index + 1).padStart(2, "0")}
                </span>
                <div className="min-w-0">
                  <h2 className="text-sm font-medium leading-6">
                    {question.title}
                  </h2>
                  <div className="mt-2 flex flex-wrap items-center gap-2">
                    <span className="text-xs text-muted-foreground">
                      难度 {question.difficulty}
                    </span>
                    {question.tags.slice(0, 3).map((tag) => (
                      <span
                        key={tag}
                        className="rounded-sm bg-secondary px-1.5 py-0.5 text-xs text-secondary-foreground"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
                <span className="col-start-2 text-xs text-muted-foreground sm:col-start-auto">
                  {question.view_count} 次浏览
                </span>
              </li>
            ))}
          </ol>
        ) : null}
      </main>
    </div>
  )
}

function QuestionListSkeleton() {
  return (
    <div className="divide-y divide-border border-y border-border" aria-label="题目加载中">
      {[0, 1, 2, 3].map((item) => (
        <div key={item} className="grid min-h-24 grid-cols-[2.5rem_1fr] gap-3 py-5 sm:grid-cols-[3rem_1fr_auto]">
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
