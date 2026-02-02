<script setup lang="ts">
import MarkdownEditor from '@/components/markdown/MarkdownEditor.vue'
import { useRoute } from 'vue-router'
import * as postApi from '@/api/content/post'
import { usePostHook } from '../utils/hook'
import { MdCatalog, MdPreview } from 'md-editor-v3'
import {register} from "@hoshikuzu/md-kit/src/index"

const { settingPost, savePost, publishPost, unpublishPost, importPost } = usePostHook()

register()

const route = useRoute()
const editorRef = shallowRef()
const showPreview = ref(false)

const publishStatus = ref(false)

const editorState = reactive({
  id:'post-editor'
})

const valueHtml = ref('<p>hello</p>')
const valueMarkdown = ref('')

const openSettingDialog = () => {
  settingPost({ id: route.query.id })
}

const handlePreview = () => {
  showPreview.value = true
}

const handleSave = () => {
  savePost({
    id: route.query.id,
    content: valueHtml.value,
    md_content: valueMarkdown.value,
    html_content: valueHtml.value,
  }).then(() => {
    getPostData()
  })
}

const handlePublish = () => {
  publishPost({ id: route.query.id }).then(() => {
    publishStatus.value = true
  })
}

const handleUnpublish = () => {
  unpublishPost({ id: route.query.id }).then(() => {
    publishStatus.value = false
  })
}

const handleImport = () =>{
  importPost({
    id:route.query.id
  }).then((res:any)=>{
    valueHtml.value = res.importContent
  })
}

const handleHtmlChange = (h:string)=>{
  valueHtml.value = h
}

const getPostData = () => {
  const id = route.query.id
  if (id) {
    postApi.queryPost(id + '').then((res) => {
      valueMarkdown.value = res.data.md_content
      if (res.data.status == 'draft') {
        publishStatus.value = false
      } else if (res.data.status == 'published') {
        publishStatus.value = true
      }
    })
  }
}

onMounted(() => {
  getPostData()
})

</script>
<template>
  <div class="p-6 flex gap-1">
    <div class="editor-wrapper w-[calc(100%-400px)] border border-[#ccc]">
      <MarkdownEditor :ref="editorRef" v-model:modelValue="valueMarkdown" :id="editorState.id" @onHtmlChanged="handleHtmlChange" />
    </div>
    <div class="w-[400px]">
      <div class="border border-[#ccc] bg-white flex flex-col">
        <div class="pt-6 flex justify-end pr-6 items-center">
          <n-button style="margin-right: 10px"> 历史版本 </n-button>
          <n-button style="margin-right: 10px" @click="handlePreview"> 预览 </n-button>
          <n-button style="margin-right: 10px" @click="openSettingDialog"> 设置 </n-button>
          <n-button @click="handleSave"> 保存 </n-button>
        </div>
        <div class="pt-6 flex justify-end pr-6 items-center">
          <n-button style="margin-right: 10px;" @click="handleImport">导入</n-button>
          <n-button v-if="!publishStatus" type="primary" @click="handlePublish">
            发布
          </n-button>
          <n-button v-else type="primary" @click="handleUnpublish">
            取消发布
          </n-button>
        </div>
        <n-divider />
        <n-tabs type="segment" animated>
          <n-tab-pane name="chap1" tab="大纲">
            <!-- <ul v-html="catalogHtml" class="p-2"></ul> -->
             <MdCatalog :editorId="editorState.id" />
          </n-tab-pane>
          <n-tab-pane name="chap2" tab="详情"></n-tab-pane>
        </n-tabs>
      </div>
    </div>
    <n-modal v-model:show="showPreview" preset="card" style="height: 100vh">
      <div class="w-full h-full"  tabindex="1">
        <n-scrollbar style="height: calc(100vh - 80px);">
          <MdPreview  :model-value="valueHtml"/>
      </n-scrollbar>
      </div>
    </n-modal>
  </div>
</template>
<style scoped>
.editor-wrapper {
  --editor-wrapper-height: calc(
    100vh - var(--pro-layout-footer-height) - var(--pro-layout-nav-height) - 48px
  );
  --editor-height: calc(var(--editor-wrapper-height) - 83px);
  height: var(--editor-wrapper-height);
}
</style>
