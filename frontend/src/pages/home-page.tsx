import { Link } from 'react-router-dom'

import { Header } from '@/components/Header'
import { Button } from '@/components/ui/button'

export function HomePage() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <Header />
      <main className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center px-6 text-center">
        <p className="text-sm text-muted-foreground">Offer Hub</p>
        <h1 className="mt-3 text-2xl font-semibold">首页</h1>
        <Button asChild className="mt-6">
          <Link to="/questions-collection">进入题库</Link>
        </Button>
      </main>
    </div>
  )
}
