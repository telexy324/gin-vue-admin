import { useEffect, useMemo, useRef, useState } from 'react'
import { useLocation, useParams, useSearchParams } from 'react-router-dom'
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

type ConnForm = {
  manageIp: string
  sshPort: number
  username: string
  password: string
}

export default function SshPage() {
  const location = useLocation()
  const params = useParams()
  const [searchParams] = useSearchParams()
  const [form, setForm] = useState<ConnForm>({ manageIp: '', sshPort: 22, username: '', password: '' })
  const [connecting, setConnecting] = useState(false)
  const [connected, setConnected] = useState(false)
  const terminalBoxRef = useRef<HTMLDivElement | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const termRef = useRef<Terminal | null>(null)
  const fitRef = useRef<FitAddon | null>(null)
  const dataDisposerRef = useRef<{ dispose: () => void } | null>(null)
  const autoConnectedRef = useRef(false)

  const wsUrl = useMemo(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const wsPort = import.meta.env.VITE_WS_PORT
    const host = wsPort ? `${window.location.hostname}:${wsPort}` : window.location.host
    return `${protocol}://${host}/ssh/run`
  }, [])

  const writeLine = (msg: string) => {
    termRef.current?.writeln(msg)
  }

  const sendResize = () => {
    const ws = wsRef.current
    const term = termRef.current
    if (!ws || ws.readyState !== WebSocket.OPEN || !term) return
    ws.send(
      JSON.stringify({
        type: 'resize',
        cols: term.cols,
        rows: term.rows
      })
    )
  }

  const initTerm = () => {
    if (termRef.current || !terminalBoxRef.current) return
    const term = new Terminal({
      cursorBlink: true,
      cursorStyle: 'bar',
      scrollback: 1000,
      disableStdin: false,
      convertEol: true
    })
    const fitAddon = new FitAddon()
    term.loadAddon(fitAddon)
    term.open(terminalBoxRef.current)
    fitAddon.fit()
    term.write('正在连接...\r\n')
    termRef.current = term
    fitRef.current = fitAddon
  }

  const disconnect = () => {
    wsRef.current?.close()
    wsRef.current = null
    dataDisposerRef.current?.dispose()
    dataDisposerRef.current = null
    setConnected(false)
    setConnecting(false)
  }

  const connect = (input?: ConnForm) => {
    const conn = input || form
    if (!conn.manageIp || !conn.sshPort || !conn.username || !conn.password) {
      toast.error('请填写完整连接参数')
      return
    }

    disconnect()
    initTerm()
    setConnecting(true)

    const ws = new WebSocket(wsUrl)
    ws.binaryType = 'arraybuffer'
    wsRef.current = ws

    ws.onopen = () => {
      setConnecting(false)
      setConnected(true)
      writeLine('连接成功...')

      dataDisposerRef.current?.dispose()
      dataDisposerRef.current = termRef.current?.onData((data) => {
        if (ws.readyState === WebSocket.OPEN) ws.send(data)
      }) || null

      const authPayload = window.btoa(
        JSON.stringify({
          manageIp: conn.manageIp,
          sshPort: Number(conn.sshPort),
          username: conn.username,
          password: conn.password
        })
      )
      ws.send(authPayload)
      fitRef.current?.fit()
      sendResize()
    }

    ws.onmessage = (recv) => {
      if (typeof recv.data === 'string') {
        termRef.current?.write(recv.data)
        return
      }
      if (recv.data instanceof ArrayBuffer) {
        const text = new TextDecoder().decode(new Uint8Array(recv.data))
        termRef.current?.write(text)
      }
    }

    ws.onerror = () => {
      writeLine('读取数据异常')
      toast.error('SSH websocket 连接异常')
    }

    ws.onclose = () => {
      setConnecting(false)
      setConnected(false)
      writeLine('连接关闭')
    }
  }

  useEffect(() => {
    initTerm()
    const onResize = () => {
      fitRef.current?.fit()
      sendResize()
    }
    window.addEventListener('resize', onResize)
    return () => {
      window.removeEventListener('resize', onResize)
      disconnect()
      termRef.current?.dispose()
      termRef.current = null
      fitRef.current = null
    }
  }, [])

  useEffect(() => {
    if (autoConnectedRef.current) return
    const state = (location.state || {}) as Partial<ConnForm>
    const fromParams: Partial<ConnForm> = {
      manageIp: params.manageIp ? decodeURIComponent(params.manageIp) : '',
      username: params.username ? decodeURIComponent(params.username) : '',
      password: params.password ? decodeURIComponent(params.password) : '',
      sshPort: params.sshPort ? Number(decodeURIComponent(params.sshPort)) : 0
    }
    const fromRoute: ConnForm = {
      manageIp: String(fromParams.manageIp || state.manageIp || searchParams.get('manageIp') || ''),
      sshPort: Number(fromParams.sshPort || state.sshPort || searchParams.get('sshPort') || 22),
      username: String(fromParams.username || state.username || searchParams.get('username') || ''),
      password: String(fromParams.password || state.password || searchParams.get('password') || '')
    }
    if (!fromRoute.manageIp || !fromRoute.username || !fromRoute.password) return
    setForm(fromRoute)
    autoConnectedRef.current = true
    setTimeout(() => connect(fromRoute), 0)
  }, [location.state, params.manageIp, params.password, params.sshPort, params.username, searchParams])

  return (
    <Card>
      <CardHeader>
        <CardTitle>SSH 终端</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="grid gap-3 md:grid-cols-4">
          <Input placeholder="管理IP" value={form.manageIp} onChange={(e) => setForm((f) => ({ ...f, manageIp: e.target.value }))} />
          <Input type="number" placeholder="SSH端口" value={form.sshPort} onChange={(e) => setForm((f) => ({ ...f, sshPort: Number(e.target.value) }))} />
          <Input placeholder="用户名" value={form.username} onChange={(e) => setForm((f) => ({ ...f, username: e.target.value }))} />
          <Input type="password" placeholder="密码" value={form.password} onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))} />
        </div>
        <div className="flex gap-2">
          <Button onClick={() => connect()} disabled={connecting}>{connecting ? '连接中...' : connected ? '重新连接' : '连接'}</Button>
          <Button variant="outline" onClick={disconnect} disabled={!connected}>断开</Button>
        </div>
        <div ref={terminalBoxRef} className="h-[72vh] rounded-md border bg-black p-1" />
      </CardContent>
    </Card>
  )
}
