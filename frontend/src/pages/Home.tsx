import {
  ArrowRight,
  BookOpenCheck,
  CheckCircle2,
  FileCheck2,
  MessageCircleMore,
  Target,
} from 'lucide-react'
import { Link, useNavigate } from 'react-router-dom'

import { Header } from '@/components/Header'
import { HotQuestions } from '@/components/HotContent'
import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'

const stats = [
  { value: '200+', label: '面试题目' },
  { value: '1,000+', label: '学习用户' },
  { value: '16', label: '精选题库' },
]

const features = [
  {
    icon: BookOpenCheck,
    title: '海量真题',
    description: '覆盖后端、前端与 AI 等热门方向，按题库系统练习。',
  },
  {
    icon: FileCheck2,
    title: '精选解析',
    description: '从核心思路到延伸追问，帮助你建立完整的知识脉络。',
  },
  {
    icon: Target,
    title: '掌握状态追踪',
    description: '标记已掌握、待复习和薄弱题目，让每次练习更有重点。',
  },
  {
    icon: MessageCircleMore,
    title: '交流讨论',
    description: '围绕每道题分享思路，在讨论中补齐理解盲区。',
  },
]

export function HomePage() {
  const navigate = useNavigate()
  const { isLoggedIn, setShowLoginDialog } = useLogin()

  const startPracticing = () => {
    navigate('/questions-collection')

    if (!isLoggedIn) {
      setShowLoginDialog(true)
    }
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Header />
      <main>
        <HeroSection />
        <HotContentSection />
        <FeaturesSection />
        <CTASection onStart={startPracticing} />
      </main>
      <FooterSection />
    </div>
  )
}

function HeroSection() {
  return (
    <section className="relative isolate flex min-h-120 items-end overflow-hidden bg-primary text-primary-foreground md:min-h-136">
      <img
        src="/images/home-hero.jpg"
        alt="摆放着代码编辑器与学习资料的开发者工作台"
        className="absolute inset-0 -z-20 size-full object-cover object-center"
      />
      <div className="absolute inset-0 -z-10 bg-black/60" aria-hidden="true" />

      <div className="mx-auto w-full max-w-6xl px-4 pb-8 pt-20 sm:px-6 sm:pb-10">
        <div className="max-w-2xl">
          <p className="text-sm font-medium text-white/75">Offer Hub 技术面试题库</p>
          <h1 className="mt-4 text-4xl font-semibold leading-tight text-white sm:text-5xl">
            成为技术面试高手
          </h1>
          <p className="mt-5 max-w-xl text-base leading-7 text-white/80 sm:text-lg">
            用体系化题库、精选解析和学习进度追踪，拆解面试重点，把零散知识变成真正可表达的能力。
          </p>
          <div className="mt-7 flex flex-wrap gap-3">
            <Button asChild className="h-11 bg-white px-5 text-primary hover:bg-white/90">
              <Link to="/questions-collection">
                开始刷题
                <ArrowRight aria-hidden="true" />
              </Link>
            </Button>
            <Button
              asChild
              variant="outline"
              className="h-11 border-white/40 bg-black/10 px-5 text-white hover:bg-white/10 hover:text-white"
            >
              <Link to="/hot">查看热门榜单</Link>
            </Button>
          </div>
        </div>

        <StatsSection />
      </div>
    </section>
  )
}

function StatsSection() {
  return (
    <dl className="mt-10 grid max-w-2xl grid-cols-3 border-t border-white/25 pt-5 sm:mt-12">
      {stats.map((stat, index) => (
        <div key={stat.label} className={index > 0 ? 'border-l border-white/25 pl-5' : ''}>
          <dt className="text-xs text-white/65 sm:text-sm">{stat.label}</dt>
          <dd className="mt-1 text-xl font-semibold text-white sm:text-2xl">{stat.value}</dd>
        </div>
      ))}
    </dl>
  )
}

function HotContentSection() {
  return (
    <section
      className="border-b border-border bg-background py-14 sm:py-16"
      aria-labelledby="hot-title"
    >
      <div className="mx-auto max-w-5xl px-4 sm:px-6">
        <div className="flex items-end justify-between gap-4 border-b border-border pb-5">
          <div>
            <p className="text-sm text-muted-foreground">近期高热度题目</p>
            <h2 id="hot-title" className="mt-1 text-2xl font-semibold">
              热门榜单
            </h2>
          </div>
          <Link
            to="/hot"
            className="shrink-0 text-sm font-medium text-muted-foreground outline-none hover:text-foreground hover:underline focus-visible:ring-2 focus-visible:ring-ring"
          >
            查看完整榜单
          </Link>
        </div>
        <HotQuestions limit={6} showHeader={false} variant="section" />
      </div>
    </section>
  )
}

function FeaturesSection() {
  return (
    <section className="bg-muted/35 py-14 sm:py-16" aria-labelledby="features-title">
      <div className="mx-auto max-w-6xl px-4 sm:px-6">
        <div className="max-w-xl">
          <p className="text-sm text-muted-foreground">高效准备每一次面试</p>
          <h2 id="features-title" className="mt-1 text-2xl font-semibold">
            从练习到复盘，形成完整闭环
          </h2>
        </div>
        <ul className="mt-8 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
          {features.map(feature => {
            const Icon = feature.icon

            return (
              <li key={feature.title} className="rounded-md border border-border bg-card p-5">
                <span className="flex size-9 items-center justify-center rounded-md bg-secondary text-secondary-foreground">
                  <Icon className="size-4" aria-hidden="true" />
                </span>
                <h3 className="mt-5 text-sm font-semibold">{feature.title}</h3>
                <p className="mt-2 text-sm leading-6 text-muted-foreground">
                  {feature.description}
                </p>
              </li>
            )
          })}
        </ul>
      </div>
    </section>
  )
}

function CTASection({ onStart }: { onStart: () => void }) {
  return (
    <section className="bg-primary py-12 text-primary-foreground" aria-labelledby="cta-title">
      <div className="mx-auto flex max-w-5xl flex-col gap-6 px-4 sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <div>
          <h2 id="cta-title" className="text-2xl font-semibold">
            准备好开始了吗？
          </h2>
          <p className="mt-2 text-sm leading-6 text-primary-foreground/70">
            选择你的技术方向，从第一道题开始建立面试信心。
          </p>
        </div>
        <Button
          type="button"
          className="h-11 w-full bg-white px-5 text-primary hover:bg-white/90 sm:w-auto"
          onClick={onStart}
        >
          <CheckCircle2 aria-hidden="true" />
          立即开始刷题
        </Button>
      </div>
    </section>
  )
}

function FooterSection() {
  return (
    <footer className="border-t border-border bg-background">
      <div className="mx-auto flex max-w-6xl flex-col gap-2 px-4 py-6 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <p>© 2026 Offer Hub. 保留所有权利。</p>
        <p>专注技术面试练习与知识复盘</p>
      </div>
    </footer>
  )
}
