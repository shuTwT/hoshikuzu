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
const permissionStore = usePermissionStore()

const message = useMessage()
window.$message = message
const dialog = useDialog()
window.$dialog = dialog
const notification = useNotification()
window.$notification = notification








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
  left: 64px;
  right: 0;
  transition: left 0.3s ease-in-out;
}
.nav-right-show ~ .main-container{
  left: 280px;
}
</style>
