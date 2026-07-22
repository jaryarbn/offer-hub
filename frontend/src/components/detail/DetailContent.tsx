import { QuestionMarkdown } from '@/components/QuestionMarkdown'
import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'

interface DetailContentProps {
  content: string
}

export function DetailContent({ content }: DetailContentProps) {
  const { isLoggedIn, setShowLoginDialog } = useLogin()

  return (
    <div className="relative">
      <QuestionMarkdown content={content} />

      {!isLoggedIn ? (
        <div className="pointer-events-none absolute inset-x-0 bottom-0 flex h-[200px] max-h-full items-end justify-center bg-gradient-to-b from-transparent via-background/90 to-background pb-6">
          <Button
            type="button"
            className="pointer-events-auto shadow-sm"
            onClick={() => setShowLoginDialog(true)}
          >
            登录后查看完整内容
          </Button>
        </div>
      ) : null}
    </div>
  )
}
