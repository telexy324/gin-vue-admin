import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { delSysHistory, getSysHistory, rollback } from '@/api/autoCode'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type History = { ID: number; structName?: string; tableName?: string; packageName?: string; createdAt?: string }

export default function AutoCodeAdminPage() {
  const [items, setItems] = useState<History[]>([])
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getSysHistory({ page, pageSize })
    if (res?.code === 0) {
      setItems(res.data?.list || [])
      setTotal(res.data?.total || 0)
    }
  }

  useEffect(() => {
    query()
  }, [page, pageSize])

  const onRollback = async (row: History) => {
    if (!window.confirm(`确认回滚 #${row.ID} ?`)) return
    const res = await rollback({ ID: row.ID })
    if (res?.code === 0) {
      toast.success('回滚成功')
      query()
    }
  }

  const onDelete = async (row: History) => {
    if (!window.confirm(`确认删除 #${row.ID} ?`)) return
    const res = await delSysHistory({ ID: row.ID })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>自动代码历史</CardTitle>
        <CardDescription>支持回滚与删除历史记录。</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>Struct</TableHead><TableHead>Table</TableHead><TableHead>Package</TableHead><TableHead>时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>{row.ID}</TableCell><TableCell>{row.structName}</TableCell><TableCell>{row.tableName}</TableCell><TableCell>{row.packageName}</TableCell><TableCell>{row.createdAt || '-'}</TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="secondary" onClick={() => onRollback(row)}>回滚</Button>
                  <Button size="sm" variant="destructive" onClick={() => onDelete(row)}>删除</Button>
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
