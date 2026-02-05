<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NPopconfirm } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline } from '@vicons/ionicons5'
import * as categoryApi from '@/api/content/category'
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
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  },
})

const dataList = ref<any>([])
const loading = ref(false)
const editFormRef = ref<any>(null)
const currentCategoryId = ref<number | undefined>(undefined)

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '分类名称',
    key: 'name',
    width: 200,
  },
  {
    title: '分类别名',
    key: 'slug',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
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
    title: '排序',
    key: 'sort_order',
    width: 100,
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
    width: 150,
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
              default: () => '确定删除该分类吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  categoryApi.getCategoryPage({
    page: pagination.page,
    page_size: pagination.pageSize,
  }).then(res => {
    dataList.value = res.data.records || []
    pagination.total = res.data.total || 0
  }).finally(() => {
    loading.value = false
  })
}

const openEditDialog = (title = '新增', row?: any) => {
  currentCategoryId.value = row?.id
  addDialog({
    title: `${title}分类`,
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        slug: row?.slug || '',
        description: row?.description || '',
        sort_order: row?.sort_order || 0,
        active: row?.active !== undefined ? row.active : true,
      }
    },
    contentRenderer: ({ options }) => h(EditForm, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        
        if (currentCategoryId.value) {
          await categoryApi.updateCategory(currentCategoryId.value, data)
          window.$message?.success('更新成功')
        } else {
          await categoryApi.createCategory(data)
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
    await categoryApi.deleteCategory(id)
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
    <n-card title="分类管理" class="category-card">
      <div class="header-section">
        <div class="search-section">
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openEditDialog('新增')">
            <template #icon>
              <n-icon><pencil /></n-icon>
            </template>
            添加分类
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
.category-card {
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
