import axios, { type AxiosRequestConfig } from 'axios'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/useAuthStore'

export type ApiResponse<T = any> = {
  code: number
  data?: T
  msg?: string
  success?: boolean
  [key: string]: any
}

type HttpClient = {
  <T = any>(config: AxiosRequestConfig): Promise<ApiResponse<T>>
}

const instance = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 60000
})

instance.interceptors.request.use(
  (config) => {
    const { token, userInfo } = useAuthStore.getState()
    const headers = {
      'Content-Type': 'application/json',
      'x-token': token || '',
      'x-user-id': userInfo?.ID || ''
    }
    ;(config as any).headers = { ...((config as any).headers || {}), ...headers }
    if (config.data && typeof config.data !== 'string') {
      config.data = JSON.stringify(config.data)
    }
    return config
  },
  (error) => Promise.reject(error)
)

instance.interceptors.response.use(
  (response) => {
    const newToken = response.headers['new-token']
    if (newToken) {
      useAuthStore.getState().setToken(newToken)
    }

    if (response.data?.code === 0 || response.headers.success === 'true') {
      return response.data
    }

    if (response.data?.msg) {
      toast.error(response.data.msg)
    }

    if (response.data?.data?.reload) {
      useAuthStore.getState().clearAuth()
      window.location.hash = '#/login'
    }

    return response.data
  },
  (error) => {
    const message = error?.response?.data?.msg || error?.message || '网络请求失败'
    toast.error(message)
    if (error?.response?.status === 401) {
      useAuthStore.getState().clearAuth()
      window.location.hash = '#/login'
    }
    return Promise.reject(error)
  }
)

const service = ((config: AxiosRequestConfig) => instance(config)) as HttpClient

export default service
