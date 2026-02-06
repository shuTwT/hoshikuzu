<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NPopconfirm, NInput, NSelect } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline } from '@vicons/ionicons5'
import * as couponUsageApi from '@/api/mall/couponUsage'
import { addDialog } from '@/components/dialog'
import EditForm from './editForm.vue'

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

const dataList = ref<any>([])
const loading = ref(false)
const checkedRowKeys = ref<number[]>([])
const searchCouponCode = ref('')
const searchUserId = ref<number | null>(null)
const searchStatus = ref<number | null>(null)
const editFormRef = ref()

const statusOptions = [
  { label: '全部', value: '' },
  { label: '未使用', value: '0' },
  { label: '已使用', value: '1' },
  { label: '已过期', value: '2' },
]

const columns: DataTableColumns<any> = [
  {
    type: 'selection',
    width: 50,
  },
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '优惠券代码',
    key: 'coupon_code',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '用户ID',
    key: 'user_id',
    width: 100,
  },
  {
    title: '订单ID',
    key: 'order_id',
    width: 100,
    render: (row) => row.order_id || '-',
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const statusMap:Record<string,string> = { '0': '未使用', '1': '已使用', '2': '已过期' }
      const typeMap:Record<string,string> = { '0': 'default', '1': 'success', '2': 'error' }
      return h(NTag, { type: typeMap[row.status] as any }, () => statusMap[row.status] || '-')
    },
  },
  {
    title: '使用时间',
    key: 'used_at',
    width: 180,
    render: (row) => row.used_at ? new Date(row.used_at).toLocaleString() : '-',
  },
  {
    title: '优惠金额(分)',
    key: 'discount_amount',
    width: 120,
  },
  {
    title: '过期时间',
    key: 'expire_at',
    width: 180,
    render: (row) => new Date(row.expire_at).toLocaleString(),
  },
  {
    title: '备注',
    key: 'remark',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row) => {
      return h(
        'div',
        { style: { display: 'flex', gap: '8px' } },
        [
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              quaternary: true,
            },
            {
              icon: () => h(NIcon, {}, () => h(Pencil)),
              default: () => '编辑',
            },
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDelete(row.id),
            },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: 'small',
                    type: 'error',
                    quaternary: true,
                  },
                  {
                    icon: () => h(NIcon, {}, () => h(TrashOutline)),
                    default: () => '删除',
                  },
                ),
              default: () => '确定删除该使用记录吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  couponUsageApi.getCouponUsagePage({
    page: pagination.page,
    page_size: pagination.pageSize,
    coupon_code: searchCouponCode.value || undefined,
    user_id: searchUserId.value ?? undefined,
    status: searchStatus.value ?? undefined,
  }).then(res => {
    dataList.value = res.data.records || []
    pagination.total = res.data.total || 0
  }).finally(() => {
    loading.value = false
  })
}



const handleDelete = async (id: number) => {
  try {
    await couponUsageApi.deleteCouponUsage(id)
    window.$message?.success('删除成功')
    onSearch()
  } catch (error) {
    console.error('删除失败:', error)
    window.$message?.error('删除失败')
  }
}

const handleBatchDelete = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要删除的使用记录')
    return
  }
  try {
    await couponUsageApi.batchDeleteCouponUsages(checkedRowKeys.value)
    window.$message?.success('删除成功')
    checkedRowKeys.value = []
    onSearch()
  } catch (error) {
    console.error('删除失败:', error)
    window.$message?.error('删除失败')
  }
}

onMounted(() => {
  onSearch()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="优惠券使用记录" class="coupon-usage-card">
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchCouponCode"
            placeholder="搜索优惠券代码"
            clearable
            style="width: 200px; margin-right: 12px"
            @keyup.enter="onSearch"
          />
          <n-input-number
            v-model:value="searchUserId"
            placeholder="用户ID"
            clearable
            style="width: 150px; margin-right: 12px"
            @update:value="onSearch"
          />
          <n-select
            v-model:value="searchStatus"
            :options="statusOptions"
            placeholder="选择状态"
            style="width: 150px; margin-right: 12px"
            @update:value="onSearch"
          />
          <n-button type="error" @click="handleBatchDelete">
            <template #icon>
              <n-icon><trash-outline /></n-icon>
            </template>
            删除选中
          </n-button>
        </div>
        <div class="action-section">
          <n-button @click="onSearch()">
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </div>

      <n-data-table
        v-model:checked-row-keys="checkedRowKeys"
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
.coupon-usage-card {
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
