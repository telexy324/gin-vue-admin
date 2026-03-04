import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { initDB } from '@/api/initdb'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function InitPage() {
  const navigate = useNavigate()
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({
    sqlType: 'mysql',
    host: '127.0.0.1',
    port: '3306',
    userName: 'root',
    password: '',
    dbName: 'gva'
  })

  const setField = (k, v) => setForm((prev) => ({ ...prev, [k]: v }))

  const onSubmit = async (e) => {
    e.preventDefault()
    setSubmitting(true)
    try {
      const res = await initDB(form)
      if (res?.code === 0) {
        toast.success(res.msg || '初始化成功')
        navigate('/login', { replace: true })
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="grid min-h-screen place-items-center p-6">
      <Card className="w-full max-w-xl">
        <CardHeader>
          <CardTitle>初始化数据库</CardTitle>
          <CardDescription>首次部署时填写数据库连接参数</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="grid gap-3" onSubmit={onSubmit}>
            <Input value={form.sqlType} disabled />
            <Input value={form.host} onChange={(e) => setField('host', e.target.value)} placeholder="host" />
            <Input value={form.port} onChange={(e) => setField('port', e.target.value)} placeholder="port" />
            <Input value={form.userName} onChange={(e) => setField('userName', e.target.value)} placeholder="userName" />
            <Input value={form.password} onChange={(e) => setField('password', e.target.value)} placeholder="password" />
            <Input value={form.dbName} onChange={(e) => setField('dbName', e.target.value)} placeholder="dbName" />
            <Button disabled={submitting}>{submitting ? '初始化中...' : '立即初始化'}</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
