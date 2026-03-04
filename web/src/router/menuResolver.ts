const pageModules = {
  ...import.meta.glob([
    '../pages/**/*.tsx',
    '!../pages/auth/LoginPage.tsx',
    '!../pages/system/InitPage.tsx',
    '!../pages/dashboard/index.tsx',
    '!../pages/person/index.tsx',
    '!../pages/fallback/NotFoundPage.tsx',
    '!../pages/fallback/NotMigratedPage.tsx'
  ]),
  ...import.meta.glob([
    '../pages/**/*.jsx',
    '!../pages/auth/LoginPage.jsx',
    '!../pages/system/InitPage.jsx',
    '!../pages/dashboard/index.jsx',
    '!../pages/person/index.jsx',
    '!../pages/fallback/NotFoundPage.jsx',
    '!../pages/fallback/NotMigratedPage.jsx'
  ])
}

export type MenuItem = {
  path?: string
  name?: string
  hidden?: boolean
  component?: string
  meta?: {
    title?: string
    icon?: string
    keepAlive?: boolean
  }
  children?: MenuItem[]
  [key: string]: any
}

export type RouteEntry = {
  path: string
  name?: string
  title: string
  icon?: string
  keepAlive: boolean
  componentPath?: string
  loader: any
  raw: MenuItem
}

const normalizeSegment = (v = '') => String(v).replace(/^\/+|\/+$/g, '')

const resolveComponentModule = (componentPath?: string) => {
  if (!componentPath) return null
  const normalized = String(componentPath).replace(/^\/+/, '')
  const basePath = normalized.replace(/^view\//, '').replace(/\.vue$/, '')
  const tsxPath = `../pages/${basePath}.tsx`
  const jsxPath = `../pages/${basePath}.jsx`
  return pageModules[tsxPath] || pageModules[jsxPath] || null
}

const walkMenus = (menus: MenuItem[], parentPath = ''): RouteEntry[] => {
  const result: RouteEntry[] = []
  ;(menus || []).forEach((menu) => {
    const seg = normalizeSegment(menu.path)
    const currentPath = [parentPath, seg].filter(Boolean).join('/')
    const fullPath = currentPath.replace(/^layout\/?/, '')

    if (menu.component) {
      result.push({
        path: fullPath || 'dashboard',
        name: menu.name,
        title: menu.meta?.title || menu.name,
        icon: menu.meta?.icon,
        keepAlive: Boolean(menu.meta?.keepAlive),
        componentPath: menu.component,
        loader: resolveComponentModule(menu.component),
        raw: menu
      })
    }

    if (menu.children?.length) {
      result.push(...walkMenus(menu.children, currentPath))
    }
  })
  return result
}

export const createRouteEntriesFromMenus = (menus: MenuItem[] = []): RouteEntry[] => {
  const entries = walkMenus(menus)
  const dedup = new Map()
  entries.forEach((entry) => {
    if (!entry.path || entry.path === '404') return
    dedup.set(entry.path, entry)
  })
  return [...dedup.values()]
}

export const findPathByRouteName = (menus: MenuItem[] = [], routeName?: string): string | null => {
  if (!routeName) return null
  const entries = createRouteEntriesFromMenus(menus)
  const matched = entries.find((v) => v.name === routeName)
  return matched ? `/layout/${matched.path}` : null
}
