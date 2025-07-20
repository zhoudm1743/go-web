import type { MenuOption } from 'naive-ui'
import { router } from '@/router'
import { staticRoutes } from '@/router/routes.static'
import { fetchUserRoutes } from '@/service'
import { useAuthStore } from '@/store/auth'
import { $t, local } from '@/utils'
import { createMenus, createRoutes, generateCacheRoutes } from './helper'

interface RoutesStatus {
  isInitAuthRoute: boolean
  menus: MenuOption[]
  rowRoutes: AppRoute.RowRoute[]
  activeMenu: string | null
  cacheRoutes: string[]
}
export const useRouteStore = defineStore('route-store', {
  state: (): RoutesStatus => {
    return {
      isInitAuthRoute: false,
      activeMenu: null,
      menus: [],
      rowRoutes: [],
      cacheRoutes: [],
    }
  },
  actions: {
    resetRouteStore() {
      this.resetRoutes()
      this.$reset()
    },
    resetRoutes() {
      if (router.hasRoute('appRoot'))
        router.removeRoute('appRoot')
    },
    // set the currently highlighted menu key
    setActiveMenu(key: string) {
      this.activeMenu = key
    },

    async initRouteInfo() {
      if (import.meta.env.VITE_ROUTE_LOAD_MODE === 'dynamic') {
        try {
          const userInfo = local.get('userInfo')
          if (!userInfo || !userInfo.id) {
            console.log('用户信息不完整，使用静态路由')
            return staticRoutes
          }
          
          // Get user's route
          const { data } = await fetchUserRoutes({
            id: userInfo.id,
          })

          if (!data) {
            console.log('获取用户路由数据失败，使用静态路由')
            return staticRoutes
          }

          return data
        } catch (error) {
          console.error('获取路由时出错:', error)
          return staticRoutes
        }
      }
      else {
        this.rowRoutes = staticRoutes
        return staticRoutes
      }
    },
    async initAuthRoute() {
      try {
        this.isInitAuthRoute = false

        // Initialize route information
        const rowRoutes = await this.initRouteInfo()
        if (!rowRoutes) {
          window.$message.error($t(`app.getRouteError`))
          return
        }
        this.rowRoutes = rowRoutes

        // Generate actual route and insert
        const routes = createRoutes(rowRoutes)
        
        // 先检查appRoot路由是否已存在，如果存在则先移除
        if (router.hasRoute('appRoot')) {
          router.removeRoute('appRoot')
        }
        
        router.addRoute(routes)

        // Generate side menu
        this.menus = createMenus(rowRoutes)

        // Generate the route cache
        this.cacheRoutes = generateCacheRoutes(rowRoutes)

        this.isInitAuthRoute = true
      } catch (error) {
        console.error('初始化路由失败:', error)
        // 如果初始化失败，也标记为已初始化，避免无限循环
        this.isInitAuthRoute = true
        window.$message.error($t(`app.getRouteError`))
      }
    },
  },
})
