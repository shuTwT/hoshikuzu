<script lang="ts" setup>
import { NButton, NIcon, NPopconfirm, NSpace, NInput, NSelect } from 'naive-ui'
import { RefreshOutline, Link, Search } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import { addDialog } from '@/components/dialog'
import uploadForm from './uploadForm.vue'
import * as fileApi from '@/api/infra/file'

const searchForm = reactive({
  name: '',
  type: '',
  storage_strategy_id: null as number | null,
})

const storageStrategyOptions = ref<any[]>([])

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

const dataList = ref<any[]>([])
const loading = ref(false)

const columns: DataTableColumns<any> = [
  {
    title: '编号',
    key: 'id',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title:"名称",
    key:"name",
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
    render:(row)=>{
      return row.size ? `${(row.size / 1024 / 1024).toFixed(2)} MB` : '0 KB'
    }
  },
  {
    title: '上传时间',
    key: 'created_at',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: () => {
      return h(
        NSpace,
        {},
        {
          default: ()=>[
            h(
              NButton,
              {
                size: 'small',
                type: 'primary',
                quaternary: true,
              },
              {
                icon: () => h(NIcon, {}, () => h(Link)),
                default: () => '复制链接',
              },
            ),
            h(
              NPopconfirm,
              {},
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
                      default: () => '删除',
                    },
                  ),
                default:()=>"确认删除吗？"
              },
            ),
          ],
        },
      )
    },
  },
]

const openUploadDialog = (title = '新增', row?: any) => {
  addDialog({
    title: '上传文件',
    contentRenderer: () => h(uploadForm),
  })
}

const onSearch = async () => {
  loading.value = true
  const res = await fileApi.getFilePage({
    page: pagination.page,
    page_size: pagination.pageSize,
    name: searchForm.name || undefined,
    type: searchForm.type || undefined,
    storage_strategy_id: searchForm.storage_strategy_id || undefined,
  })
  if (res.code === 200) {
    dataList.value = res.data.records || []
    pagination.itemCount = res.data.total || 0
  }
  loading.value = false
}

const onReset = () => {
  searchForm.name = ''
  searchForm.type = ''
  searchForm.storage_strategy_id = null
  pagination.page = 1
  onSearch()
}

onMounted(() => {
  onSearch()
})
</script>
<template>
  <div class="container-fluid p-6">
    <n-card title="文件管理" class="file-card">
      <!-- 头部操作栏 -->
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchForm.name"
            placeholder="文件名称"
            clearable
            style="width: 200px"
          />
          <n-input
            v-model:value="searchForm.type"
            placeholder="文件类型"
            clearable
            style="width: 150px"
          />
          <n-select
            v-model:value="searchForm.storage_strategy_id"
            placeholder="存储策略"
            clearable
            :options="storageStrategyOptions"
            style="width: 180px"
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
          <n-button type="primary" style="margin-right: 12px" @click="openUploadDialog('新增')">
            <i class="bi bi-plus"></i> 上传文件
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
.file-card {
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
