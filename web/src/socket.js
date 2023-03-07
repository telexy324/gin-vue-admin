import Socket from '@/utils/Socket'

const socket = new Socket(() => {
  // const baseURI = `ws${document.baseURI.substr(4)}`
  return new WebSocket('ws://' + location.hostname + ':8888/task/ws')
})

export default socket
