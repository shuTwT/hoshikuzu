<script setup lang="ts">
import { h, ref, reactive, onMounted } from 'vue'
import { NCard, NDataTable, NButton, NSpace, NSwitch, NTag, NIcon, NPopconfirm, useMessage, NInput, NSelect } from 'naive-ui'
import { RefreshOutline, Play, Stop, Refresh as RefreshIcon, TrashOutline, Search } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import {
  getPluginPage,
  deletePlugin,
  startPlugin,
  stopPlugin,
  restartPlugin,
  type Plugin,
} from '@/api/infra/plugin'
import { addDialog } from '@/components/dialog'
import FormComponent from './form.vue'

const message = useMessage()

const searchForm = reactive({
  name: '',
  key: '',
  status: '',
  enabled: null as boolean | null,
  auto_start: null as boolean | null,
})

const statusOptions = [
  { label: '全部', value: '' },
  { label: '已停止', value: 'stopped' },
  { label: '运行中', value: 'running' },
  { label: '错误', value: 'error' },
  { label: '加载中', value: 'loading' },
]

const enabledOptions = [
  { label: '是', value: true },
  { label: '否', value: false },
]

const autoStartOptions = [
  { label: '是', value: true },
  { label: '否', value: false },
]

const loading = ref(false)
const data = ref<Plugin[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  total: 0,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    pagination.page = page
    onSearch()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    onSearch()
  },
})

const columns: DataTableColumns<Plugin> = [
  {
    title: '插件标识',
    key: 'key',
  },
  {
    title: '插件名称',
    key: 'name',
  },
  {
    title: '版本',
    key: 'version',
  },
  {
    title: '状态',
    key: 'status',
    render(row) {
      const statusMap = {
        stopped: { type: 'default', label: '已停止' },
        running: { type: 'success', label: '运行中' },
        error: { type: 'error', label: '错误' },
        loading: { type: 'warning', label: '加载中' },
      }
      const status = statusMap[row.status] || { type: 'default', label: row.status }
      return h(NTag, { type: status.type as any }, { default: () => status.label })
    },
  },
  {
    title: '是否启用',
    key: 'enabled',
    render(row) {
      return h(NSwitch, {
        value: row.enabled,
        disabled: true,
      })
    },
  },
  {
    title: '是否自动启动',
    key: 'auto_start',
    render(row) {
      return h(NSwitch, {
        value: row.auto_start,
        disabled: true,
      })
    },
  },
  {
    title: '依赖插件',
    key: 'dependencies',
    render(row) {
      if (!row.dependencies || row.dependencies.length === 0) {
        return h('span', { style: { color: '#999' } }, '-')
      }
      return h('div', { style: { display: 'flex', gap: '4px', flexWrap: 'wrap' } }, 
        row.dependencies.map(dep => 
          h(NTag, { type: 'info', size: 'small' }, { default: () => dep })
        )
      )
    },
  },
  {
    title: '最后启动时间',
    key: 'last_started_at',
    render(row) {
      return row.last_started_at || '-'
    },
  },
  {
    title: '操作',
    key: 'actions',
    render(row) {
      return h(
        NSpace,
        {},
        {
          default: () => [
            row.status === 'running'
              ? h(
                  NButton,
                  {
                    size: 'small',
                    quaternary: true,
                    type: 'warning',
                    onClick: () => handleStop(row),
                  },
                  {
                    icon: () => h(NIcon, null, { default: () => h(Stop) }),
                    default: () => '停止',
                  },
                )
              : h(
                  NButton,
                  {
                    size: 'small',
                    quaternary: true,
                    type: 'success',
                    disabled: !row.enabled,
                    onClick: () => handleStart(row),
                  },
                  {
                    icon: () => h(NIcon, null, { default: () => h(Play) }),
                    default: () => '启动',
                  },
                ),
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                type: 'info',
                disabled: row.status === 'running',
                onClick: () => handleRestart(row),
              },
              {
                icon: () => h(NIcon, null, { default: () => h(RefreshIcon) }),
                default: () => '重启',
              },
            ),
            h(
              NPopconfirm,
              {
                onPositiveClick: () => handleDelete(row),
              },
              {
                trigger: () =>
                  h(
                    NButton,
                    {
                      size: 'small',
                      quaternary: true,
                      type: 'error',
                    },
                    {
                      icon: () => h(NIcon, null, { default: () => h(TrashOutline) }),
                      default: () => '删除',
                    },
                  ),
                default: () => '确定要删除该插件吗？删除后无法恢复！',
              },
            ),
          ],
        },
      )
    },
  },
]

const onSearch = async () => {
  loading.value = true
  try {
    const res = await getPluginPage({ 
      page: pagination.page, 
      page_size: pagination.pageSize,
      name: searchForm.name || undefined,
      key: searchForm.key || undefined,
      status: searchForm.status || undefined,
      enabled: searchForm.enabled ?? undefined,
      auto_start: searchForm.auto_start ?? undefined,
    })
    data.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    message.error('获取插件列表失败：' + (error as Error).message)
  } finally {
    loading.value = false
  }
}

const onReset = () => {
  searchForm.name = ''
  searchForm.key = ''
  searchForm.status = ''
  searchForm.enabled = null
  searchForm.auto_start = null
  pagination.page = 1
  onSearch()
}

const handleAddPlugin = () => {
  addDialog({
    title: '添加插件',
    contentRenderer: ({ options }) =>
      h(FormComponent, { ...options.props, onSuccess: () => onSearch() }),
  })
}

const handleStart = async (row: Plugin) => {
  try {
    await startPlugin(row.id)
    message.success('插件启动成功')
    onSearch()
  } catch (error) {
    message.error('插件启动失败：' + (error as Error).message)
  }
}

const handleStop = async (row: Plugin) => {
  try {
    await stopPlugin(row.id)
    message.success('插件停止成功')
    onSearch()
  } catch (error) {
    message.error('插件停止失败：' + (error as Error).message)
  }
}

const handleRestart = async (row: Plugin) => {
  try {
    await restartPlugin(row.id)
    message.success('插件重启成功')
    onSearch()
  } catch (error) {
    message.error('插件重启失败：' + (error as Error).message)
  }
}

const handleDelete = async (row: Plugin) => {
  try {
    await deletePlugin(row.id)
    message.success('插件删除成功')
    onSearch()
  } catch (error) {
    message.error('插件删除失败：' + (error as Error).message)
  }
}

onMounted(() => {
  onSearch()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="插件管理" class="plugin-card">
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchForm.name"
            placeholder="插件名称"
            clearable
            style="width: 200px"
          />
          <n-input
            v-model:value="searchForm.key"
            placeholder="插件标识"
            clearable
            style="width: 180px"
          />
          <n-select
            v-model:value="searchForm.status"
            placeholder="状态"
            clearable
            :options="statusOptions"
            style="width: 120px"
          />
          <n-select
            v-model:value="searchForm.enabled"
            placeholder="是否启用"
            clearable
            :options="enabledOptions"
            style="width: 120px"
          />
          <n-select
            v-model:value="searchForm.auto_start"
            placeholder="是否自动启动"
            clearable
            :options="autoStartOptions"
            style="width: 130px"
          />
          <n-button type="primary" @click="onSearch">
            <template #icon>
              <n-icon><search /></n-icon>
            </template>
            搜索
          </n-button>
          <n-button @click="onReset">重置</n-button>
        </div>
        <div class="action-section">
          <n-button type="primary" @click="handleAddPlugin" style="margin-right: 12px">添加插件</n-button>
          <n-button>
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </div>
      <n-data-table
        :loading="loading"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :remote="true"
      />
    </n-card>
  </div>
</template>

<style scoped>
.plugin-card {
  max-width: 1600px;
  margin: 0 auto;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.action-section {
  display: flex;
  align-items: center;
}
</style>
