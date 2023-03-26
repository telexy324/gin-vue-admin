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
  data() {
    return {
      terminalBox: ref(null),
      term: null,
      socket: null,
      manageIp: '',
      username: '',
      password: '',
      sshPort: ''
    }
  },
  mounted() {
    const manageIp = this.$route.params.manageIp
    const username = this.$route.params.username
    const password = this.$route.params.password
    const sshPort = this.$route.params.sshPort
    this.manageIp = manageIp
    this.username = username
    this.password = password
    this.sshPort = sshPort
    this.initTerm()
    this.initSocket()
  },
  beforeDestroy() {
    this.socket.close()
    this.term.dispose()
  },
  methods: {
    initTerm() {
      const term = new Terminal({
        rendererType: 'canvas', cursorBlink: true, cursorStyle: 'bar'
      })
      const fitAddon = new FitAddon()
      term.loadAddon(fitAddon)
      console.log(this.manageIp)
      term.open(this.$refs.terminalBox)
      fitAddon.fit()
      this.term = term
      this.term.write('正在连接...\r\n')
    },
    initSocket() {
      this.socket = new WebSocket('ws://' + location.hostname + ':8888/ssh/run')
      this.socket.binaryType = 'arraybuffer'
      this.socketOnClose()
      this.socketOnOpen()
      this.socketOnError()
      this.socketOnMessage()
    },
    socketOnOpen() {
      const _this = this
      this.socket.onopen = () => {
        // this.initTerm()
        this.term.write('连接成功...\r\n')
        // fitAddon.fit()
        this.term.onData(function(data) {
          // socket.send(JSON.stringify({ type: "stdin", data: data }))
          // console.log(data)
          _this.socket.send(data)
          // console.log(data)
        })
        // ElMessage.success("会话成功连接！")
        var jsonStr = `{"manageIp":"${this.manageIp}", "sshPort":${this.sshPort}, "username":"${this.username}", "password":"${this.password}"}`
        console.log(jsonStr)
        var datMsg = window.btoa(jsonStr)
        // socket.send(JSON.stringify({ ip: ip.value, name: name.value, password: password.value }))
        this.socket.send(datMsg)
      }
    },
    socketOnClose() {
      this.socket.onclose = () => {
        this.term.writeln('连接关闭')
      }
    },
    socketOnError() {
      this.socket.onerror = err => {
        // console.log(err)
        this.term.writeln('读取数据异常：', err)
      }
    },
    socketOnMessage() {
      // 接收数据
      this.socket.onmessage = recv => {
        try {
          this.term.write(recv.data)
        } catch (e) {
          this.console.log('unsupport data', recv.data)
        }
      }
    },
  }
}
</script>

<style lang="scss" scoped>
.upload {
  min-height: 100px;
}

.term1 {
  margin-left: 60px;
}

.go_out {
  margin-left: -89%;
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
