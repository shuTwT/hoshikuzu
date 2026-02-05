<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NPopconfirm, NUpload } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline, CloudUploadOutline, EyeOutline } from '@vicons/ionicons5'
import * as themeApi from '@/api/infra/theme'
import { addDialog } from '@/components/dialog'
import EditForm from './editForm.vue'
import Detail from './detail.vue'

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
const editFormRef = ref<any>(null)
const currentThemeId = ref<number | undefined>(undefined)

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '主题类型',
    key: 'type',
    width: 100,
    render: (row) => {
      return row.type === 'internal' ? '内部主题' : '外部主题'
    },
  },
  {
    title: '主题名称',
    key: 'name',
    width: 150,
  },
  {
    title: '显示名称',
    key: 'display_name',
    width: 150,
  },
  {
    title: '版本',
    key: 'version',
    width: 100,
  },
  {
    title: '作者',
    key: 'author_name',
    width: 120,
  },
  {
    title: '描述',
    key: 'description',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '外部URL',
    key: 'external_url',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '状态',
    key: 'enabled',
    width: 100,
    render: (row) => {
      return h(NTag, { type: row.enabled ? 'success' : 'default' }, () => row.enabled ? '已启用' : '未启用')
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
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
              type: 'info',
              quaternary: true,
              onClick: () => handleDetail(row),
            },
            {
              icon: () => h(NIcon, {}, () => h(EyeOutline)),
              default: () => '详情',
            },
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              quaternary: true,
              onClick: () => {
                if (row.enabled) {
                  handleDisable(row.id)
                } else {
                  handleEnable(row.id)
                }
              },
            },
            {
              default: () => row.enabled ? '禁用' : '启用',
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
              default: () => '确定删除该主题吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  themeApi.getThemePage({
    page: pagination.page,
    page_size: pagination.pageSize,
  }).then(res => {
    dataList.value = res.data.records || []
    pagination.total = res.data.total || 0
  }).finally(() => {
    loading.value = false
  })
}

const handleUpload = () => {
  addDialog({
    title: '上传主题',
    props: {
      formInline: {
        type: 'internal',
        file_path: '',
        name: '',
        display_name: '',
        version: '',
        external_url: '',
        description: ''
      }
    },
    contentRenderer: ({ options }) => h(EditForm, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        await themeApi.createTheme(data)
        window.$message?.success('主题创建成功')
        onSearch()
        done()
      } catch (error) {
        console.error('提交失败:', error)
        window.$message?.error('提交失败')
      }
    },
  })
}

const handleDetail = (row: any) => {
  addDialog({
    title: '主题详情',
    contentRenderer: () => h(Detail, {
      theme: row
    }),
    beforeSure: (done) => {
      done()
    },
  })
}

const handleEnable = async (id: number) => {
  try {
    await themeApi.enableTheme(id)
    window.$message?.success('启用成功')
    onSearch()
  } catch (error) {
    console.error('启用失败:', error)
    window.$message?.error('启用失败')
  }
}

const handleDisable = async (id: number) => {
  try {
    await themeApi.disableTheme(id)
    window.$message?.success('禁用成功')
    onSearch()
  } catch (error) {
    console.error('禁用失败:', error)
    window.$message?.error('禁用失败')
  }
}

const handleDelete = async (id: number) => {
  try {
    await themeApi.deleteTheme(id)
    window.$message?.success('删除成功')
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
    <n-card title="主题管理" class="theme-card">
      <div class="header-section">
        <div class="search-section">
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="handleUpload">
            <template #icon>
              <n-icon><cloud-upload-outline /></n-icon>
            </template>
            上传主题
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
      />
    </n-card>
  </div>
</template>

<style scoped>
.theme-card {
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
