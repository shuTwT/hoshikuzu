<script setup lang="ts">
import { findRouteByPath, getParentPaths } from '@/router/utils'
import type { RouteRecordRaw } from 'vue-router'

const route = useRoute()
const router = useRouter()
const routes: any = router.options.routes
const levelList = ref<any[]>([])

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
</script>
<template>
    <div class="header-nav">
        <div class="nav-left">
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
        </div>
        
    </div>
</template>
<style scoped>
.header-nav {
  height: 50px;
  background-color: #f5f7fa;
}
</style>