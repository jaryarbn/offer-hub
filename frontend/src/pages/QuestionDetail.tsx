import { AlertCircle, ArrowLeft } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'

import { BankQuestionSidebar } from '@/components/BankQuestionSidebar'
import { Header } from '@/components/Header'
import { Comments } from '@/components/comment/Comments'
import { DetailContent } from '@/components/detail/DetailContent'
import { Button } from '@/components/ui/button'
import { useQuestionDetail } from '@/hooks/useQuestionQueries'
import { TargetType } from '@/types/comment'

export function QuestionDetail() {
  const { question_id: questionId = '' } = useParams<{ question_id: string }>()
  const { data: question, isPending, isError, error, refetch } = useQuestionDetail(questionId)

  if (!questionId || (isError && error.code === 404)) {
    return <QuestionNotFound />
  }

  if (isPending) {
    return <QuestionDetailSkeleton />
  }

  if (isError || !question) {
    return (
      <QuestionDetailLayout>
        <div
          className="flex min-h-[60vh] flex-col items-center justify-center text-center"
          role="alert"
        >
          <AlertCircle className="size-7 text-destructive" aria-hidden="true" />
          <h1 className="mt-4 text-lg font-semibold">题目加载失败</h1>
          <p className="mt-1 max-w-md text-sm text-muted-foreground">
            {error?.message || '暂时无法获取题目详情，请稍后重试。'}
          </p>
          <Button variant="outline" className="mt-5" onClick={() => void refetch()}>
            重新加载
          </Button>
        </div>
      </QuestionDetailLayout>
    )
  }

  const bankId = question.bank_list[0] ?? ''

  return (
    <QuestionDetailLayout bankId={bankId}>
      <div className="grid gap-8 lg:grid-cols-[17rem_minmax(0,1fr)] lg:gap-10">
        {bankId ? (
          <div className="lg:sticky lg:top-6 lg:self-start">
            <BankQuestionSidebar bankId={bankId} activeQuestionId={question.question_id} />
          </div>
        ) : null}

        <article className="min-w-0 lg:border-l lg:border-border lg:pl-10">
          <p className="font-mono text-xs text-muted-foreground">{question.question_id}</p>
          <h1 className="mt-3 text-2xl font-semibold leading-tight sm:text-3xl">
            {question.title}
          </h1>

          <div className="mt-5 flex flex-wrap items-center gap-2">
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

          <div className="mt-8 border-t border-border pt-4 text-base">
            <DetailContent content={question.content} />
          </div>

          <div className="mt-10">
            <Comments targetType={TargetType.QUESTION} targetId={questionId} />
          </div>
        </article>
      </div>
    </QuestionDetailLayout>
  )
}

function QuestionDetailLayout({
  bankId,
  children,
}: {
  bankId?: string
  children: React.ReactNode
}) {
  const listUrl = bankId ? `/questions?bank_id=${encodeURIComponent(bankId)}` : '/questions'

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Header
        action={
          <Link
            to={listUrl}
            className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground hover:underline"
          >
            <ArrowLeft className="size-4" aria-hidden="true" />
            返回题目列表
          </Link>
        }
      />
      <main className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6">{children}</main>
    </div>
  )
}

function QuestionNotFound() {
  return (
    <QuestionDetailLayout>
      <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
        <p className="font-mono text-sm text-muted-foreground">404</p>
        <h1 className="mt-3 text-xl font-semibold">题目不存在</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          该题目可能已被删除，或链接中的题目编号不正确。
        </p>
        <Button asChild className="mt-5">
          <Link to="/questions">返回题目列表</Link>
        </Button>
      </div>
    </QuestionDetailLayout>
  )
}

function QuestionDetailSkeleton() {
  return (
    <QuestionDetailLayout>
      <div
        className="grid gap-8 lg:grid-cols-[17rem_minmax(0,1fr)] lg:gap-10"
        aria-busy="true"
        aria-label="题目详情加载中"
      >
        <div className="space-y-4 border-b border-border pb-8 lg:border-b-0">
          <div className="h-4 w-20 animate-pulse rounded-sm bg-muted" />
          <div className="h-9 w-full animate-pulse rounded-md bg-muted" />
          {[0, 1, 2, 3, 4].map(item => (
            <div key={item} className="h-12 w-full animate-pulse rounded-sm bg-muted" />
          ))}
        </div>
        <div className="lg:border-l lg:border-border lg:pl-10">
          <div className="h-3 w-24 animate-pulse rounded-sm bg-muted" />
          <div className="mt-4 h-8 w-4/5 animate-pulse rounded-sm bg-muted" />
          <div className="mt-5 h-6 w-48 animate-pulse rounded-sm bg-muted" />
          <div className="mt-10 space-y-3">
            <div className="h-4 w-full animate-pulse rounded-sm bg-muted" />
            <div className="h-4 w-full animate-pulse rounded-sm bg-muted" />
            <div className="h-4 w-3/4 animate-pulse rounded-sm bg-muted" />
          </div>
        </div>
      </div>
    </QuestionDetailLayout>
  )
}

function getDifficultyLabel(difficulty: number): string {
  return ['', '简单', '中等', '困难'][difficulty] ?? `难度 ${difficulty}`
}
