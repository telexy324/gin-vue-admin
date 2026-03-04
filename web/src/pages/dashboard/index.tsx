import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore } from '@/stores/useAuthStore'

export default function DashboardPage() {
  const userInfo = useAuthStore((s) => s.userInfo)

  return (
    <div className="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader>
          <CardTitle>欢迎回来</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">当前用户：{userInfo?.nickName || '未知'}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>当前角色</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">{userInfo?.authority?.authorityName || '-'}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>迁移状态</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">已切换到 React + shadcn/ui 架构，页面可按菜单持续迁移。</p>
        </CardContent>
      </Card>
    </div>
  )
}
