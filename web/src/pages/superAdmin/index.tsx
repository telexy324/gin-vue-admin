import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function SuperAdminIndexPage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>系统管理</CardTitle>
        <CardDescription>已迁移到 React + shadcn/ui。请选择左侧菜单进入具体功能。</CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">支持模块：用户、API、菜单、角色、字典、日志。</CardContent>
    </Card>
  )
}
