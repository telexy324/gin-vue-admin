import Listenable from '@/utils/Listenable'

export default class Socket<T = any> extends Listenable<T> {
  private websocketCreator: () => WebSocket
  private ws: WebSocket | null = null

  constructor(websocketCreator: () => WebSocket) {
    super()
    this.websocketCreator = websocketCreator
  }

  start() {
    if (this.ws) {
      throw new Error('Websocket already started. Please stop it before starting.')
    }
    this.ws = this.websocketCreator()
    this.ws.onclose = () => {
      if (!this.isRunning()) return
      this.ws = null
      setTimeout(() => this.start(), 2000)
    }
    this.ws.onmessage = ({ data }) => {
      try {
        this.callListeners(JSON.parse(String(data)))
      } catch (_) {
        this.callListeners(data as any)
      }
    }
  }

  isRunning() {
    return this.ws != null
  }

  stop() {
    if (!this.ws) return
    this.ws.close()
    this.ws = null
  }
}

