import { FormEvent, useEffect, useMemo, useState } from 'react'
import { toast } from 'sonner'
import { addTask } from '@/api/task'
import { addTemplate, deleteTemplate, getTemplateById, getTemplateList, updateTemplate } from '@/api/template'
import { getSystemList } from '@/api/cmdb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'

const empty = {
  name: '',
  description: '',
  systemId: 0,
  executeType: 1,
  shellType: 1,
  deployType: 1,
  sysUser: '',
  becomeUser: '',
  command: '',
  commandVarNumbers: 0,
  targetIds: [] as number[],
  localPath: '',
  remotePath: '',
  logPath: ''
}

const executeTypeLabel = (v: number) => {
  if (v === 2) return '日志提取'
  if (v === 3) return '上传部署'
  return '通用执行'
}

export default function TemplateListPage() {
  const [items, setItems] = useState<any[]>([])
  const [systems, setSystems] = useState<any[]>([])
  const [form, setForm] = useState<any>(empty)
  const [editing, setEditing] = useState(false)
  const [search, setSearch] = useState({ name: '', executeType: 0 })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const filteredSystems = useMemo(() => systems, [systems])

  const query = async () => {
    const params: any = { page, pageSize, ...search }
    if (!params.executeType) delete params.executeType
    const res = await getTemplateList(params)
    if (res?.code === 0) {
      setItems(res.data?.list || [])
      setTotal(res.data?.total || 0)
    }
  }

  useEffect(() => {
    query()
  }, [page, pageSize])

  useEffect(() => {
    const load = async () => {
      const res = await getSystemList({ page: 1, pageSize: 99999 })
      if (res?.code === 0) setSystems(res.data?.list || [])
    }
    load()
  }, [])

  const save = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.systemId) return toast.error('请填写模板名与所属系统')
    if (!form.command && form.executeType === 1) return toast.error('通用执行模板必须填写命令/脚本')
    if (form.executeType === 2 && !form.logPath) return toast.error('日志提取模板必须填写日志路径')
    if (form.executeType === 3 && (!form.localPath || !form.remotePath)) return toast.error('上传部署模板必须填写本地与远端路径')
    const payload = {
      ...form,
      commandVarNumbers: Number(form.commandVarNumbers || 0),
      targetIds: (form.targetIds || []).map((x: any) => Number(x)).filter((x: number) => x > 0)
    }
    const res = editing ? await updateTemplate(payload) : await addTemplate(payload)
    if (res?.code === 0) {
      toast.success(editing ? '编辑成功' : '新增成功')
      setEditing(false)
      setForm(empty)
      query()
    }
  }

  const edit = async (row: any) => {
    const res = await getTemplateById({ id: row.ID })
    if (res?.code === 0) {
      setForm({ ...empty, ...(res.data?.template || row) })
      setEditing(true)
    }
  }

  const remove = async (row: any) => {
    if (!window.confirm(`确认删除模板 ${row.name}?`)) return
    const res = await deleteTemplate({ ID: row.ID })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  const run = async (row: any) => {
    const res = await addTask({ templateId: row.ID })
    if (res?.code === 0) toast.success('任务已创建')
  }

  const systemName = (id: number) => systems.find((s) => s.ID === id)?.name || '-'

  return (
    <Card>
      <CardHeader>
        <CardTitle>模板列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-4">
          <Input placeholder="模板名" value={search.name} onChange={(e) => setSearch((s) => ({ ...s, name: e.target.value }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={search.executeType} onChange={(e) => setSearch((s) => ({ ...s, executeType: Number(e.target.value) }))}>
            <option value={0}>全部类型</option>
            <option value={1}>通用执行</option>
            <option value={2}>日志提取</option>
            <option value={3}>上传部署</option>
          </select>
          <div className="flex gap-2 md:col-span-2">
            <Button onClick={() => { setPage(1); query() }}>查询</Button>
            <Button variant="outline" onClick={() => { setSearch({ name: '', executeType: 0 }); setPage(1) }}>重置</Button>
          </div>
        </div>

        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-4" onSubmit={save}>
          <Input placeholder="模板名" value={form.name || ''} onChange={(e) => setForm((f: any) => ({ ...f, name: e.target.value }))} />
          <Input placeholder="描述" value={form.description || ''} onChange={(e) => setForm((f: any) => ({ ...f, description: e.target.value }))} />
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.systemId || 0} onChange={(e) => setForm((f: any) => ({ ...f, systemId: Number(e.target.value) }))}>
            <option value={0}>选择系统</option>
            {filteredSystems.map((s) => <option key={s.ID} value={s.ID}>{s.name}</option>)}
          </select>
          <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.executeType || 1} onChange={(e) => setForm((f: any) => ({ ...f, executeType: Number(e.target.value) }))}>
            <option value={1}>通用执行</option>
            <option value={2}>日志提取</option>
            <option value={3}>上传部署</option>
          </select>

          <Input placeholder="连接用户" value={form.sysUser || ''} onChange={(e) => setForm((f: any) => ({ ...f, sysUser: e.target.value }))} />
          <Input placeholder="执行用户" value={form.becomeUser || ''} onChange={(e) => setForm((f: any) => ({ ...f, becomeUser: e.target.value }))} />
          <Input placeholder="目标ID(逗号分隔)" value={(form.targetIds || []).join(',')} onChange={(e) => setForm((f: any) => ({ ...f, targetIds: e.target.value.split(',').map((x) => Number(x.trim())).filter((n: number) => n > 0) }))} />
          <Input type="number" placeholder="参数个数" value={form.commandVarNumbers || 0} onChange={(e) => setForm((f: any) => ({ ...f, commandVarNumbers: Number(e.target.value) }))} />

          {form.executeType === 1 && (
            <>
              <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.shellType || 1} onChange={(e) => setForm((f: any) => ({ ...f, shellType: Number(e.target.value) }))}>
                <option value={1}>shell</option>
                <option value={2}>python</option>
                <option value={3}>powershell</option>
              </select>
              <Input className="md:col-span-3" placeholder="命令/脚本内容" value={form.command || ''} onChange={(e) => setForm((f: any) => ({ ...f, command: e.target.value }))} />
            </>
          )}

          {form.executeType === 2 && (
            <Input className="md:col-span-4" placeholder="日志路径 logPath" value={form.logPath || ''} onChange={(e) => setForm((f: any) => ({ ...f, logPath: e.target.value }))} />
          )}

          {form.executeType === 3 && (
            <>
              <select className="h-10 rounded-md border bg-background px-3 text-sm" value={form.deployType || 1} onChange={(e) => setForm((f: any) => ({ ...f, deployType: Number(e.target.value) }))}>
                <option value={1}>文件上传</option>
                <option value={2}>部署触发</option>
              </select>
              <Input placeholder="本地路径 localPath" value={form.localPath || ''} onChange={(e) => setForm((f: any) => ({ ...f, localPath: e.target.value }))} />
              <Input className="md:col-span-2" placeholder="远端路径 remotePath" value={form.remotePath || ''} onChange={(e) => setForm((f: any) => ({ ...f, remotePath: e.target.value }))} />
            </>
          )}

          <div className="flex gap-2 md:col-span-4">
            <Button type="submit">{editing ? '保存编辑' : '新增模板'}</Button>
            {editing && <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(empty) }}>取消</Button>}
          </div>
        </form>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>模板名</TableHead>
              <TableHead>系统</TableHead>
              <TableHead>类型</TableHead>
              <TableHead>最近任务</TableHead>
              <TableHead>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>{row.ID}</TableCell>
                <TableCell>{row.name}</TableCell>
                <TableCell>{systemName(row.systemId)}</TableCell>
                <TableCell><Badge variant="secondary">{executeTypeLabel(row.executeType || 1)}</Badge></TableCell>
                <TableCell>{row.lastTask?.ID ? `#${row.lastTask.ID}` : '-'}</TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="outline" onClick={() => edit(row)}>编辑</Button>
                  <Button size="sm" variant="destructive" onClick={() => remove(row)}>删除</Button>
                  <Button size="sm" variant="secondary" onClick={() => run(row)}>构建</Button>
                </TableCell>
              </TableRow>
            ))}
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
