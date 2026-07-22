import { AlertCircle, ArrowRight, BookOpen, Layers3, RotateCcw } from 'lucide-react'
import { Link, useSearchParams } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { useQuestionBanks } from '@/hooks/useQuestionQueries'
import { cn } from '@/lib/utils'
import type { QuestionBankGroup } from '@/types/question'

export function QuestionCollections() {
  const [searchParams, setSearchParams] = useSearchParams()
  const { data, isPending, isError, error, refetch } = useQuestionBanks()

  if (isPending) {
    return <QuestionCollectionsSkeleton />
  }

  if (isError) {
    return (
      <PageLayout>
        <div
          className="flex min-h-72 flex-col items-center justify-center border-y border-border px-4 text-center"
          role="alert"
        >
          <AlertCircle className="size-6 text-destructive" aria-hidden="true" />
          <h2 className="mt-3 text-sm font-semibold">题库加载失败</h2>
          <p className="mt-1 max-w-md text-sm text-muted-foreground">
            {error.message || '暂时无法获取题库，请稍后重试。'}
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
      </PageLayout>
    )
  }

  if (!data.length) {
    return (
      <PageLayout>
        <div className="flex min-h-72 flex-col items-center justify-center border-y border-border px-4 text-center">
          <Layers3 className="size-6 text-muted-foreground" aria-hidden="true" />
          <h2 className="mt-3 text-sm font-semibold">暂无可用题库</h2>
          <p className="mt-1 text-sm text-muted-foreground">题库数据准备完成后会显示在这里。</p>
        </div>
      </PageLayout>
    )
  }

  const activeIndex = getActiveIndex(searchParams.get('activeIndex'), data.length)
  const activeGroup = data[activeIndex]

  const selectGroup = (nextIndex: number) => {
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set('activeIndex', String(nextIndex))
    setSearchParams(nextParams, { replace: true })
  }

  return (
    <PageLayout>
      <nav aria-label="后端领域" className="overflow-x-auto border-b border-border">
        <ul className="flex min-w-max gap-1">
          {data.map((group, index) => {
            const isActive = index === activeIndex

            return (
              <li key={group.job_name}>
                <button
                  type="button"
                  className={cn(
                    'relative h-11 px-4 text-sm font-medium text-muted-foreground outline-none transition-colors',
                    'hover:text-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring',
                    isActive &&
                      'text-foreground after:absolute after:inset-x-4 after:bottom-0 after:h-0.5 after:bg-foreground'
                  )}
                  aria-current={isActive ? 'page' : undefined}
                  onClick={() => selectGroup(index)}
                >
                  {group.job_name}
                </button>
              </li>
            )
          })}
        </ul>
      </nav>

      <QuestionBankGroupContent group={activeGroup} activeIndex={activeIndex} />
    </PageLayout>
  )
}

function QuestionBankGroupContent({
  group,
  activeIndex,
}: {
  group: QuestionBankGroup
  activeIndex: number
}) {
  return (
    <div className="space-y-9 pt-8">
      {group.series_list.map(series => (
        <section key={series.series_id} aria-labelledby={`series-${series.series_id}`}>
          <div className="mb-4 flex items-baseline justify-between gap-4">
            <h2 id={`series-${series.series_id}`} className="text-base font-semibold">
              {series.series_name}
            </h2>
            <span className="shrink-0 text-xs text-muted-foreground">
              {series.bank_list.length} 个题库
            </span>
          </div>

          <ul className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {series.bank_list.map(bank => (
              <li key={bank.bank_id}>
                <Link
                  to={getQuestionListUrl(bank.bank_id, activeIndex)}
                  className="group flex min-h-36 flex-col rounded-md border border-border bg-card p-4 text-card-foreground outline-none transition-colors hover:border-foreground/30 hover:bg-muted/40 focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex size-9 shrink-0 items-center justify-center rounded-md bg-secondary text-secondary-foreground">
                      <BookOpen className="size-4" aria-hidden="true" />
                    </div>
                    <ArrowRight
                      className="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5 group-hover:text-foreground"
                      aria-hidden="true"
                    />
                  </div>
                  <h3 className="mt-4 text-sm font-semibold leading-5">{bank.bank_name}</h3>
                  <p className="mt-1 line-clamp-2 text-sm leading-5 text-muted-foreground">
                    {bank.desc || '暂无题库描述'}
                  </p>
                  <span className="mt-auto pt-4 text-xs text-muted-foreground">
                    {bank.count} 道题目
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      ))}
    </div>
  )
}

function PageLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b border-border bg-background">
        <div className="mx-auto flex h-14 max-w-6xl items-center gap-3 px-4 sm:px-6">
          <div className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <BookOpen className="size-4" aria-hidden="true" />
          </div>
          <span className="text-sm font-semibold">Offer Hub</span>
          <span className="h-4 w-px bg-border" aria-hidden="true" />
          <span className="text-sm text-muted-foreground">题库分类</span>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6">
        <div className="mb-7">
          <h1 className="text-2xl font-semibold">选择题库</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            按技术领域浏览题目，选择一个题库开始练习。
          </p>
        </div>
        {children}
      </main>
    </div>
  )
}

function QuestionCollectionsSkeleton() {
  return (
    <PageLayout>
      <div aria-busy="true" aria-label="题库加载中">
        <div className="flex h-11 gap-3 border-b border-border">
          {[0, 1, 2].map(item => (
            <div key={item} className="mt-3 h-4 w-20 animate-pulse rounded-sm bg-muted" />
          ))}
        </div>
        <div className="pt-8">
          <div className="h-5 w-28 animate-pulse rounded-sm bg-muted" />
          <div className="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {[0, 1, 2, 3, 4, 5].map(item => (
              <div
                key={item}
                className="min-h-36 animate-pulse rounded-md border border-border p-4"
              >
                <div className="size-9 rounded-md bg-muted" />
                <div className="mt-4 h-4 w-2/3 rounded-sm bg-muted" />
                <div className="mt-2 h-3 w-5/6 rounded-sm bg-muted" />
              </div>
            ))}
          </div>
        </div>
      </div>
    </PageLayout>
  )
}

function getActiveIndex(value: string | null, groupCount: number): number {
  if (value === null || value.trim() === '') return 0

  const index = Number(value)
  return Number.isInteger(index) && index >= 0 && index < groupCount ? index : 0
}

function getQuestionListUrl(bankId: string, activeIndex: number): string {
  const params = new URLSearchParams({
    bank_id: bankId,
    activeIndex: String(activeIndex),
  })
  return `/questions?${params.toString()}`
}
