<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import { ProLayout, useLayoutMenu, type ProLayoutProps } from 'pro-naive-ui'
import { RouterLink, useRoute, useRouter, type RouteRecordRaw } from 'vue-router'
import { useAppStore } from '@/stores/modules/app'
import { storeToRefs } from 'pinia'
import { usePermissionStore } from '@/stores/modules/permission'
import { useUserStore } from '@/stores/modules/user'
import { remainingPaths } from '@/router'
import { findRouteByPath, getParentPaths } from '@/router/utils'
import LaySidebar from './components/lay-sidebar/NavDouble.vue'
import LayHeader from './components/lay-header/index.vue'
import LayContent from './components/lay-content/index.vue'

defineOptions({
  name: "AdminLayout"
})

const route = useRoute()
const router = useRouter()
const routes: any = router.options.routes
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const userStore = useUserStore()
const message = useMessage()
window.$message = message
const dialog = useDialog()
window.$dialog = dialog
const notification = useNotification()
window.$notification = notification

const levelList = ref<any[]>([])

type LayoutThemeOverrides = NonNullable<ProLayoutProps['builtinThemeOverrides']>

const layoutThemeOverrides: LayoutThemeOverrides = {
  color: '#f5f7fa',
}

const menuData = computed(() => permissionStore.wholeMenus)

const menuOptions = computed<MenuOption[]>(() =>
  menuData.value.map((item: any) => {
    if (item.children && item.children.length > 1) {
      return {
        label: item.meta?.title || item.name,
        key: item.name,
        children: item.children.map((child: any) => ({
          label: () => {
            return h(
              RouterLink,
              {
                to: {
                  name: child.name,
                },
              },
              {
                default: () => child.meta?.title || child.name,
              },
            )
          },
          key: child.name,
        })),
      }
    }
    return {
      label: () => {
        return h(
          RouterLink,
          {
            to: {
              name: item.name,
            },
          },
          {
            default: () => item.meta?.title || item.name,
          },
        )
      },
      key: item.name,
    }
  }),
)

const {
  isMobile,
  showFooter,
  showTabbar,
  showLogo,
  showSidebar,
  showNav,
  collapsed,
  navFixed,
  footerFixed,
  layoutMode,
} = storeToRefs(appStore)
const navHeight = ref(50)
const sidebarWidth = ref(224)
const tabbarHeight = ref(38)
const footerHeight = ref(50)
const sidebarCollapsedWidth = ref(58)
const loading = computed(() => menuData.value.length === 0)

const { layout, verticalLayout, activeKey } = useLayoutMenu({
  mode: layoutMode,
  menus: menuOptions,
})

const hasHorizontalMenu = computed(() =>
  ['horizontal', 'mixed-two-column', 'mixed-sidebar'].includes(layoutMode.value),
)

// function updateMode(v: ProLayoutMode) {
//   mode.value = v
//   isMobile.value = v === 'mobile'
// }

function logout() {
  dialog.create({
    type: 'info',
    title: '提示',
    content: '确定要登出吗',
    positiveText: '确定',
    negativeText: '不确定',
    onPositiveClick: () => {
      // 发送登出请求
      userStore.logOut().then(() => {
        router.push('/login').then(() => {
          message.success('登出成功')
        })
      })
    },
    onNegativeClick: () => { },
  })
}

/** 判断路径是否参与菜单 */
function isRemaining(path: string) {
  return remainingPaths.includes(path)
}

function menuSelect(path: string) {
  if (permissionStore.wholeMenus.length === 0 || isRemaining(path)) return
}

const getBreadcrumb = () => {
  // 当前路由信息
  const currentRoute: any = findRouteByPath(router.currentRoute.value.path, routes)

  // 当前路由的父级路径组成的数组
  const parentRoutes = getParentPaths(router.currentRoute.value.name as string, routes, 'name')

  // 存放组成面包屑的数组
  const matched: (RouteRecordRaw | undefined)[] = []

  // 获取每个父级路径对应的路由信息
  parentRoutes.forEach((path) => {
    if (path !== '/') matched.push(findRouteByPath(path, routes))
  })

  matched.push(currentRoute)

  matched.forEach((item, index) => {
    if (currentRoute?.query || currentRoute?.params) return
    if (item?.children) {
      item.children.forEach((v) => {
        if (v?.meta?.title === item?.meta?.title) {
          matched.splice(index, 1)
        }
      })
    }
  })

  levelList.value = matched.filter((item) => item?.meta && item?.meta.title !== false)
}

const gotoUserCenter = () => {
  router.push({ name: 'UserCenter' })
}

watch(
  () => [route.path, usePermissionStore().wholeMenus],
  () => {
    if (route.path.includes('/redirect')) return
    menuSelect(route.path)
  },
)

watch(
  () => route.path,
  () => {
    getBreadcrumb()
  },
)

onMounted(() => {
  getBreadcrumb()
  activeKey.value = route.name?.toString() || ''
})
</script>
<template>
  <div class="app-wrapper h-dvh w-dvw">
    <lay-sidebar></lay-sidebar>
    
     <div class="main-container">
      <!-- 头部 -->
      <lay-header></lay-header>
      <!-- 主体内容 -->
       <lay-content></lay-content>
     </div>
    <pro-layout v-model:collapsed="collapsed" :mode="layoutMode" :show-nav="showNav" :show-logo="showLogo"
      :is-mobile="isMobile" :nav-fixed="navFixed" :nav-height="navHeight" :show-footer="showFooter"
      :show-tabbar="showTabbar" :show-sidebar="showSidebar" :footer-fixed="footerFixed" :footer-height="footerHeight"
      :sidebar-width="sidebarWidth" :tabbar-height="tabbarHeight" :sidebar-collapsed-width="sidebarCollapsedWidth"
      :builtin-theme-overrides="layoutThemeOverrides" logo-class="flex justify-center">
      <template #logo>
        <n-image src="/console/logo.png" :preview-disabled="true" width="64" />
      </template>
      <template #nav-left>
        <template v-if="!isMobile">
          <div class="flex items-center h-full pl-6">
            <template v-for="(item, index) in levelList" :key="index">

              <span>
                {{ item.meta?.title }}
              </span>
              <template v-if="index == 0">
                <n-divider vertical />
              </template>
            </template>

          </div>
        </template>
        <n-popover v-if="isMobile" trigger="click" style="padding: 0">
          <template #trigger>
            <n-button type="primary" size="small"> 菜单 </n-button>
          </template>
          <n-scrollbar class="flex-[1_0_0]">
            <n-spin v-if="loading" />
            <n-menu v-else v-bind="verticalLayout.verticalMenuProps" :collapsed="false" />
          </n-scrollbar>
        </n-popover>
      </template>
      <template #nav-center>
        <n-menu v-if="hasHorizontalMenu" v-bind="layout.horizontalMenuProps" />
      </template>
      <template #footer>
        <div class="w-full h-full text-center leading-[var(--pro-layout-footer-height)] text-gray-500">
          <span>Copyright ©2025</span>
        </div>
      </template>
      <template #sidebar>
        <n-scrollbar class="flex-[1_0_0]">
          <n-menu v-bind="layout.verticalMenuProps" :collapsed-width="sidebarCollapsedWidth" />
        </n-scrollbar>
        <n-divider />

        <!-- 用户信息卡片 -->
        <div class="p-2">
          <div class="bg-gray-100 dark:bg-gray-800 rounded-lg py-4 px-2 shadow-sm">
            <div class="flex items-center space-x-3">
              <n-avatar class="cursor-pointer" :style="{
                color: 'yellow',
                backgroundColor: 'red',
              }" @click="gotoUserCenter">
                M
              </n-avatar>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900 dark:text-white truncate cursor-pointer"
                  @click="gotoUserCenter">
                  {{ userStore.username }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
                  {{ userStore.username }}
                </p>
              </div>
              <button @click="logout()"
                class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1">
                  </path>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </template>
      <div class="app-main app-main--vertical">
        <slot v-if="$slots.default" />
        <router-view v-else />
      </div>
    </pro-layout>
  </div>
</template>
<style scoped>
.app-main {
  position: relative;
  width: 100%;
  overflow-x: hidden;
  /* background-color: #f5f7fa; */
}

.app-main--vertical {
  display: flex;
  flex-direction: column;
}

.main-container{
  position: absolute;
  top: 0;
  bottom: 0;
  left: 280px;
  right: 0;
}
</style>
