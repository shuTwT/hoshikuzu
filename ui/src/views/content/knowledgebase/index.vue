<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  NButton,
  NSpace,
  NModal,
  NInput,
  NSelect,
  NInputNumber,
  NForm,
  NFormItem,
  NPagination,
  NSpin,
  NPopconfirm,
  NTag,
  NCard,
  NCheckbox,
  NEmpty,
  useMessage,
  type DataTableColumns,
} from 'naive-ui'
import * as knowledgebaseApi from '@/api/content/knowledgebase'
import dayjs from 'dayjs'

const message = useMessage()

interface KnowledgeBase {
  id: number
  name: string
  model_provider: string
  model: string
  vector_dimension: number
  max_batch_document_count: number
  created_at: string
  updated_at: string
}

interface KnowledgeBaseFormData {
  name: string
  model_provider: string
  model: string
  vector_dimension: number
  max_batch_document_count: number
}

const loading = ref(false)
const data = ref<KnowledgeBase[]>([])
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
const currentId = ref<number | null>(null)
const formData = ref<KnowledgeBaseFormData>({
  name: '',
  model_provider: '',
  model: '',
  vector_dimension: 1536,
  max_batch_document_count: 100,
})

const modelProviderOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Google', value: 'google' },
  { label: 'Azure', value: 'azure' },
  { label: 'Cohere', value: 'cohere' },
  { label: 'HuggingFace', value: 'huggingface' },
  { label: 'Local', value: 'local' },
]

const modelProviderConfig: Record<string, { type: any; label: string; icon: string; color: string }> = {
  openai: { type: 'info', label: 'OpenAI', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#10a37f' },
  anthropic: { type: 'success', label: 'Anthropic', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#d97706' },
  google: { type: 'warning', label: 'Google', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#4285f4' },
  azure: { type: 'default', label: 'Azure', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#0078d4' },
  cohere: { type: 'error', label: 'Cohere', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#f59e0b' },
  huggingface: { type: 'info', label: 'HuggingFace', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#ffd21e' },
  local: { type: 'success', label: 'Local', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z', color: '#52c41a' },
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
    message.warning('请先选择要删除的知识库')
    return
  }

  try {
    const deletePromises = selectedIds.value.map(async id => {
      try {
        await knowledgebaseApi.deleteKnowledgeBase(id)
        return { success: true, id }
      } catch (error) {
        console.error(`删除知识库 ${id} 失败:`, error)
        return { success: false, id }
      }
    })
    
    const results = await Promise.all(deletePromises)
    const successCount = results.filter(r => r.success).length
    const failedCount = results.filter(r => !r.success).length
    
    if (successCount > 0) {
      message.success(`成功删除 ${successCount} 个知识库${failedCount > 0 ? `，失败 ${failedCount} 个` : ''}`)
    } else {
      message.error('删除失败，请稍后重试')
    }
    
    selectedIds.value = []
    fetchKnowledgeBases()
  } catch (error: any) {
    console.error('批量删除失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || '批量删除失败，请稍后重试'
    message.error(errorMessage)
  }
}

const fetchKnowledgeBases = async () => {
  loading.value = true
  try {
    const res = await knowledgebaseApi.getKnowledgeBasePage({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
    })
    data.value = res.data.records
    pagination.value.total = res.data.total
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    message.error('获取知识库列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  formData.value = {
    name: '',
    model_provider: '',
    model: '',
    vector_dimension: 1536,
    max_batch_document_count: 100,
  }
  showModal.value = true
}

const handleEdit = (row: KnowledgeBase) => {
  isEdit.value = true
  currentId.value = row.id
  formData.value = {
    name: row.name,
    model_provider: row.model_provider,
    model: row.model,
    vector_dimension: row.vector_dimension,
    max_batch_document_count: row.max_batch_document_count,
  }
  showModal.value = true
}

const handleDelete = async (id: number) => {
  try {
    await knowledgebaseApi.deleteKnowledgeBase(id)
    message.success('删除成功')
    selectedIds.value = selectedIds.value.filter(selectedId => selectedId !== id)
    fetchKnowledgeBases()
  } catch (error: any) {
    console.error('删除失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || '删除失败，请稍后重试'
    message.error(errorMessage)
  }
}

const handleSubmit = async () => {
  if (!formData.value.name.trim()) {
    message.error('请输入知识库名称')
    return
  }
  if (!formData.value.model_provider) {
    message.error('请选择模型供应商')
    return
  }
  if (!formData.value.model.trim()) {
    message.error('请输入模型名称')
    return
  }
  if (!formData.value.vector_dimension || formData.value.vector_dimension <= 0) {
    message.error('请输入有效的向量维度')
    return
  }
  if (!formData.value.max_batch_document_count || formData.value.max_batch_document_count <= 0) {
    message.error('请输入有效的最大批量处理文档数量')
    return
  }

  try {
    if (isEdit.value && currentId.value) {
      await knowledgebaseApi.updateKnowledgeBase(currentId.value, formData.value)
      message.success('更新成功')
    } else {
      await knowledgebaseApi.createKnowledgeBase(formData.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchKnowledgeBases()
  } catch (error: any) {
    console.error('操作失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || (isEdit.value ? '更新失败，请稍后重试' : '创建失败，请稍后重试')
    message.error(errorMessage)
  }
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  fetchKnowledgeBases()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  fetchKnowledgeBases()
}

const handleRefresh = () => {
  fetchKnowledgeBases()
  message.success('刷新成功')
}

const handleModalClose = () => {
  showModal.value = false
  isEdit.value = false
  currentId.value = null
  formData.value = {
    name: '',
    model_provider: '',
    model: '',
    vector_dimension: 1536,
    max_batch_document_count: 100,
  }
}

onMounted(() => {
  fetchKnowledgeBases()
})
</script>

<template>
  <div class="knowledgebase-container">
    <NCard :bordered="false" class="knowledgebase-card">
      <div class="knowledgebase-header">
        <h1 class="page-title">知识库管理</h1>
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
            新建知识库
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
          <NEmpty description="暂无知识库数据" />
        </div>
        
        <div v-else class="card-grid">
          <div
            v-for="kb in data"
            :key="kb.id"
            class="kb-card"
            :class="{ 'selected': selectedIds.includes(kb.id) }"
          >
            <div class="card-header">
              <NCheckbox
                :checked="selectedIds.includes(kb.id)"
                @update:checked="(checked: boolean) => handleCardSelect(kb.id, checked)"
                class="card-checkbox"
              />
              <div class="kb-icon" :style="{ backgroundColor: modelProviderConfig[kb.model_provider]?.color + '15' }">
                <svg width="32" height="32" viewBox="0 0 24 24" :fill="modelProviderConfig[kb.model_provider]?.color">
                  <path :d="modelProviderConfig[kb.model_provider]?.icon" />
                </svg>
              </div>
            </div>
            
            <div class="card-content">
              <h3 class="kb-name" :title="kb.name">{{ kb.name }}</h3>
              <p class="kb-model" :title="kb.model">{{ kb.model }}</p>
              
              <div class="kb-meta">
                <NTag :type="modelProviderConfig[kb.model_provider]?.type" size="small" :bordered="false">
                  {{ modelProviderConfig[kb.model_provider]?.label }}
                </NTag>
              </div>
              
              <div class="kb-info">
                <div class="info-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 2L2 7l10 5v10l-10 5z"></path>
                    <path d="M12 2L22 7l-10 5v10l10 5z"></path>
                  </svg>
                  <span>向量维度: {{ kb.vector_dimension }}</span>
                </div>
                <div class="info-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h8a2 2 0 002-2V4a2 2 0 00-2-2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"></path>
                  </svg>
                  <span>批量处理: {{ kb.max_batch_document_count }}</span>
                </div>
              </div>
              
              <div class="kb-date">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
                <span>{{ dayjs(kb.created_at).format('YYYY-MM-DD HH:mm') }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <NButton size="small" type="primary" ghost @click="handleEdit(kb)">
                <template #icon>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                </template>
                编辑
              </NButton>
              <NPopconfirm @positive-click="() => handleDelete(kb.id)">
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
                确定要删除这个知识库吗？
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
      :title="isEdit ? '编辑知识库' : '新建知识库'"
      style="width: 600px"
      :mask-closable="false"
    >
      <NForm :model="formData" label-placement="left" label-width="140px">
        <NFormItem label="知识库名称" required>
          <NInput v-model:value="formData.name" placeholder="请输入知识库名称" />
        </NFormItem>
        <NFormItem label="模型供应商" required>
          <NSelect
            v-model:value="formData.model_provider"
            :options="modelProviderOptions"
            placeholder="请选择模型供应商"
          />
        </NFormItem>
        <NFormItem label="模型名称" required>
          <NInput v-model:value="formData.model" placeholder="请输入模型名称" />
        </NFormItem>
        <NFormItem label="向量维度" required>
          <NInputNumber
            v-model:value="formData.vector_dimension"
            :min="1"
            :max="10000"
            placeholder="请输入向量维度"
          />
        </NFormItem>
        <NFormItem label="最大批量处理" required>
          <NInputNumber
            v-model:value="formData.max_batch_document_count"
            :min="1"
            :max="1000"
            placeholder="请输入最大批量处理文档数量"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="handleModalClose">取消</NButton>
          <NButton type="primary" @click="handleSubmit">
            {{ isEdit ? '保存' : '创建' }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.knowledgebase-container {
  padding: 24px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.knowledgebase-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.knowledgebase-header {
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

.kb-card {
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

.kb-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #18a058;
}

.kb-card.selected {
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

.kb-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.3s ease;
}

.kb-card:hover .kb-icon {
  transform: scale(1.1);
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.kb-name {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kb-model {
  font-size: 14px;
  color: #64748b;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kb-meta {
  display: flex;
  justify-content: flex-start;
  margin-top: 8px;
}

.kb-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #64748b;
}

.info-item svg {
  flex-shrink: 0;
}

.kb-date {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #94a3b8;
  margin-top: auto;
}

.kb-date svg {
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
  .knowledgebase-container {
    padding: 16px;
  }
  
  .knowledgebase-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
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
