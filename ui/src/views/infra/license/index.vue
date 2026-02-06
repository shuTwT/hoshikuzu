<script setup lang="ts">
import { h, ref, reactive, onMounted } from 'vue'
import { NCard, NDataTable, NButton, NSpace, NTag, NIcon, NPopconfirm, useMessage, NInput, NSelect } from 'naive-ui'
import { RefreshOutline, AddOutline, TrashOutline, CheckmarkCircleOutline, Search } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import {
  getLicensePage,
  deleteLicense,
  verifyLicense,
  createLicense,
  updateLicense,
  type License,
} from '@/api/infra/license'
import { addDialog } from '@/components/dialog'
import FormComponent from './form.vue'
import type { FormProps } from './utils/types'

const message = useMessage()

const searchForm = reactive({
  domain: '',
  customer_name: '',
  status: null as number | null,
})

const statusOptions = [
  { label: '全部', value: null },
  { label: '有效', value: 1 },
  { label: '过期', value: 2 },
  { label: '禁用', value: 3 },
]

const loading = ref(false)
const data = ref<License[]>([])
const editFormRef = ref<any>(null)
const currentLicenseId = ref<number | undefined>(undefined)

const pagination = reactive({
  page:1,
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

const columns: DataTableColumns<License> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '域名',
    key: 'domain',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '授权密钥',
    key: 'license_key',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '客户名称',
    key: 'customer_name',
  },
  {
    title: '过期日期',
    key: 'expire_date',
  },
  {
    title: '状态',
    key: 'status',
    render(row) {
      const statusMap: Record<number, { type: string; label: string }> = {
        1: { type: 'success', label: '有效' },
        2: { type: 'warning', label: '过期' },
        3: { type: 'error', label: '禁用' },
      }
      const status = statusMap[row.status] || { type: 'default', label: '未知' }
      return h(NTag, { type: status.type as any }, { default: () => status.label })
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render(row) {
      return h(
        'div',
        { style: { display: 'flex', gap: '8px' } },
        [
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              type: 'info',
              onClick: () => handleVerify(row),
            },
            {
              icon: () => h(NIcon, null, { default: () => h(CheckmarkCircleOutline) }),
              default: () => '验证',
            },
          ),
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              type: 'primary',
              onClick: () => openEditDialog('编辑', row),
            },
            {
              default: () => '编辑',
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
              default: () => '确定要删除该授权吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = async () => {
  loading.value = true
  try {
    const res = await getLicensePage({
      page: pagination.page,
      page_size: pagination.pageSize,
      domain: searchForm.domain || undefined,
      customer_name: searchForm.customer_name || undefined,
      status: searchForm.status ?? undefined,
    })
    data.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    message.error('获取授权列表失败：' + (error as Error).message)
  } finally {
    loading.value = false
  }
}

const onReset = () => {
  searchForm.domain = ''
  searchForm.customer_name = ''
  searchForm.status = null
  pagination.page = 1
  onSearch()
}

const openEditDialog = (title = '新增', row?: License) => {
  currentLicenseId.value = row?.id
  addDialog({
    title: `${title}授权`,
    props: {
      formInline: {
        id: row?.id || undefined,
        domain: row?.domain || '',
        license_key: row?.license_key,
        customer_name: row?.customer_name || '',
        expire_date: row?.expire_date ? new Date(row.expire_date).getTime() : undefined,
        status: row?.status || 1,
      }
    },
    contentRenderer: ({ options }) => h(FormComponent, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        
        if (currentLicenseId.value) {
          await updateLicense(currentLicenseId.value, data)
          message.success('更新成功')
        } else {
          await createLicense(data)
          message.success('创建成功')
        }
        
        onSearch()
        done()
      } catch (error) {
        console.error('提交失败:', error)
        message.error('提交失败')
      }
    },
  })
}

const handleDelete = async (row: License) => {
  try {
    await deleteLicense(row.id)
    message.success('授权删除成功')
    onSearch()
  } catch (error) {
    message.error('授权删除失败：' + (error as Error).message)
  }
}

const handleVerify = async (row: License) => {
  try {
    const res = await verifyLicense({ domain: row.domain })
    if (res.data.valid) {
      message.success(`授权验证成功：${res.data.message}，客户：${res.data.customer_name}，过期时间：${res.data.expire_date}`)
    } else {
      message.warning(`授权验证失败：${res.data.message}`)
    }
  } catch (error) {
    message.error('授权验证失败：' + (error as Error).message)
  }
}

onMounted(() => {
  onSearch()
})
</script>
<template>
  <div class="container-fluid p-6">
    <n-card title="授权管理" class="license-card">
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchForm.domain"
            placeholder="域名"
            clearable
            style="width: 200px"
          />
          <n-input
            v-model:value="searchForm.customer_name"
            placeholder="客户名称"
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
          <n-button type="primary" @click="onSearch">
            <template #icon>
              <n-icon><search /></n-icon>
            </template>
            搜索
          </n-button>
          <n-button @click="onReset">重置</n-button>
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openEditDialog('新增')">
            <template #icon>
              <n-icon><add-outline /></n-icon>
            </template>
            添加授权
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
        :columns="columns"
        :data="data"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row) => row.id"
      />
    </n-card>
  </div>
</template>

<style scoped>
.license-card {
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
