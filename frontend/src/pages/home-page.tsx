import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'

export function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-background px-6 text-center text-foreground">
      <p className="text-sm text-muted-foreground">Offer Hub</p>
      <h1 className="mt-3 text-2xl font-semibold">首页</h1>
      <Button asChild className="mt-6">
        <Link to="/questions-collection">进入题库</Link>
      </Button>
    </main>
  )
}
