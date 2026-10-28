<script setup lang="ts">
import type { FormInst, FormRules, SelectOption, UploadCustomRequestOptions } from 'naive-ui'
import { CloudUploadOutline, Trash } from '@vicons/ionicons5'
import type { PromptFormProps } from './utils/types'
import type { EntPromptTag, ModelStorageStrategyListResp } from '@hoshikuzu/api-client'
import { apiClient, useApi } from '@/api'

const props = defineProps<PromptFormProps>()

const message = useMessage()

const formRef = ref<FormInst | null>(null)
const formData = ref(props.formInline)
const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入提示词正文', trigger: 'blur' }],
}

const tagOptions = ref<SelectOption[]>([])
const storageStrategyOptions = ref<SelectOption[]>([])
const storageStrategyId = ref<number | null>(null)
const uploading = ref(false)

// 复用文件上传接口，图片写入所选存储策略（默认为主存储策略）
const handleUpload = async ({ file, onFinish, onError }: UploadCustomRequestOptions) => {
  if (!storageStrategyId.value) {
    message.warning('请先在存储策略中配置并选择一个存储策略')
    onError()
    return
  }
  uploading.value = true
  try {
    const res = await useApi(apiClient.api.v1FileUploadCreate, {
      file: file.file as File,
      storage_strategy: storageStrategyId.value,
    })
    const url = res.data?.[0]?.url
    if (!url) {
      throw new Error('上传失败')
    }
    formData.value.images.push({
      image_url: url,
      model_name: '',
      sort: formData.value.images.length,
    })
    onFinish()
  } catch {
    onError()
  } finally {
    uploading.value = false
  }
}

const removeImage = (index: number) => {
  formData.value.images.splice(index, 1)
}

const getData = () => {
  return new Promise((resolve, reject) => {
    if (formRef.value) {
      formRef.value?.validate((errors) => {
        if (!errors) {
          resolve(toRaw(formData.value))
        } else {
          reject(errors)
        }
      })
    } else {
      reject(new Error('表单实例不存在'))
    }
  })
}

onMounted(async () => {
  try {
    const [tagRes, strategyRes] = await Promise.all([
      useApi(apiClient.api.v1PromptTagListList),
      useApi(apiClient.api.v1StorageStrategyListList),
    ])
    tagOptions.value = ((tagRes.data || []) as EntPromptTag[]).map((item) => ({
      label: item.name,
      value: item.name,
    }))
    const strategies = (strategyRes.data || []) as ModelStorageStrategyListResp[]
    storageStrategyOptions.value = strategies.map((item) => ({
      label: item.master ? `${item.name}（默认）` : item.name,
      value: item.id,
    }))
    storageStrategyId.value = (strategies.find((item) => item.master) ?? strategies[0])?.id ?? null
  } catch (e) {
    console.error(e)
  }
})

defineExpose({ getData })
</script>
<template>
  <n-form ref="formRef" :rules="rules" :model="formData" label-placement="top">
    <n-form-item label="标题" path="title">
      <n-input v-model:value="formData.title" placeholder="请输入标题" />
    </n-form-item>
    <n-form-item label="提示词" path="content">
      <n-input
        v-model:value="formData.content"
        type="textarea"
        :autosize="{ minRows: 4, maxRows: 10 }"
        placeholder="请输入提示词正文"
      />
    </n-form-item>
    <n-form-item label="负向提示词" path="negative_prompt">
      <n-input
        v-model:value="formData.negative_prompt"
        type="textarea"
        :autosize="{ minRows: 2, maxRows: 6 }"
        placeholder="可选"
      />
    </n-form-item>
    <n-form-item label="描述" path="description">
      <n-input v-model:value="formData.description" placeholder="可选" />
    </n-form-item>
    <n-grid :cols="2" x-gap="12">
      <n-form-item-gi label="标签" path="tags">
        <n-select
          v-model:value="formData.tags"
          multiple
          filterable
          tag
          :options="tagOptions"
          placeholder="选择或输入后回车新建标签"
        />
      </n-form-item-gi>
      <n-form-item-gi label="排序" path="sort">
        <n-input-number v-model:value="formData.sort" :min="0" class="w-full" />
      </n-form-item-gi>
    </n-grid>
    <n-form-item label="模型图片" path="images">
      <div class="w-full">
        <div class="flex items-center gap-2 mb-3">
          <n-select
            v-model:value="storageStrategyId"
            :options="storageStrategyOptions"
            placeholder="存储策略"
            style="width: 200px"
          />
          <n-upload
            accept="image/*"
            multiple
            :show-file-list="false"
            :custom-request="handleUpload"
          >
            <n-button :loading="uploading">
              <template #icon>
                <n-icon><CloudUploadOutline /></n-icon>
              </template>
              上传图片
            </n-button>
          </n-upload>
        </div>
        <n-empty v-if="formData.images.length === 0" description="暂无图片" class="py-4" />
        <div v-else class="space-y-2">
          <div
            v-for="(item, index) in formData.images"
            :key="item.image_url + index"
            class="flex items-center gap-3 p-2 bg-gray-50 dark:bg-gray-700 rounded-md"
          >
            <n-image :src="item.image_url" width="64" height="64" object-fit="cover" />
            <n-input v-model:value="item.model_name" placeholder="生成所用模型名" />
            <n-input-number
              v-model:value="item.sort"
              :min="0"
              placeholder="排序"
              style="width: 110px"
            />
            <n-button quaternary type="error" @click="removeImage(index)">
              <template #icon>
                <n-icon><Trash /></n-icon>
              </template>
            </n-button>
          </div>
        </div>
      </div>
    </n-form-item>
  </n-form>
</template>
