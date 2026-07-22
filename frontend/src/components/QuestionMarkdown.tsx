import type { ComponentPropsWithoutRef } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

interface QuestionMarkdownProps {
  content: string
}

export function QuestionMarkdown({ content }: QuestionMarkdownProps) {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      components={{
        h1: props => <h2 className="mb-4 mt-8 text-xl font-semibold" {...props} />,
        h2: props => <h2 className="mb-3 mt-7 text-lg font-semibold" {...props} />,
        h3: props => <h3 className="mb-2 mt-6 text-base font-semibold" {...props} />,
        p: props => <p className="my-4 leading-8 text-foreground/90" {...props} />,
        ul: props => <ul className="my-4 list-disc space-y-2 pl-6" {...props} />,
        ol: props => <ol className="my-4 list-decimal space-y-2 pl-6" {...props} />,
        li: props => <li className="pl-1 leading-7" {...props} />,
        blockquote: props => (
          <blockquote
            className="my-5 border-l-2 border-foreground/30 pl-4 text-muted-foreground"
            {...props}
          />
        ),
        a: ({ href, ...props }) => (
          <a
            href={href}
            className="font-medium underline underline-offset-4"
            target={href?.startsWith('http') ? '_blank' : undefined}
            rel={href?.startsWith('http') ? 'noreferrer' : undefined}
            {...props}
          />
        ),
        code: MarkdownCode,
        pre: props => (
          <pre
            className="my-5 overflow-x-auto rounded-md bg-primary p-4 text-sm leading-6 text-primary-foreground"
            {...props}
          />
        ),
        table: props => (
          <div className="my-5 overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm" {...props} />
          </div>
        ),
        th: props => (
          <th className="border border-border bg-muted px-3 py-2 font-semibold" {...props} />
        ),
        td: props => <td className="border border-border px-3 py-2 align-top" {...props} />,
        hr: props => <hr className="my-8 border-border" {...props} />,
      }}
    >
      {content}
    </ReactMarkdown>
  )
}

function MarkdownCode({ className, children, ...props }: ComponentPropsWithoutRef<'code'>) {
  const isBlock = className?.startsWith('language-')

  return (
    <code
      className={
        isBlock
          ? className
          : 'rounded-sm bg-muted px-1.5 py-0.5 font-mono text-[0.9em] text-foreground'
      }
      {...props}
    >
      {children}
    </code>
  )
}
