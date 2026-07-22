import { ChevronLeft, ChevronRight } from 'lucide-react'

import { Button } from '@/components/ui/button'

interface QuestionPaginationProps {
  page: number
  totalPages: number
  onChange: (page: number) => void
}

export function QuestionPagination({ page, totalPages, onChange }: QuestionPaginationProps) {
  return (
    <nav aria-label="题目分页" className="mt-6 flex items-center justify-between gap-3">
      <span className="text-sm text-muted-foreground">
        第 {page} / {totalPages} 页
      </span>
      <div className="flex items-center gap-1">
        <PageArrow label="上一页" disabled={page <= 1} onClick={() => onChange(page - 1)}>
          <ChevronLeft aria-hidden="true" />
        </PageArrow>
        {getVisiblePages(page, totalPages).map((item, index) =>
          typeof item === 'number' ? (
            <Button
              key={item}
              type="button"
              variant={item === page ? 'default' : 'outline'}
              size="icon"
              className={item === page ? undefined : 'hidden sm:inline-flex'}
              aria-label={`第 ${item} 页`}
              aria-current={item === page ? 'page' : undefined}
              onClick={() => onChange(item)}
            >
              {item}
            </Button>
          ) : (
            <span
              key={`${item}-${index}`}
              className="hidden size-9 items-center justify-center text-sm text-muted-foreground sm:flex"
              aria-hidden="true"
            >
              ...
            </span>
          )
        )}
        <PageArrow label="下一页" disabled={page >= totalPages} onClick={() => onChange(page + 1)}>
          <ChevronRight aria-hidden="true" />
        </PageArrow>
      </div>
    </nav>
  )
}

function PageArrow({
  label,
  disabled,
  onClick,
  children,
}: {
  label: string
  disabled: boolean
  onClick: () => void
  children: React.ReactNode
}) {
  return (
    <Button
      type="button"
      variant="outline"
      size="icon"
      disabled={disabled}
      aria-label={label}
      title={label}
      onClick={onClick}
    >
      {children}
    </Button>
  )
}

function getVisiblePages(current: number, total: number): Array<number | 'start' | 'end'> {
  if (total <= 7) return Array.from({ length: total }, (_, index) => index + 1)

  const items: Array<number | 'start' | 'end'> = [1]
  if (current > 4) items.push('start')
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)
  for (let page = start; page <= end; page += 1) items.push(page)
  if (current < total - 3) items.push('end')
  items.push(total)
  return items
}
