import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { addSchedule, deleteSchedule, getScheduleById, getTemplateScheduleList, updateSchedule, validateScheduleCronFormat } from '@/api/schedule'
import { getTemplateList } from '@/api/template'
import { getAdminSystems, getSystemServerIds } from '@/api/cmdb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const empty = { templateId: 0, cronFormat: '', valid: 0, systemId: 0, commandVars: [] as string[], targetIds: [] as number[] }

export default function SchedulePage() {
  const [items, setItems] = useState<any[]>([])
  const [systems, setSystems] = useState<any[]>([])
  const [templates, setTemplates] = useState<any[]>([])
  const [form, setForm] = useState<any>(empty)
  const [editing, setEditing] = useState(false)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getTemplateScheduleList({ page, pageSize })
    if (res?.code === 0) { setItems(res.data?.list || []); setTotal(res.data?.total || 0) }
  }

  useEffect(() => { query() }, [page, pageSize])
  useEffect(() => {
    const load = async () => {
      const [sr, tr] = await Promise.all([
        getAdminSystems({}),
        getTemplateList({ page: 1, pageSize: 99999 })
      ])
      if (sr?.code === 0) setSystems(sr.data || sr.data?.list || [])
      if (tr?.code === 0) setTemplates(tr.data?.list || [])
    }
    load()
  }, [])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.systemId || !form.templateId || !form.cronFormat) return toast.error('请填写系统、模板和cron')
    const res = editing ? await updateSchedule(form) : await addSchedule(form)
    if (res?.code === 0) { toast.success(editing ? '编辑成功' : '新增成功'); setEditing(false); setForm(empty); query() }
  }

  const edit = async (row: any) => {
    const res = await getScheduleById({ id: row.ID })
    if (res?.code === 0) { setForm(res.data?.schedule || row); setEditing(true) }
  }

  const remove = async (row: any) => {
    if (!window.confirm('确认删除该调度?')) return
    const res = await deleteSchedule({ ID: row.ID })
    if (res?.code === 0) { toast.success('删除成功'); query() }
  }

  const checkCron = async () => {
    const res = await validateScheduleCronFormat({ cronFormat: form.cronFormat })
    if (res?.code === 0) toast.success('cron 格式校验通过')
  }

  const changeSystem = async (id: number) => {
    setForm((f: any) => ({ ...f, systemId: id, targetIds: [] }))
    const r = await getSystemServerIds({ ID: id })
    const ids = (r?.data?.[0]?.children || []).map((x: any) => x.ID)
    setForm((f: any) => ({ ...f, targetIds: ids }))
  }

  const tName = (id: number) => templates.find((x) => x.ID === id)?.name || '-'

  return (
    <Card>
      <CardHeader><CardTitle>定时任务</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-4" onSubmit={save}>
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.systemId || 0} onChange={(e) => changeSystem(Number(e.target.value))}><option value={0}>选择系统</option>{systems.map((s: any) => <option key={s.ID} value={s.ID}>{s.name}</option>)}</select>
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.templateId || 0} onChange={(e) => setForm((f: any) => ({ ...f, templateId: Number(e.target.value) }))}><option value={0}>选择模板</option>{templates.filter((t) => !form.systemId || t.systemId === form.systemId).map((t) => <option key={t.ID} value={t.ID}>{t.name}</option>)}</select>
          <Input placeholder="cronFormat" value={form.cronFormat || ''} onChange={(e) => setForm((f: any) => ({ ...f, cronFormat: e.target.value }))} />
          <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={form.valid === 1} onChange={(e) => setForm((f: any) => ({ ...f, valid: e.target.checked ? 1 : 0 }))} />启用</label>
          <Input className="md:col-span-2" placeholder="参数，逗号分隔" value={(form.commandVars || []).join(',')} onChange={(e) => setForm((f: any) => ({ ...f, commandVars: e.target.value.split(',').map((x) => x.trim()).filter(Boolean) }))} />
          <Input className="md:col-span-2" placeholder="目标ID，逗号分隔" value={(form.targetIds || []).join(',')} onChange={(e) => setForm((f: any) => ({ ...f, targetIds: e.target.value.split(',').map((x) => Number(x.trim())).filter((n) => Number.isFinite(n) && n > 0) }))} />
          <div className="flex gap-2 md:col-span-4"><Button type="button" variant="outline" onClick={checkCron}>检查cron</Button><Button type="submit">{editing ? '保存编辑' : '新增调度'}</Button>{editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(empty) }}>取消</Button>}</div>
        </form>

        <Table><TableHeader><TableRow><TableHead>ID</TableHead><TableHead>模板名</TableHead><TableHead>启用</TableHead><TableHead>cron</TableHead><TableHead>操作</TableHead></TableRow></TableHeader><TableBody>{items.map((row) => <TableRow key={row.ID}><TableCell>{row.ID}</TableCell><TableCell>{tName(row.templateId)}</TableCell><TableCell>{row.valid === 1 ? '是' : '否'}</TableCell><TableCell>{row.cronFormat}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button></TableCell></TableRow>)}</TableBody></Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
