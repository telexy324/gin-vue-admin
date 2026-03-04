import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/stores/useAuthStore'

export default function PersonPage() {
  const userInfo = useAuthStore((s) => s.userInfo)

  return (
    <Card>
      <CardHeader>
        <CardTitle>个人信息</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3 text-sm">
        <div>昵称：{userInfo?.nickName || '-'}</div>
        <div>UUID：{userInfo?.uuid || '-'}</div>
        <div className="flex items-center gap-2">
          角色：<Badge variant="secondary">{userInfo?.authority?.authorityName || '-'}</Badge>
        </div>
      </CardContent>
    </Card>
  )
}
