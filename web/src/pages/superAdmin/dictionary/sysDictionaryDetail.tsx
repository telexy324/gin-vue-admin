import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getSysDictionaryDetailList } from '@/api/sysDictionaryDetail'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Detail = { ID: number; label: string; value: number; status: boolean; sort: number }

export default function SysDictionaryDetailPage() {
  const { id } = useParams()
  const [items, setItems] = useState<Detail[]>([])

  useEffect(() => {
    const query = async () => {
      if (!id) return
      const res = await getSysDictionaryDetailList({ page: 1, pageSize: 999, sysDictionaryID: Number(id) })
      if (res?.code === 0) setItems(res.data?.list || [])
    }
    query()
  }, [id])

  return (
    <Card>
      <CardHeader>
        <CardTitle>字典详情</CardTitle>
        <CardDescription>Dictionary ID: {id || '-'}</CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>展示值</TableHead><TableHead>字典值</TableHead><TableHead>排序</TableHead><TableHead>状态</TableHead></TableRow></TableHeader>
          <TableBody>
            {items.map((row) => (
              <TableRow key={row.ID}><TableCell>{row.ID}</TableCell><TableCell>{row.label}</TableCell><TableCell>{row.value}</TableCell><TableCell>{row.sort}</TableCell><TableCell>{row.status ? '启用' : '停用'}</TableCell></TableRow>
            ))}
            {!items.length && <TableRow><TableCell colSpan={5} className="text-center text-muted-foreground">暂无数据</TableCell></TableRow>}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
