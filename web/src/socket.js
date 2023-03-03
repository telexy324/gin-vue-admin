import Socket from '@/utils/Socket'

const socket = new Socket(() => {
  const baseURI = `ws${document.baseURI.substr(4)}`
  return new WebSocket(`${baseURI}api/ws`)
})

export default socket
