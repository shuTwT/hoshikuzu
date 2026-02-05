<script setup lang="ts">
import { ref, reactive, h } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NTag,
  NPopconfirm,
  NSpace,
  NSelect,
  NIcon,
  useDialog,
  useMessage,
  type DataTableColumns,
  type PaginationProps,
} from 'naive-ui'
import { Checkmark, Close, RefreshOutline } from '@vicons/ionicons5'
import { getFlinkApplicationPage, approveFlinkApplication } from '@/api/content/flinkApplication'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const dataList = ref<any[]>([])

const statusFilter = ref<number | null>(null)
const applicationTypeFilter = ref<string | null>(null)

const statusOptions = [
  { label: '待审批', value: 0 },
  { label: '已通过', value: 1 },
  { label: '已拒绝', value: 2 },
]

const applicationTypeOptions = [
  { label: '新增', value: 'create' },
  { label: '修改', value: 'update' },
]

const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    pagination.page = page
    loadData()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    loadData()
  },
})

const columns: DataTableColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '网站名称',
    key: 'website_name',
    width: 150,
  },
  {
    title: '网站链接',
    key: 'website_url',
    width: 200,
    render: (row: any) => {
      return h('a', { href: row.website_url, target: '_blank', class: 'text-blue-500 hover:underline' }, row.website_url)
    },
  },
  {
    title: 'Logo',
    key: 'website_logo',
    width: 100,
    render: (row: any) => {
      return row.website_logo ? h('img', { src: row.website_logo, class: 'w-10 h-10 object-cover rounded' }) : '暂无'
    },
  },
  {
    title: '网站描述',
    key: 'website_description',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '联系邮箱',
    key: 'contact_email',
    width: 180,
  },
  {
    title: '申请类型',
    key: 'application_type',
    width: 100,
    render: (row: any) => {
      const typeMap: Record<string, { text: string; type: 'success' | 'info' }> = {
        create: { text: '新增', type: 'success' },
        update: { text: '修改', type: 'info' },
      }
      const type = typeMap[row.application_type] || { text: row.application_type, type: 'info' }
      return h(NTag, { type: type.type }, { default: () => type.text })
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const statusMap: Record<number, { text: string; type: 'warning' | 'success' | 'error' }> = {
        0: { text: '待审批', type: 'warning' },
        1: { text: '已通过', type: 'success' },
        2: { text: '已拒绝', type: 'error' },
      }
      const status = statusMap[row.status] || { text: '未知', type: 'warning' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render: (row: any) => {
      return h(
        NSpace,
        {},
        {
          default: () => [
            row.status === 0 &&
              h(
                NPopconfirm,
                {
                  onPositiveClick: () => handleApprove(row),
                },
                {
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: 'small',
                        type: 'success',
                        quaternary: true,
                      },
                      {
                        icon: () => h(NIcon, {}, () => h(Checkmark)),
                        default: () => '通过',
                      },
                    ),
                  default: () => '确定通过该申请吗？',
                },
              ),
            row.status === 0 &&
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  quaternary: true,
                  onClick: () => handleReject(row),
                },
                {
                  icon: () => h(NIcon, {}, () => h(Close)),
                  default: () => '拒绝',
                },
              ),
          ],
        },
      )
    },
  },
]

const loadData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize,
    }
    if (statusFilter.value !== null) {
      params.status = statusFilter.value
    }
    if (applicationTypeFilter.value !== null) {
      params.application_type = applicationTypeFilter.value
    }
    const res = await getFlinkApplicationPage(params)
    dataList.value = res.data.records
    pagination.itemCount = res.data.total
  } catch (error) {
    window.$message?.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleApprove = async (row: any) => {
  try {
    await approveFlinkApplication(row.id, { status: 1, reject_reason: '' })
    window.$message?.success('审批通过')
    await loadData()
  } catch (error) {
    window.$message?.error('操作失败')
  }
}

const handleReject = (row: any) => {
  dialog.warning({
    title: '拒绝申请',
    content: '请输入拒绝原因',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      const reason = prompt('请输入拒绝原因：')
      if (reason && reason.trim()) {
        confirmReject(row.id, reason)
      } else {
        window.$message?.warning('请输入拒绝原因')
      }
    },
  })
}

const confirmReject = async (id: number, reason: string) => {
  try {
    await approveFlinkApplication(id, {
      status: 2,
      reject_reason: reason,
    })
    window.$message?.success('已拒绝申请')
    await loadData()
  } catch (error) {
    window.$message?.error('操作失败')
  }
}

const handleFilterChange = () => {
  pagination.page = 1
  loadData()
}

const handleRefresh = () => {
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="友链申请管理" class="flink-application-card">
      <div class="header-section">
        <div class="search-section">
          <n-select
            v-model:value="statusFilter"
            placeholder="状态筛选"
            :options="statusOptions"
            clearable
            style="width: 120px"
            @update:value="handleFilterChange"
          />
          <n-select
            v-model:value="applicationTypeFilter"
            placeholder="申请类型"
            :options="applicationTypeOptions"
            clearable
            style="width: 120px"
            @update:value="handleFilterChange"
          />
        </div>
        <div class="action-section">
          <n-button @click="handleRefresh">
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
        :data="dataList"
        :pagination="pagination"
        :remote="true"
        :row-key="(row) => row.id"
      />
    </n-card>
  </div>
</template>

<style scoped>
.flink-application-card {
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
