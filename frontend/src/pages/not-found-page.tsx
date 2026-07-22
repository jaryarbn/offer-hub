import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'

export function NotFoundPage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-background px-6 text-center text-foreground">
      <p className="font-mono text-sm text-muted-foreground">404</p>
      <h1 className="mt-3 text-xl font-semibold">页面不存在</h1>
      <Button asChild className="mt-5">
        <Link to="/questions">返回题目列表</Link>
      </Button>
    </main>
  )
}
