import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

const tasks = [
  { name: '资产巡检', status: 'success', owner: 'ops-1' },
  { name: '日志归档', status: 'running', owner: 'ops-2' },
  { name: '数据库备份', status: 'failed', owner: 'dba-1' },
  { name: '证书续期', status: 'success', owner: 'sec-1' }
] as const

function statusBadge(status: string) {
  if (status === 'success') return <Badge className="bg-emerald-500 hover:bg-emerald-500">成功</Badge>
  if (status === 'running') return <Badge className="bg-blue-500 hover:bg-blue-500">执行中</Badge>
  return <Badge className="bg-rose-500 hover:bg-rose-500">失败</Badge>
}

export default function Dashboard2Table() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>最近任务</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {tasks.map((task) => (
            <div key={task.name} className="flex items-center justify-between rounded-md border p-3">
              <div>
                <p className="text-sm font-medium">{task.name}</p>
                <p className="text-xs text-muted-foreground">负责人：{task.owner}</p>
              </div>
              {statusBadge(task.status)}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
