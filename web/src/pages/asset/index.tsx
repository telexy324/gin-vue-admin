import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function AssetIndexPage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>资产管理</CardTitle>
        <CardDescription>服务器、系统关系、SSH 执行、日志上传配置已迁移到 React。</CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">从左侧菜单进入具体资产页面。</CardContent>
    </Card>
  )
}
