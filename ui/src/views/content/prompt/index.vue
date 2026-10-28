<script setup lang="ts">
import { Pencil, Trash } from '@vicons/ionicons5'
import type { SelectOption } from 'naive-ui'
import type { EntPrompt, EntPromptImage, EntPromptTag } from '@hoshikuzu/api-client'
import { addDialog } from '@/components/dialog'
import { apiClient, useApi } from '@/api'
import promptForm from './promptForm.vue'
import type { PromptFormProps } from './utils/types'

const message = useMessage()

const loading = ref(false)
const promptList = ref<EntPrompt[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const keyword = ref('')
const tagId = ref<number | null>(null)
const tagOptions = ref<SelectOption[]>([])

const onSearchTag = async () => {
  const res = await useApi(apiClient.api.v1PromptTagListList)
  tagOptions.value = ((res.data || []) as EntPromptTag[]).map((item) => ({
    label: item.name,
    value: item.id,
  }))
}

const onSearchPrompt = async () => {
  loading.value = true
  try {
    const res = await useApi(apiClient.api.v1PromptPageList, {
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
      tag_id: tagId.value ?? undefined,
    })
    promptList.value = res.data?.records || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const onRefresh = () => {
  onSearchTag()
  onSearchPrompt()
}

const onFilterChange = () => {
  page.value = 1
  onSearchPrompt()
}

const openPromptDialog = async (title = '新增', row?: EntPrompt) => {
  // 编辑时拉取详情，确保图片与标签为最新数据
  const detail: EntPrompt | undefined = row?.id
    ? (await useApi(apiClient.api.v1PromptQueryDetail, String(row.id))).data
    : undefined
  const formRef = ref()
  addDialog<PromptFormProps>({
    title: `${title}提示词`,
    width: '800px',
    props: {
      formInline: {
        title: detail?.title ?? '',
        content: detail?.content ?? '',
        negative_prompt: detail?.negative_prompt ?? '',
        description: detail?.description ?? '',
        sort: detail?.sort ?? 0,
        images: (detail?.edges?.images || []).map((img: EntPromptImage) => ({
          image_url: img.image_url ?? '',
          model_name: img.model_name ?? '',
          sort: img.sort ?? 0,
        })),
        tags: (detail?.edges?.tags || []).map((tag: EntPromptTag) => tag.name ?? ''),
      },
    },
    contentRenderer: ({ options }) =>
      h(promptForm, { ref: formRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await formRef.value.getData()
        if (title == '新增') {
          await useApi(apiClient.api.v1PromptCreateCreate, data)
        } else {
          await useApi(apiClient.api.v1PromptUpdateUpdate, String(row?.id), data)
        }
        message.success('操作成功')
        done()
        onRefresh()
      } catch (e) {
        console.error(e)
      }
    },
  })
}

const handleDelete = async (row: EntPrompt) => {
  try {
    await useApi(apiClient.api.v1PromptDeleteDelete, String(row.id))
    message.success('删除成功')
    if (promptList.value.length === 1 && page.value > 1) {
      page.value -= 1
    }
    await onRefresh()
  } catch (e) {
    console.error(e)
  }
}

const copyPrompt = async (row: EntPrompt) => {
  try {
    await navigator.clipboard.writeText(row.content ?? '')
    message.success('已复制提示词')
  } catch {
    message.error('复制失败')
  }
}

watch([page, pageSize], () => {
  onSearchPrompt()
})

onMounted(() => {
  onRefresh()
})
</script>
<template>
  <div class="p-6">
    <div class="prompt-card">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-6">提示词画廊</h2>
      <n-card title="提示词列表">
        <template #header-extra>
          <n-button type="primary" style="margin-right: 5px" @click="openPromptDialog('新增')">
            <i class="fas fa-plus mr-2"></i>新增提示词
          </n-button>
          <n-button @click="onRefresh()">刷新</n-button>
        </template>
        <div class="flex flex-wrap gap-3 mb-4">
          <n-input
            v-model:value="keyword"
            placeholder="搜索标题或提示词"
            clearable
            style="width: 240px"
            @keyup.enter="onFilterChange"
            @clear="onFilterChange"
          />
          <n-select
            v-model:value="tagId"
            :options="tagOptions"
            placeholder="按标签筛选"
            clearable
            style="width: 200px"
            @update:value="onFilterChange"
          />
          <n-button @click="onFilterChange">搜索</n-button>
        </div>
        <n-spin :show="loading">
          <n-empty v-if="promptList.length === 0" description="暂无提示词" class="py-16" />
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div
              v-for="item in promptList"
              :key="item.id"
              class="bg-gray-50 dark:bg-gray-700 rounded-lg shadow-sm overflow-hidden flex flex-col"
            >
              <n-image-group v-if="item.edges?.images?.length">
                <div class="grid grid-cols-2 gap-px">
                  <div
                    v-for="img in item.edges.images.slice(0, 4)"
                    :key="img.id"
                    class="relative h-28 overflow-hidden"
                  >
                    <n-image
                      :src="img.image_url"
                      :alt="img.model_name"
                      object-fit="cover"
                      width="100%"
                      class="h-full w-full"
                    />
                    <span v-if="img.model_name" class="model-badge">{{ img.model_name }}</span>
                  </div>
                </div>
              </n-image-group>
              <div v-else class="h-28 flex items-center justify-center text-sm text-gray-400">
                暂无图片
              </div>
              <div class="p-3 flex flex-col gap-2 flex-1">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-200 truncate">
                  {{ item.title }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-3 break-all">
                  {{ item.content }}
                </p>
                <div v-if="item.edges?.tags?.length" class="flex flex-wrap gap-1">
                  <n-tag v-for="tag in item.edges.tags" :key="tag.id" size="small" round>
                    {{ tag.name }}
                  </n-tag>
                </div>
                <div class="flex justify-end gap-1 mt-auto">
                  <n-button size="small" quaternary @click="copyPrompt(item)">复制</n-button>
                  <n-button
                    size="small"
                    type="primary"
                    quaternary
                    @click="openPromptDialog('编辑', item)"
                  >
                    <template #icon>
                      <n-icon><Pencil /></n-icon>
                    </template>
                    编辑
                  </n-button>
                  <n-popconfirm @positive-click="handleDelete(item)">
                    <template #trigger>
                      <n-button size="small" type="error" quaternary>
                        <template #icon>
                          <n-icon><Trash /></n-icon>
                        </template>
                        删除
                      </n-button>
                    </template>
                    确认删除吗
                  </n-popconfirm>
                </div>
              </div>
            </div>
          </div>
        </n-spin>
        <div class="mt-4 flex justify-end">
          <n-pagination
            v-model:page="page"
            v-model:page-size="pageSize"
            :item-count="total"
            show-size-picker
            :page-sizes="[12, 24, 48]"
          />
        </div>
      </n-card>
    </div>
  </div>
</template>
<style scoped>
.prompt-card {
  max-width: 1600px;
  margin: 0 auto;
}
.model-badge {
  position: absolute;
  left: 4px;
  bottom: 4px;
  padding: 0 6px;
  font-size: 12px;
  line-height: 18px;
  color: #fff;
  background: rgba(0, 0, 0, 0.55);
  border-radius: 4px;
}
</style>
