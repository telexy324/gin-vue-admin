import { create } from 'zustand'

type MenuStoreState = {
  menus: any[]
  routeEntries: any[]
  initialized: boolean
  setMenus: (menus: any[]) => void
  setRouteEntries: (routeEntries: any[]) => void
  setInitialized: (initialized: boolean) => void
  clearMenus: () => void
}

export const useMenuStore = create<MenuStoreState>()((set) => ({
  menus: [],
  routeEntries: [],
  initialized: false,
  setMenus: (menus) => set({ menus }),
  setRouteEntries: (routeEntries) => set({ routeEntries }),
  setInitialized: (initialized) => set({ initialized }),
  clearMenus: () => set({ menus: [], routeEntries: [], initialized: false })
}))
