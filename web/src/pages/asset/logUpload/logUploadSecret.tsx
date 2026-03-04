import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { addSecret, deleteSecret, getSecretById, getSecretList, getServerList, updateSecret } from '@/api/logUpload'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Secret = { ID: number; name: string; serverId: number }
const empty = { name: '', password: '', serverId: 0 }

export default function LogUploadSecretPage() {
  const [items, setItems] = useState<Secret[]>([])
  const [servers, setServers] = useState<any[]>([])
  const [form, setForm] = useState<any>(empty)
  const [editing, setEditing] = useState(false)
  const [search, setSearch] = useState({ name: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getSecretList({ page, pageSize, ...search })
    if (res?.code === 0) { setItems(res.data?.list || []); setTotal(res.data?.total || 0) }
  }

  useEffect(() => { query() }, [page, pageSize])
  useEffect(() => {
    const loadServers = async () => {
      const res = await getServerList({ page: 1, pageSize: 99999 })
      if (res?.code === 0) setServers(res.data?.list || [])
    }
    loadServers()
  }, [])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.password || !form.serverId) return toast.error('请填写完整信息')
    const res = editing ? await updateSecret(form) : await addSecret(form)
    if (res?.code === 0) { toast.success(editing ? '编辑成功' : '新增成功'); setEditing(false); setForm(empty); query() }
  }

  const edit = async (row: Secret) => {
    const res = await getSecretById({ id: row.ID })
    if (res?.code === 0) { setForm(res.data?.secret || row); setEditing(true) }
  }

  const remove = async (row: Secret) => {
    if (!window.confirm(`确认删除密钥 ${row.name}?`)) return
    const res = await deleteSecret({ ID: row.ID })
    if (res?.code === 0) { toast.success('删除成功'); query() }
  }

  const serverName = (id: number) => servers.find((s) => s.ID === id)?.hostname || '-'

  return (
    <Card>
      <CardHeader><CardTitle>日志上传密钥</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-3"><Input placeholder="密钥名" value={search.name} onChange={(e) => setSearch({ name: e.target.value })} /><div className="flex gap-2"><Button onClick={() => { setPage(1); query() }}>查询</Button><Button variant="outline" onClick={() => { setSearch({ name: '' }); setPage(1) }}>重置</Button></div></div>

        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-4" onSubmit={save}>
          <Input placeholder="密钥名" value={form.name || ''} onChange={(e) => setForm((f: any) => ({ ...f, name: e.target.value }))} />
          <Input placeholder="密码" type="password" value={form.password || ''} onChange={(e) => setForm((f: any) => ({ ...f, password: e.target.value }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.serverId || 0} onChange={(e) => setForm((f: any) => ({ ...f, serverId: Number(e.target.value) }))}><option value={0}>选择服务器</option>{servers.map((s) => <option key={s.ID} value={s.ID}>{s.hostname}</option>)}</select>
          <div className="flex gap-2"><Button type="submit">{editing ? '保存编辑' : '新增'}</Button>{editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(empty) }}>取消</Button>}</div>
        </form>

        <Table><TableHeader><TableRow><TableHead>ID</TableHead><TableHead>密钥名</TableHead><TableHead>所属服务器</TableHead><TableHead>操作</TableHead></TableRow></TableHeader><TableBody>{items.map((row) => <TableRow key={row.ID}><TableCell>{row.ID}</TableCell><TableCell>{row.name}</TableCell><TableCell>{serverName(row.serverId)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button></TableCell></TableRow>)}</TableBody></Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
