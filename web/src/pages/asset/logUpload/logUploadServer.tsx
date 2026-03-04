import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { addServer, deleteServer, getServerById, getServerList, updateServer } from '@/api/logUpload'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Server = { ID: number; hostname: string; manageIp: string; port: number; mode: number }
const empty = { hostname: '', manageIp: '', port: 22, mode: 1 }

export default function LogUploadServerPage() {
  const [items, setItems] = useState<Server[]>([])
  const [form, setForm] = useState<any>(empty)
  const [editing, setEditing] = useState(false)
  const [search, setSearch] = useState({ hostname: '', manageIp: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getServerList({ page, pageSize, ...search })
    if (res?.code === 0) { setItems(res.data?.list || []); setTotal(res.data?.total || 0) }
  }
  useEffect(() => { query() }, [page, pageSize])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.hostname || !form.manageIp || !form.port || !form.mode) return toast.error('请填写完整信息')
    const res = editing ? await updateServer(form) : await addServer({ ...form, port: Number(form.port) })
    if (res?.code === 0) { toast.success(editing ? '编辑成功' : '新增成功'); setEditing(false); setForm(empty); query() }
  }

  const edit = async (row: Server) => {
    const res = await getServerById({ id: row.ID })
    if (res?.code === 0) { setForm(res.data?.server || row); setEditing(true) }
  }

  const remove = async (row: Server) => {
    if (!window.confirm(`确认删除服务器 ${row.hostname}?`)) return
    const res = await deleteServer({ ID: row.ID })
    if (res?.code === 0) { toast.success('删除成功'); query() }
  }

  return (
    <Card>
      <CardHeader><CardTitle>日志上传服务器</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-3"><Input placeholder="服务器名" value={search.hostname} onChange={(e) => setSearch((s) => ({ ...s, hostname: e.target.value }))} /><Input placeholder="管理IP" value={search.manageIp} onChange={(e) => setSearch((s) => ({ ...s, manageIp: e.target.value }))} /><div className="flex gap-2"><Button onClick={() => { setPage(1); query() }}>查询</Button><Button variant="outline" onClick={() => { setSearch({ hostname: '', manageIp: '' }); setPage(1) }}>重置</Button></div></div>

        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-5" onSubmit={save}>
          <Input placeholder="服务器名" value={form.hostname || ''} onChange={(e) => setForm((f: any) => ({ ...f, hostname: e.target.value }))} />
          <Input placeholder="管理IP" value={form.manageIp || ''} onChange={(e) => setForm((f: any) => ({ ...f, manageIp: e.target.value }))} />
          <Input type="number" placeholder="端口" value={form.port} onChange={(e) => setForm((f: any) => ({ ...f, port: Number(e.target.value) }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.mode} onChange={(e) => setForm((f: any) => ({ ...f, mode: Number(e.target.value) }))}><option value={1}>ftp</option><option value={2}>sftp</option></select>
          <div className="flex gap-2"><Button type="submit">{editing ? '保存编辑' : '新增'}</Button>{editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(empty) }}>取消</Button>}</div>
        </form>

        <Table><TableHeader><TableRow><TableHead>ID</TableHead><TableHead>服务器</TableHead><TableHead>IP</TableHead><TableHead>端口</TableHead><TableHead>方式</TableHead><TableHead>操作</TableHead></TableRow></TableHeader><TableBody>{items.map((row) => <TableRow key={row.ID}><TableCell>{row.ID}</TableCell><TableCell>{row.hostname}</TableCell><TableCell>{row.manageIp}</TableCell><TableCell>{row.port}</TableCell><TableCell>{row.mode === 2 ? 'sftp' : 'ftp'}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button></TableCell></TableRow>)}</TableBody></Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
