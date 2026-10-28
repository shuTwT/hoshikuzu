<script lang="ts" setup>
import {
  NButton,
  NIcon,
  NDataTable,
  type DataTableColumns,
  NTag,
  NPopconfirm,
  NInput,
} from 'naive-ui'
import { Pencil, RefreshOutline, SearchOutline, TrashOutline } from '@vicons/ionicons5'
import { addDialog } from '@/components/dialog'
import EditForm from './editForm.vue'
import { apiClient, useApi } from '@/api'

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
const currentId = ref<number | undefined>(undefined)
const searchName = ref('')
const searchCode = ref('')

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 70,
  },
  {
    title: '模板名称',
    key: 'name',
    width: 140,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '模板编码',
    key: 'code',
    width: 160,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '短信签名',
    key: 'sign_name',
    width: 110,
    render: (row) => row.sign_name || '-',
  },
  {
    title: '服务商模板ID',
    key: 'provider_template_id',
    width: 130,
    ellipsis: {
      tooltip: true,
    },
    render: (row) => row.provider_template_id || '-',
  },
  {
    title: '模板内容',
    key: 'content',
    minWidth: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '状态',
    key: 'enabled',
    width: 90,
    render: (row) => {
      return h(NTag, { type: row.enabled ? 'success' : 'default' }, () =>
        row.enabled ? '启用' : '禁用',
      )
    },
  },
  {
    title: '备注',
    key: 'remark',
    width: 140,
    ellipsis: {
      tooltip: true,
    },
    render: (row) => row.remark || '-',
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: (row) => (row.created_at ? new Date(row.created_at).toLocaleString() : '-'),
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render: (row) => {
      return h('div', { style: { display: 'flex', gap: '8px' } }, [
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
            default: () => '确定删除该短信模板吗？',
          },
        ),
      ])
    },
  },
]

const onSearch = () => {
  loading.value = true
  useApi(apiClient.api.v1SmsTemplatePageList, {
    page: pagination.page,
    page_size: pagination.pageSize,
    name: searchName.value || undefined,
    code: searchCode.value || undefined,
  } as any)
    .then((res) => {
      dataList.value = res.data.records || []
      pagination.total = res.data.total || 0
    })
    .finally(() => {
      loading.value = false
    })
}

const openEditDialog = (title = '新增', row?: any) => {
  currentId.value = row?.id
  addDialog({
    title: `${title}短信模板`,
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        code: row?.code || '',
        content: row?.content || '',
        sign_name: row?.sign_name || '',
        provider_template_id: row?.provider_template_id || '',
        enabled: row?.enabled !== undefined ? row.enabled : true,
        remark: row?.remark || '',
      },
    },
    contentRenderer: ({ options }) =>
      h(EditForm, { ref: editFormRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const data = await editFormRef.value?.getData()

        if (currentId.value) {
          await useApi(apiClient.api.v1SmsTemplateUpdateUpdate, currentId.value, data)
          window.$message?.success('更新成功')
        } else {
          await useApi(apiClient.api.v1SmsTemplateCreateCreate, data)
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
    await useApi(apiClient.api.v1SmsTemplateDeleteDelete, id)
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
    <n-card title="短信模板" class="sms-template-card">
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchName"
            clearable
            placeholder="模板名称"
            style="width: 180px"
            @keyup.enter="onSearch()"
          />
          <n-input
            v-model:value="searchCode"
            clearable
            placeholder="模板编码"
            style="width: 180px"
            @keyup.enter="onSearch()"
          />
          <n-button type="primary" quaternary @click="onSearch()">
            <template #icon>
              <n-icon><search-outline /></n-icon>
            </template>
            搜索
          </n-button>
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openEditDialog('新增')">
            <template #icon>
              <n-icon><pencil /></n-icon>
            </template>
            添加模板
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
        :scroll-x="1400"
      />
    </n-card>
  </div>
</template>

<style scoped>
.sms-template-card {
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
