import Socket from '@/utils/Socket'
import { useAuthStore } from '@/stores/useAuthStore'

const socket = new Socket<any>(() => {
  const token = useAuthStore.getState().token || ''
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const wsPort = import.meta.env.VITE_WS_PORT
  const host = wsPort ? `${window.location.hostname}:${wsPort}` : window.location.host
  return new WebSocket(`${protocol}://${host}/task/ws?x-token=${encodeURIComponent(token)}`)
})

export default socket

