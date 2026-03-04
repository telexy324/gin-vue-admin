import { useEffect, useState } from 'react'
import { getTaskList } from '@/api/task'
import { getTemplateList } from '@/api/template'
import { getUserList } from '@/api/user'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

export default function TaskListPage() {
  const [items, setItems] = useState<any[]>([])
  const [templates, setTemplates] = useState<any[]>([])
  const [users, setUsers] = useState<any[]>([])
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const query = async () => {
    const res = await getTaskList({ page, pageSize })
    if (res?.code === 0) {
      setItems(res.data?.list || [])
      setTotal(res.data?.total || 0)
    }
  }

  useEffect(() => { query() }, [page, pageSize])
  useEffect(() => {
    const load = async () => {
      const [tr, ur] = await Promise.all([
        getTemplateList({ page: 1, pageSize: 99999 }),
        getUserList({ page: 1, pageSize: 99999 })
      ])
      if (tr?.code === 0) setTemplates(tr.data?.list || [])
      if (ur?.code === 0) setUsers(ur.data?.list || [])
    }
    load()
  }, [])

  const templateName = (id: number) => templates.find((x) => x.ID === id)?.name || '-'
  const userName = (id: number) => id === 999999 ? '定时任务' : (users.find((x) => x.ID === id)?.userName || '-')

  return (
    <Card>
      <CardHeader><CardTitle>任务列表</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>模板名</TableHead><TableHead>状态</TableHead><TableHead>创建人</TableHead><TableHead>开始时间</TableHead><TableHead>结束时间</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>#{row.ID}</TableCell>
                <TableCell>{templateName(row.templateId)}</TableCell>
                <TableCell>{row.status || '-'}</TableCell>
                <TableCell>{userName(row.systemUserId)}</TableCell>
                <TableCell>{row.beginTime?.Time || '-'}</TableCell>
                <TableCell>{row.endTime?.Time || '-'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground"><span>共 {total} 条</span><div className="flex items-center gap-2"><select className="h-8 rounded-md border bg-background px-2" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)) }}>{[10, 30, 50, 100].map((v) => <option key={v} value={v}>{v}/页</option>)}</select><Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</Button><span>第 {page} 页</span><Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>下一页</Button></div></div>
      </CardContent>
    </Card>
  )
}
