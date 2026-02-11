<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import type { FormProps } from './utils/types'

const props = defineProps<FormProps>()

const formRef = ref<FormInst | null>(null)
const formData = ref(props.formInline)

const rules: FormRules = {
  name: {
    required: true,
    message: '请输入商品名称',
    trigger: 'blur',
  },
  sku: {
    required: true,
    message: '请输入商品SKU',
    trigger: 'blur',
  },
  price: {
    required: true,
    type: 'number',
    message: '请输入商品价格',
    trigger: 'blur',
  },
  stock: {
    required: true,
    type: 'number',
    message: '请输入库存数量',
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
            sku: data.sku,
            price: Math.round(data.price * 100),
            stock: data.stock,
            original_price: data.original_price ? Math.round(data.original_price * 100) : undefined,
            cost_price: data.cost_price ? Math.round(data.cost_price * 100) : undefined,
            min_stock: data.min_stock,
            category_id: data.category_id,
            brand: data.brand,
            unit: data.unit,
            weight: data.weight,
            volume: data.volume,
            description: data.description,
            short_description: data.short_description,
            images: data.images,
            attributes: data.attributes,
            tags: data.tags,
            active: data.active,
            featured: data.featured,
            digital: data.digital,
            meta_title: data.meta_title,
            meta_description: data.meta_description,
            meta_keywords: data.meta_keywords,
            sort_order: data.sort_order
          }

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
    <n-tabs type="segment" animated>
      <n-tab-pane name="base" tab="基本信息">
        <n-form-item label="商品名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入商品名称" />
        </n-form-item>
        <n-form-item label="商品价格(元)" path="price">
          <n-input-number v-model:value="formData.price" :min="0" :precision="2" placeholder="请输入商品价格"
            style="width: 100%" />
        </n-form-item>
        <n-form-item label="原价(元)" path="original_price">
          <n-input-number v-model:value="formData.original_price" :min="0" :precision="2" placeholder="请输入原价"
            style="width: 100%" />
        </n-form-item>
        <n-form-item label="成本价(元)" path="cost_price">
          <n-input-number v-model:value="formData.cost_price" :min="0" :precision="2" placeholder="请输入成本价"
            style="width: 100%" />
        </n-form-item>
        <n-form-item label="库存数量" path="stock">
          <n-input-number v-model:value="formData.stock" :min="0" placeholder="请输入库存数量" style="width: 100%" />
        </n-form-item>
        <n-form-item label="简短描述" path="short_description">
          <n-input v-model:value="formData.short_description" type="textarea" placeholder="请输入简短描述" :rows="2" />
        </n-form-item>

        <n-form-item label="商品标签" path="tags">
          <n-dynamic-tags v-model:value="formData.tags" />
        </n-form-item>
        <n-form-item label="排序" path="sort_order">
          <n-input-number v-model:value="formData.sort_order" :min="0" placeholder="请输入排序" style="width: 100%" />
        </n-form-item>
        <n-form-item label="是否上架" path="active">
          <n-switch v-model:value="formData.active" />
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="sku" tab="商品SKU">
        <n-form-item label="商品SKU" path="sku">
          <n-input v-model:value="formData.sku" placeholder="请输入商品SKU" />
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="description" tab="商品介绍">
        <n-form-item path="description" label-placement="top">
          <n-input v-model:value="formData.description" type="textarea" placeholder="请输入商品描述" :rows="4" />
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="restrictions" tab="商品限制">
        <n-form-item label="是否推荐" path="featured">
          <n-switch v-model:value="formData.featured" />
        </n-form-item>
        <n-form-item label="是否数字商品" path="digital">
          <n-switch v-model:value="formData.digital" />
        </n-form-item>
        <n-form-item label="最低库存预警" path="min_stock">
          <n-input-number v-model:value="formData.min_stock" :min="0" placeholder="请输入最低库存预警" style="width: 100%" />
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="seo" tab="SEO设置">
        <n-form-item label="SEO标题" path="meta_title">
          <n-input v-model:value="formData.meta_title" placeholder="请输入SEO标题" />
        </n-form-item>
        <n-form-item label="SEO描述" path="meta_description">
          <n-input v-model:value="formData.meta_description" type="textarea" placeholder="请输入SEO描述" :rows="2" />
        </n-form-item>
        <n-form-item label="SEO关键词" path="meta_keywords">
          <n-input v-model:value="formData.meta_keywords" placeholder="请输入SEO关键词" />
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="config" tab="配置参数">
        <n-form-item label="品牌" path="brand">
          <n-input v-model:value="formData.brand" placeholder="请输入品牌" />
        </n-form-item>
        <n-form-item label="单位" path="unit">
          <n-input v-model:value="formData.unit" placeholder="请输入单位" />
        </n-form-item>
        <n-form-item label="重量(kg)" path="weight">
          <n-input-number v-model:value="formData.weight" :min="0" :precision="2" placeholder="请输入重量"
            style="width: 100%" />
        </n-form-item>
        <n-form-item label="体积(立方米)" path="volume">
          <n-input-number v-model:value="formData.volume" :min="0" :precision="3" placeholder="请输入体积"
            style="width: 100%" />
        </n-form-item>
      </n-tab-pane>
    </n-tabs>






  </n-form>
</template>
