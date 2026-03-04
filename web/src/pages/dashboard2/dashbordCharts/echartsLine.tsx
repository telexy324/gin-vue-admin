import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

const values = [14, 18, 26, 22, 30, 37, 34, 42]

export default function Dashboard2Line() {
  const max = Math.max(...values)

  return (
    <Card>
      <CardHeader>
        <CardTitle>任务执行趋势</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {values.map((value, index) => (
            <div key={index} className="flex items-center gap-2">
              <span className="w-20 text-xs text-muted-foreground">第 {index + 1} 周</span>
              <div className="h-2 flex-1 rounded bg-slate-100">
                <div className="h-2 rounded bg-emerald-500" style={{ width: `${(value / max) * 100}%` }} />
              </div>
              <span className="w-10 text-right text-xs">{value}</span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
