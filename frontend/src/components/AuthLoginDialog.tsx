import { BookOpenCheck, ChartNoAxesCombined, ListChecks } from 'lucide-react'

import { UsernamePasswordAuthForms } from '@/components/auth/UsernamePasswordAuthForms'
import { useLogin } from '@/components/provider/LoginProvider'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const platformFeatures = [
  {
    icon: BookOpenCheck,
    title: '精选面试题库',
    description: '按技术方向和岗位系统整理',
  },
  {
    icon: ListChecks,
    title: '结构化练习',
    description: '从基础到进阶逐步巩固知识',
  },
  {
    icon: ChartNoAxesCombined,
    title: '持续记录成长',
    description: '登录后保存练习与学习进度',
  },
]

export function AuthLoginDialog() {
  const { showLoginDialog, setShowLoginDialog } = useLogin()

  return (
    <Dialog open={showLoginDialog} onOpenChange={setShowLoginDialog}>
      <DialogContent className="max-h-[calc(100vh-2rem)] max-w-md gap-0 overflow-y-auto border-0 p-0 md:max-w-4xl">
        <div className="grid min-h-[34rem] md:grid-cols-[0.9fr_1.1fr]">
          <aside className="hidden bg-gradient-to-br from-emerald-700 via-emerald-600 to-teal-600 p-10 text-white md:flex md:flex-col">
            <div>
              <p className="text-sm font-medium text-emerald-100">Offer Hub</p>
              <h2 className="mt-3 text-3xl font-semibold tracking-tight">让每次练习都有收获</h2>
              <p className="mt-3 max-w-sm text-sm leading-6 text-emerald-50/90">
                聚合高质量技术题目，帮助你更有节奏地准备下一场面试。
              </p>
            </div>

            <ul className="mt-10 space-y-6" aria-label="平台特性">
              {platformFeatures.map(feature => {
                const Icon = feature.icon
                return (
                  <li key={feature.title} className="flex gap-3">
                    <span className="flex size-9 shrink-0 items-center justify-center rounded-md bg-white/15">
                      <Icon className="size-4" aria-hidden="true" />
                    </span>
                    <div>
                      <p className="text-sm font-semibold">{feature.title}</p>
                      <p className="mt-1 text-sm text-emerald-50/80">{feature.description}</p>
                    </div>
                  </li>
                )
              })}
            </ul>

            <p className="mt-auto pt-10 text-xs text-emerald-100/80">
              为技术学习者打造的题库与成长平台
            </p>
          </aside>

          <section className="flex flex-col justify-center bg-background px-6 py-10 sm:px-10 md:px-12">
            <DialogHeader className="mb-7 text-left">
              <DialogTitle className="text-2xl">欢迎使用 Offer Hub</DialogTitle>
              <DialogDescription>登录或创建账号，开始你的面试准备。</DialogDescription>
            </DialogHeader>
            <UsernamePasswordAuthForms />
          </section>
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default AuthLoginDialog
