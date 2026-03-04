import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function SystemToolsIndexPage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>系统工具</CardTitle>
        <CardDescription>已迁移 React 版本，支持系统配置与自动代码核心功能。</CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">请在左侧继续进入具体工具页面。</CardContent>
    </Card>
  )
}
