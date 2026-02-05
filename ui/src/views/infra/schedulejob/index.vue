<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { NButton, NIcon, NDataTable, NInput, NSelect, NPopconfirm, NTag, NSwitch, useMessage, useDialog } from 'naive-ui'
import { RefreshOutline, Pencil, TrashOutline, PlayCircleOutline } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import {
  getScheduleJobPage,
  createScheduleJob,
  updateScheduleJob,
  deleteScheduleJob,
  executeScheduleJob,
  type ScheduleJob,
  type CreateScheduleJobParams,
} from '@/api/infra/schedulejob'
import { addDialog } from '@/components/dialog'
import FormComponent from './form.vue'
import dayjs from 'dayjs'

const dialog = useDialog()

const loading = ref(false)
const data = ref<ScheduleJob[]>([])

const searchForm = reactive({
  name: '',
  type: '',
  enabled: undefined as string | undefined,
})

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

const typeOptions = [
  { label: '全部', value: '' },
  { label: 'Cron', value: 'cron' },
  { label: '间隔', value: 'interval' },
]

const enabledOptions = [
  { label: '全部', value: undefined },
  { label: '启用', value: 'true' },
  { label: '禁用', value: 'false' },
]

const columns: DataTableColumns<ScheduleJob> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '任务名称',
    key: 'name',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '任务类型',
    key: 'type',
    width: 100,
    render(row) {
      const typeMap = {
        cron: 'Cron',
        interval: '间隔',
      }
      return h(NTag, { type: 'info' }, { default: () => typeMap[row.type] || row.type })
    },
  },
  {
    title: '调度表达式',
    key: 'expression',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '内部任务',
    key: 'job_name',
    width: 120,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '状态',
    key: 'enabled',
    width: 80,
    render(row) {
      return h(NSwitch, {
        value: row.enabled,
        disabled: true,
      })
    },
  },
  {
    title: '下次执行时间',
    key: 'next_run_time',
    ellipsis: {
      tooltip: true,
    },
    render(row) {
      return row.next_run_time ? dayjs(row.next_run_time).format('YYYY-MM-DD HH:mm:ss') : '-'
    },
  },
  {
    title: '上次执行时间',
    key: 'last_run_time',
    ellipsis: {
      tooltip: true,
    },
    render(row) {
      return row.last_run_time ? dayjs(row.last_run_time).format('YYYY-MM-DD HH:mm:ss') : '-'
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    ellipsis: {
      tooltip: true,
    },
    render(row){
      return dayjs(row.created_at).format('YYYY-MM-DD HH:mm:ss')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 300,
    fixed: 'right',
    render(row) {
      return h(
        'div',
        { style: { display: 'flex', gap: '8px' } },
        [
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleExecute(row),
            },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: 'small',
                    type: 'success',
                    quaternary: true,
                  },
                  {
                    icon: () => h(NIcon, null, { default: () => h(PlayCircleOutline) }),
                    default: () => '立即执行',
                  },
                ),
              default: () => `确定要立即执行定时任务"${row.name}"吗？`,
            },
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              quaternary: true,
              onClick: () => openEditDialog('编辑', row),
            },
            {
              icon: () => h(NIcon, null, { default: () => h(Pencil) }),
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
                    type: 'error',
                    quaternary: true,
                  },
                  {
                    icon: () => h(NIcon, null, { default: () => h(TrashOutline) }),
                    default: () => '删除',
                  },
                ),
              default: () => `确定要删除定时任务"${row.name}"吗？`,
            },
          ),
        ],
      )
    },
  },
]

const onSearch = () => {
  loading.value = true
  getScheduleJobPage({
    page: pagination.page,
    page_size: pagination.pageSize,
  }).then(res => {
    if (res.code === 200) {
      data.value = res.data.records || []
      pagination.total = res.data.total || 0
    }
  }).finally(() => {
    loading.value = false
  })
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.type = ''
  searchForm.enabled = undefined
  pagination.page = 1
  onSearch()
}

const openEditDialog = (title = '新增', row?: ScheduleJob) => {
  const formRef = ref<any>()
  
  addDialog({
    title: `${title}定时任务`,
    props: {
      formInline: {
        id: row?.id || undefined,
        name: row?.name || '',
        type: row?.type || 'cron',
        expression: row?.expression || '',
        description: row?.description || '',
        enabled: row?.enabled !== undefined ? row.enabled : true,
        job_name: row?.job_name || 'friendCircle',
        max_retries: row?.max_retries || 3,
        failure_notification: row?.failure_notification || false,
      },
    },
    contentRenderer: ({ options }) =>
      h(FormComponent, { ref: formRef, formInline: options.props!.formInline }),
    beforeSure: async (done) => {
      try {
        const curData = await formRef.value?.getData()
        const chores = () => {
          window.$message?.success(`${title}成功`)
          done()
          onSearch()
        }
        if (title === '新增') {
          await createScheduleJob(curData as CreateScheduleJobParams)
          chores()
        } else {
          await updateScheduleJob(curData)
          chores()
        }
      } catch (error) {
        console.error(`${title}失败:`, error)
        window.$message?.error(`${title}失败`)
      }
    },
  })
}

const handleDelete = (row: ScheduleJob) => {
  deleteScheduleJob(row.id).then(() => {
    window.$message?.success('删除成功')
    onSearch()
  }).catch(error => {
    console.error('删除失败:', error)
    window.$message?.error('删除失败')
  })
}

const handleExecute = (row: ScheduleJob) => {
  executeScheduleJob(row.id).then(() => {
    window.$message?.success('定时任务执行成功')
    onSearch()
  }).catch(error => {
    console.error('执行失败:', error)
    window.$message?.error('执行失败')
  })
}

onMounted(() => {
  onSearch()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="定时任务管理" class="schedule-job-card">
      <template #header-extra>
        <n-space>
          <n-button type="primary" @click="openEditDialog('新增')">
            <i class="bi bi-plus"></i>
            添加定时任务
          </n-button>
          <n-button @click="onSearch()">
            <template #icon>
              <n-icon>
                <refresh-outline />
              </n-icon>
            </template>
            刷新
          </n-button>
        </n-space>
      </template>

      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchForm.name"
            placeholder="请输入任务名称"
            clearable
            style="width: 200px"
            @keyup.enter="onSearch"
          />
          <n-select
            v-model:value="searchForm.type"
            :options="typeOptions"
            placeholder="任务类型"
            clearable
            style="width: 150px"
            @update:value="onSearch"
          />
          <n-select
            v-model:value="searchForm.enabled"
            :options="enabledOptions"
            placeholder="状态"
            clearable
            style="width: 120px"
            @update:value="onSearch"
          />
          <n-button @click="resetSearch">
            重置
          </n-button>
        </div>
      </div>

      <n-data-table
        :columns="columns"
        :data="data"
        :loading="loading"
        :pagination="pagination"
        :remote="true"
        :row-key="(row) => row.id"
      />
    </n-card>
  </div>
</template>

<style scoped>
.schedule-job-card {
  max-width: 1600px;
  margin: 0 auto;
}

.header-section {
  display: flex;
  justify-content: flex-start;
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
</style>
