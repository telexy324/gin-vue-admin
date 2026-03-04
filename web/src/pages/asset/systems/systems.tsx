import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { addSystem, deleteSystem, getSystemById, getSystemList, updateSystem } from '@/api/cmdb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useNavigate } from 'react-router-dom'

type Sys = { ID: number; name: string }

export default function SystemsPage() {
  const navigate = useNavigate()
  const [items, setItems] = useState<Sys[]>([])
  const [form, setForm] = useState<any>({ name: '' })
  const [editing, setEditing] = useState(false)
  const [search, setSearch] = useState({ name: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getSystemList({ page, pageSize, ...search })
    if (res?.code === 0) {
      setItems((res.data?.list || []).map((x: any) => x.system || x))
      setTotal(res.data?.total || 0)
    }
  }
  useEffect(() => { query() }, [page, pageSize])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.name) return toast.error('请输入系统名')
    const res = editing ? await updateSystem(form) : await addSystem(form)
    if (res?.code === 0) {
      toast.success(editing ? '编辑成功' : '新增成功')
      setEditing(false)
      setForm({ name: '' })
      query()
    }
  }

  const edit = async (row: Sys) => {
    const res = await getSystemById({ id: row.ID })
    if (res?.code === 0) {
      setForm(res.data?.system || row)
      setEditing(true)
    }
  }

  const remove = async (row: Sys) => {
    if (!window.confirm(`确认删除系统 ${row.name}?`)) return
    const res = await deleteSystem({ ID: row.ID })
    if (res?.code === 0) { toast.success('删除成功'); query() }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>系统管理</CardTitle>
        <CardDescription>可进入系统关系图页面查看拓扑。</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-3">
          <Input placeholder="系统名" value={search.name} onChange={(e) => setSearch((s) => ({ ...s, name: e.target.value }))} />
          <div className="flex gap-2"><Button onClick={() => { setPage(1); query() }}>查询</Button><Button variant="outline" onClick={() => { setSearch({ name: '' }); setPage(1) }}>重置</Button></div>
        </div>

        <form className="flex gap-2 rounded-md border p-3" onSubmit={save}>
          <Input placeholder="系统名" value={form.name || ''} onChange={(e) => setForm((f: any) => ({ ...f, name: e.target.value }))} />
          <Button type="submit">{editing ? '保存编辑' : '新增系统'}</Button>
          {editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm({ name: '' }) }}>取消</Button>}
        </form>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>系统名</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>{row.ID}</TableCell><TableCell>{row.name}</TableCell>
                <TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button><Button size="sm" variant="secondary" onClick={() => navigate('/layout/asset/graph', { state: { systemId: row.ID, name: row.name } })}>关系图</Button></TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
