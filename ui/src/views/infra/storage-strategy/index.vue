<script setup lang="ts">
import { h, ref, reactive, onMounted } from 'vue'
import { NCard, NDataTable, NButton, NSpace, NSwitch, NIcon, useMessage, NInput, NSelect } from 'naive-ui'
import { RefreshOutline,Pencil, Search } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import {
  getStorageStrategyList,
  createStorageStrategy,
  updateStorageStrategy,
  deleteStorageStrategy,
  setDefaultStorageStrategy,
  type StorageStrategy,
} from '@/api/infra/storage'
import { addDialog } from '@/components/dialog'
import FormComponent from './form.vue'
import type { FormItemProps } from './utils/types'

const message = useMessage()

const searchForm = reactive({
  name: '',
  type: '',
  master: null as boolean | null,
})

const typeOptions = [
  { label: '本地存储', value: 'local' },
  { label: 'S3存储', value: 's3' },
]

const masterOptions = [
  { label: '是', value: true },
  { label: '否', value: false },
]

const loading = ref(false)
const data = ref<StorageStrategy[]>([])

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

const columns: DataTableColumns<StorageStrategy> = [
  {
    title: '策略名称',
    key: 'name',
  },
  {
    title: '存储类型',
    key: 'type',
    render(row) {
      const typeMap = {
        local: '本地存储',
        s3: 'S3存储',
      }
      return typeMap[row.type] || row.type
    },
  },
  {
    title: '访问域名',
    key: 'domain',
  },
  {
    title: '是否默认',
    key: 'is_default',
    render(row) {
      return h(NSwitch, {
        value: row.master,
        disabled: row.master,
        onUpdateValue: () => handleSetDefault(row),
      })
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
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                type: 'primary',
                onClick: () => openEditDialog('编辑',row),
              },
              {
                icon: () => h(NIcon, null, { default: () => h(Pencil) }),
                default: () => '编辑',
              },
            ),
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                type: 'error',
                disabled: row.master,
                onClick: () => handleDelete(row),
              },
              { default: () => '删除' },
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
    const res = await getStorageStrategyList({
      page: pagination.page,
      page_size: pagination.pageSize,
      name: searchForm.name || undefined,
      type: searchForm.type || undefined,
      master: searchForm.master ?? undefined,
    })
    data.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    message.error('获取存储策略列表失败：' + (error as Error).message)
  } finally {
    loading.value = false
  }
}

const onReset = () => {
  searchForm.name = ''
  searchForm.type = ''
  searchForm.master = null
  pagination.page = 1
  onSearch()
}

const openEditDialog = (title='新增',row?:FormItemProps) => {
  const formRef = ref<any>()
  addDialog({
    title:`${title}存储策略`,
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        type: row?.type || 's3',
        endpoint:row?.endpoint||'',
        base_path: row?.base_path || '',
        node_id: row?.node_id || '',
        access_key: row?.access_key || '',
        secret_key: row?.secret_key || '',
        bucket: row?.bucket || '',
        region: row?.region||'',
        domain: row?.domain || '',
        master: row?.master || false,
      },
    },
    contentRenderer: ({ options }) =>
      h(FormComponent, { ref: formRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const curData = await formRef.value?.getData()
        console.log(curData)
        const chores = () => {
          message.success('创建成功喵~')
          done()
          onSearch()
        }
        if(title=='新增'){
          createStorageStrategy(curData).then(() => {
            chores()
          })
        }else{
          updateStorageStrategy(curData).then(() => {
            chores()
          })
        }

      } catch {}
    },
  })
}


const handleDelete = async (row: StorageStrategy) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除该存储策略吗？删除后无法恢复！',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteStorageStrategy(row.id)
        message.success('删除成功喵~')
        onSearch()
      } catch (error) {
        message.error('删除失败：' + (error as Error).message)
      }
    },
  })
}

const handleSetDefault = async (row: StorageStrategy) => {
  try {
    await setDefaultStorageStrategy(row.id)
    message.success('设置默认策略成功喵~')
    onSearch()
  } catch (error) {
    message.error('设置默认策略失败：' + (error as Error).message)
  }
}

onMounted(() => {
  onSearch()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="存储策略管理" class="storage-strategy-card">
      <!-- 头部操作栏 -->
       <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchForm.name"
            placeholder="策略名称"
            clearable
            style="width: 200px"
          />
          <n-select
            v-model:value="searchForm.type"
            placeholder="存储类型"
            clearable
            :options="typeOptions"
            style="width: 150px"
          />
          <n-select
            v-model:value="searchForm.master"
            placeholder="是否默认"
            clearable
            :options="masterOptions"
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
          <n-button type="primary" @click="openEditDialog('新增')" style="margin-right: 12px">添加策略</n-button>
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
.storage-strategy-card {
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
