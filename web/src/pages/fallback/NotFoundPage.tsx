import { Link } from 'react-router-dom'
import { buttonVariants } from '@/components/ui/button'

export default function NotFoundPage() {
  return (
    <div className="grid min-h-[60vh] place-items-center">
      <div className="space-y-3 text-center">
        <h1 className="text-2xl font-semibold">404</h1>
        <p className="text-muted-foreground">页面不存在或无访问权限。</p>
        <Link className={buttonVariants()} to="/layout">
          返回首页
        </Link>
      </div>
    </div>
  )
}
