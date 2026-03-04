import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

export type UserInfo = {
  ID?: number
  uuid?: string
  nickName?: string
  headerImg?: string
  authority?: {
    defaultRouter?: string
    authorityName?: string
  }
  authorityId?: string | number
  authorities?: Array<{
    authorityId: string
    authorityName: string
  }>
  [key: string]: any
}

type AuthState = {
  token: string
  userInfo: UserInfo | null
  setToken: (token: string) => void
  setUserInfo: (userInfo: UserInfo | null) => void
  clearAuth: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: '',
      userInfo: null,
      setToken: (token: string) => set({ token }),
      setUserInfo: (userInfo: UserInfo | null) => set({ userInfo }),
      clearAuth: () => set({ token: '', userInfo: null })
    }),
    {
      name: 'gva-auth',
      storage: createJSONStorage(() => localStorage)
    }
  )
)
