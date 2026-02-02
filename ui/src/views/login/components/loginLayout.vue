<script lang="ts" setup>
import { Language } from '@vicons/ionicons5'
const currentLanguage = ref('zh-CN')
const languageOptions = [
  { label: '简体中文', key: 'zh-CN' },
  { label: 'English', key: 'en-US' },
]

const handleLanguageChange = (key: string) => {
  currentLanguage.value = key
  // 这里可以添加语言切换逻辑
}
</script>
<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <!-- 语言切换 -->
    <div class="flex justify-end mb-4 absolute top-4 right-4">
      <n-dropdown trigger="click" :options="languageOptions" @select="handleLanguageChange">
        <n-button text>
          <template #icon>
            <n-icon>
              <Language />
            </n-icon>
          </template>
          {{ currentLanguage === 'zh-CN' ? '简体中文' : 'English' }}
        </n-button>
      </n-dropdown>
    </div>
    <div class="flex h-[580px]">

      <!-- 登录卡片左侧的图片 -->
      <div class="w-[300px] h-full">
        <div class="w-[300px] h-full translate-x-[10px]">
          <n-image src="/console/login-form-bg.webp" :preview-disabled="true" class="w-full h-full object-cover"></n-image>
        </div>
      </div>

      <!-- 登录卡片 -->
      <n-card class="shadow-xl rounded-2xl" style="width:400px">
        <div class="p-6 relative">
          <slot />
        </div>
      </n-card>


    </div>
    <!-- 底部信息 -->
    <div class="text-center mt-6 text-gray-500 text-sm absolute bottom-4">
      <p>© 2026 Hoshikuzu. All rights reserved.</p>
    </div>
  </div>
</template>
<style scoped>
/* 自定义样式 */
:deep(.n-card) {
  border-radius: 1rem;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

:deep(.n-button--primary:hover) {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
}
</style>
