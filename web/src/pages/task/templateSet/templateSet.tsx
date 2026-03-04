import { FormEvent, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { addSet, addSetTask, deleteSet, getSetById, getSetList, updateSet } from '@/api/template'
import { getSystemList } from '@/api/cmdb'
import { getTemplateList } from '@/api/template'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const empty = { name: '', systemId: 0, templates: [{ templateIds: [] as number[], seq: 1 }] }

export default function TemplateSetPage() {
  const navigate = useNavigate()
  const [items, setItems] = useState<any[]>([])
  const [systems, setSystems] = useState<any[]>([])
  const [templates, setTemplates] = useState<any[]>([])
  const [form, setForm] = useState<any>(empty)
  const [editing, setEditing] = useState(false)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getSetList({ page, pageSize })
    if (res?.code === 0) { setItems(res.data?.list || []); setTotal(res.data?.total || 0) }
  }
  useEffect(() => { query() }, [page, pageSize])
  useEffect(() => {
    const load = async () => {
      const [sr, tr] = await Promise.all([
        getSystemList({ page: 1, pageSize: 99999 }),
        getTemplateList({ page: 1, pageSize: 99999 })
      ])
      if (sr?.code === 0) setSystems(sr.data?.list || [])
      if (tr?.code === 0) setTemplates(tr.data?.list || [])
    }
    load()
  }, [])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.systemId || !form.templates?.length) return toast.error('请填写完整信息')
    const res = editing ? await updateSet(form) : await addSet(form)
    if (res?.code === 0) { toast.success(editing ? '编辑成功' : '新增成功'); setEditing(false); setForm(empty); query() }
  }

  const edit = async (row: any) => {
    const res = await getSetById({ id: row.ID })
    if (res?.code === 0) { setForm(res.data || row); setEditing(true) }
  }

  const remove = async (row: any) => {
    if (!window.confirm('确认删除模板集?')) return
    const res = await deleteSet({ ID: row.ID })
    if (res?.code === 0) { toast.success('删除成功'); query() }
  }

  const run = async (row: any) => {
    const res = await addSetTask({ setId: row.ID })
    if (res?.code === 0) {
      toast.success('模板集任务已创建')
      const setTaskId = res.data?.setTask?.ID
      if (setTaskId) navigate(`/layout/task/templateSet/detail/${setTaskId}`)
    }
  }

  const addItem = () => setForm((f: any) => ({ ...f, templates: [...(f.templates || []), { templateIds: [], seq: (f.templates?.length || 0) + 1 }] }))
  const removeItem = (idx: number) => setForm((f: any) => ({ ...f, templates: f.templates.filter((_: any, i: number) => i !== idx) }))
  const updateItem = (idx: number, patch: any) => setForm((f: any) => ({ ...f, templates: f.templates.map((x: any, i: number) => i === idx ? { ...x, ...patch } : x) }))

  const systemName = (id: number) => systems.find((x) => x.ID === id)?.name || '-'

  return (
    <Card>
      <CardHeader><CardTitle>模板集</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <form className="space-y-2 rounded-md border p-3" onSubmit={save}>
          <div className="grid gap-2 md:grid-cols-3">
            <Input placeholder="模板集名称" value={form.name || ''} onChange={(e) => setForm((f: any) => ({ ...f, name: e.target.value }))} />
            <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.systemId || 0} onChange={(e) => setForm((f: any) => ({ ...f, systemId: Number(e.target.value) }))}><option value={0}>选择系统</option>{systems.map((s) => <option key={s.ID} value={s.ID}>{s.name}</option>)}</select>
            <Button type="button" variant="secondary" onClick={addItem}>新增步骤</Button>
          </div>
          {(form.templates || []).map((item: any, idx: number) => (
            <div className="grid gap-2 md:grid-cols-6" key={idx}>
              <select className="h-10 rounded-md border bg-background px-3 text-sm md:col-span-4" multiple value={item.templateIds || []} onChange={(e) => updateItem(idx, { templateIds: [...e.target.selectedOptions].map((o) => Number(o.value)) })}>
                {templates.filter((t) => !form.systemId || t.systemId === form.systemId).map((t) => <option key={t.ID} value={t.ID}>{t.name}</option>)}
              </select>
              <Input type="number" value={item.seq} onChange={(e) => updateItem(idx, { seq: Number(e.target.value) })} />
              <Button type="button" variant="destructive" onClick={() => removeItem(idx)}>删除步骤</Button>
            </div>
          ))}
          <div className="flex gap-2"><Button type="submit">{editing ? '保存编辑' : '新增模板集'}</Button>{editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(empty) }}>取消</Button>}</div>
        </form>

        <Table><TableHeader><TableRow><TableHead>ID</TableHead><TableHead>名称</TableHead><TableHead>系统</TableHead><TableHead>步骤数</TableHead><TableHead>操作</TableHead></TableRow></TableHeader><TableBody>{items.map((row) => <TableRow key={row.ID}><TableCell>{row.ID}</TableCell><TableCell>{row.name}</TableCell><TableCell>{systemName(row.systemId)}</TableCell><TableCell>{row.templates?.length || 0}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button><Button size="sm" variant="secondary" onClick={() => run(row)}>执行</Button></TableCell></TableRow>)}</TableBody></Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
