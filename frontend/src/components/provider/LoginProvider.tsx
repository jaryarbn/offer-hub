import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react'
import { useQueryClient } from '@tanstack/react-query'

import { queryKeys } from '@/hooks/useQuestionQueries'
import {
  getUserInfo,
  mapPasswordAuthUserToUserInfo,
  type PasswordAuthUserDTO,
  type UserInfo,
} from '@/services/auth'

const TOKEN_KEY = 'token'

export interface LoginContextValue {
  isLoggedIn: boolean
  userInfo: UserInfo | null
  login: (userInfo: PasswordAuthUserDTO, token: string) => void
  logout: () => void
  showLoginDialog: boolean
  setShowLoginDialog: (show: boolean) => void
  refreshUserInfo: () => Promise<void>
}

interface LoginProviderProps {
  children: ReactNode
}

export const LoginContext = createContext<LoginContextValue | undefined>(undefined)

export function LoginProvider({ children }: LoginProviderProps) {
  const queryClient = useQueryClient()
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_KEY))
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null)
  const [showLoginDialog, setShowLoginDialogState] = useState(false)

  // 用版本号阻止登出后仍在执行的旧请求重新写入登录状态。
  const sessionVersionRef = useRef(0)
  const initializedRef = useRef(false)
  const mountedRef = useRef(true)

  const login = useCallback(
    (passwordAuthUser: PasswordAuthUserDTO, nextToken: string) => {
      sessionVersionRef.current += 1
      localStorage.setItem(TOKEN_KEY, nextToken)
      setToken(nextToken)
      setUserInfo(mapPasswordAuthUserToUserInfo(passwordAuthUser))

      // token 已写入后再刷新详情，确保重新请求携带登录凭证。
      void queryClient.invalidateQueries({ queryKey: queryKeys.details })
    },
    [queryClient]
  )

  const clearSession = useCallback(() => {
    sessionVersionRef.current += 1
    localStorage.removeItem(TOKEN_KEY)

    if (mountedRef.current) {
      setToken(null)
      setUserInfo(null)
    }

    // 详情响应包含登录用户的点赞和掌握状态，清空后由活跃查询以游客身份重新拉取。
    void queryClient.resetQueries({ queryKey: queryKeys.details })
  }, [queryClient])

  const logout = clearSession

  const setShowLoginDialog = useCallback((show: boolean) => {
    setShowLoginDialogState(show)
  }, [])

  const refreshUserInfo = useCallback(async () => {
    const currentToken = localStorage.getItem(TOKEN_KEY)

    if (!currentToken) {
      clearSession()
      return
    }

    const requestVersion = sessionVersionRef.current

    try {
      const refreshedUserInfo = await getUserInfo()

      if (
        mountedRef.current &&
        sessionVersionRef.current === requestVersion &&
        localStorage.getItem(TOKEN_KEY) === currentToken
      ) {
        setToken(currentToken)
        setUserInfo(refreshedUserInfo)
      }
    } catch (error) {
      if (
        mountedRef.current &&
        sessionVersionRef.current === requestVersion &&
        localStorage.getItem(TOKEN_KEY) === currentToken
      ) {
        clearSession()
      }

      throw error
    }
  }, [clearSession])

  useEffect(() => {
    mountedRef.current = true

    return () => {
      mountedRef.current = false
    }
  }, [])

  useEffect(() => {
    if (initializedRef.current) {
      return
    }
    initializedRef.current = true

    if (localStorage.getItem(TOKEN_KEY)) {
      // 初始化失败时 refreshUserInfo 已清理无效登录态，无需产生未处理的 Promise。
      void refreshUserInfo().catch(() => undefined)
    }
  }, [refreshUserInfo])

  const contextValue = useMemo<LoginContextValue>(
    () => ({
      isLoggedIn: Boolean(token && userInfo),
      userInfo,
      login,
      logout,
      showLoginDialog,
      setShowLoginDialog,
      refreshUserInfo,
    }),
    [token, userInfo, login, logout, showLoginDialog, setShowLoginDialog, refreshUserInfo]
  )

  return <LoginContext.Provider value={contextValue}>{children}</LoginContext.Provider>
}

export function useLogin(): LoginContextValue {
  const context = useContext(LoginContext)

  if (!context) {
    throw new Error('useLogin 必须在 LoginProvider 内使用')
  }

  return context
}

export default LoginProvider
