import { useState, type ReactNode } from 'react'
import { BookOpen, ChevronDown, LoaderCircle, LogOut } from 'lucide-react'
import { Link } from 'react-router-dom'

import { useLogin } from '@/components/provider/LoginProvider'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { postAuthLogout } from '@/services/auth'

export interface HeaderProps {
  sectionLabel?: string
  action?: ReactNode
}

export function Header({ sectionLabel, action }: HeaderProps) {
  const { isLoggedIn, userInfo, logout, setShowLoginDialog } = useLogin()
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const [logoutError, setLogoutError] = useState<string | null>(null)

  const displayName = userInfo?.nick_name.trim() || userInfo?.username || '用户'
  const avatarUrl = userInfo?.avatar_url?.trim() || userInfo?.avatar.trim() || undefined
  const avatarFallback = displayName.slice(0, 1).toUpperCase()

  const handleLogout = async () => {
    if (isLoggingOut) {
      return
    }

    setIsLoggingOut(true)
    setLogoutError(null)

    try {
      // 本地 token 必须在请求完成前保留，axios 拦截器才能添加 Authorization Header。
      await postAuthLogout()
    } catch {
      setLogoutError('服务器登出失败，本地登录状态已清除')
    } finally {
      logout()
      setIsLoggingOut(false)
    }
  }

  const openLoginDialog = () => {
    setLogoutError(null)
    setShowLoginDialog(true)
  }

  return (
    <header className="border-b border-border bg-background">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-4 px-4 sm:px-6">
        <div className="flex min-w-0 items-center gap-3">
          <Link
            to="/"
            className="flex shrink-0 items-center gap-3 rounded-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <span className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <BookOpen className="size-4" aria-hidden="true" />
            </span>
            <span className="hidden text-sm font-semibold sm:inline">Offer Hub</span>
          </Link>

          {sectionLabel && (
            <>
              <span className="hidden h-4 w-px bg-border sm:block" aria-hidden="true" />
              <span className="truncate text-sm text-muted-foreground">{sectionLabel}</span>
            </>
          )}
        </div>

        <div className="flex shrink-0 items-center gap-3">
          {action && (
            <>
              <div className="text-sm">{action}</div>
              <span className="h-4 w-px bg-border" aria-hidden="true" />
            </>
          )}

          {isLoggedIn && userInfo ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  className="flex items-center gap-2 rounded-md p-1 outline-none transition-colors hover:bg-muted focus-visible:ring-2 focus-visible:ring-ring"
                  aria-label={`${displayName}的用户菜单`}
                >
                  <Avatar className="size-8 border border-border">
                    <AvatarImage src={avatarUrl} alt={`${displayName}的头像`} />
                    <AvatarFallback>{avatarFallback}</AvatarFallback>
                  </Avatar>
                  <span className="hidden max-w-28 truncate text-sm font-medium sm:inline">
                    {displayName}
                  </span>
                  <ChevronDown
                    className="hidden size-4 text-muted-foreground sm:block"
                    aria-hidden="true"
                  />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel className="max-w-48 truncate">{displayName}</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  className="text-destructive focus:bg-destructive/10 focus:text-destructive"
                  disabled={isLoggingOut}
                  onSelect={() => void handleLogout()}
                >
                  {isLoggingOut ? (
                    <LoaderCircle className="animate-spin" aria-hidden="true" />
                  ) : (
                    <LogOut aria-hidden="true" />
                  )}
                  {isLoggingOut ? '正在退出…' : '退出登录'}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <div className="flex items-center gap-3">
              {logoutError && (
                <span className="hidden text-xs text-destructive md:inline" role="status">
                  {logoutError}
                </span>
              )}
              <Button type="button" size="sm" onClick={openLoginDialog}>
                登录
              </Button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header
