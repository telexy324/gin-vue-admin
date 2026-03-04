import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

const points = [62, 78, 66, 92, 88, 97, 110]

function toPolyline(values: number[]) {
  const max = Math.max(...values)
  const min = Math.min(...values)
  const gap = max - min || 1
  return values
    .map((value, index) => {
      const x = (index / (values.length - 1)) * 560 + 20
      const y = 160 - ((value - min) / gap) * 120
      return `${x},${y}`
    })
    .join(' ')
}

export default function DashboardLineChart() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>访问趋势</CardTitle>
      </CardHeader>
      <CardContent>
        <svg viewBox="0 0 600 180" className="h-[220px] w-full rounded-md bg-slate-50">
          <polyline points={toPolyline(points)} fill="none" stroke="#2563eb" strokeWidth="3" strokeLinecap="round" />
          {points.map((v, i) => {
            const x = (i / (points.length - 1)) * 560 + 20
            const max = Math.max(...points)
            const min = Math.min(...points)
            const y = 160 - ((v - min) / (max - min || 1)) * 120
            return <circle key={i} cx={x} cy={y} r="4" fill="#1d4ed8" />
          })}
        </svg>
      </CardContent>
    </Card>
  )
}
