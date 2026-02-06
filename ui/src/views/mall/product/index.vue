<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NPopconfirm } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline, CloudUploadOutline, CloudDownloadOutline } from '@vicons/ionicons5'
import * as productApi from '@/api/mall/product'
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
const checkedRowKeys = ref<number[]>([])
const editFormRef = ref()
const currentProductId = ref<number | undefined>(undefined)

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
    title: '商品名称',
    key: 'name',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: 'SKU',
    key: 'sku',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '价格',
    key: 'price',
    width: 120,
    render: (row) => `${(row.price / 100).toFixed(2)}元`,
  },
  {
    title: '原价',
    key: 'original_price',
    width: 120,
    render: (row) => row.original_price ? `${(row.original_price / 100).toFixed(2)}元` : '-',
  },
  {
    title: '库存',
    key: 'stock',
    width: 100,
  },
  {
    title: '销量',
    key: 'sales',
    width: 100,
  },
  {
    title: '品牌',
    key: 'brand',
    width: 120,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '状态',
    key: 'active',
    width: 100,
    render: (row) => {
      return h(NTag, { type: row.active ? 'success' : 'default' }, () => row.active ? '已上架' : '已下架')
    },
  },
  {
    title: '推荐',
    key: 'featured',
    width: 100,
    render: (row) => {
      return h(NTag, { type: row.featured ? 'warning' : 'default' }, () => row.featured ? '是' : '否')
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
              default: () => '确定删除该商品吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  productApi.getProductPage({
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
  currentProductId.value = row?.id
  addDialog({
    title: `${title}商品`,
    scroll: true,
    scrollbarHeight: '600px',
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        sku: row?.sku || '',
        price: row?.price || 0,
        original_price: row?.original_price || 0,
        cost_price: row?.cost_price || 0,
        stock: row?.stock || 0,
        min_stock: row?.min_stock || 0,
        category_id: row?.category_id || null,
        brand: row?.brand || '',
        unit: row?.unit || '',
        weight: row?.weight || 0,
        volume: row?.volume || 0,
        description: row?.description || '',
        short_description: row?.short_description || '',
        images: row?.images || [],
        attributes: row?.attributes || {},
        tags: row?.tags || [],
        active: row?.active !== undefined ? row.active : true,
        featured: row?.featured || false,
        digital: row?.digital || false,
        meta_title: row?.meta_title || '',
        meta_description: row?.meta_description || '',
        meta_keywords: row?.meta_keywords || '',
        sort_order: row?.sort_order || 0,
      }
    },
    contentRenderer: ({ options }) => h(EditForm, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        
        if (currentProductId.value) {
          await productApi.updateProduct(currentProductId.value, data)
          window.$message?.success('更新成功')
        } else {
          await productApi.createProduct(data)
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
    await productApi.deleteProduct(id)
    window.$message?.success('删除成功')
    onSearch()
  } catch (error) {
    console.error('删除失败:', error)
    window.$message?.error('删除失败')
  }
}

const handleBatchOnline = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要上架的商品')
    return
  }
  try {
    await productApi.batchUpdateProducts(checkedRowKeys.value, { active: true })
    window.$message?.success('上架成功')
    checkedRowKeys.value = []
    onSearch()
  } catch (error) {
    console.error('上架失败:', error)
    window.$message?.error('上架失败')
  }
}

const handleBatchOffline = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要下架的商品')
    return
  }
  try {
    await productApi.batchUpdateProducts(checkedRowKeys.value, { active: false })
    window.$message?.success('下架成功')
    checkedRowKeys.value = []
    onSearch()
  } catch (error) {
    console.error('下架失败:', error)
    window.$message?.error('下架失败')
  }
}

const handleBatchDelete = async () => {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning('请先选择要删除的商品')
    return
  }
  try {
    await productApi.batchDeleteProducts(checkedRowKeys.value)
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
    <n-card title="商品管理" class="product-card">
      <div class="header-section">
        <div class="search-section">
          <n-button type="success" @click="handleBatchOnline">
            <template #icon>
              <n-icon><cloud-upload-outline /></n-icon>
            </template>
            上架选中商品
          </n-button>
          <n-button type="warning" @click="handleBatchOffline">
            <template #icon>
              <n-icon><cloud-download-outline /></n-icon>
            </template>
            下架选中商品
          </n-button>
          <n-button type="error" @click="handleBatchDelete">
            <template #icon>
              <n-icon><trash-outline /></n-icon>
            </template>
            删除选中商品
          </n-button>
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openEditDialog('新增')">
            <template #icon>
              <n-icon><pencil /></n-icon>
            </template>
            添加商品
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
.product-card {
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
