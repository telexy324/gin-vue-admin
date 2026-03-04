export default class Listenable<T = any> {
  private listeners: Record<symbol, (data: T) => void> = {}

  addListener(callback: (data: T) => void) {
    const id = Symbol('listener')
    this.listeners[id] = callback
    return id
  }

  removeListener(id: symbol) {
    if (this.listeners[id] == null) return false
    delete this.listeners[id]
    return true
  }

  callListeners(data: T) {
    Object.getOwnPropertySymbols(this.listeners).forEach((id) => {
      this.listeners[id]?.(data)
    })
  }

  hasListeners() {
    return Object.getOwnPropertySymbols(this.listeners).length > 0
  }
}

