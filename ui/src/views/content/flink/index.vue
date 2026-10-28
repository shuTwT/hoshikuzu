<script setup lang="ts">
import { ref } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  useMessage,
  NImage,
  NIcon,
  NSpace,
  NPopconfirm,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { Pencil, Trash } from '@vicons/ionicons5'
import { addDialog } from '@/components/dialog'
import type {
  FlinkFormItemProps,
  FlinkFormProps,
  FlinkGroupFormItemProps,
  FlinkGroupFormPorps,
} from './utils/types'
import flinkForm from './flinkForm.vue'
import { apiClient, useApi } from '@/api'
import FlinkGroupForm from './flinkGroupForm.vue'
import type { DropdownMixedOption } from 'naive-ui/es/dropdown/src/interface'

const message = useMessage()

const currentFlinkGroup = ref<number>(0)

// 友链表格数据
const loading = ref(false)
const flinkGroups = ref<any[]>([])
const dataList = ref<any[]>([])
const columns: DataTableColumns = [
  {
    title: '网站名称',
    key: 'name',
  },
  {
    title: 'Logo',
    key: 'logo',
    render: (row: any) => {
      return row.logo ? h(NImage, { src: row.logo }) : '暂无'
    },
  },
  {
    title: '网站地址',
    key: 'url',
  },
  {
    title: '描述',
    key: 'description',
  },
  {
    title: '操作',
    key: 'actions',
    render: (row: any) => {
      return h(
        NSpace,
        {},
        {
          default: () => [
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
                onPositiveClick: () => handleDelete(row),
              },
              {
                default: () => '确认删除吗',
                trigger: () =>
                  h(
                    NButton,
                    {
                      size: 'small',
                      type: 'error',
                      quaternary: true,
                    },
                    {
                      icon: () => h(NIcon, {}, () => h(Trash)),
                      default: () => '删除',
                    },
                  ),
              },
            ),
          ],
        },
      )
    },
  },
]
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const showDropdownRef = ref(false)
const dropdownRow = ref<any>()

const handleFlinkGroupContextMenu = (e: MouseEvent, row: any) => {
  e.preventDefault()
  dropdownRow.value = row
  showDropdownRef.value = false
  nextTick().then(() => {
    contextMenuX.value = e.clientX
    contextMenuY.value = e.clientY
    showDropdownRef.value = true
  })
}

const openGroupEditDialog = (title = '新增', row?: FlinkGroupFormItemProps) => {
  const formRef = ref()
  addDialog<FlinkGroupFormPorps>({
    title: `${title}分组`,
    props: {
      formInline: {
        name: row?.name ?? '',
        description: row?.description ?? '',
      },
    },
    contentRenderer: ({ options }) =>
      h(FlinkGroupForm, { formInline: options.props!.formInline, ref: formRef }),
    beforeSure: async (done) => {
      try {
        const curData = await formRef.value.getData()
        const chores = () => {
          message.success('操作成功')
          done()
          onSearchCategory()
        }
        if (title == '新增') {
          await useApi(apiClient.api.v1FlinkGroupCreateCreate, curData)
          chores()
        } else {
          await useApi(apiClient.api.v1FlinkGroupUpdateUpdate, row?.id!, curData)
          chores()
        }
      } catch { }
    },
  })
}

const dropdownOptions: DropdownMixedOption[] = [
  {
    label: '编辑',
    key: 'edit',
  },
  {
    label: '删除',
    key: 'delete',
  },
]

const openEditDialog = (title = '新增', row?: FlinkFormItemProps) => {
  const formRef = ref()
  addDialog<FlinkFormProps>({
    title: `${title}友链`,
    props: {
      formInline: {
        id: row?.id ?? undefined,
        name: row?.name ?? '',
        url: row?.url ?? '',
        avatar_url: row?.avatar_url ?? '',
        description: row?.description ?? '',
        cover_url: row?.cover_url ?? '',
        snapshot_url: row?.snapshot_url ?? '',
        email: row?.email ?? '',
        enable_friend_circle: row?.enable_friend_circle ?? true,
        friend_circle_rule_id: row?.friend_circle_rule_id ?? undefined,
        group_id: currentFlinkGroup.value ?? 0,
      },
    },
    contentRenderer: ({ options }) =>
      h(flinkForm, { formInline: options.props!.formInline, ref: formRef }),
    beforeSure: async (done) => {
      try {
        const curData = await formRef.value.getData()
        const chores = () => {
          message.success('操作成功')
          done()
          onSearchFlink()
        }
        if (title == '新增') {
          await useApi(apiClient.api.v1FlinkCreateCreate, curData)
          chores()
        } else {
          await useApi(apiClient.api.v1FlinkUpdateUpdate, row?.id!, curData)
          chores()
        }
      } catch (e) {
        console.error(e)
        done()
      }
    },
  })
}

const handleDeleteGroup = async (row: any) => {
  // TODO: 实现删除逻辑
  try {
    await useApi(apiClient.api.v1FlinkGroupDeleteDelete, row.id!)
    message.success('删除成功喵~')
    await onSearchCategory()
  } catch { }
}

const handleDelete = async (row: any) => {
  // TODO: 实现删除逻辑
  try {
    await useApi(apiClient.api.v1FlinkDeleteDelete, row.id!)
    message.success('删除成功喵~')
    await onSearchFlink()
  } catch { }
}

const onSearchCategory = async () => {
  const res = await useApi(apiClient.api.v1FlinkGroupListList)
  flinkGroups.value = res.data || []
}
const onSearchFlink = async () => {
  try {
    const res = await useApi(apiClient.api.v1FlinkPageList, {
      page: 1,
      size: 5,
      group_id: currentFlinkGroup.value,
    })
    dataList.value = res.data.records || []
  } catch { }
}

const handleContextMenuClickoutside = async () => {
  showDropdownRef.value = false
}

const handleContextMenuSelect = async (key: string | number) => {
  showDropdownRef.value = false
  switch (key) {
    case 'edit':
      openGroupEditDialog('编辑', dropdownRow.value)
      break
    case 'delete':
      handleDeleteGroup(dropdownRow.value)
      break
  }
}

const switchCurrentFlinkGroup = (id: number) => {
  currentFlinkGroup.value = id
  onSearchFlink()
}

onMounted(async () => {
  try {
    await onSearchCategory()
    if (flinkGroups.value.length > 0) {
      currentFlinkGroup.value = flinkGroups.value[0].id!
    }
    await onSearchFlink()

  } catch (e) {
    console.error(e)
  }
})
</script>

<template>
  <div class="container-fluid p-6">
    <div class="flink-card">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-6">友链管理</h2>
      <n-grid x-gap="6" cols="3">
        <n-gi span="1">
          <n-card title="友链分类">
            <template #header-extra>
              <n-button type="primary" style="margin-right: 5px" @click="openGroupEditDialog('新增')">
                <i class="fas fa-plus mr-2"></i>新增分类
              </n-button>
              <n-button @click="onSearchCategory()">刷新</n-button>
            </template>
            <n-empty v-if="flinkGroups.length === 0" description="暂无分类" class="py-10" />
            <ul v-else class="space-y-2">
              <!-- 示例相册 -->
              <li v-for="(item, index) in flinkGroups" :key="index" :class="{ active: item.id == currentFlinkGroup }"
                class="flex justify-between items-center p-2 bg-gray-50 dark:bg-gray-700 rounded-md cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600 flink-group-item"
                @click="switchCurrentFlinkGroup(item.id!)" @contextmenu="(e) => handleFlinkGroupContextMenu(e, item)">
                <span class="text-gray-800 dark:text-gray-200">{{ item.name }}</span>
                <span class="text-sm text-gray-500 dark:text-gray-400">{{ item.count }} 个</span>
              </li>
            </ul>
            <n-dropdown placement="bottom-start" trigger="manual" :x="contextMenuX" :y="contextMenuY"
              :options="dropdownOptions" :show="showDropdownRef" :on-clickoutside="handleContextMenuClickoutside"
              @select="handleContextMenuSelect" />
          </n-card>
        </n-gi>
        <n-gi span="2">
          <n-card>
            <template #header>
              <div class="header-section">
                <div class="title">友链管理</div>
                <div>
                  <n-button type="primary" @click="openEditDialog('新增')" style="margin-right: 5px">添加友链</n-button>
                  <n-button @click="onSearchFlink()">刷新</n-button>
                </div>
              </div>
            </template>
            <n-data-table :loading="loading" :columns="columns" :data="dataList" :pagination="{ pageSize: 10 }"
              :remote="true">
              <template #empty>
                <n-empty description="暂无友链" />
              </template>
            </n-data-table>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </div>
</template>

<style scoped>
.flink-card {
  max-width: 1600px;
  margin: 0 auto;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 16px;
  font-weight: 500;
}

.flink-group-item.active {
  background-color: var(--color-gray-100);
}
</style>
