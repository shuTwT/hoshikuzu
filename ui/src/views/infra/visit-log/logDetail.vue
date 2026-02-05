<script setup lang="ts">
/**
 * 访问日志详情组件
 * 用于显示访问日志的详细信息
 */
import type { LogDetailProps } from './utils/types'
import dayjs from 'dayjs';

const { data } = defineProps<LogDetailProps>()
</script>

<template>
  <div style="padding: 20px">
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">编号：</span>
      <span>{{ data.id }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">访问IP：</span>
      <span>{{ data.ip }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">访问路径：</span>
      <span>{{ data.path }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">操作系统：</span>
      <span>{{ data.os || '-' }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">浏览器：</span>
      <span>{{ data.browser || '-' }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">设备：</span>
      <span>{{ data.device || '-' }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">用户代理：</span>
      <span style="word-break: break-all">{{ data.user_agent || '-' }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">创建时间：</span>
      <span>{{ dayjs(data.created_at).format('YYYY-MM-DD HH:mm:ss') }}</span>
    </div>
    <div style="margin-bottom: 12px">
      <span style="font-weight: bold; width: 100px; display: inline-block">更新时间：</span>
      <span>{{ dayjs(data.updated_at).format('YYYY-MM-DD HH:mm:ss') }}</span>
    </div>
  </div>
</template>
