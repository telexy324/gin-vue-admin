import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { addServer, deleteServer, getServerById, getServerList, getSystemList, updateServer } from '@/api/cmdb'
import { runSsh } from '@/api/ssh'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Server = {
  ID: number
  hostname: string
  displayName?: string
  architecture?: string
  manageIp: string
  sshPort: number
  os?: string
  osVersion?: string
  systemId?: number
}

const emptyForm = { hostname: '', displayName: '', architecture: '', manageIp: '', sshPort: 22, os: '', osVersion: '', systemId: 0 }

export default function ServerPage() {
  const [items, setItems] = useState<Server[]>([])
  const [systems, setSystems] = useState<any[]>([])
  const [search, setSearch] = useState({ hostname: '', manageIp: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState<any>(emptyForm)

  const query = async () => {
    const res = await getServerList({ page, pageSize, ...search })
    if (res?.code === 0) {
      setItems(res.data?.list || [])
      setTotal(res.data?.total || 0)
    }
  }

  useEffect(() => { query() }, [page, pageSize])
  useEffect(() => {
    const loadSystems = async () => {
      const r = await getSystemList({ page: 1, pageSize: 99999 })
      if (r?.code === 0) setSystems(r.data?.list || [])
    }
    loadSystems()
  }, [])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.hostname || !form.manageIp || !form.sshPort || !form.systemId) {
      toast.error('请填写完整信息')
      return
    }
    const payload = { ...form, sshPort: Number(form.sshPort) }
    const res = editing ? await updateServer(payload) : await addServer(payload)
    if (res?.code === 0) {
      toast.success(editing ? '编辑成功' : '新增成功')
      setEditing(false)
      setForm(emptyForm)
      query()
    }
  }

  const edit = async (row: Server) => {
    const res = await getServerById({ id: row.ID })
    if (res?.code === 0) {
      setForm(res.data?.server)
      setEditing(true)
    }
  }

  const remove = async (row: Server) => {
    if (!window.confirm(`确认删除服务器 ${row.hostname}?`)) return
    const res = await deleteServer({ ID: row.ID })
    if (res?.code === 0) { toast.success('删除成功'); query() }
  }

  const execute = async (row: Server) => {
    const username = window.prompt(`输入 ${row.manageIp} 的 SSH 用户名`)
    if (!username) return
    const password = window.prompt('输入 SSH 密码')
    if (!password) return
    const res = await runSsh({ server: { manageIp: row.manageIp, sshPort: row.sshPort }, username, password })
    if (res?.code === 0) toast.success('已发起执行')
  }

  const systemName = (id?: number) => systems.find((s) => s.ID === id)?.name || '-'

  return (
    <Card>
      <CardHeader><CardTitle>服务器管理</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-3">
          <Input placeholder="服务器名" value={search.hostname} onChange={(e) => setSearch((s) => ({ ...s, hostname: e.target.value }))} />
          <Input placeholder="管理IP" value={search.manageIp} onChange={(e) => setSearch((s) => ({ ...s, manageIp: e.target.value }))} />
          <div className="flex gap-2"><Button onClick={() => { setPage(1); query() }}>查询</Button><Button variant="outline" onClick={() => { setSearch({ hostname: '', manageIp: '' }); setPage(1) }}>重置</Button></div>
        </div>

        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-5" onSubmit={save}>
          <Input placeholder="服务器名" value={form.hostname} onChange={(e) => setForm((f: any) => ({ ...f, hostname: e.target.value }))} />
          <Input placeholder="展示名" value={form.displayName || ''} onChange={(e) => setForm((f: any) => ({ ...f, displayName: e.target.value }))} />
          <Input placeholder="管理IP" value={form.manageIp} onChange={(e) => setForm((f: any) => ({ ...f, manageIp: e.target.value }))} />
          <Input placeholder="SSH端口" type="number" value={form.sshPort} onChange={(e) => setForm((f: any) => ({ ...f, sshPort: Number(e.target.value) }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.systemId || 0} onChange={(e) => setForm((f: any) => ({ ...f, systemId: Number(e.target.value) }))}>
            <option value={0}>选择所属系统</option>
            {systems.map((s) => <option key={s.ID} value={s.ID}>{s.name}</option>)}
          </select>
          <Input placeholder="架构" value={form.architecture || ''} onChange={(e) => setForm((f: any) => ({ ...f, architecture: e.target.value }))} />
          <Input placeholder="操作系统" value={form.os || ''} onChange={(e) => setForm((f: any) => ({ ...f, os: e.target.value }))} />
          <Input placeholder="版本" value={form.osVersion || ''} onChange={(e) => setForm((f: any) => ({ ...f, osVersion: e.target.value }))} />
          <div className="flex gap-2 md:col-span-2">
            <Button type="submit">{editing ? '保存编辑' : '新增'}</Button>
            {editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(emptyForm) }}>取消</Button>}
          </div>
        </form>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>服务器</TableHead><TableHead>管理IP</TableHead><TableHead>端口</TableHead><TableHead>系统</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>{row.ID}</TableCell><TableCell>{row.hostname}</TableCell><TableCell>{row.manageIp}</TableCell><TableCell>{row.sshPort}</TableCell><TableCell>{systemName(row.systemId)}</TableCell>
                <TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button><Button size="sm" variant="secondary" onClick={() => execute(row)}>执行</Button></TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
