import { useState, type ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'

import { AuthLoginDialog } from '@/components/AuthLoginDialog'
import { LoginProvider } from '@/components/provider/LoginProvider'

interface AppProviderProps {
  children: ReactNode
}

export function AppProvider({ children }: AppProviderProps) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            refetchOnWindowFocus: false,
            retry: 1,
            staleTime: 5 * 60 * 1_000,
          },
        },
      })
  )

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
        <LoginProvider>
          {children}
          <AuthLoginDialog />
        </LoginProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
