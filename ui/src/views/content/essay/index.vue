<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import {
  NButton,
  NSpace,
  NModal,
  NInput,
  NSwitch,
  NUpload,
  NImage,
  NTag,
  useMessage,
  useDialog,
  NCard,
  NForm,
  NFormItem,
  NPagination,
  NSpin,
  NAvatar,
  NPopconfirm,
} from 'naive-ui'
import type { UploadFileInfo, MenuOption } from 'naive-ui'
import * as essayApi from '@/api/content/essay'
import dayjs from 'dayjs'

const message = useMessage()
const dialog = useDialog()

// 用户信息（模拟数据，可替换为真实用户信息）
const currentUser = ref({
  name: 'Administrator',
  avatar: 'https://picsum.photos/id/1/200/200',
  status: 'online'
})

interface Essay {
  id: number
  content: string
  draft: boolean
  images: string[]
  create_at: string
  user?: {
    name: string
    avatar: string
  }
}

interface EssayFormData {
  content: string
  draft: boolean
  images: string[]
}

const loading = ref(false)
const data = ref<Essay[]>([])
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 40],
})

const showModal = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const formData = ref<EssayFormData>({
  content: '',
  draft: false,
  images: [],
})

const fileList = ref<UploadFileInfo[]>([])

// 格式化时间
const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 计算图片显示样式
const getImageGridClass = (count: number) => {
  if (count === 1) return 'image-single'
  if (count === 2) return 'image-double'
  if (count <= 4) return 'image-grid-2'
  return 'image-grid-3'
}

// 获取说说列表
const fetchEssays = async () => {
  loading.value = true
  try {
    const res = await essayApi.getEssayPage({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
    })
    
    // 为每个说说添加模拟用户信息
    data.value = res.data.records.map((essay: any) => ({
      ...essay,
      user: {
        name: currentUser.value.name,
        avatar: currentUser.value.avatar
      }
    }))
    
    pagination.value.total = res.data.total
  } catch (error) {
    console.error('获取说说列表失败:', error)
    message.error('获取说说列表失败')
  } finally {
    loading.value = false
  }
}

// 新建说说
const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  formData.value = {
    content: '',
    draft: false,
    images: [],
  }
  fileList.value = []
  showModal.value = true
}

// 编辑说说
const handleEdit = (row: Essay) => {
  isEdit.value = true
  currentId.value = row.id
  formData.value = {
    content: row.content,
    draft: row.draft,
    images: row.images || [],
  }
  fileList.value = (row.images || []).map((url, index) => ({
    id: String(index),
    name: `image-${index}`,
    status: 'finished',
    url: url,
  }))
  showModal.value = true
}

// 删除说说
const handleDelete = async (id: number) => {
  try {
    await essayApi.deleteEssay(id)
    message.success('删除成功')
    fetchEssays()
  } catch (error) {
    console.error('删除失败:', error)
    message.error('删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formData.value.content.trim()) {
    message.error('请输入说说内容')
    return
  }

  try {
    if (isEdit.value && currentId.value) {
      await essayApi.updateEssay(currentId.value, formData.value)
      message.success('更新成功')
    } else {
      await essayApi.createEssay(formData.value)
      message.success('创建成功')
    }
    showModal.value = false
    fetchEssays()
  } catch (error) {
    console.error('操作失败:', error)
    message.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

// 上传文件变化处理
const handleUploadChange = (options: { fileList: UploadFileInfo[] }) => {
  fileList.value = options.fileList
  formData.value.images = fileList.value
    .filter((file) => file.status === 'finished' && file.url)
    .map((file) => file.url!)
}

// 页面变化处理
const handlePageChange = (page: number) => {
  pagination.value.page = page
  fetchEssays()
}

// 每页数量变化处理
const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  fetchEssays()
}

// 初始化
onMounted(() => {
  fetchEssays()
})
</script>

<template>
  <div class="container-fluid p-6">
    <div class="max-w-[1600px] mx-auto">
    <!-- 顶部导航栏 -->
    <div class="essay-header">
      <div class="header-left">
        <h1 class="page-title">说说管理</h1>
      </div>
      <div class="header-right">
        <n-button 
          type="primary" 
          @click="handleCreate"
          round
          size="large"
        >
          <template #icon>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
          </template>
          发布说说
        </n-button>
      </div>
    </div>

    <!-- 说说列表 -->
    <div class="essay-list">
      <n-spin :show="loading" size="large">
        <!-- 说说卡片 -->
        <n-card 
          v-for="item in data" 
          :key="item.id"
          class="essay-card"
          :bordered="false"
          :shadow="'hover'"
        >
          <!-- 卡片头部：用户信息 -->
          <div class="card-header">
            <div class="user-info">
              <n-avatar 
                :src="item.user?.avatar || currentUser.avatar"
                size="large"
                class="user-avatar"
              />
              <div class="user-details">
                <div class="user-name">{{ item.user?.name || currentUser.name }}</div>
                <div class="post-time">{{ formatTime(item.create_at) }}</div>
              </div>
            </div>
            
            <!-- 状态标签 -->
            <n-tag 
              :type="item.draft ? 'warning' : 'success'"
              round
              size="small"
            >
              {{ item.draft ? '草稿' : '已发布' }}
            </n-tag>
          </div>

          <!-- 卡片内容 -->
          <div class="card-content">
            <div class="essay-content">{{ item.content }}</div>
            
            <!-- 图片展示 -->
            <div 
              v-if="item.images && item.images.length > 0" 
              class="essay-images"
              :class="getImageGridClass(item.images.length)"
            >
              <n-image 
                v-for="(image, index) in item.images" 
                :key="index"
                :src="image"
                class="essay-image"
                :preview-src-list="item.images"
                :initial-index="index"
              />
            </div>
          </div>

          <!-- 卡片底部：操作按钮 -->
          <div class="card-footer">
            <n-space>
              <n-popconfirm @positive-click="handleEdit(item)" placement="top">
                <template #trigger>
                  <n-button size="small" quaternary>
                    <template #icon>
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                      </svg>
                    </template>
                    编辑
                  </n-button>
                </template>
                <template #default>
                  <span>确定要编辑这条说说吗？</span>
                </template>
              </n-popconfirm>
              
              <n-popconfirm @positive-click="handleDelete(item.id)" placement="top">
                <template #trigger>
                  <n-button size="small" quaternary type="error">
                    <template #icon>
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      </svg>
                    </template>
                    删除
                  </n-button>
                </template>
                <template #default>
                  <span>确定要删除这条说说吗？</span>
                </template>
              </n-popconfirm>
            </n-space>
          </div>
        </n-card>
      </n-spin>

      <!-- 空状态 -->
      <div v-if="data.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <div class="empty-text">暂无说说</div>
        <n-button type="primary" @click="handleCreate" style="margin-top: 16px;">
          发布第一条说说
        </n-button>
      </div>

      <!-- 分页 -->
      <div v-if="pagination.total > 0" class="pagination-container">
        <n-pagination
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
    </div>
    </div>
  </div>

  <!-- 新建/编辑模态框 -->
  <n-modal
    v-model:show="showModal"
    preset="card"
    :title="isEdit ? '编辑说说' : '发布说说'"
    style="width: 90%; max-width: 800px;"
    :mask-closable="false"
    :close-on-esc="false"
    :destroy-on-close="true"
  >
    <n-form 
      :model="formData" 
      label-placement="left" 
      label-width="80px"
      class="essay-form"
    >
      <n-form-item label="内容" required>
        <n-input
          v-model:value="formData.content"
          type="textarea"
          placeholder="分享你的想法..."
          :autosize="{ minRows: 6, maxRows: 12 }"
          class="content-input"
        />
      </n-form-item>

      <n-form-item label="图片">
        <n-upload
          v-model:file-list="fileList"
          multiple
          list-type="image-card"
          :max="9"
          @change="handleUploadChange"
          class="image-upload"
        >
          <div class="upload-placeholder">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            <div>点击上传图片</div>
            <div class="upload-hint">最多上传9张</div>
          </div>
        </n-upload>
      </n-form-item>

      <n-form-item label="状态" class="status-item">
        <div class="status-switch">
          <span class="status-label">{{ formData.draft ? '保存为草稿' : '立即发布' }}</span>
          <n-switch 
            v-model:value="formData.draft"
            size="large"
            class="status-toggle"
          />
        </div>
      </n-form-item>
    </n-form>

    <template #footer>
      <n-space justify="end" size="large">
        <n-button @click="showModal = false" size="large">取消</n-button>
        <n-button 
          type="primary" 
          @click="handleSubmit"
          size="large"
          round
        >
          {{ isEdit ? '保存修改' : '发布' }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<style scoped>
/* 全局容器 */
.essay-container {
  background-color: #f5f7fa;
  padding: 24px;
}

/* 头部样式 */
.essay-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0 8px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

/* 列表容器 */
.essay-list {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
}

/* 说说卡片样式 */
.essay-card {
  border-radius: 16px;
  background-color: #ffffff;
  padding: 20px;
  margin: 0.75rem 0;
  transition: all 0.3s ease;
}

.essay-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  border: 2px solid #e8eaed;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-weight: 600;
  color: #1a1a1a;
  font-size: 16px;
}

.post-time {
  font-size: 14px;
  color: #666;
}

/* 卡片内容 */
.card-content {
  margin-bottom: 16px;
}

.essay-content {
  font-size: 16px;
  line-height: 1.6;
  color: #333;
  margin-bottom: 16px;
  white-space: pre-wrap;
  word-break: break-word;
}

/* 图片网格 */
.essay-images {
  display: grid;
  gap: 8px;
  border-radius: 8px;
  overflow: hidden;
}

.image-single {
  grid-template-columns: 1fr;
}

.image-double {
  grid-template-columns: repeat(2, 1fr);
}

.image-grid-2 {
  grid-template-columns: repeat(2, 1fr);
}

.image-grid-3 {
  grid-template-columns: repeat(3, 1fr);
}

.essay-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.essay-image:hover {
  transform: scale(1.02);
}

/* 卡片底部 */
.card-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #999;
  text-align: center;
}

.empty-icon {
  opacity: 0.3;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  margin-bottom: 8px;
}

/* 分页容器 */
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 20px;
}

/* 表单样式 */
.essay-form {
  padding: 16px 0;
}

.content-input {
  font-size: 16px;
  resize: vertical;
  border-radius: 12px;
}

.image-upload {
  margin-top: 8px;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #999;
}

.upload-hint {
  font-size: 12px;
  color: #bbb;
}

.status-item {
  margin-top: 8px;
}

.status-switch {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-label {
  font-size: 14px;
  color: #666;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .essay-container {
    padding: 16px;
  }
  
  .essay-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .essay-list {
    max-width: 100%;
  }
  
  .essay-card {
    padding: 16px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .essay-content {
    font-size: 14px;
  }
  
  .image-grid-3 {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .essay-container {
    padding: 12px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .essay-card {
    padding: 12px;
  }
  
  .essay-images {
    gap: 4px;
  }
  
  .image-upload {
    margin-top: 4px;
  }
}
</style>
