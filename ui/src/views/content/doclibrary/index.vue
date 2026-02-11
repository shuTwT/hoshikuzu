<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton,
  NSpace,
  NModal,
  NInput,
  NSelect,
  NForm,
  NFormItem,
  NPagination,
  NSpin,
  NPopconfirm,
  NTag,
  NCard,
  NCheckbox,
  NCheckboxGroup,
  NEmpty,
  useMessage,
  type DataTableColumns,
} from 'naive-ui'
import * as doclibraryApi from '@/api/content/doclibrary'
import dayjs from 'dayjs'

const router = useRouter()
const message = useMessage()

interface DocLibrary {
  id: number
  name: string
  alias: string
  description: string
  source: string
  url: string
  created_at: string
  updated_at: string
}

interface DocLibraryFormData {
  name: string
  alias: string
  description: string
  source: string
  url: string
}

const loading = ref(false)
const data = ref<DocLibrary[]>([])
const selectedIds = ref<number[]>([])
const pagination = ref({
  page: 1,
  pageSize: 12,
  total: 0,
  showSizePicker: true,
  pageSizes: [12, 24, 36, 48],
})

const showModal = ref(false)
const isEdit = ref(false)
const isSettings = ref(false)
const currentId = ref<number | null>(null)
const formData = ref<DocLibraryFormData>({
  name: '',
  alias: '',
  description: '',
  source: '',
  url: '',
})

const sourceOptions = [
  { label: 'Git', value: 'git' },
  { label: 'OpenAPI', value: 'openapi' },
  { label: 'LLMs Txt', value: 'llms_txt' },
  { label: 'Website', value: 'website' },
]

const sourceConfig: Record<string, { type: any; label: string; icon: string; color: string }> = {
  git: { type: 'info', label: 'Git', icon: 'M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z', color: '#24292e' },
  openapi: { type: 'success', label: 'OpenAPI', icon: 'M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5', color: '#6BA539' },
  llms_txt: { type: 'warning', label: 'LLMs Txt', icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 2l5 5h-5V4zM6 20V4h6v6h6v10H6z', color: '#FF6B35' },
  website: { type: 'default', label: 'Website', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#4285F4' },
}

const isAllSelected = computed(() => data.value.length > 0 && selectedIds.value.length === data.value.length)

const isIndeterminate = computed(() => selectedIds.value.length > 0 && selectedIds.value.length < data.value.length)

const handleSelectAll = (checked: boolean) => {
  if (checked) {
    selectedIds.value = data.value.map(item => item.id)
  } else {
    selectedIds.value = []
  }
}

const handleCardSelect = (id: number, checked: boolean) => {
  if (checked) {
    if (!selectedIds.value.includes(id)) {
      selectedIds.value.push(id)
    }
  } else {
    const index = selectedIds.value.indexOf(id)
    if (index > -1) {
      selectedIds.value.splice(index, 1)
    }
  }
}

const handleBatchDelete = async () => {
  if (selectedIds.value.length === 0) {
    message.warning('请先选择要删除的文档库')
    return
  }

  try {
    const deletePromises = selectedIds.value.map(async id => {
      try {
        await doclibraryApi.deleteDocLibrary(id)
        return { success: true, id }
      } catch (error) {
        console.error(`删除文档库 ${id} 失败:`, error)
        return { success: false, id }
      }
    })
    
    const results = await Promise.all(deletePromises)
    const successCount = results.filter(r => r.success).length
    const failedCount = results.filter(r => !r.success).length
    
    if (successCount > 0) {
      message.success(`成功删除 ${successCount} 个文档库${failedCount > 0 ? `，失败 ${failedCount} 个` : ''}`)
    } else {
      message.error('删除失败，请稍后重试')
    }
    
    selectedIds.value = []
    fetchDocLibraries()
  } catch (error: any) {
    console.error('批量删除失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || '批量删除失败，请稍后重试'
    message.error(errorMessage)
  }
}

const fetchDocLibraries = async () => {
  loading.value = true
  try {
    const res = await doclibraryApi.getDocLibraryPage({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
    })
    data.value = res.data.records
    pagination.value.total = res.data.total
  } catch (error) {
    console.error('获取文档库列表失败:', error)
    message.error('获取文档库列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  formData.value = {
    name: '',
    alias: '',
    description: '',
    source: '',
    url: '',
  }
  showModal.value = true
}

const handleEdit = (row: DocLibrary) => {
  router.push({
    name: 'DocLibraryDetailManagement',
    query: { docLibraryId: row.id }
  })
}

const handleSettings = (row: DocLibrary) => {
  isEdit.value = true
  isSettings.value = true
  currentId.value = row.id
  formData.value = {
    name: row.name,
    alias: row.alias,
    description: row.description,
    source: row.source,
    url: row.url,
  }
  showModal.value = true
}

const handleDelete = async (id: number) => {
  try {
    await doclibraryApi.deleteDocLibrary(id)
    message.success('删除成功')
    selectedIds.value = selectedIds.value.filter(selectedId => selectedId !== id)
    fetchDocLibraries()
  } catch (error: any) {
    console.error('删除失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || '删除失败，请稍后重试'
    message.error(errorMessage)
  }
}

const handleSubmit = async () => {
  if (!formData.value.name.trim()) {
    message.error('请输入文档库名称')
    return
  }
  if (!formData.value.alias.trim()) {
    message.error('请输入文档库别名')
    return
  }
  if (!formData.value.source) {
    message.error('请选择文档库来源')
    return
  }
  if (!formData.value.url.trim()) {
    message.error('请输入文档库URL')
    return
  }

  try {
    if (isEdit.value && currentId.value) {
      await doclibraryApi.updateDocLibrary(currentId.value, formData.value)
      message.success('更新成功')
    } else {
      await doclibraryApi.createDocLibrary(formData.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchDocLibraries()
  } catch (error: any) {
    console.error('操作失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || (isEdit.value ? '更新失败，请稍后重试' : '创建失败，请稍后重试')
    message.error(errorMessage)
  }
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  fetchDocLibraries()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  fetchDocLibraries()
}

const handleRefresh = () => {
  fetchDocLibraries()
  message.success('刷新成功')
}

const handleModalClose = () => {
  showModal.value = false
  isEdit.value = false
  isSettings.value = false
  currentId.value = null
  formData.value = {
    name: '',
    alias: '',
    description: '',
    source: '',
    url: '',
  }
}

const getDocSize = (doc: DocLibrary): string => {
  const sizes = ['1.2 MB', '2.5 MB', '856 KB', '3.4 MB', '1.8 MB', '5.2 MB']
  return sizes[doc.id % sizes.length]
}

onMounted(() => {
  fetchDocLibraries()
})
</script>

<template>
  <div class="container-fluid p-6">
    <div class="max-w-[1600px] mx-auto">
    <NCard :bordered="false" class="doclibrary-card">
      <div class="doclibrary-header">
        <h1 class="page-title">文档库管理</h1>
        <NSpace>
          <NButton 
            v-if="selectedIds.length > 0"
            type="error" 
            @click="handleBatchDelete"
          >
            <template #icon>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
            </template>
            批量删除 ({{ selectedIds.length }})
          </NButton>
          <NButton type="primary" @click="handleCreate">
            <template #icon>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
            </template>
            新建文档库
          </NButton>
          <NButton @click="handleRefresh">
            <template #icon>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 4v6h-6"></path>
                <path d="M1 20v-6h6"></path>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
              </svg>
            </template>
            刷新
          </NButton>
        </NSpace>
      </div>

      <NSpin :show="loading">
        <div v-if="data.length === 0 && !loading" class="empty-container">
          <NEmpty description="暂无文档库数据" />
        </div>
        
        <div v-else class="card-grid">
          <div
            v-for="doc in data"
            :key="doc.id"
            class="doc-card"
            :class="{ 'selected': selectedIds.includes(doc.id) }"
          >
            <div class="card-header">
              <NCheckbox
                :checked="selectedIds.includes(doc.id)"
                @update:checked="(checked: boolean) => handleCardSelect(doc.id, checked)"
                class="card-checkbox"
              />
              <div class="doc-icon" :style="{ backgroundColor: sourceConfig[doc.source]?.color + '15' }">
                <svg width="32" height="32" viewBox="0 0 24 24" :fill="sourceConfig[doc.source]?.color">
                  <path :d="sourceConfig[doc.source]?.icon" />
                </svg>
              </div>
            </div>
            
            <div class="card-content">
              <h3 class="doc-name" :title="doc.name">{{ doc.name }}</h3>
              <p class="doc-alias" :title="doc.alias">{{ doc.alias }}</p>
              <p class="doc-description" :title="doc.description">{{ doc.description || '暂无描述' }}</p>
              
              <div class="doc-meta">
                <NTag :type="sourceConfig[doc.source]?.type" size="small" :bordered="false">
                  {{ sourceConfig[doc.source]?.label }}
                </NTag>
                <span class="doc-size">{{ getDocSize(doc) }}</span>
              </div>
              
              <div class="doc-url" :title="doc.url">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
                </svg>
                <span>{{ doc.url }}</span>
              </div>
              
              <div class="doc-date">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
                <span>{{ dayjs(doc.created_at).format('YYYY-MM-DD HH:mm') }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <NButton size="small" type="primary" ghost @click="handleEdit(doc)">
                <template #icon>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                </template>
                编辑
              </NButton>
              <NButton size="small" type="info" ghost @click="handleSettings(doc)">
                <template #icon>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="3"></circle>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                  </svg>
                </template>
                设置
              </NButton>
              <NPopconfirm @positive-click="() => handleDelete(doc.id)">
                <template #trigger>
                  <NButton size="small" type="error" ghost>
                    <template #icon>
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      </svg>
                    </template>
                    删除
                  </NButton>
                </template>
                确定要删除这个文档库吗？
              </NPopconfirm>
            </div>
          </div>
        </div>
      </NSpin>

      <div class="pagination-container">
        <div class="select-all-container">
          <NCheckbox
            :checked="isAllSelected"
            :indeterminate="isIndeterminate"
            @update:checked="handleSelectAll"
          >
            全选当前页
          </NCheckbox>
        </div>
        <NPagination
          v-model:page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-count="Math.ceil(pagination.total / pagination.pageSize)"
          :page-sizes="pagination.pageSizes"
          :show-size-picker="pagination.showSizePicker"
          :show-quick-jumper="true"
          :show-total="true"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </NCard>

    <NModal
      v-model:show="showModal"
      preset="card"
      :title="isSettings ? '文档库设置' : (isEdit ? '编辑文档库' : '新建文档库')"
      style="width: 600px"
      :mask-closable="false"
    >
      <NForm :model="formData" label-placement="left" label-width="100px">
        <NFormItem label="文档库名称" required>
          <NInput v-model:value="formData.name" placeholder="请输入文档库名称" />
        </NFormItem>
        <NFormItem label="别名" required>
          <NInput v-model:value="formData.alias" placeholder="请输入文档库别名" />
        </NFormItem>
        <NFormItem label="描述">
          <NInput
            v-model:value="formData.description"
            type="textarea"
            placeholder="请输入文档库描述"
            :autosize="{ minRows: 3, maxRows: 6 }"
          />
        </NFormItem>
        <NFormItem label="来源" required>
          <NSelect
            v-model:value="formData.source"
            :options="sourceOptions"
            placeholder="请选择文档库来源"
          />
        </NFormItem>
        <NFormItem label="URL" required>
          <NInput v-model:value="formData.url" placeholder="请输入文档库URL" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="handleModalClose">取消</NButton>
          <NButton type="primary" @click="handleSubmit">
            {{ isEdit || isSettings ? '保存' : '创建' }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
  </div>
</template>

<style scoped>

.doclibrary-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.doclibrary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.doc-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  border: 2px solid transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.doc-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #18a058;
}

.doc-card.selected {
  border-color: #18a058;
  background-color: #f0fdf4;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.card-checkbox {
  flex-shrink: 0;
}

.doc-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.3s ease;
}

.doc-card:hover .doc-icon {
  transform: scale(1.1);
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.doc-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.doc-alias {
  font-size: 14px;
  color: #64748b;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.doc-description {
  font-size: 13px;
  color: #94a3b8;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
  min-height: 39px;
}

.doc-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.doc-size {
  font-size: 12px;
  color: #64748b;
}

.doc-url {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.doc-url svg {
  flex-shrink: 0;
}

.doc-date {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #94a3b8;
}

.doc-date svg {
  flex-shrink: 0;
}

.card-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.card-actions .n-button {
  flex: 1;
}

.pagination-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding: 16px 0;
  flex-wrap: wrap;
  gap: 16px;
}

.select-all-container {
  font-size: 14px;
  color: #64748b;
}

@media (max-width: 768px) {
  .doclibrary-container {
    padding: 16px;
  }
  
  .doclibrary-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .card-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .pagination-container {
    flex-direction: column;
    align-items: stretch;
  }
  
  .select-all-container {
    text-align: center;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .card-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (min-width: 1025px) and (max-width: 1440px) {
  .card-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}
</style>
