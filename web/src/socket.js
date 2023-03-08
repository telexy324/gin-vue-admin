import Socket from '@/utils/Socket'
import { store } from '@/store/index'

const socket = new Socket(() => {
  // const baseURI = `ws${document.baseURI.substr(4)}`
  const token = store.getters['user/token']
  return new WebSocket('ws://' + location.hostname + ':8888/task/ws' + '?x-token=' + token)
})

export default socket
