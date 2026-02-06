<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NPopconfirm, NInput, NSelect } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline, CheckmarkOutline, CloseOutline } from '@vicons/ionicons5'
import * as couponApi from '@/api/mall/coupon'
import { addDialog } from '@/components/dialog'
import EditForm from './editForm.vue'
import type { FormProps } from './utils/types'
import type { VNodeRef } from 'vue'

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
const searchKeyword = ref('')
const searchType = ref<number>()
const searchActive = ref<string >()


const couponTypeOptions = [
  { label: '全部', value: '' },
  { label: '满减', value: '1' },
  { label: '折扣', value: '2' },
  { label: '无门槛', value:'3' },
]

const activeOptions = [
  { label: '全部', value: '' },
  { label: '启用', value: '1' },
  { label: '禁用', value: '0' },
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
    title: '优惠券名称',
    key: 'name',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '优惠券代码',
    key: 'code',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row:any) => {
      const typeMap:Record<string,string> = {'1': '满减','2': '折扣','3': '无门槛' }
      return h(NTag, { type: 'info' }, () => typeMap[row.type as string] || '-')
    },
  },
  {
    title: '优惠券值',
    key: 'value',
    width: 120,
    render: (row) => `${row.value}分`,
  },
  {
    title: '最低消费',
    key: 'min_amount',
    width: 120,
    render: (row) => `${row.min_amount}分`,
  },
  {
    title: '发放总数',
    key: 'total_count',
    width: 100,
  },
  {
    title: '已使用',
    key: 'used_count',
    width: 100,
  },
  {
    title: '每用户限领',
    key: 'per_user_limit',
    width: 120,
  },
  {
    title: '开始时间',
    key: 'start_time',
    width: 180,
    render: (row) => new Date(row.start_time).toLocaleString(),
  },
  {
    title: '结束时间',
    key: 'end_time',
    width: 180,
    render: (row) => new Date(row.end_time).toLocaleString(),
  },
  {
    title: '状态',
    key: 'active',
    width: 100,
    render: (row) => {
      return h(NTag, { type: row.active ? 'success' : 'default' }, () => row.active ? '启用' : '禁用')
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
              onClick: () => {
                openEditDialog('编辑', row)
              },
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
              default: () => '确定删除该优惠券吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  couponApi.getCouponPage({
    page: pagination.page,
    page_size: pagination.pageSize,
    keyword: searchKeyword.value || undefined,
    type: searchType.value ?? undefined,
    active: searchActive.value ?? undefined,
  }).then(res => {
    dataList.value = res.data.records || []
    pagination.total = res.data.total || 0
  }).finally(() => {
    loading.value = false
  })
}

const openEditDialog = (title = '新增', row?: any) => {
  const editFormRef = ref()
  addDialog<FormProps>({
    title: `${title}优惠券`,
    scroll:true,
    scrollbarHeight:"600px",
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        code: row?.code || '',
        type: row?.type || 0,
        value: row?.value || 0,
        min_amount: row?.min_amount || 0,
        max_discount: row?.max_discount || 0,
        total_count: row?.total_count || 0,
        per_user_limit: row?.per_user_limit || 1,
        start_time: row?.start_time || 0,
        end_time: row?.end_time || 0,
        active: row?.active || true,
        image: row?.image || '',
        product_ids: row?.product_ids || [],
        category_ids: row?.category_ids || [],
      },
    },
    contentRenderer: ({ options }) => h(EditForm, { ref: editFormRef.value,formInline:options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        
        if (row?.id) {
          await couponApi.updateCoupon(row.id, data)
          window.$message?.success('更新成功')
        } else {
          await couponApi.createCoupon(data)
          window.$message?.success('创建成功')
        }
        
        onSearch()
        done()
      } catch (error) {
        console.error('提交失败:', error)
        window.$message?.error('提交失败')
      }
    },
  })
}

const handleDelete = async (id: number) => {
  try {
    await couponApi.deleteCoupon(id)
    window.$message?.success('删除成功')
    onSearch()
  } catch (error) {
    console.error('删除失败:', error)
    window.$message?.error('删除失败')
  }
}

const handleBatchEnable = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要启用的优惠券')
    return
  }
  try {
    await couponApi.batchUpdateCoupons(checkedRowKeys.value, { active: true })
    window.$message?.success('启用成功')
    checkedRowKeys.value = []
    onSearch()
  } catch (error) {
    console.error('启用失败:', error)
    window.$message?.error('启用失败')
  }
}

const handleBatchDisable = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要禁用的优惠券')
    return
  }
  try {
    await couponApi.batchUpdateCoupons(checkedRowKeys.value, { active: false })
    window.$message?.success('禁用成功')
    checkedRowKeys.value = []
    onSearch()
  } catch (error) {
    console.error('禁用失败:', error)
    window.$message?.error('禁用失败')
  }
}

const handleBatchDelete = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要删除的优惠券')
    return
  }
  try {
    await couponApi.batchDeleteCoupons(checkedRowKeys.value)
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
    <n-card title="优惠券管理" class="coupon-card">
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索优惠券名称或代码"
            clearable
            style="width: 200px; margin-right: 12px"
            @keyup.enter="onSearch"
          />
          <n-select
            v-model:value="searchType"
            :options="couponTypeOptions"
            placeholder="选择类型"
            style="width: 150px; margin-right: 12px"
            @update:value="onSearch"
          />
          <n-select
            v-model:value="searchActive"
            :options="activeOptions"
            placeholder="选择状态"
            style="width: 150px; margin-right: 12px"
            @update:value="onSearch"
          />
          <n-button type="success" @click="handleBatchEnable">
            <template #icon>
              <n-icon><checkmark-outline /></n-icon>
            </template>
            启用选中
          </n-button>
          <n-button type="warning" @click="handleBatchDisable">
            <template #icon>
              <n-icon><close-outline /></n-icon>
            </template>
            禁用选中
          </n-button>
          <n-button type="error" @click="handleBatchDelete">
            <template #icon>
              <n-icon><trash-outline /></n-icon>
            </template>
            删除选中
          </n-button>
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openEditDialog('新增')">
            <template #icon>
              <n-icon><pencil /></n-icon>
            </template>
            添加优惠券
          </n-button>
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
.coupon-card {
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
