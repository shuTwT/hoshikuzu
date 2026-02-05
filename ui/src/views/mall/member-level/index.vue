<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NPopconfirm } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline } from '@vicons/ionicons5'
import * as memberLevelApi from '@/api/mall/memberLevel'
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
const currentMemberLevelId = ref<number | undefined>(undefined)

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '等级名称',
    key: 'name',
    width: 150,
  },
  {
    title: '等级级别',
    key: 'level',
    width: 120,
  },
  {
    title: '最低积分',
    key: 'min_points',
    width: 120,
  },
  {
    title: '折扣率(%)',
    key: 'discount_rate',
    width: 120,
    render: (row) => `${row.discount_rate}%`,
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
              default: () => '确定删除该会员等级吗？',
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  memberLevelApi.getMemberLevelPage({
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
  currentMemberLevelId.value = row?.id
  addDialog({
    title: `${title}会员等级`,
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        level: row?.level || 1,
        min_points: row?.min_points || 0,
        discount_rate: row?.discount_rate || 100,
        description: row?.description || '',
        privileges: row?.privileges || [],
        icon: row?.icon || '',
        active: row?.active !== undefined ? row.active : true,
        sort_order: row?.sort_order || 0,
      }
    },
    contentRenderer: ({ options }) => h(EditForm, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        
        if (currentMemberLevelId.value) {
          await memberLevelApi.updateMemberLevel(currentMemberLevelId.value, data)
          window.$message?.success('更新成功')
        } else {
          await memberLevelApi.createMemberLevel(data)
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
    await memberLevelApi.deleteMemberLevel(id)
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
    <n-card title="会员等级管理" class="member-level-card">
      <div class="header-section">
        <div class="search-section">
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openEditDialog('新增')">
            <template #icon>
              <n-icon><pencil /></n-icon>
            </template>
            添加会员等级
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
.member-level-card {
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
