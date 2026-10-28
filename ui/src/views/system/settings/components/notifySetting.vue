<script setup lang="ts">
import type { FormInst } from 'naive-ui'
import { apiClient, useApi } from '@/api'

const message = useMessage()

const notificationFormRef = ref<FormInst | null>(null)

const encryptionOptions = [
  { label: 'TLS', value: 'tls' },
  { label: 'SSL', value: 'ssl' },
  { label: '无', value: 'none' },
]

// 通知设置（存储于 notify 组）
const defaultNotifyForm = {
  enableEmailNotification: true,
  enableSmsNotification: false,
}
const notificationForm = ref({ ...defaultNotifyForm })

// 邮件配置（存储于 email 组，启用邮件通知时显示）
const defaultEmailForm = {
  smtpHost: '',
  smtpPort: 587,
  smtpUsername: '',
  smtpPassword: '',
  smtpEncryption: 'tls',
  senderEmail: '',
  senderName: '',
}
const emailForm = ref({ ...defaultEmailForm })

const notificationLoading = ref(false)
const testMailLoading = ref(false)
const testMailTo = ref('')

// 保存通知设置（通知与邮件配置分两组保存，保持 email 组数据兼容）
const saveNotificationSettings = async () => {
  notificationLoading.value = true
  try {
    await Promise.all([
      useApi(apiClient.api.v1SettingsJsonSaveCreate, 'notify', notificationForm.value),
      useApi(apiClient.api.v1SettingsJsonSaveCreate, 'email', emailForm.value),
    ])
    await new Promise((resolve) => setTimeout(resolve, 1000))
    onSearch()
    message.success('通知设置保存成功')
  } catch {
    message.error('通知设置保存失败')
  } finally {
    notificationLoading.value = false
  }
}

const onSearch = async () => {
  const [notifyRes, emailRes] = await Promise.all([
    useApi(apiClient.api.v1SettingsJsonDetail, 'notify'),
    useApi(apiClient.api.v1SettingsJsonDetail, 'email'),
  ])
  notificationForm.value = Object.assign({}, defaultNotifyForm, notifyRes.data)
  emailForm.value = Object.assign({}, defaultEmailForm, emailRes.data)
}

// 发送测试邮件，验证 SMTP 配置是否正确
const sendTestMail = async () => {
  if (!testMailTo.value) {
    message.warning('请输入收件人邮箱')
    return
  }
  testMailLoading.value = true
  try {
    await useApi(apiClient.api.v1MailTestCreate, { to: testMailTo.value })
    message.success('测试邮件发送成功')
  } catch {
    message.error('测试邮件发送失败，请检查 SMTP 配置')
  } finally {
    testMailLoading.value = false
  }
}

onMounted(() => {
  onSearch()
})
</script>
<template>
  <n-form
    ref="notificationFormRef"
    :model="notificationForm"
    label-placement="left"
    label-width="auto"
    require-mark-placement="right-hanging"
    class="settings-form"
  >
    <n-form-item label="启用邮件通知" path="enableEmailNotification">
      <n-switch v-model:value="notificationForm.enableEmailNotification" />
    </n-form-item>
    <n-form-item label="启用短信通知" path="enableSmsNotification">
      <n-switch v-model:value="notificationForm.enableSmsNotification" />
    </n-form-item>
    <template v-if="notificationForm.enableEmailNotification">
      <n-divider>邮件配置</n-divider>
      <n-form-item label="SMTP服务器" path="smtpHost">
        <n-input v-model:value="emailForm.smtpHost" placeholder="例如: smtp.gmail.com" />
      </n-form-item>
      <n-form-item label="SMTP端口" path="smtpPort">
        <n-input-number v-model:value="emailForm.smtpPort" :min="1" :max="65535" />
      </n-form-item>
      <n-form-item label="SMTP用户名" path="smtpUsername">
        <n-input v-model:value="emailForm.smtpUsername" placeholder="请输入SMTP用户名" />
      </n-form-item>
      <n-form-item label="SMTP密码" path="smtpPassword">
        <n-input
          v-model:value="emailForm.smtpPassword"
          type="password"
          placeholder="请输入SMTP密码"
          show-password-on="mousedown"
        />
      </n-form-item>
      <n-form-item label="加密方式" path="smtpEncryption">
        <n-select v-model:value="emailForm.smtpEncryption" :options="encryptionOptions" />
      </n-form-item>
      <n-form-item label="发件人邮箱" path="senderEmail">
        <n-input v-model:value="emailForm.senderEmail" placeholder="请输入发件人邮箱地址" />
      </n-form-item>
      <n-form-item label="发件人名称" path="senderName">
        <n-input v-model:value="emailForm.senderName" placeholder="请输入发件人名称" />
      </n-form-item>
    </template>
    <n-form-item>
      <n-button type="primary" @click="saveNotificationSettings" :loading="notificationLoading">
        保存通知设置
      </n-button>
      <n-input
        v-model:value="testMailTo"
        placeholder="收件人邮箱"
        style="width: 240px; margin-left: 12px"
      />
      <n-button
        @click="sendTestMail"
        :loading="testMailLoading"
        :disabled="!notificationForm.enableEmailNotification"
        style="margin-left: 12px"
      >
        发送测试邮件
      </n-button>
    </n-form-item>
  </n-form>
</template>
