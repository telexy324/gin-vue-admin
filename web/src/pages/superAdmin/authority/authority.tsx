import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { copyAuthority, createAuthority, deleteAuthority, getAuthorityList, updateAuthority } from '@/api/authority'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Authority = {
  authorityId: string
  authorityName: string
  parentId?: string
  dataAuthorityId?: string[]
}

type AuthorityForm = { authorityId: string; authorityName: string; parentId: string }
const emptyForm: AuthorityForm = { authorityId: '', authorityName: '', parentId: '0' }

const flatten = (list: any[] = []) => {
  const result: Authority[] = []
  const walk = (arr: any[]) => {
    arr.forEach((item) => {
      result.push(item)
      if (item.children?.length) walk(item.children)
    })
  }
  walk(list)
  return result
}

export default function AuthorityPage() {
  const [items, setItems] = useState<Authority[]>([])
  const [form, setForm] = useState<AuthorityForm>(emptyForm)
  const [mode, setMode] = useState<'add' | 'edit' | 'copy'>('add')
  const [copyFrom, setCopyFrom] = useState<Authority | null>(null)

  const query = async () => {
    const res = await getAuthorityList({ page: 1, pageSize: 999 })
    if (res?.code === 0) {
      setItems(flatten(res.data?.list || []))
    }
  }

  useEffect(() => {
    query()
  }, [])

  const onSave = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.authorityId || !form.authorityName || !form.parentId) {
      toast.error('请填写完整字段')
      return
    }
    if (form.authorityId === '0') {
      toast.error('角色ID不能为0')
      return
    }

    let res
    if (mode === 'add') res = await createAuthority(form)
    if (mode === 'edit') res = await updateAuthority(form)
    if (mode === 'copy' && copyFrom) {
      res = await copyAuthority({
        authority: {
          authorityId: form.authorityId,
          authorityName: form.authorityName,
          parentId: form.parentId,
          dataAuthorityId: copyFrom.dataAuthorityId || []
        },
        oldAuthorityId: Number(copyFrom.authorityId)
      })
    }
    if (res?.code === 0) {
      toast.success('保存成功')
      setForm(emptyForm)
      setMode('add')
      setCopyFrom(null)
      query()
    }
  }

  const onDelete = async (row: Authority) => {
    if (!window.confirm(`确认删除角色 ${row.authorityName} ?`)) return
    const res = await deleteAuthority({ authorityId: row.authorityId })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>角色管理</CardTitle>
        <CardDescription>角色菜单/API/数据权限配置页将在下一批迁移，此页先支持角色增删改拷贝。</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <form className="grid gap-2 rounded-md border p-3 md:grid-cols-4" onSubmit={onSave}>
          <Input placeholder="角色ID" value={form.authorityId} disabled={mode === 'edit'} onChange={(e) => setForm((f) => ({ ...f, authorityId: e.target.value }))} />
          <Input placeholder="角色名称" value={form.authorityName} onChange={(e) => setForm((f) => ({ ...f, authorityName: e.target.value }))} />
          <Input placeholder="父级角色ID" value={form.parentId} onChange={(e) => setForm((f) => ({ ...f, parentId: e.target.value }))} />
          <div className="flex gap-2">
            <Button type="submit">{mode === 'add' ? '新增' : mode === 'edit' ? '保存编辑' : '执行拷贝'}</Button>
            <Button type="button" variant="outline" onClick={() => { setMode('add'); setForm(emptyForm); setCopyFrom(null) }}>重置</Button>
          </div>
        </form>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>角色ID</TableHead>
              <TableHead>角色名称</TableHead>
              <TableHead>父ID</TableHead>
              <TableHead>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.authorityId}>
                <TableCell>{row.authorityId}</TableCell>
                <TableCell>{row.authorityName}</TableCell>
                <TableCell>{row.parentId || '-'}</TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="outline" onClick={() => { setMode('add'); setForm({ ...emptyForm, parentId: row.authorityId }) }}>子角色</Button>
                  <Button size="sm" variant="outline" onClick={() => { setMode('edit'); setForm({ authorityId: row.authorityId, authorityName: row.authorityName, parentId: row.parentId || '0' }) }}>编辑</Button>
                  <Button size="sm" variant="secondary" onClick={() => { setMode('copy'); setCopyFrom(row); setForm({ authorityId: '', authorityName: `${row.authorityName}-copy`, parentId: row.parentId || '0' }) }}>拷贝</Button>
                  <Button size="sm" variant="destructive" onClick={() => onDelete(row)}>删除</Button>
                </TableCell>
              </TableRow>
            ))}
            {!items.length && <TableRow><TableCell colSpan={4} className="text-center text-muted-foreground">暂无数据</TableCell></TableRow>}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
