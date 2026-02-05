<script setup lang="ts">
import { h, ref, onMounted, reactive } from 'vue'
import { NCard, NDataTable, NButton, NSpace, useMessage, NIcon, useDialog, NPopconfirm, NInput } from 'naive-ui'
import { RefreshOutline, TrashOutline, EyeOutline } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import {
  getVisitLogPage,
  deleteVisitLog,
  batchDeleteVisitLog,
  type VisitLogResponse,
} from '@/api/infra/visitLog'
import { addDialog } from '@/components/dialog'
import LogDetail from './logDetail.vue'
import dayjs from 'dayjs'

const message = useMessage()
const dialog = useDialog()

const searchParams = reactive({
  ip: '',
  path: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  itemCount: 0,
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

const dataList = ref<VisitLogResponse[]>([])
const loading = ref(false)
const selectedRowKeys = ref<number[]>([])

const columns: DataTableColumns<VisitLogResponse> = [
  {
    type: 'selection',
    width: 50,
  },
  {
    title: '编号',
    key: 'id',
    width: 80,
  },
  {
    title: '访问IP',
    key: 'ip',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '访问路径',
    key: 'path',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作系统',
    key: 'os',
    width: 120,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '浏览器',
    key: 'browser',
    width: 120,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '设备',
    key: 'device',
    width: 120,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '访问时间',
    key: 'created_at',
    width: 180,
    render(row){
      return dayjs(row.created_at).format('YYYY-MM-DD HH:mm:ss')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render(row) {
      return h(
        NSpace,
        {},
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                type: 'primary',
                onClick: () => openDetailDialog(row),
              },
              {
                icon: () => h(NIcon, null, { default: () => h(EyeOutline) }),
                default: () => '详情',
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
                default: () => '确认删除该访问日志吗？',
              },
            ),
          ],
        },
      )
    },
  },
]

const openDetailDialog = (row: VisitLogResponse) => {
  addDialog({
    title: '访问日志详情',
    contentRenderer: () => h(LogDetail, { data: row }),
  })
}

const onSearch = async () => {
  loading.value = true
  try {
    const res = await getVisitLogPage({
      page: pagination.page,
      page_size: pagination.pageSize,
      ip: searchParams.ip,
      path: searchParams.path,
    })
    if (res.code === 200) {
      dataList.value = res.data.records || []
      pagination.itemCount = res.data.total || 0
    }
  } catch (error) {
    message.error('获取访问日志列表失败：' + (error as Error).message)
  } finally {
    loading.value = false
  }
}

const handleDelete = async (row: VisitLogResponse) => {
  try {
    await deleteVisitLog(row.id)
    message.success('删除成功喵~')
    onSearch()
  } catch (error) {
    message.error('删除失败：' + (error as Error).message)
  }
}

const handleBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请先选择要删除的访问日志')
    return
  }
  dialog.warning({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 条访问日志吗？删除后无法恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await batchDeleteVisitLog(selectedRowKeys.value)
        message.success('批量删除成功喵~')
        selectedRowKeys.value = []
        onSearch()
      } catch (error) {
        message.error('批量删除失败：' + (error as Error).message)
      }
    },
  })
}

const resetSearch = () => {
  searchParams.ip = ''
  searchParams.path = ''
  pagination.page = 1
  onSearch()
}

onMounted(() => {
  onSearch()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="访问日志" class="visit-log-card">
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchParams.ip"
            placeholder="请输入访问IP"
            clearable
            style="width: 200px"
            @keyup.enter="onSearch"
          />
          <n-input
            v-model:value="searchParams.path"
            placeholder="请输入访问路径"
            clearable
            style="width: 200px"
            @keyup.enter="onSearch"
          />
          <n-button type="primary" @click="onSearch">查询</n-button>
          <n-button @click="resetSearch">重置</n-button>
        </div>
        <div class="action-section">
          <n-button
            v-if="selectedRowKeys.length > 0"
            type="error"
            @click="handleBatchDelete"
            style="margin-right: 12px"
          >
            <template #icon>
              <n-icon><trash-outline /></n-icon>
            </template>
            批量删除 ({{ selectedRowKeys.length }})
          </n-button>
          <n-button @click="onSearch">
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </div>
      <n-data-table
        v-model:checked-row-keys="selectedRowKeys"
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row) => row.id"
        :remote="true"
      />
    </n-card>
  </div>
</template>

<style scoped>
.visit-log-card {
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