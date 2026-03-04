import { useState } from 'react'
import { toast } from 'sonner'
import { runSsh } from '@/api/ssh'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

export default function SshPage() {
  const [form, setForm] = useState({ manageIp: '', sshPort: 22, username: '', password: '' })
  const [running, setRunning] = useState(false)

  const execute = async () => {
    if (!form.manageIp || !form.sshPort || !form.username || !form.password) {
      toast.error('请填写完整连接参数')
      return
    }
    setRunning(true)
    try {
      const res = await runSsh({ server: { manageIp: form.manageIp, sshPort: Number(form.sshPort) }, username: form.username, password: form.password })
      if (res?.code === 0) toast.success('SSH 执行请求已发送')
    } finally {
      setRunning(false)
    }
  }

  return (
    <Card>
      <CardHeader><CardTitle>SSH 执行</CardTitle></CardHeader>
      <CardContent className="grid gap-3 md:grid-cols-2">
        <Input placeholder="管理IP" value={form.manageIp} onChange={(e) => setForm((f) => ({ ...f, manageIp: e.target.value }))} />
        <Input type="number" placeholder="SSH端口" value={form.sshPort} onChange={(e) => setForm((f) => ({ ...f, sshPort: Number(e.target.value) }))} />
        <Input placeholder="用户名" value={form.username} onChange={(e) => setForm((f) => ({ ...f, username: e.target.value }))} />
        <Input type="password" placeholder="密码" value={form.password} onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))} />
        <div><Button onClick={execute} disabled={running}>{running ? '执行中...' : '执行命令'}</Button></div>
      </CardContent>
    </Card>
  )
}
