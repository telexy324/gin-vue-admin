import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function TaskIndexPage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>任务中心</CardTitle>
        <CardDescription>任务、模板、定时调度、模板集页面已迁移基础操作。</CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">从左侧菜单选择具体任务模块。</CardContent>
    </Card>
  )
}
