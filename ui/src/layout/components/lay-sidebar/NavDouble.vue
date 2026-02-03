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
const childMenu = ref<RouteRecordRaw[]>([])
const subMenuData = ref([])

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
  childMenu.value = menu.children ?? []
  console.log(childMenu.value)
  router.push(menu.path);
  curActive.value = index;
}



</script>
<template>
  <div class="sidebar" style="--sidebar-width:280px;">
    <div class="sidebar-left">
      <div class="logo w-[64px] h-[64px]">
        <n-image src="/console/logo.png" :preview-disabled="true" width="64" />
      </div>
      <div class="h-calc(100%_-_64px)">
        <n-scrollbar>
          <ul>
            <li class="menu-item cursor-pointer select-none py-2 px-0.5" v-for="(menu, index) in menuData" :key="index"
              @click="handleChildMenu(menu, index)">
              <div class="w-full flex justify-center">
                <div class="text-xl mb-2">
                  <div></div>
                </div>
                <n-text class="text-xs">{{ menu.meta?.title || menu.name }}</n-text>
              </div>
            </li>
          </ul>
        </n-scrollbar>
      </div>

    </div>
    <div class="sidebar-right">
      <div class="logo-text"></div>
      <n-scrollbar>
        <n-menu :options="menuOptions">
        </n-menu>
      </n-scrollbar>
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
  background: #fff;
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
}
</style>