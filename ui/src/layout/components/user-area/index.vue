<script setup lang="ts">
import { useUserStore } from '@/stores/modules/user'

const router = useRouter()
const userStore = useUserStore()
const dialog = useDialog()
const message = useMessage()

const gotoUserCenter = () => {
  router.push({ name: 'UserCenter' })
}

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

</script>
<template>
    <div>
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
    </div>
</template>