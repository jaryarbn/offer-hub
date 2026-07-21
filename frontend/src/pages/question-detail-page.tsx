import { Link, useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"

export function QuestionDetailPage() {
  const { id } = useParams<{ id: string }>()

  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-background px-6 text-center text-foreground">
      <p className="font-mono text-sm text-muted-foreground">{id}</p>
      <h1 className="mt-3 text-2xl font-semibold">题目详情</h1>
      <Button asChild variant="outline" className="mt-6">
        <Link to="/questions">返回题库</Link>
      </Button>
    </main>
  )
}
