import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { getSystemConfig, setSystemConfig } from '@/api/system'
import { emailTest } from '@/api/email'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

export default function SystemConfigPage() {
  const [config, setConfig] = useState<any>(null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    const load = async () => {
      const res = await getSystemConfig()
      if (res?.code === 0) {
        setConfig(res.data?.config || res.data)
      }
    }
    load()
  }, [])

  const updateField = (path: string, value: any) => {
    setConfig((prev: any) => {
      if (!prev) return prev
      const next = structuredClone(prev)
      const keys = path.split('.')
      let obj = next
      for (let i = 0; i < keys.length - 1; i++) {
        if (!obj[keys[i]]) obj[keys[i]] = {}
        obj = obj[keys[i]]
      }
      obj[keys[keys.length - 1]] = value
      return next
    })
  }

  const save = async () => {
    if (!config) return
    setSaving(true)
    try {
      const res = await setSystemConfig(config)
      if (res?.code === 0) toast.success('配置已更新')
    } finally {
      setSaving(false)
    }
  }

  const sendEmailTest = async () => {
    if (!config?.email) return
    const res = await emailTest(config.email)
    if (res?.code === 0) toast.success('测试邮件发送成功')
  }

  if (!config) {
    return <div className="p-6 text-sm text-muted-foreground">加载配置中...</div>
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>基础配置</CardTitle>
          <CardDescription>支持在线修改并提交到后端配置文件。</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-3">
          <Input placeholder="环境 ENV" value={config.system?.env || ''} onChange={(e) => updateField('system.env', e.target.value)} />
          <Input placeholder="端口 Addr" value={config.system?.addr || ''} onChange={(e) => updateField('system.addr', Number(e.target.value))} />
          <Input placeholder="数据库类型" value={config.system?.dbType || ''} onChange={(e) => updateField('system.dbType', e.target.value)} />
          <Input placeholder="OSS类型" value={config.system?.ossType || ''} onChange={(e) => updateField('system.ossType', e.target.value)} />
          <Input placeholder="JWT SigningKey" value={config.jwt?.signingKey || ''} onChange={(e) => updateField('jwt.signingKey', e.target.value)} />
          <Input placeholder="JWT ExpiresTime" value={config.jwt?.expiresTime || ''} onChange={(e) => updateField('jwt.expiresTime', e.target.value)} />
          <Input placeholder="MySQL Path" value={config.mysql?.path || ''} onChange={(e) => updateField('mysql.path', e.target.value)} />
          <Input placeholder="MySQL DBName" value={config.mysql?.dbname || ''} onChange={(e) => updateField('mysql.dbname', e.target.value)} />
          <Input placeholder="Redis Addr" value={config.redis?.addr || ''} onChange={(e) => updateField('redis.addr', e.target.value)} />
          <Input placeholder="邮件 From" value={config.email?.from || ''} onChange={(e) => updateField('email.from', e.target.value)} />
          <Input placeholder="邮件 To" value={config.email?.to || ''} onChange={(e) => updateField('email.to', e.target.value)} />
          <Input placeholder="邮件 Host" value={config.email?.host || ''} onChange={(e) => updateField('email.host', e.target.value)} />
        </CardContent>
      </Card>

      <div className="flex gap-2">
        <Button onClick={save} disabled={saving}>{saving ? '保存中...' : '立即更新'}</Button>
        <Button variant="outline" onClick={sendEmailTest}>测试邮件</Button>
      </div>
    </div>
  )
}
