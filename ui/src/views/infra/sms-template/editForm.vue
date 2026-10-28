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
  content: {
    required: true,
    message: '请输入模板内容',
    trigger: 'blur',
  },
}

const getData = () => {
  return new Promise((resolve, reject) => {
    if (formRef.value) {
      formRef.value?.validate((errors) => {
        if (!errors) {
          const data = toRaw(formData.value)
          const submitData: any = {
            name: data.name,
            code: data.code,
            content: data.content,
            enabled: data.enabled,
          }

          if (data.sign_name) submitData.sign_name = data.sign_name
          if (data.provider_template_id) submitData.provider_template_id = data.provider_template_id
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
      <n-input v-model:value="formData.code" placeholder="如 register_verify_code，全局唯一" />
    </n-form-item>
    <n-form-item label="短信签名" path="sign_name">
      <n-input v-model:value="formData.sign_name" placeholder="请输入短信签名" />
    </n-form-item>
    <n-form-item label="服务商模板ID" path="provider_template_id">
      <n-input
        v-model:value="formData.provider_template_id"
        placeholder="请输入服务商模板ID（可选）"
      />
    </n-form-item>
    <n-form-item label="模板内容" path="content">
      <n-input
        v-model:value="formData.content"
        type="textarea"
        placeholder="支持 {变量} 占位符，如：您的验证码为{code}，{minutes}分钟内有效"
        :rows="4"
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
