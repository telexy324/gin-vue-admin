import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { deleteUser, getUserList } from '@/api/user'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

export default function UserManagePage() {
  const [loading, setLoading] = useState(false)
  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [rows, setRows] = useState([])

  const query = async () => {
    setLoading(true)
    try {
      const res = await getUserList({ page, pageSize, keyword })
      if (res?.code === 0) {
        setRows(res.data?.list || [])
        setTotal(res.data?.total || 0)
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    query()
  }, [page])

  const onDelete = async (row) => {
    if (!window.confirm(`确认删除用户 ${row.userName || row.nickName} ?`)) return
    const res = await deleteUser({ ID: row.ID })
    if (res?.code === 0) {
      toast.success('删除成功')
      query()
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>用户管理</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Input value={keyword} onChange={(e) => setKeyword(e.target.value)} placeholder="昵称/用户名关键词" />
          <Button
            onClick={() => {
              setPage(1)
              query()
            }}
          >
            查询
          </Button>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>用户名</TableHead>
              <TableHead>昵称</TableHead>
              <TableHead>邮箱</TableHead>
              <TableHead>角色</TableHead>
              <TableHead className="w-[100px]">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>{row.userName}</TableCell>
                <TableCell>{row.nickName}</TableCell>
                <TableCell>{row.email}</TableCell>
                <TableCell>{row.authority?.authorityName || '-'}</TableCell>
                <TableCell>
                  <Button variant="destructive" size="sm" onClick={() => onDelete(row)}>
                    删除
                  </Button>
                </TableCell>
              </TableRow>
            ))}
            {!rows.length && !loading && (
              <TableRow>
                <TableCell colSpan={5} className="text-center text-muted-foreground">
                  暂无数据
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>

        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <span>共 {total} 条</span>
          <div className="space-x-2">
            <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
              上一页
            </Button>
            <span>第 {page} 页</span>
            <Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p) => p + 1)}>
              下一页
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
