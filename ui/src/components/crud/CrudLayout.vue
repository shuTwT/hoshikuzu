<script setup lang="ts">
import { ProSearchForm, type ProSearchFormColumns } from 'pro-naive-ui';
import type { CrudLayoutProps } from './types';

const props = defineProps<CrudLayoutProps>()



const searchColumns:ProSearchFormColumns<any>=[] as any[]
</script>
<template>
  <div>
    <!-- 头部操作栏 -->
    <div class="header-section">
      <pro-search-form :form="form" :loading="false" :columns="searchColumns"/>
    </div>
    <!-- 内容 -->
     <n-data-table/>
  </div>
</template>
<style scoped>
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
