import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { createSysDictionary, deleteSysDictionary, findSysDictionary, getSysDictionaryList, updateSysDictionary } from '@/api/sysDictionary'
import { createSysDictionaryDetail, deleteSysDictionaryDetail, findSysDictionaryDetail, getSysDictionaryDetailList, updateSysDictionaryDetail } from '@/api/sysDictionaryDetail'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Dictionary = { ID: number; name: string; type: string; status: boolean; desc: string }
type Detail = { ID: number; label: string; value: number; status: boolean; sort: number; sysDictionaryID: number }

const emptyDict = { name: '', type: '', status: true, desc: '' }
const emptyDetail = { label: '', value: 0, status: true, sort: 0 }

export default function SysDictionaryPage() {
  const [items, setItems] = useState<Dictionary[]>([])
  const [search, setSearch] = useState({ name: '', type: '', status: '', desc: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const [dictForm, setDictForm] = useState<any>(emptyDict)
  const [dictEditing, setDictEditing] = useState(false)

  const [selectedDict, setSelectedDict] = useState<Dictionary | null>(null)
  const [detailItems, setDetailItems] = useState<Detail[]>([])
  const [detailForm, setDetailForm] = useState<any>(emptyDetail)
  const [detailEditing, setDetailEditing] = useState(false)

  const query = async () => {
    const params: any = { page, pageSize, ...search }
    if (params.status === '') delete params.status
    const res = await getSysDictionaryList(params)
    if (res?.code === 0) {
      setItems(res.data?.list || [])
      setTotal(res.data?.total || 0)
    }
  }

  useEffect(() => {
    query()
  }, [page, pageSize])

  const loadDetails = async (dictionaryId: number) => {
    const res = await getSysDictionaryDetailList({ page: 1, pageSize: 999, sysDictionaryID: dictionaryId })
    if (res?.code === 0) {
      setDetailItems(res.data?.list || [])
    }
  }

  const onSelectDict = async (row: Dictionary) => {
    setSelectedDict(row)
    setDetailEditing(false)
    setDetailForm(emptyDetail)
    loadDetails(row.ID)
  }

  const onSaveDict = async (e: FormEvent) => {
    e.preventDefault()
    if (!dictForm.name || !dictForm.type || !dictForm.desc) {
      toast.error('请填写完整字典信息')
      return
    }
    const res = dictEditing ? await updateSysDictionary(dictForm) : await createSysDictionary(dictForm)
    if (res?.code === 0) {
      toast.success(dictEditing ? '更新成功' : '创建成功')
      setDictEditing(false)
      setDictForm(emptyDict)
      query()
    }
  }

  const onEditDict = async (row: Dictionary) => {
    const res = await findSysDictionary({ ID: row.ID })
    if (res?.code === 0) {
      setDictForm(res.data?.resysDictionary)
      setDictEditing(true)
    }
  }

  const onDeleteDict = async (row: Dictionary) => {
    if (!window.confirm(`确认删除字典 ${row.name} ?`)) return
    const res = await deleteSysDictionary({ ID: row.ID })
    if (res?.code === 0) {
      toast.success('删除成功')
      if (selectedDict?.ID === row.ID) {
        setSelectedDict(null)
        setDetailItems([])
      }
      query()
    }
  }

  const onSaveDetail = async (e: FormEvent) => {
    e.preventDefault()
    if (!selectedDict) {
      toast.error('请先选择字典')
      return
    }
    if (!detailForm.label) {
      toast.error('请输入展示值')
      return
    }
    const payload = { ...detailForm, sysDictionaryID: selectedDict.ID }
    const res = detailEditing ? await updateSysDictionaryDetail(payload) : await createSysDictionaryDetail(payload)
    if (res?.code === 0) {
      toast.success(detailEditing ? '字典项更新成功' : '字典项创建成功')
      setDetailEditing(false)
      setDetailForm(emptyDetail)
      loadDetails(selectedDict.ID)
    }
  }

  const onEditDetail = async (row: Detail) => {
    const res = await findSysDictionaryDetail({ ID: row.ID })
    if (res?.code === 0) {
      setDetailForm(res.data?.resysDictionaryDetail)
      setDetailEditing(true)
    }
  }

  const onDeleteDetail = async (row: Detail) => {
    if (!window.confirm('确认删除该字典项?')) return
    const res = await deleteSysDictionaryDetail({ ID: row.ID })
    if (res?.code === 0 && selectedDict) {
      toast.success('删除成功')
      loadDetails(selectedDict.ID)
    }
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>字典管理</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-2 md:grid-cols-5">
            <Input placeholder="字典名(中)" value={search.name} onChange={(e) => setSearch((s) => ({ ...s, name: e.target.value }))} />
            <Input placeholder="字典名(英)" value={search.type} onChange={(e) => setSearch((s) => ({ ...s, type: e.target.value }))} />
            <select className="h-10 rounded-md border bg-background px-3 text-sm" value={search.status} onChange={(e) => setSearch((s) => ({ ...s, status: e.target.value }))}>
              <option value="">全部状态</option><option value="true">启用</option><option value="false">停用</option>
            </select>
            <Input placeholder="描述" value={search.desc} onChange={(e) => setSearch((s) => ({ ...s, desc: e.target.value }))} />
            <div className="flex gap-2">
              <Button onClick={() => { setPage(1); query() }}>查询</Button>
              <Button variant="outline" onClick={() => { setSearch({ name: '', type: '', status: '', desc: '' }); setPage(1); }}>重置</Button>
            </div>
          </div>

          <form className="grid gap-2 rounded-md border p-3 md:grid-cols-5" onSubmit={onSaveDict}>
            <Input placeholder="字典名(中)" value={dictForm.name || ''} onChange={(e) => setDictForm((f: any) => ({ ...f, name: e.target.value }))} />
            <Input placeholder="字典名(英)" value={dictForm.type || ''} onChange={(e) => setDictForm((f: any) => ({ ...f, type: e.target.value }))} />
            <Input placeholder="描述" value={dictForm.desc || ''} onChange={(e) => setDictForm((f: any) => ({ ...f, desc: e.target.value }))} />
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!dictForm.status} onChange={(e) => setDictForm((f: any) => ({ ...f, status: e.target.checked }))} />启用</label>
            <div className="flex gap-2">
              <Button type="submit">{dictEditing ? '保存编辑' : '新增字典'}</Button>
              {dictEditing && <Button type="button" variant="outline" onClick={() => { setDictEditing(false); setDictForm(emptyDict) }}>取消</Button>}
            </div>
          </form>

          <Table>
            <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>中文名</TableHead><TableHead>英文名</TableHead><TableHead>状态</TableHead><TableHead>描述</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
            <TableBody>
              {items.map((row) => (
                <TableRow key={row.ID} className={selectedDict?.ID === row.ID ? 'bg-muted/60' : ''}>
                  <TableCell>{row.ID}</TableCell><TableCell>{row.name}</TableCell><TableCell>{row.type}</TableCell><TableCell>{row.status ? '启用' : '停用'}</TableCell><TableCell>{row.desc}</TableCell>
                  <TableCell className="space-x-2">
                    <Button size="sm" variant="secondary" onClick={() => onSelectDict(row)}>详情</Button>
                    <Button size="sm" variant="outline" onClick={() => onEditDict(row)}>编辑</Button>
                    <Button size="sm" variant="destructive" onClick={() => onDeleteDict(row)}>删除</Button>
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

      <Card>
        <CardHeader>
          <CardTitle>字典项管理</CardTitle>
          <CardDescription>{selectedDict ? `当前字典：${selectedDict.name} (${selectedDict.type})` : '请先在上方选择字典'}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <form className="grid gap-2 rounded-md border p-3 md:grid-cols-5" onSubmit={onSaveDetail}>
            <Input placeholder="展示值" value={detailForm.label || ''} onChange={(e) => setDetailForm((f: any) => ({ ...f, label: e.target.value }))} />
            <Input type="number" placeholder="字典值" value={detailForm.value ?? 0} onChange={(e) => setDetailForm((f: any) => ({ ...f, value: Number(e.target.value) }))} />
            <Input type="number" placeholder="排序" value={detailForm.sort ?? 0} onChange={(e) => setDetailForm((f: any) => ({ ...f, sort: Number(e.target.value) }))} />
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!detailForm.status} onChange={(e) => setDetailForm((f: any) => ({ ...f, status: e.target.checked }))} />启用</label>
            <div className="flex gap-2">
              <Button type="submit" disabled={!selectedDict}>{detailEditing ? '保存编辑' : '新增字典项'}</Button>
              {detailEditing && <Button type="button" variant="outline" onClick={() => { setDetailEditing(false); setDetailForm(emptyDetail) }}>取消</Button>}
            </div>
          </form>

          <Table>
            <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>展示值</TableHead><TableHead>字典值</TableHead><TableHead>排序</TableHead><TableHead>状态</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
            <TableBody>
              {detailItems.map((row) => (
                <TableRow key={row.ID}>
                  <TableCell>{row.ID}</TableCell><TableCell>{row.label}</TableCell><TableCell>{row.value}</TableCell><TableCell>{row.sort}</TableCell><TableCell>{row.status ? '启用' : '停用'}</TableCell>
                  <TableCell className="space-x-2">
                    <Button size="sm" variant="outline" onClick={() => onEditDetail(row)}>编辑</Button>
                    <Button size="sm" variant="destructive" onClick={() => onDeleteDetail(row)}>删除</Button>
                  </TableCell>
                </TableRow>
              ))}
              {!detailItems.length && <TableRow><TableCell colSpan={6} className="text-center text-muted-foreground">暂无字典项</TableCell></TableRow>}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}
