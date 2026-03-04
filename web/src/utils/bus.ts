type Handler = (payload?: any) => void

class TinyEmitter {
  private map = new Map<string, Set<Handler>>()

  on(event: string, handler: Handler) {
    const set = this.map.get(event) || new Set<Handler>()
    set.add(handler)
    this.map.set(event, set)
  }

  off(event: string, handler: Handler) {
    this.map.get(event)?.delete(handler)
  }

  emit(event: string, payload?: any) {
    this.map.get(event)?.forEach((handler) => handler(payload))
  }
}

export const emitter = new TinyEmitter()

