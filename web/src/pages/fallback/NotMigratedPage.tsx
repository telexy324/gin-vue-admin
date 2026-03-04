import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function NotMigratedPage({ entry }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{entry?.title || '未迁移页面'}</CardTitle>
        <CardDescription>该路由已连通，但页面尚未迁移到 React。</CardDescription>
      </CardHeader>
      <CardContent>
        <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs text-muted-foreground">{JSON.stringify(entry?.raw || entry, null, 2)}</pre>
      </CardContent>
    </Card>
  )
}
