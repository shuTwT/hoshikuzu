<script lang="ts" setup>
import { NButton, NIcon, NDataTable, type DataTableColumns, NTag, NInputNumber } from 'naive-ui'
import { Pencil, RefreshOutline } from '@vicons/ionicons5'
import * as memberApi from '@/api/mall/member'
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
const memberLevelMap = ref<Record<number, any>>({})
const editFormRef = ref()

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '会员编号',
    key: 'member_no',
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
    title: '会员等级',
    key: 'member_level',
    width: 150,
    render: (row) => {
      const level = memberLevelMap.value[row.member_level]
      return h(NTag, { type: 'primary' }, () => level?.name || `等级${row.member_level}`)
    },
  },
  {
    title: '会员积分',
    key: 'points',
    width: 120,
  },
  {
    title: '累计消费',
    key: 'total_spent',
    width: 120,
    render: (row) => `${(row.total_spent / 100).toFixed(2)}元`,
  },
  {
    title: '订单数量',
    key: 'order_count',
    width: 100,
  },
  {
    title: '入会时间',
    key: 'join_time',
    width: 180,
    render: (row) => row.join_time ? new Date(row.join_time).toLocaleString() : '-',
  },
  {
    title: '到期时间',
    key: 'expire_time',
    width: 180,
    render: (row) => row.expire_time ? new Date(row.expire_time).toLocaleString() : '-',
  },
  {
    title: '状态',
    key: 'active',
    width: 100,
    render: (row) => {
      return h(NTag, { type: row.active ? 'success' : 'default' }, () => row.active ? '激活' : '未激活')
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row) => {
      return h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          quaternary: true,
          onClick: () => {
            openEditDialog(row)
          },
        },
        {
          icon: () => h(NIcon, {}, () => h(Pencil)),
          default: () => '编辑',
        },
      )
    },
  },
]

const getMemberLevels = async () => {
  try {
    const res = await memberLevelApi.getMemberLevelPage({ page: 1, page_size: 100 })
    const levels = res.data.records || []
    const map: Record<number, any> = {}
    levels.forEach((level: any) => {
      map[level.id] = level
    })
    memberLevelMap.value = map
  } catch (error) {
    console.error('获取会员等级失败:', error)
  }
}

const onSearch = () => {
  loading.value = true
  memberApi.getMemberPage({
    page: pagination.page,
    page_size: pagination.pageSize,
  }).then(res => {
    dataList.value = res.data.records || []
    pagination.total = res.data.total || 0
  }).finally(() => {
    loading.value = false
  })
}

const openEditDialog = (row: any) => {
  addDialog({
    title: '编辑会员',
    props: {
      formInline: {
        id: row.id,
        member_level: row.member_level,
        balance: 0,
        discount_rate: 100,
      },
    },
    contentRenderer: ({ options }) => h(EditForm, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()
        
        await memberApi.updateMember(row.id, data)
        window.$message?.success('更新成功')
        
        onSearch()
        done()
      } catch (error) {
        console.error('提交失败:', error)
        window.$message?.error('提交失败')
      }
    },
  })
}

onMounted(() => {
  getMemberLevels()
  onSearch()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="会员管理" class="member-card">
      <div class="header-section">
        <div class="search-section">
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
.member-card {
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
