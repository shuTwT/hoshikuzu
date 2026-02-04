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
import UserArea from './components/user-area/index.vue'

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
  layoutMode,
} = storeToRefs(appStore)

const { activeKey } = useLayoutMenu({
  mode: layoutMode,
  menus: menuOptions,
})




/** 判断路径是否参与菜单 */
function isRemaining(path: string) {
  return remainingPaths.includes(path)
}

function menuSelect(path: string) {
  if (permissionStore.wholeMenus.length === 0 || isRemaining(path)) return
}


watch(
  () => [route.path, usePermissionStore().wholeMenus],
  () => {
    if (route.path.includes('/redirect')) return
    menuSelect(route.path)
  },
)


onMounted(() => {
  activeKey.value = route.name?.toString() || ''
})
</script>
<template>
  <div class="app-wrapper">
    <lay-sidebar>
      <template #sidebar-bottom>
        <user-area/>
      </template>
    
    </lay-sidebar>
    
     <div class="main-container">
      <!-- 头部 -->
      <lay-header></lay-header>
      <!-- 主体内容 -->
       <lay-content></lay-content>
     </div>
  </div>
</template>
<style scoped>
.app-wrapper{
  position: relative;
  height: 100%;
  width: 100%;
}


.main-container{
  position: absolute;
  top: 0;
  bottom: 0;
  left: 280px;
  right: 0;
}
</style>
