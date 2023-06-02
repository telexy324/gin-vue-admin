import Socket from '@/utils/Socket'
import { store } from '@/store'

const socket = new Socket(() => {
  // const baseURI = `ws${document.baseURI.substr(4)}`
  const token = store.getters['user/token']
  return new WebSocket('ws://' + location.hostname + ':' + import.meta.env.VITE_WS_PORT + '/task/ws' + '?x-token=' + token)
})

export default socket
