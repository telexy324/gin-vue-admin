import { lazy, Suspense, useEffect, useMemo, useState } from 'react'
import { Navigate, Route, Routes, useLocation, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { ChevronDown, LogOut, Menu, UserCircle2 } from 'lucide-react'
import { getUserInfo, setUserAuthority } from '@/api/user'
import { asyncMenu } from '@/api/menu'
import { jsonInBlacklist } from '@/api/jwt'
import { useAuthStore } from '@/stores/useAuthStore'
import { useMenuStore } from '@/stores/useMenuStore'
import { createRouteEntriesFromMenus, findPathByRouteName, type MenuItem, type RouteEntry } from '@/router/menuResolver'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import LoginPage from '@/pages/auth/LoginPage'
import InitPage from '@/pages/system/InitPage'
import DashboardPage from '@/pages/dashboard/index'
import PersonPage from '@/pages/person/index'
import NotFoundPage from '@/pages/fallback/NotFoundPage'
import NotMigratedPage from '@/pages/fallback/NotMigratedPage'
import nfhLogo from '@/assets/nfhlogo.jpg'

const Loading = () => <div className="grid h-screen place-items-center text-sm text-muted-foreground">加载中...</div>

const resolveElement = (entry: RouteEntry) => {
  if (!entry?.loader) {
    return <NotMigratedPage entry={entry} />
  }
  const Dynamic = lazy(entry.loader)
  return (
    <Suspense fallback={<div className="p-6 text-sm text-muted-foreground">页面加载中...</div>}>
      <Dynamic />
    </Suspense>
  )
}

function MenuNode({ item, parentPath = '', collapsed }: { item: MenuItem; parentPath?: string; collapsed: boolean }) {
  const navigate = useNavigate()
  const location = useLocation()
  const seg = String(item.path || '').replace(/^\/+|\/+$/g, '')
  const current = [parentPath, seg].filter(Boolean).join('/')
  const to = `/layout/${current.replace(/^layout\/?/, '')}`
  const active = location.pathname === to

  if (item.hidden) return null

  return (
    <div className="space-y-1">
      {item.component && (
        <button
          onClick={() => navigate(to)}
          className={`w-full rounded-md px-3 py-2 text-left text-sm transition ${
            active ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
          }`}
          title={collapsed ? item.meta?.title : ''}
        >
          {collapsed ? (item.meta?.title || item.name || '?').slice(0, 1) : item.meta?.title || item.name}
        </button>
      )}
      {item.children?.length ? (
        <div className={collapsed ? 'pl-0' : 'pl-3'}>
          {item.children.filter((child) => !child.hidden).map((child) => (
            <MenuNode key={`${child.name}-${child.path}`} item={child} parentPath={current} collapsed={collapsed} />
          ))}
        </div>
      ) : null}
    </div>
  )
}

function LayoutPage() {
  const [loading, setLoading] = useState(true)
  const [collapsed, setCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)
  const token = useAuthStore((s) => s.token)
  const userInfo = useAuthStore((s) => s.userInfo)
  const setUserInfo = useAuthStore((s) => s.setUserInfo)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const menus = useMenuStore((s) => s.menus)
  const routeEntries = useMenuStore((s) => s.routeEntries)
  const initialized = useMenuStore((s) => s.initialized)
  const setMenus = useMenuStore((s) => s.setMenus)
  const setRouteEntries = useMenuStore((s) => s.setRouteEntries)
  const setInitialized = useMenuStore((s) => s.setInitialized)
  const clearMenus = useMenuStore((s) => s.clearMenus)
  const navigate = useNavigate()
  const location = useLocation()

  useEffect(() => {
    if (!token) {
      navigate('/login', { replace: true })
      return
    }

    const bootstrap = async () => {
      if (initialized) {
        setLoading(false)
        return
      }
      try {
        const [userRes, menuRes] = await Promise.all([getUserInfo(), asyncMenu()])
        if (userRes?.code !== 0 || menuRes?.code !== 0) {
          throw new Error('初始化用户信息失败')
        }
        const nextMenus = menuRes.data?.menus || []
        setUserInfo(userRes.data?.userInfo || null)
        setMenus(nextMenus)
        setRouteEntries(createRouteEntriesFromMenus(nextMenus))
        setInitialized(true)
      } catch (err) {
        clearAuth()
        clearMenus()
        toast.error(err.message || '登录状态失效，请重新登录')
        navigate('/login', { replace: true })
      } finally {
        setLoading(false)
      }
    }

    bootstrap()
  }, [token, initialized, navigate, clearAuth, clearMenus, setInitialized, setMenus, setRouteEntries, setUserInfo])

  const defaultPath = useMemo(() => {
    const byName = findPathByRouteName(menus, userInfo?.authority?.defaultRouter)
    if (byName) return byName
    if (routeEntries[0]?.path) return `/layout/${routeEntries[0].path}`
    return '/layout/dashboard'
  }, [menus, routeEntries, userInfo])

  const currentTitle = useMemo(() => {
    const pathname = location.pathname.replace(/^\/layout\/?/, '')
    const matched = routeEntries.find((entry) => entry.path === pathname)
    return matched?.title || '控制台'
  }, [location.pathname, routeEntries])

  const logout = async () => {
    try {
      await jsonInBlacklist()
    } catch (_) {
      // ignore server error and continue local logout
    }
    clearAuth()
    clearMenus()
    navigate('/login', { replace: true })
  }

  const switchAuthority = async (authorityId: string) => {
    const res = await setUserAuthority({ authorityId })
    if (res?.code === 0) {
      clearMenus()
      window.location.reload()
    }
  }

  if (loading) return <Loading />

  return (
    <div className="flex h-screen bg-[radial-gradient(circle_at_10%_10%,#d9efff_0,#f8fbff_35%,#f5f7fa_100%)]">
      <aside
        className={`fixed z-30 h-full border-r bg-white/90 backdrop-blur md:static ${
          collapsed ? 'w-[84px]' : 'w-[280px]'
        } ${mobileOpen ? 'left-0' : '-left-[280px] md:left-0'} transition-all`}
      >
        <div className="flex h-16 items-center gap-3 border-b px-4">
          <img src={nfhLogo} alt="logo" className="h-9 w-9 rounded-md object-cover" />
          {!collapsed && <span className="text-sm font-semibold tracking-wide">Gin Admin React</span>}
        </div>
        <div className="h-[calc(100%-64px)] overflow-y-auto p-3">
          {menus.map((menu) => (
            <MenuNode key={`${menu.name}-${menu.path}`} item={menu} collapsed={collapsed} />
          ))}
        </div>
      </aside>

      <div className={`flex min-w-0 flex-1 flex-col ${collapsed ? 'md:ml-0' : 'md:ml-0'}`}>
        <header className="flex h-16 items-center justify-between border-b bg-white/80 px-4 backdrop-blur md:px-6">
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm" onClick={() => setMobileOpen((v) => !v)} className="md:hidden">
              <Menu className="h-4 w-4" />
            </Button>
            <Button variant="outline" size="sm" onClick={() => setCollapsed((v) => !v)} className="hidden md:inline-flex">
              <Menu className="h-4 w-4" />
            </Button>
            <div>
              <p className="text-xs text-muted-foreground">当前页面</p>
              <p className="text-sm font-semibold">{currentTitle}</p>
            </div>
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button className="flex items-center gap-2 rounded-md border px-2 py-1 hover:bg-muted">
                <Avatar>
                  <AvatarImage src={userInfo?.headerImg} />
                  <AvatarFallback>{(userInfo?.nickName || 'U').slice(0, 1)}</AvatarFallback>
                </Avatar>
                <span className="max-w-[120px] truncate text-sm">{userInfo?.nickName || '未登录用户'}</span>
                <ChevronDown className="h-4 w-4 text-muted-foreground" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => navigate('/layout/person')}>
                <UserCircle2 className="mr-2 h-4 w-4" />个人信息
              </DropdownMenuItem>
              {userInfo?.authorities?.length > 1 && <DropdownMenuSeparator className="my-1 h-px bg-border" />}
              {(userInfo?.authorities || [])
                .filter((a) => a.authorityId !== userInfo?.authorityId)
                .map((a) => (
                  <DropdownMenuItem key={a.authorityId} onClick={() => switchAuthority(a.authorityId)}>
                    切换到 {a.authorityName}
                  </DropdownMenuItem>
                ))}
              <DropdownMenuSeparator className="my-1 h-px bg-border" />
              <DropdownMenuItem onClick={logout}>
                <LogOut className="mr-2 h-4 w-4" />退出登录
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </header>

        <main className="min-h-0 flex-1 overflow-y-auto p-4 md:p-6">
          <Routes>
            <Route path="/" element={<Navigate to={defaultPath} replace />} />
            <Route path="person" element={<PersonPage />} />
            <Route path="dashboard" element={<DashboardPage />} />
            {routeEntries.map((entry) => (
              <Route key={entry.path} path={entry.path} element={resolveElement(entry)} />
            ))}
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </main>
      </div>

      {mobileOpen && <div className="fixed inset-0 z-20 bg-black/20 md:hidden" onClick={() => setMobileOpen(false)} />}
    </div>
  )
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/login" replace />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/init" element={<InitPage />} />
      <Route path="/layout/*" element={<LayoutPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}
