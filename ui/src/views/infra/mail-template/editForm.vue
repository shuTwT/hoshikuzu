<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import type { FormProps } from './utils/types'

const props = defineProps<FormProps>()

const formRef = ref<FormInst | null>(null)
const formData = ref(props.formInline)

const rules: FormRules = {
  name: {
    required: true,
    message: '请输入模板名称',
    trigger: 'blur',
  },
  code: {
    required: true,
    message: '请输入模板编码',
    trigger: 'blur',
  },
  subject: {
    required: true,
    message: '请输入邮件主题',
    trigger: 'blur',
  },
  content: {
    required: true,
    message: '请输入邮件正文',
    trigger: 'blur',
  },
}

const contentTypeOptions = [
  { label: 'HTML', value: 'html' },
  { label: '纯文本', value: 'text' },
]

const getData = () => {
  return new Promise((resolve, reject) => {
    if (formRef.value) {
      formRef.value?.validate((errors) => {
        if (!errors) {
          const data = toRaw(formData.value)
          const submitData: any = {
            name: data.name,
            code: data.code,
            subject: data.subject,
            content: data.content,
            content_type: data.content_type || 'html',
            enabled: data.enabled,
          }

          if (data.remark) submitData.remark = data.remark

          resolve(submitData)
        } else {
          reject(errors)
        }
      })
    } else {
      reject(new Error('表单实例不存在'))
    }
  })
}

defineExpose({ getData })
</script>

<template>
  <n-form :model="formData" :rules="rules" ref="formRef" label-placement="left" label-width="120px">
    <n-form-item label="模板名称" path="name">
      <n-input v-model:value="formData.name" placeholder="请输入模板名称" />
    </n-form-item>
    <n-form-item label="模板编码" path="code">
      <n-input v-model:value="formData.code" placeholder="如 register_welcome，全局唯一" />
    </n-form-item>
    <n-form-item label="邮件主题" path="subject">
      <n-input
        v-model:value="formData.subject"
        placeholder="支持 {变量} 占位符，如：欢迎加入星屑，{nickname}"
      />
    </n-form-item>
    <n-form-item label="内容类型" path="content_type">
      <n-select v-model:value="formData.content_type" :options="contentTypeOptions" />
    </n-form-item>
    <n-form-item label="邮件正文" path="content">
      <n-input
        v-model:value="formData.content"
        type="textarea"
        placeholder="支持 {变量} 占位符与 HTML，如：<h1>你好 {nickname}</h1><p>感谢你的注册。</p>"
        :rows="6"
      />
    </n-form-item>
    <n-form-item label="是否启用" path="enabled">
      <n-switch v-model:value="formData.enabled" />
    </n-form-item>
    <n-form-item label="备注" path="remark">
      <n-input v-model:value="formData.remark" type="textarea" placeholder="请输入备注" :rows="2" />
    </n-form-item>
  </n-form>
</template>
