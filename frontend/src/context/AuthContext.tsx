import { createContext, useContext, useState } from 'react'

type AuthContextType = {
  accessToken: string
  login: (accessToken: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [accessToken, setAccessToken] = useState('')

  function login(accessToken: string) {
    setAccessToken(accessToken)
  }

  function logout() {
    setAccessToken('')
  }

  return (
    <AuthContext.Provider
      value={{
        accessToken,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }

  return context
}