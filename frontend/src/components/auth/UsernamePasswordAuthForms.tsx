import { useMemo, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { LoaderCircle } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

import { useLogin } from '@/components/provider/LoginProvider'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { postAuthLogin, postAuthRegister } from '@/services/auth'
import { cn } from '@/lib/utils'

type AuthMode = 'login' | 'register'

const baseAuthSchema = z.object({
  username: z.string().trim().min(1, '请输入用户名').max(50, '用户名不能超过 50 个字符'),
  password: z.string().min(6, '密码至少需要 6 位'),
  confirmPassword: z.string().optional(),
})

type AuthFormValues = z.infer<typeof baseAuthSchema>

function createAuthSchema(mode: AuthMode) {
  return baseAuthSchema.superRefine((values, context) => {
    if (mode !== 'register') {
      return
    }

    if (!values.confirmPassword) {
      context.addIssue({
        code: 'custom',
        message: '请再次输入密码',
        path: ['confirmPassword'],
      })
      return
    }

    if (values.password !== values.confirmPassword) {
      context.addIssue({
        code: 'custom',
        message: '两次输入的密码不一致',
        path: ['confirmPassword'],
      })
    }
  })
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error && error.message.trim()) {
    return error.message
  }

  // axios 拦截器可能直接抛出后端的统一错误响应对象。
  if (typeof error === 'object' && error !== null) {
    const candidate = error as { msg?: unknown }
    if (typeof candidate.msg === 'string' && candidate.msg.trim()) {
      return candidate.msg
    }
  }

  return '请求失败，请稍后重试'
}

export function UsernamePasswordAuthForms() {
  const [mode, setMode] = useState<AuthMode>('login')
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const { login, setShowLoginDialog } = useLogin()
  const authSchema = useMemo(() => createAuthSchema(mode), [mode])

  const {
    register,
    handleSubmit,
    reset,
    getValues,
    formState: { errors, isSubmitting },
  } = useForm<AuthFormValues>({
    resolver: zodResolver(authSchema),
    defaultValues: {
      username: '',
      password: '',
      confirmPassword: '',
    },
  })

  const switchMode = (nextMode: AuthMode) => {
    if (nextMode === mode || isSubmitting) {
      return
    }

    const username = getValues('username')
    setMode(nextMode)
    setSubmitError(null)
    setSuccessMessage(null)
    reset({ username, password: '', confirmPassword: '' })
  }

  const onSubmit = handleSubmit(async values => {
    setSubmitError(null)
    setSuccessMessage(null)

    try {
      const username = values.username.trim()

      if (mode === 'login') {
        const result = await postAuthLogin(username, values.password)
        login(result.userInfo, result.token)
        reset()
        setShowLoginDialog(false)
        return
      }

      await postAuthRegister(username, values.password)
      setMode('login')
      reset({ username, password: '', confirmPassword: '' })
      setSuccessMessage('注册成功，请使用新账号登录')
    } catch (error) {
      setSubmitError(getErrorMessage(error))
    }
  })

  const isLoginMode = mode === 'login'

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 rounded-md bg-muted p-1" aria-label="登录方式">
        <button
          type="button"
          className={cn(
            'rounded-sm px-3 py-2 text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
            isLoginMode
              ? 'bg-background text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          )}
          aria-pressed={isLoginMode}
          disabled={isSubmitting}
          onClick={() => switchMode('login')}
        >
          登录
        </button>
        <button
          type="button"
          className={cn(
            'rounded-sm px-3 py-2 text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
            !isLoginMode
              ? 'bg-background text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          )}
          aria-pressed={!isLoginMode}
          disabled={isSubmitting}
          onClick={() => switchMode('register')}
        >
          注册
        </button>
      </div>

      <div>
        <h3 className="text-xl font-semibold tracking-tight">
          {isLoginMode ? '登录账号' : '创建账号'}
        </h3>
        <p className="mt-1 text-sm text-muted-foreground">
          {isLoginMode ? '继续你的面试题练习进度。' : '注册后即可保存练习记录与学习进度。'}
        </p>
      </div>

      <form className="space-y-4" noValidate onSubmit={onSubmit}>
        <div className="space-y-1.5">
          <label className="text-sm font-medium" htmlFor="auth-username">
            用户名
          </label>
          <Input
            id="auth-username"
            autoComplete="username"
            placeholder="请输入用户名"
            aria-invalid={Boolean(errors.username)}
            aria-describedby={errors.username ? 'auth-username-error' : undefined}
            disabled={isSubmitting}
            {...register('username')}
          />
          {errors.username && (
            <p id="auth-username-error" className="text-sm text-destructive" role="alert">
              {errors.username.message}
            </p>
          )}
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium" htmlFor="auth-password">
            密码
          </label>
          <Input
            id="auth-password"
            type="password"
            autoComplete={isLoginMode ? 'current-password' : 'new-password'}
            placeholder="请输入至少 6 位密码"
            aria-invalid={Boolean(errors.password)}
            aria-describedby={errors.password ? 'auth-password-error' : undefined}
            disabled={isSubmitting}
            {...register('password')}
          />
          {errors.password && (
            <p id="auth-password-error" className="text-sm text-destructive" role="alert">
              {errors.password.message}
            </p>
          )}
        </div>

        {!isLoginMode && (
          <div className="space-y-1.5">
            <label className="text-sm font-medium" htmlFor="auth-confirm-password">
              确认密码
            </label>
            <Input
              id="auth-confirm-password"
              type="password"
              autoComplete="new-password"
              placeholder="请再次输入密码"
              aria-invalid={Boolean(errors.confirmPassword)}
              aria-describedby={errors.confirmPassword ? 'auth-confirm-password-error' : undefined}
              disabled={isSubmitting}
              {...register('confirmPassword')}
            />
            {errors.confirmPassword && (
              <p id="auth-confirm-password-error" className="text-sm text-destructive" role="alert">
                {errors.confirmPassword.message}
              </p>
            )}
          </div>
        )}

        {submitError && (
          <div
            className="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive"
            role="alert"
          >
            {submitError}
          </div>
        )}

        {successMessage && (
          <div
            className="rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-800"
            role="status"
          >
            {successMessage}
          </div>
        )}

        <Button
          type="submit"
          className="w-full bg-emerald-600 text-white hover:bg-emerald-700"
          disabled={isSubmitting}
        >
          {isSubmitting && <LoaderCircle className="animate-spin" aria-hidden="true" />}
          {isSubmitting ? (isLoginMode ? '正在登录…' : '正在注册…') : isLoginMode ? '登录' : '注册'}
        </Button>
      </form>

      <p className="text-center text-sm text-muted-foreground">
        {isLoginMode ? '还没有账号？' : '已经有账号？'}
        <button
          type="button"
          className="ml-1 font-medium text-emerald-700 underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          disabled={isSubmitting}
          onClick={() => switchMode(isLoginMode ? 'register' : 'login')}
        >
          {isLoginMode ? '立即注册' : '返回登录'}
        </button>
      </p>
    </div>
  )
}

export default UsernamePasswordAuthForms
