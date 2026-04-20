<script lang="ts" setup>
import { NButton, NIcon, NDataTable, NInput, NInputNumber, NSwitch, NTag, NPopconfirm, NSelect, type DataTableColumns } from 'naive-ui'
import { Pencil, RefreshOutline, TrashOutline, AddOutline } from '@vicons/ionicons5'
import * as menuApi from '@/api/content/menu'

const pagination = reactive({
  page: 1,
  pageSize: 100,
  total: 0,
})

const dataList = ref<any[]>([])
const loading = ref(false)
const currentMenuId = ref<number | undefined>(undefined)
const menuFormVisible = ref(false)
const formRef = ref<any>(null)
const formData = ref({
  name: '',
  title: '',
  path: '',
  icon: '',
  parent_id: 0,
  sort_order: 0,
  visible: true,
  target: '_self',
})

const parentMenuOptions = ref<{ label: string; value: number }[]>([{ label: '顶级菜单', value: 0 }])

const targetOptions = [
  { label: '当前窗口', value: '_self' },
  { label: '新窗口', value: '_blank' },
]

const buildMenuTree = (menus: any[]): any[] => {
  const menuMap = new Map<number, any>()
  const tree: any[] = []

  menus.forEach(menu => {
    menuMap.set(menu.id, { ...menu, children: [] })
  })

  menus.forEach(menu => {
    const node = menuMap.get(menu.id)!
    if (menu.parent_id === 0) {
      tree.push(node)
    } else {
      const parent = menuMap.get(menu.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        tree.push(node)
      }
    }
  })

  tree.sort((a, b) => a.sort_order - b.sort_order)
  const sortChildren = (nodes: any[]) => {
    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        node.children.sort((a: any, b: any) => a.sort_order - b.sort_order)
        sortChildren(node.children)
      }
    })
  }
  sortChildren(tree)

  return tree
}

const columns: DataTableColumns<any> = [
  {
    title: '菜单名称',
    key: 'name',
    width: 250,
  },
  {
    title: '菜单标题',
    key: 'title',
    width: 150,
    render: (row) => row.title || row.name,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '路径',
    key: 'path',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '图标',
    key: 'icon',
    width: 100,
  },
  {
    title: '排序',
    key: 'sort_order',
    width: 100,
  },
  {
    title: '可见',
    key: 'visible',
    width: 100,
    render: (row) => {
      return h(NTag, { type: row.visible ? 'success' : 'default' }, () => row.visible ? '可见' : '隐藏')
    },
  },
  {
    title: '目标',
    key: 'target',
    width: 100,
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
              onClick: () => openAddChildDialog(row),
            },
            {
              icon: () => h(NIcon, {}, () => h(AddOutline)),
              default: () => '添加子菜单',
            },
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              quaternary: true,
              onClick: () => openEditDialog(row),
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
              default: () => '确定删除该菜单吗？',
            },
          ),
        ],
      )
    },
  },
]

const loadMenuList = async () => {
  loading.value = true
  try {
    const res = await menuApi.getMenuList()
    const flatData = res.data || []
    dataList.value = buildMenuTree(flatData)
    updateParentMenuOptions()
  } catch (error) {
    console.error('加载菜单列表失败:', error)
    window.$message?.error('加载菜单列表失败')
  } finally {
    loading.value = false
  }
}

const flattenMenuTree = (menus: any[]): any[] => {
  const result: any[] = []
  const flatten = (nodes: any[]) => {
    nodes.forEach(node => {
      result.push(node)
      if (node.children && node.children.length > 0) {
        flatten(node.children)
      }
    })
  }
  flatten(menus)
  return result
}

const updateParentMenuOptions = () => {
  parentMenuOptions.value = [
    { label: '顶级菜单', value: 0 },
    ...flattenMenuTree(dataList.value).map((item: any) => ({
      label: item.name,
      value: item.id,
    })),
  ]
}

const resetFormData = () => {
  formData.value = {
    name: '',
    title: '',
    path: '',
    icon: '',
    parent_id: 0,
    sort_order: 0,
    visible: true,
    target: '_self',
  }
}

const openAddDialog = () => {
  currentMenuId.value = undefined
  resetFormData()
  menuFormVisible.value = true
}

const openAddChildDialog = (row: any) => {
  currentMenuId.value = undefined
  resetFormData()
  formData.value.parent_id = row.id
  menuFormVisible.value = true
}

const openEditDialog = (row: any) => {
  currentMenuId.value = row.id
  formData.value = {
    name: row.name || '',
    title: row.title || '',
    path: row.path || '',
    icon: row.icon || '',
    parent_id: row.parent_id || 0,
    sort_order: row.sort_order || 0,
    visible: row.visible !== undefined ? row.visible : true,
    target: row.target || '_self',
  }
  menuFormVisible.value = true
}

const handleSubmit = async () => {
  if (!formData.value.name) {
    window.$message?.warning('请输入菜单名称')
    return
  }

  try {
    if (currentMenuId.value) {
      await menuApi.updateMenu(currentMenuId.value, formData.value)
      window.$message?.success('更新成功')
    } else {
      await menuApi.createMenu(formData.value)
      window.$message?.success('创建成功')
    }
    
    menuFormVisible.value = false
    loadMenuList()
  } catch (error) {
    console.error('提交失败:', error)
    window.$message?.error('提交失败')
  }
}

const handleDelete = async (id: number) => {
  try {
    await menuApi.deleteMenu(id)
    window.$message?.success('删除成功')
    loadMenuList()
  } catch (error) {
    console.error('删除失败:', error)
    window.$message?.error('删除失败')
  }
}

onMounted(() => {
  loadMenuList()
})
</script>

<template>
  <div class="container-fluid p-6">
    <n-card title="前台菜单管理" class="menu-card">
      <div class="header-section">
        <div class="search-section">
        </div>
        <div class="action-section">
          <n-button type="primary" style="margin-right: 12px" @click="openAddDialog">
            <template #icon>
              <n-icon><add-outline /></n-icon>
            </template>
            添加菜单
          </n-button>
          <n-button @click="loadMenuList">
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
        :row-key="(row) => row.id"
        :children-key="'children'"
        default-expand-all
      />
    </n-card>

    <n-modal v-model:show="menuFormVisible" preset="card" :title="currentMenuId ? '编辑菜单' : '添加菜单'" style="width: 600px">
      <n-form ref="formRef" :model="formData" label-placement="left" label-width="80px">
        <n-form-item label="菜单名称" path="name" :rule="{ required: true, message: '请输入菜单名称', trigger: 'blur' }">
          <n-input v-model:value="formData.name" placeholder="请输入菜单名称" />
        </n-form-item>
        <n-form-item label="菜单标题" path="title">
          <n-input v-model:value="formData.title" placeholder="请输入菜单标题" />
        </n-form-item>
        <n-form-item label="路径" path="path">
          <n-input v-model:value="formData.path" placeholder="请输入菜单路径" />
        </n-form-item>
        <n-form-item label="图标" path="icon">
          <n-input v-model:value="formData.icon" placeholder="请输入图标类名" />
        </n-form-item>
        <n-form-item label="父菜单" path="parent_id">
          <n-select v-model:value="formData.parent_id" :options="parentMenuOptions" />
        </n-form-item>
        <n-form-item label="排序" path="sort_order">
          <n-input-number v-model:value="formData.sort_order" :min="0" style="width: 100%" />
        </n-form-item>
        <n-form-item label="是否可见" path="visible">
          <n-switch v-model:value="formData.visible" />
        </n-form-item>
        <n-form-item label="打开方式" path="target">
          <n-select v-model:value="formData.target" :options="targetOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="menuFormVisible = false">取消</n-button>
          <n-button type="primary" @click="handleSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.menu-card {
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
