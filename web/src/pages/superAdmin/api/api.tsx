import { FormEvent, useEffect, useMemo, useState } from 'react'
import { toast } from 'sonner'
import { createApi, deleteApi, deleteApisByIds, getApiById, getApiList, updateApi } from '@/api/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type ApiItem = {
  ID: number
  path: string
  apiGroup: string
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | string
  description: string
}

type ApiForm = Pick<ApiItem, 'path' | 'apiGroup' | 'method' | 'description'> & { ID?: number }

const emptyForm: ApiForm = { path: '', apiGroup: '', method: 'GET', description: '' }

export default function ApiManagePage() {
  const [loading, setLoading] = useState(false)
  const [items, setItems] = useState<ApiItem[]>([])
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const [search, setSearch] = useState({ path: '', description: '', apiGroup: '', method: '' })
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState<ApiForm>(emptyForm)

  const query = async () => {
    setLoading(true)
    try {
      const res = await getApiList({ page, pageSize, ...search })
      if (res?.code === 0) {
        setItems(res.data?.list || [])
        setTotal(res.data?.total || 0)
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    query()
  }, [page, pageSize])

  const methods = useMemo(() => ['GET', 'POST', 'PUT', 'DELETE'], [])

  const onEdit = async (row: ApiItem) => {
    const res = await getApiById({ id: row.ID })
    if (res?.code === 0) {
      setForm({ ...res.data.api })
      setEditing(true)
    }
  }

  const onSave = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.path || !form.apiGroup || !form.method || !form.description) {
      toast.error('请填写完整字段')
      return
    }
    const res = editing ? await updateApi(form) : await createApi(form)
    if (res?.code === 0) {
      toast.success(editing ? '编辑成功' : '新增成功')
      setForm(emptyForm)
      setEditing(false)
      query()
    }
  }

  const onDeleteOne = async (row: ApiItem) => {
    if (!window.confirm(`确认删除 API ${row.path}?`)) return
    const res = await deleteApi(row)
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  const onBatchDelete = async () => {
    if (!selectedIds.length) return
    if (!window.confirm(`确认删除选中的 ${selectedIds.length} 条 API?`)) return
    const res = await deleteApisByIds({ ids: selectedIds })
    if (res?.code === 0) {
      toast.success('批量删除成功')
      setSelectedIds([])
      query()
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>API 管理</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-5">
          <Input placeholder="路径" value={search.path} onChange={(e) => setSearch((s) => ({ ...s, path: e.target.value }))} />
          <Input placeholder="描述" value={search.description} onChange={(e) => setSearch((s) => ({ ...s, description: e.target.value }))} />
          <Input placeholder="API组" value={search.apiGroup} onChange={(e) => setSearch((s) => ({ ...s, apiGroup: e.target.value }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={search.method} onChange={(e) => setSearch((s) => ({ ...s, method: e.target.value }))}>
            <option value="">全部请求</option>
            {methods.map((m) => (
              <option key={m} value={m}>{m}</option>
            ))}
          </select>
          <div className="flex gap-2">
            <Button onClick={() => { setPage(1); query() }}>查询</Button>
            <Button variant="outline" onClick={() => { setSearch({ path: '', description: '', apiGroup: '', method: '' }); setPage(1); }}>重置</Button>
          </div>
        </div>

        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-6" onSubmit={onSave}>
          <Input className="md:col-span-2" placeholder="路径" value={form.path} onChange={(e) => setForm((f) => ({ ...f, path: e.target.value }))} />
          <Input placeholder="组" value={form.apiGroup} onChange={(e) => setForm((f) => ({ ...f, apiGroup: e.target.value }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.method} onChange={(e) => setForm((f) => ({ ...f, method: e.target.value }))}>
            {methods.map((m) => (
              <option key={m} value={m}>{m}</option>
            ))}
          </select>
          <Input placeholder="描述" value={form.description} onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))} />
          <div className="flex gap-2">
            <Button type="submit">{editing ? '保存编辑' : '新增'}</Button>
            {editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(emptyForm) }}>取消</Button>}
          </div>
        </form>

        <div className="flex gap-2">
          <Button variant="destructive" size="sm" onClick={onBatchDelete} disabled={!selectedIds.length}>批量删除</Button>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-10">选</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>路径</TableHead>
              <TableHead>组</TableHead>
              <TableHead>方法</TableHead>
              <TableHead>描述</TableHead>
              <TableHead className="w-40">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>
                  <input type="checkbox" checked={selectedIds.includes(row.ID)} onChange={(e) => setSelectedIds((ids) => e.target.checked ? [...new Set([...ids, row.ID])] : ids.filter((id) => id !== row.ID))} />
                </TableCell>
                <TableCell>{row.ID}</TableCell>
                <TableCell>{row.path}</TableCell>
                <TableCell>{row.apiGroup}</TableCell>
                <TableCell>{row.method}</TableCell>
                <TableCell>{row.description}</TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="outline" onClick={() => onEdit(row)}>编辑</Button>
                  <Button size="sm" variant="destructive" onClick={() => onDeleteOne(row)}>删除</Button>
                </TableCell>
              </TableRow>
            ))}
            {!items.length && !loading && (
              <TableRow><TableCell colSpan={7} className="text-center text-muted-foreground">暂无数据</TableCell></TableRow>
            )}
          </TableBody>
        </Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <span>共 {total} 条</span>
          <div className="flex items-center gap-2">
            <select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>
              {[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}
            </select>
            <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button>
            <span>第 {page} 页</span>
            <Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
