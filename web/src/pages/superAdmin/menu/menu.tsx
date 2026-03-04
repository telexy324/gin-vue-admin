import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { addBaseMenu, deleteBaseMenu, getBaseMenuById, getMenuList, updateBaseMenu } from '@/api/menu'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type MenuMeta = { title?: string; icon?: string; keepAlive?: boolean; closeTab?: boolean }
type MenuItem = {
  ID: number
  parentId: string
  name: string
  path: string
  component: string
  hidden: boolean
  sort?: number
  meta: MenuMeta
}

type MenuForm = {
  ID?: number
  parentId: string
  name: string
  path: string
  component: string
  hidden: boolean
  sort: number
  meta: MenuMeta
}

const emptyForm: MenuForm = {
  parentId: '0',
  name: '',
  path: '',
  component: '',
  hidden: false,
  sort: 0,
  meta: { title: '', icon: '', keepAlive: false, closeTab: false }
}

const flatten = (list: any[] = []): MenuItem[] => {
  const result: MenuItem[] = []
  const walk = (arr: any[]) => {
    arr.forEach((item) => {
      result.push(item)
      if (item.children?.length) walk(item.children)
    })
  }
  walk(list)
  return result
}

export default function MenuManagePage() {
  const [items, setItems] = useState<MenuItem[]>([])
  const [loading, setLoading] = useState(false)
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState<MenuForm>(emptyForm)

  const query = async () => {
    setLoading(true)
    try {
      const res = await getMenuList({ page: 1, pageSize: 999 })
      if (res?.code === 0) {
        setItems(flatten(res.data?.list || []))
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    query()
  }, [])

  const onEdit = async (id: number) => {
    const res = await getBaseMenuById({ id })
    if (res?.code === 0) {
      const menu = res.data?.menu
      setForm({
        ID: menu.ID,
        parentId: String(menu.parentId ?? '0'),
        name: menu.name || '',
        path: menu.path || '',
        component: menu.component || '',
        hidden: Boolean(menu.hidden),
        sort: Number(menu.sort || 0),
        meta: {
          title: menu.meta?.title || '',
          icon: menu.meta?.icon || '',
          keepAlive: Boolean(menu.meta?.keepAlive),
          closeTab: Boolean(menu.meta?.closeTab)
        }
      })
      setEditing(true)
    }
  }

  const onSave = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.path || !form.component || !form.meta.title) {
      toast.error('请填写必填项')
      return
    }
    const payload = { ...form, parentId: String(form.parentId) }
    const res = editing ? await updateBaseMenu(payload) : await addBaseMenu(payload)
    if (res?.code === 0) {
      toast.success(editing ? '编辑成功' : '新增成功')
      setForm(emptyForm)
      setEditing(false)
      query()
    }
  }

  const onDelete = async (id: number) => {
    if (!window.confirm('确认删除该菜单？')) return
    const res = await deleteBaseMenu({ ID: id })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  const onAddChild = (parentId: number | string) => {
    setEditing(false)
    setForm({ ...emptyForm, parentId: String(parentId) })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>菜单管理</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-4" onSubmit={onSave}>
          <Input placeholder="路由Name" value={form.name} onChange={(e) => setForm((f) => ({ ...f, name: e.target.value, path: e.target.value }))} />
          <Input placeholder="路由Path" value={form.path} onChange={(e) => setForm((f) => ({ ...f, path: e.target.value }))} />
          <Input placeholder="父ID（0为根）" value={form.parentId} onChange={(e) => setForm((f) => ({ ...f, parentId: e.target.value }))} />
          <Input placeholder="文件路径，例如 view/superAdmin/api/api.vue" value={form.component} onChange={(e) => setForm((f) => ({ ...f, component: e.target.value }))} />
          <Input placeholder="展示名称" value={form.meta.title || ''} onChange={(e) => setForm((f) => ({ ...f, meta: { ...f.meta, title: e.target.value } }))} />
          <Input placeholder="图标名" value={form.meta.icon || ''} onChange={(e) => setForm((f) => ({ ...f, meta: { ...f.meta, icon: e.target.value } }))} />
          <Input type="number" placeholder="排序" value={form.sort} onChange={(e) => setForm((f) => ({ ...f, sort: Number(e.target.value) }))} />
          <div className="flex items-center gap-4 text-sm">
            <label className="flex items-center gap-1"><input type="checkbox" checked={form.hidden} onChange={(e) => setForm((f) => ({ ...f, hidden: e.target.checked }))} />隐藏</label>
            <label className="flex items-center gap-1"><input type="checkbox" checked={!!form.meta.keepAlive} onChange={(e) => setForm((f) => ({ ...f, meta: { ...f.meta, keepAlive: e.target.checked } }))} />KeepAlive</label>
            <label className="flex items-center gap-1"><input type="checkbox" checked={!!form.meta.closeTab} onChange={(e) => setForm((f) => ({ ...f, meta: { ...f.meta, closeTab: e.target.checked } }))} />CloseTab</label>
          </div>
          <div className="flex gap-2">
            <Button type="submit">{editing ? '保存编辑' : '新增菜单'}</Button>
            <Button type="button" variant="outline" onClick={() => { setEditing(false); setForm(emptyForm) }}>重置</Button>
            <Button type="button" variant="secondary" onClick={() => onAddChild('0')}>新增根菜单</Button>
          </div>
        </form>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Path</TableHead>
              <TableHead>父ID</TableHead>
              <TableHead>组件</TableHead>
              <TableHead>标题</TableHead>
              <TableHead>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>{row.ID}</TableCell>
                <TableCell>{row.name}</TableCell>
                <TableCell>{row.path}</TableCell>
                <TableCell>{row.parentId}</TableCell>
                <TableCell className="max-w-[320px] truncate">{row.component}</TableCell>
                <TableCell>{row.meta?.title}</TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="outline" onClick={() => onAddChild(row.ID)}>子菜单</Button>
                  <Button size="sm" variant="outline" onClick={() => onEdit(row.ID)}>编辑</Button>
                  <Button size="sm" variant="destructive" onClick={() => onDelete(row.ID)}>删除</Button>
                </TableCell>
              </TableRow>
            ))}
            {!items.length && !loading && <TableRow><TableCell colSpan={7} className="text-center text-muted-foreground">暂无数据</TableCell></TableRow>}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
