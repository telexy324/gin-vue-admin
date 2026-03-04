import { useEffect, useMemo, useRef, useState } from 'react'
import { useLocation } from 'react-router-dom'
import { getServerRelations } from '@/api/graph'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

type GraphNode = {
  id?: string | number
  name?: string
  category?: string | number
  [key: string]: any
}

type GraphLink = {
  source?: string | number
  target?: string | number
  [key: string]: any
}

const normalizeGraph = (raw: any): { nodes: GraphNode[]; links: GraphLink[] } => {
  if (!raw) return { nodes: [], links: [] }
  if (Array.isArray(raw.nodes) || Array.isArray(raw.links)) {
    return { nodes: raw.nodes || [], links: raw.links || raw.edges || [] }
  }
  if (raw.graph && (Array.isArray(raw.graph.nodes) || Array.isArray(raw.graph.links))) {
    return { nodes: raw.graph.nodes || [], links: raw.graph.links || raw.graph.edges || [] }
  }
  if (raw.data && (Array.isArray(raw.data.nodes) || Array.isArray(raw.data.links))) {
    return { nodes: raw.data.nodes || [], links: raw.data.links || raw.data.edges || [] }
  }
  return { nodes: [], links: [] }
}

const getNodeId = (n: GraphNode, idx: number) => String(n.id ?? n.ID ?? n.name ?? idx)

export default function GraphPage() {
  const location = useLocation()
  const [systemId, setSystemId] = useState<number>((location.state as any)?.systemId || 0)
  const [result, setResult] = useState<any>(null)
  const [scale, setScale] = useState(1)
  const [offset, setOffset] = useState({ x: 0, y: 0 })
  const [dragging, setDragging] = useState(false)
  const [selectedNodeId, setSelectedNodeId] = useState<string>('')
  const dragStart = useRef({ x: 0, y: 0 })

  const query = async () => {
    if (!systemId) return
    const res = await getServerRelations({ id: Number(systemId) })
    if (res?.code === 0) setResult(res.data)
  }

  useEffect(() => {
    if (systemId) query()
  }, [])

  const { nodes, links } = normalizeGraph(result)

  const width = 980
  const height = 560
  const nodePositions = useMemo(() => {
    const map = new Map<string, { x: number; y: number }>()
    nodes.forEach((n, i) => {
      const col = i % 5
      const row = Math.floor(i / 5)
      const x = 100 + col * 180
      const y = 80 + row * 110
      map.set(getNodeId(n, i), { x, y })
    })
    return map
  }, [nodes])

  const selectedNode = useMemo(() => {
    if (!selectedNodeId) return null
    return nodes.find((n, i) => getNodeId(n, i) === selectedNodeId) || null
  }, [nodes, selectedNodeId])

  const onWheel: React.WheelEventHandler<SVGSVGElement> = (e) => {
    e.preventDefault()
    const next = Math.max(0.4, Math.min(2.5, scale + (e.deltaY < 0 ? 0.1 : -0.1)))
    setScale(next)
  }

  const onMouseDown: React.MouseEventHandler<SVGSVGElement> = (e) => {
    setDragging(true)
    dragStart.current = { x: e.clientX - offset.x, y: e.clientY - offset.y }
  }

  const onMouseMove: React.MouseEventHandler<SVGSVGElement> = (e) => {
    if (!dragging) return
    setOffset({ x: e.clientX - dragStart.current.x, y: e.clientY - dragStart.current.y })
  }

  const onMouseUp: React.MouseEventHandler<SVGSVGElement> = () => {
    setDragging(false)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>系统关系图</CardTitle>
        <CardDescription>支持缩放、拖拽和平面关系可视化。</CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex gap-2">
          <Input type="number" placeholder="系统ID" value={systemId || ''} onChange={(e) => setSystemId(Number(e.target.value))} />
          <Button onClick={query}>查询关系</Button>
          <Button variant="outline" onClick={() => { setScale(1); setOffset({ x: 0, y: 0 }) }}>重置视图</Button>
        </div>
        <div className="flex gap-2">
          <Badge variant="secondary">节点 {nodes.length}</Badge>
          <Badge variant="secondary">连线 {links.length}</Badge>
          <Badge variant="secondary">缩放 {scale.toFixed(1)}x</Badge>
        </div>
        <div className="overflow-auto rounded-md border bg-card p-2">
          <svg
            width={width}
            height={height}
            className="min-w-[980px] cursor-grab"
            onWheel={onWheel}
            onMouseDown={onMouseDown}
            onMouseMove={onMouseMove}
            onMouseUp={onMouseUp}
            onMouseLeave={onMouseUp}
          >
            <g transform={`translate(${offset.x},${offset.y}) scale(${scale})`}>
              {links.map((link, i) => {
                const s = nodePositions.get(String(link.source))
                const t = nodePositions.get(String(link.target))
                if (!s || !t) return null
                return <line key={`l-${i}`} x1={s.x} y1={s.y} x2={t.x} y2={t.y} stroke="#9ca3af" strokeWidth="1.5" />
              })}
              {nodes.map((node, i) => {
                const id = getNodeId(node, i)
                const p = nodePositions.get(id)
                if (!p) return null
                const label = node.name || id
                const active = selectedNodeId === id
                return (
                  <g key={id} onClick={() => setSelectedNodeId(id)} className="cursor-pointer">
                    <circle cx={p.x} cy={p.y} r={active ? 26 : 22} fill={active ? '#2563eb' : '#3b82f6'} opacity="0.9" />
                    <text x={p.x} y={p.y + 4} textAnchor="middle" fill="#fff" fontSize="10">
                      {String(label).slice(0, 10)}
                    </text>
                  </g>
                )
              })}
            </g>
          </svg>
        </div>

        <div className="grid gap-3 md:grid-cols-2">
          <pre className="max-h-[42vh] overflow-auto rounded-md bg-muted p-3 text-xs">{JSON.stringify(result || { tip: '请先输入系统ID并查询' }, null, 2)}</pre>
          <pre className="max-h-[42vh] overflow-auto rounded-md bg-muted p-3 text-xs">{JSON.stringify(selectedNode || { tip: '点击图中节点查看详情' }, null, 2)}</pre>
        </div>
      </CardContent>
    </Card>
  )
}
