import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { deleteSysOperationRecord, deleteSysOperationRecordByIds, getSysOperationRecordList } from '@/api/sysOperationRecord'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Log = {
  ID: number
  method: string
  path: string
  status: number
  ip: string
  body?: string
  resp?: string
  user?: { userName?: string; nickName?: string }
  CreatedAt?: string
}

export default function SysOperationRecordPage() {
  const [items, setItems] = useState<Log[]>([])
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const [search, setSearch] = useState({ method: '', path: '', status: '' })
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getSysOperationRecordList({ page, pageSize, ...search })
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
    const res = await deleteSysOperationRecord({ ID: id })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  const onDeleteBatch = async () => {
    if (!selectedIds.length) return
    if (!window.confirm(`确认删除 ${selectedIds.length} 条日志?`)) return
    const res = await deleteSysOperationRecordByIds({ ids: selectedIds })
    if (res?.code === 0) {
      toast.success('批量删除成功')
      setSelectedIds([])
      query()
    }
  }

  const pretty = (v?: string) => {
    if (!v) return '无'
    try {
      return JSON.stringify(JSON.parse(v), null, 2)
    } catch {
      return v
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>操作日志</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-4">
          <Input placeholder="请求方法" value={search.method} onChange={(e) => setSearch((s) => ({ ...s, method: e.target.value }))} />
          <Input placeholder="请求路径" value={search.path} onChange={(e) => setSearch((s) => ({ ...s, path: e.target.value }))} />
          <Input placeholder="状态码" value={search.status} onChange={(e) => setSearch((s) => ({ ...s, status: e.target.value }))} />
          <div className="flex gap-2">
            <Button onClick={() => { setPage(1); query() }}>查询</Button>
            <Button variant="outline" onClick={() => { setSearch({ method: '', path: '', status: '' }); setPage(1); }}>重置</Button>
            <Button variant="destructive" onClick={onDeleteBatch} disabled={!selectedIds.length}>删除</Button>
          </div>
        </div>

        <Table>
          <TableHeader><TableRow><TableHead>选</TableHead><TableHead>操作人</TableHead><TableHead>方法</TableHead><TableHead>路径</TableHead><TableHead>状态</TableHead><TableHead>请求IP</TableHead><TableHead>请求体</TableHead><TableHead>响应体</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell><input type="checkbox" checked={selectedIds.includes(row.ID)} onChange={(e) => setSelectedIds((ids) => e.target.checked ? [...new Set([...ids, row.ID])] : ids.filter((id) => id !== row.ID))} /></TableCell>
                <TableCell>{row.user?.userName || '-'}({row.user?.nickName || '-'})</TableCell>
                <TableCell>{row.method}</TableCell>
                <TableCell className="max-w-[200px] truncate">{row.path}</TableCell>
                <TableCell>{row.status}</TableCell>
                <TableCell>{row.ip}</TableCell>
                <TableCell><details><summary>查看</summary><pre className="max-h-52 overflow-auto text-xs">{pretty(row.body)}</pre></details></TableCell>
                <TableCell><details><summary>查看</summary><pre className="max-h-52 overflow-auto text-xs">{pretty(row.resp)}</pre></details></TableCell>
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
