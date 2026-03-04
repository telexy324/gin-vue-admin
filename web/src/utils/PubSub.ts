import Listenable from '@/utils/Listenable'

export default class PubSub<TopicData = any> {
  private topics: Record<string, Listenable<TopicData>> = {}

  subscribe(topic: string, callback: (data: TopicData) => void) {
    if (!this.topics[topic]) this.topics[topic] = new Listenable<TopicData>()
    return this.topics[topic].addListener(callback)
  }

  unsubscribe(id: symbol) {
    Object.keys(this.topics).forEach((topic) => {
      this.topics[topic].removeListener(id)
    })
  }

  publish(topic: string, data: TopicData) {
    this.topics[topic]?.callListeners(data)
  }
}

