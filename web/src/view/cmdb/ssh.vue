<template>
  <div class="term1">
    <div ref="terminalBox" style="height: 60vh;"></div>
  </div>
</template>

<script>
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { ref } from 'vue'

export default {
  name: 'SSH',
  props: {
    username: String,
    ipaddress: String,
    port: Number,
    password: String,
  },
  data() {
    return {
      terminalBox: ref(null),
      term: null,
      socket: null
    }
  },
  mounted() {
    this.initSocket()
  },
  beforeDestroy() {
    this.socket.close()
    this.term.dispose()
  },
  methods: {
    initTerm() {
      const term = new Terminal({
        rendererType: 'canvas',
        cursorBlink: true,
        cursorStyle: 'bar'
      })
      const fitAddon = new FitAddon()
      term.loadAddon(fitAddon)
      term.open(this.terminalBox.value)
      fitAddon.fit()
      this.term = term
      this.term.write('正在连接...\r\n')
    },
    initSocket() {
      this.socket = new WebSocket('ws://' + location.hostname + ':8080/ssh/run')
      this.socket.binaryType = 'arraybuffer'
      this.socketOnClose()
      this.socketOnOpen()
      this.socketOnError()
    },
    socketOnOpen() {
      this.socket.onopen = () => {
        // 链接成功后
        this.initTerm()
      }
    },
    socketOnClose() {
      this.socket.onclose = () => {
        // console.log('close socket')
      }
    },
    socketOnError() {
      this.socket.onerror = () => {
        // console.log('socket 链接失败')
      }
    }
  }
}
</script>

<style scoped>

</style>