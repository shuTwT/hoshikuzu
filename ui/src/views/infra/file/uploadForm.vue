<script lang="ts" setup>
import Uppy from '@uppy/core';
import Dashboard from '@uppy/dashboard';
import Zhcn from '@uppy/locales/lib/zh_CN';
import XHRUpload from '@uppy/xhr-upload';
import { BASE_URL } from '@/api/utils';
import { getToken } from '@/utils/auth';
import { getStorageStrategyListAll } from '@/api/infra/storage';

import '@uppy/core/css/style.min.css';
import '@uppy/dashboard/css/style.min.css';

let uppy:Uppy|null = null;


const storageStrategyList = ref<any[]>([])
const formData = reactive({
  currStorageStrategy:0
})

onMounted(()=>{
  uppy = new Uppy({
    locale:Zhcn,
  }).use(Dashboard,{
    inline:true,
    target:"#uppy-dashboard",
    height:"480px",
  })
  .use(XHRUpload,{
    endpoint:BASE_URL+"/v1/file/upload",
    onBeforeRequest:(xhr)=>{
      const data = getToken()
      if(data){
        const now = new Date().getTime()
        const expires = Number(data.expires) - now <= 0
        if(expires){

        }else{
          xhr.setRequestHeader("Authorization","Bearer "+data.accessToken)
        }
      }

    },
    getResponseData:(xhr)=>{
      const res = JSON.parse(xhr.response)
      return res.data
    }
  }).on('file-added',(file)=>{
    uppy?.setFileMeta(file.id,{
      storage_strategy:formData.currStorageStrategy
    })
  })
})

onMounted(()=>{
  getStorageStrategyListAll().then((res)=>{
    if(res.code === 200){
      for(let i=0;i<res.data.length;i++){
        if(res.data[i].master==true){
          formData.currStorageStrategy = res.data[i].id
          break;
        }
      }
      storageStrategyList.value = res.data

    }
  })
})

onBeforeUnmount(()=>{
  if(uppy){
    uppy.destroy()
    uppy = null
  }
})

</script>
<template>
  <div class="flex flex-col h-full">
    <div class="mb-2">
      <span>存储策略</span>
    </div>
    <div class="mb-2 w-full flex ">
      <div v-for="(item,index) in storageStrategyList" :key="index" class="mr-2">
        <n-button ghost :type="formData.currStorageStrategy === item.id ? 'primary' : 'default'" @click="formData.currStorageStrategy = item.id">
          {{ item.name }}
        </n-button>

      </div>
    </div>
    <div id="uppy-dashboard"></div>
  </div>
</template>
