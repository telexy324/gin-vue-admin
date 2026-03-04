import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { deleteApplicationRecord, deleteApplicationRecordByIds, exportApplicationRecord, getApplicationRecordList } from '@/api/applicationRecord'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Log = {
  ID: number
  action: string
  detail: string
  status: number
  errorMessage: string
  ip: string
  logTime?: string
  user?: { userName?: string; nickName?: string }
}

export default function ApplicationRecordPage() {
  const [items, setItems] = useState<Log[]>([])
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const [search, setSearch] = useState({ action: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getApplicationRecordList({ page, pageSize, ...search })
    if (res?.code === 0) {
      setItems(res.data?.list || [])
      setTotal(res.data?.total || 0)
    }
  }

  useEffect(() => {
    query()
  }, [page, pageSize])

  const onDeleteOne = async (id: number) => {
    if (!window.confirm('确认删除该日志?')) return
    const res = await deleteApplicationRecord({ ID: id })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  const onDeleteBatch = async () => {
    if (!selectedIds.length) return
    if (!window.confirm(`确认删除 ${selectedIds.length} 条日志?`)) return
    const res = await deleteApplicationRecordByIds({ ids: selectedIds })
    if (res?.code === 0) {
      toast.success('批量删除成功')
      setSelectedIds([])
      query()
    }
  }

  const onExport = async () => {
    if (!selectedIds.length) {
      toast.error('请先勾选要导出的日志')
      return
    }
    const fileName = `logRecord-${Date.now()}.xlsx`
    await exportApplicationRecord({ ids: selectedIds }, fileName)
    toast.success('导出任务已触发')
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>登录日志</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-3">
          <Input placeholder="请求路径" value={search.action} onChange={(e) => setSearch((s) => ({ ...s, action: e.target.value }))} />
          <div className="flex gap-2 md:col-span-2">
            <Button onClick={() => { setPage(1); query() }}>查询</Button>
            <Button variant="outline" onClick={() => { setSearch({ action: '' }); setPage(1) }}>重置</Button>
            <Button variant="destructive" onClick={onDeleteBatch} disabled={!selectedIds.length}>删除</Button>
            <Button variant="secondary" onClick={onExport} disabled={!selectedIds.length}>导出</Button>
          </div>
        </div>

        <Table>
          <TableHeader><TableRow><TableHead>选</TableHead><TableHead>操作人</TableHead><TableHead>请求路径</TableHead><TableHead>详情</TableHead><TableHead>状态</TableHead><TableHead>错误信息</TableHead><TableHead>IP</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell><input type="checkbox" checked={selectedIds.includes(row.ID)} onChange={(e) => setSelectedIds((ids) => e.target.checked ? [...new Set([...ids, row.ID])] : ids.filter((id) => id !== row.ID))} /></TableCell>
                <TableCell>{row.user?.userName || '-'}({row.user?.nickName || '-'})</TableCell>
                <TableCell className="max-w-[220px] truncate">{row.action}</TableCell>
                <TableCell className="max-w-[240px] truncate">{row.detail}</TableCell>
                <TableCell>{row.status === 0 ? '成功' : '失败'}</TableCell>
                <TableCell className="max-w-[220px] truncate">{row.errorMessage || '-'}</TableCell>
                <TableCell>{row.ip}</TableCell>
                <TableCell><Button size="sm" variant="destructive" onClick={() => onDeleteOne(row.ID)}>删除</Button></TableCell>
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
