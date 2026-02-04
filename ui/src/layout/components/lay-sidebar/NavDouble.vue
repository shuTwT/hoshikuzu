<script setup lang="ts">
import { useAppStore } from '@/stores/modules/app'
import { storeToRefs } from 'pinia'
import { useAppStoreHook } from "@/stores/modules/app";
import { usePermissionStore } from '@/stores/modules/permission'
import { useThemeVars, type MenuOption } from 'naive-ui'
import { RouterLink, type RouteRecordRaw } from 'vue-router'
import { remainingPaths } from '@/router'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const menuData = computed(() => permissionStore.wholeMenus as RouteRecordRaw[])

const activeKey = ref<string | null>(null)
const activeSubKey = ref<string | null>(null)
const curActive = ref(0);
const childMenu = ref<RouteRecordRaw[]>([])
const navRightShow = computed(() => childMenu.value.length > 0)

const menuOptions = computed<MenuOption[]>(() =>
  childMenu.value.map((item: any) => {
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

const handleChildMenu = (menu: RouteRecordRaw, index: number) => {
  console.log(menu)
  // childMenu.value = menu.children ?? []
  console.log(childMenu.value)
  router.push(menu.path);
  curActive.value = index;
}

const themeVars = useThemeVars()


/** 判断路径是否参与菜单 */
function isRemaining(path: string) {
  return remainingPaths.includes(path)
}

function menuSelect() {
  if (permissionStore.wholeMenus.length === 0 || !route.name) return
  
  // 重置默认值
  activeKey.value = null
  activeSubKey.value = null
  childMenu.value = []
  
  // 遍历所有菜单，查找匹配的路由名称
  for(const menu of menuData.value){
    // 检查一级菜单是否匹配
    if (menu.name === route.name) {
      activeKey.value = menu.name?.toString()||''
      childMenu.value = menu.children ?? []
      return
    }
    
    // 检查二级菜单是否匹配
    if (menu.children && menu.children.length > 0) {
      const matchedChild = menu.children.find(child => child.name === route.name)
      if (matchedChild) {
        activeKey.value = menu.name?.toString()||''
        activeSubKey.value = matchedChild.name?.toString()||''
        childMenu.value = menu.children ?? []
        return
      }
    }
  }
}


watch(
  () => [route.name, usePermissionStore().wholeMenus],
  () => {
    if (route.path.includes('/redirect')) return
    menuSelect()
  }
)

onMounted(()=>{
  if(!route.path.includes('/redirect')){
    menuSelect()
  }
})

</script>
<template>
  <div class="sidebar" :class="{ 'nav-right-show': navRightShow }" style="--sidebar-width:280px;">
    <div class="sidebar-left">
      <div class="logo w-[64px] h-[64px]">
        <n-image src="/console/logo.png" :preview-disabled="true" width="64" />
      </div>
      <div class="h-[calc(100%-64px)]">
        <n-scrollbar>
          <ul>
            <li class="menu-item py-1 px-0.5" :class="{active:menu.name==activeKey}" v-for="(menu, index) in menuData" :key="index"
              @click="handleChildMenu(menu, index)">
              <div class="menu-item-block w-full py-2">
                <div class="text-xl mb-2">
                  <div></div>
                </div>
                <n-text class="text-xs" style="--n-text-color:#fff">{{ menu.meta?.title || menu.name }}</n-text>
              </div>
            </li>
          </ul>
        </n-scrollbar>

      </div>
      <div class="absolute bottom-0 left-0 right-0" style="height: 64px;">
        <slot name="sidebar-bottom"></slot>
      </div>
    </div>
    <div class="sidebar-right" v-show="navRightShow">
      <div class="logo-text-block">
        <div class="logo-text">
          Hoshikuzu
        </div>
      </div>

      <div class="h-[calc(100%-64px)]">
        <n-scrollbar>
          <n-menu :options="menuOptions" v-model:value="activeSubKey">
          </n-menu>
        </n-scrollbar>
      </div>


    </div>
  </div>
</template>
<style scoped>
.sidebar {
  position: relative;
  width: var(--sidebar-width);
  height: 100%;
  background-color: #f5f5f5;
}

.sidebar-left {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 1000;
  border-right: 1px solid #e5e5e5;
  background: rgb(40, 44, 52);
  overflow-x: hidden;
  width: 64px;
}

.sidebar-right {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 64px;
  right: 0;
  z-index: 1000;
  border-right: 1px solid #e5e5e5;
  background: #fff;
  transition: all 0.3s ease-in-out;
}

.sidebar {
  transition: width 0.3s ease-in-out;
}

.sidebar:not(.nav-right-show) {
  width: 64px;
}

.logo-text-block {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 64px;
  font-size: 20px;
  font-weight: bold;
  color: #333;
  border-bottom: 1px solid #e5e5e5;
}

.menu-item {
  cursor: pointer;
  user-select: none;
}

.menu-item-block {
  --menu-item-background: transparent;
  --menu-item-color: #fff;
  display: flex;
  justify-content: center;
  background: var(--menu-item-background);
  color: var(--menu-item-color);
  border-radius: 6px;
  transition: background 0.3s ease-in-out;
}
.menu-item.active .menu-item-block,
.menu-item:not(.active) .menu-item-block:hover{
  --menu-item-background: rgb(63, 184, 132);
}
</style>