import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const rows = [
  { name: '华东机房', load: '63%', alarm: 1, updatedAt: '2026-03-05 10:15' },
  { name: '华南机房', load: '49%', alarm: 0, updatedAt: '2026-03-05 10:12' },
  { name: '研发集群', load: '81%', alarm: 2, updatedAt: '2026-03-05 10:09' },
  { name: '备份节点', load: '31%', alarm: 0, updatedAt: '2026-03-05 10:08' }
]

export default function DashboardTable() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>系统概览</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>节点</TableHead>
              <TableHead>负载</TableHead>
              <TableHead>告警数</TableHead>
              <TableHead>更新时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map((row) => (
              <TableRow key={row.name}>
                <TableCell>{row.name}</TableCell>
                <TableCell>{row.load}</TableCell>
                <TableCell>{row.alarm}</TableCell>
                <TableCell>{row.updatedAt}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
