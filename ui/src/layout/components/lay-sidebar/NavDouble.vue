<script setup lang="ts">
import { useAppStore } from '@/stores/modules/app'
import { storeToRefs } from 'pinia'
import { useAppStoreHook } from "@/stores/modules/app";
import { usePermissionStore } from '@/stores/modules/permission'
import type { MenuOption } from 'naive-ui'
import { RouterLink, type RouteRecordRaw } from 'vue-router'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const menuData = computed(() => permissionStore.wholeMenus as RouteRecordRaw[])



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

const curActive = ref(0);
const childMenu = ref()
const subMenuData =ref([])

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
  childMenu.value = menu.children
  router.push(menu.path);
  curActive.value = index;
}



</script>
<template>
  <div class="sidebar" style="--sidebar-width:280px;">
    <div class="sidebar-left">
      <div class="logo"></div>
      <n-scrollbar>
        <ul>
          <li class="menu-item" v-for="(menu, index) in menuData" :key="index" @click="handleChildMenu(menu, index)">
            <div>
              <n-text>{{ menu.meta?.title || menu.name }}</n-text>
            </div>
          </li>
        </ul>
      </n-scrollbar>
    </div>
    <div class="sidebar-right">
      <div class="logo-text"></div>
      <n-scrollbar>
        <n-menu>
          <div v-for="(menu, index) in childMenu" :key="index">
            <n-text>{{ menu.meta?.title || menu.name }}</n-text>
          </div>
        </n-menu>
      </n-scrollbar>
    </div>
  </div>
</template>
<style scoped>
.sidebar {
  width: var(--sidebar-width);
  height: 100%;
  background-color: #f5f5f5;
}
</style>