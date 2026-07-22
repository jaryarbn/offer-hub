import { useEffect, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { BookOpen, Search } from "lucide-react"
import { useForm } from "react-hook-form"
import { Link, useSearchParams } from "react-router-dom"
import { z } from "zod"

import { HotContent } from "@/components/HotContent"
import {
  QUESTION_PAGE_SIZE,
  QuestionListContent,
} from "@/components/QuestionListContent"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useQuestionList } from "@/hooks/useQuestionQueries"
import { cn } from "@/lib/utils"

const difficultyOptions = [
  { label: "全部", value: undefined },
  { label: "简单", value: 1 },
  { label: "中等", value: 2 },
  { label: "困难", value: 3 },
] as const

const searchSchema = z.object({
  keyword: z.string().trim().max(80, "关键词不能超过 80 个字符"),
})

type SearchValues = z.infer<typeof searchSchema>

export function Questions() {
  const [searchParams] = useSearchParams()
  const bankId = searchParams.get("bank_id")?.trim() || undefined
  const activeIndex = searchParams.get("activeIndex")
  const [difficulty, setDifficulty] = useState<number | undefined>()
  const [keyword, setKeyword] = useState("")
  const [page, setPage] = useState(1)
  const { data, isPending, isFetching, isError, error, refetch } = useQuestionList({
    bank_id: bankId,
    difficulty,
    keyword: keyword || undefined,
    page,
    page_size: QUESTION_PAGE_SIZE,
  })
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SearchValues>({
    resolver: zodResolver(searchSchema),
    defaultValues: { keyword: "" },
  })

  useEffect(() => setPage(1), [bankId])

  const total = data?.total ?? 0
  const totalPages = Math.ceil(total / QUESTION_PAGE_SIZE)
  const collectionUrl = activeIndex
    ? `/questions-collection?activeIndex=${encodeURIComponent(activeIndex)}`
    : "/questions-collection"

  const submitSearch = ({ keyword: nextKeyword }: SearchValues) => {
    setKeyword(nextKeyword)
    setPage(1)
  }

  const selectDifficulty = (nextDifficulty: number | undefined) => {
    setDifficulty(nextDifficulty)
    setPage(1)
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b border-border bg-background">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4 sm:px-6">
          <Link to="/" className="flex items-center gap-3 outline-none focus-visible:ring-2 focus-visible:ring-ring">
            <span className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <BookOpen className="size-4" aria-hidden="true" />
            </span>
            <span className="text-sm font-semibold">Offer Hub</span>
          </Link>
          <Link to={collectionUrl} className="text-sm text-muted-foreground hover:text-foreground hover:underline">
            返回题库分类
          </Link>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6">
        <div className="mb-6">
          <h1 className="text-2xl font-semibold">题目列表</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {data ? `共 ${total} 道题目` : "正在读取题目"}
            {isFetching && !isPending ? " · 正在更新" : ""}
          </p>
        </div>

        <section aria-label="题目筛选" className="border-y border-border py-4">
          <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
            <div>
              <span className="mb-2 block text-xs font-medium text-muted-foreground">难度</span>
              <div className="inline-flex rounded-md bg-muted p-1" aria-label="按难度筛选">
                {difficultyOptions.map((option) => (
                  <button
                    key={option.label}
                    type="button"
                    className={cn(
                      "h-8 rounded-sm px-3 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring",
                      difficulty === option.value
                        ? "bg-background font-medium text-foreground shadow-sm"
                        : "text-muted-foreground hover:text-foreground",
                    )}
                    aria-pressed={difficulty === option.value}
                    onClick={() => selectDifficulty(option.value)}
                  >
                    {option.label}
                  </button>
                ))}
              </div>
            </div>

            <form className="w-full max-w-md" onSubmit={handleSubmit(submitSearch)}>
              <label htmlFor="question-keyword" className="mb-2 block text-xs font-medium text-muted-foreground">
                关键词
              </label>
              <div className="flex items-start gap-2">
                <div className="min-w-0 flex-1">
                  <Input
                    id="question-keyword"
                    type="search"
                    placeholder="搜索标题或内容"
                    aria-invalid={Boolean(errors.keyword)}
                    {...register("keyword")}
                  />
                  {errors.keyword ? (
                    <p className="mt-1 text-xs text-destructive">{errors.keyword.message}</p>
                  ) : null}
                </div>
                <Button type="submit" disabled={isFetching}>
                  <Search aria-hidden="true" />
                  搜索
                </Button>
              </div>
            </form>
          </div>
        </section>

        <div className="mt-8 grid gap-10 lg:grid-cols-[minmax(0,1fr)_18rem]">
          <section aria-label="题目列表" className="min-w-0">
            <QuestionListContent
              questions={data?.list}
              page={page}
              totalPages={totalPages}
              isPending={isPending}
              isError={isError}
              errorMessage={error?.message}
              onRetry={() => void refetch()}
              onPageChange={setPage}
            />
          </section>

          <HotContent />
        </div>
      </main>
    </div>
  )
}
