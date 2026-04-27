<script setup lang="ts">
import * as albumApi from '@/api/content/album'
import * as albumPhotoApi from '@/api/content/albumPhoto'
import { addDialog } from '@/components/dialog'
import albumForm from "./albumForm.vue"
import albumPhotoForm from "./albumPhotoForm.vue"
import type { AlbumFormProps, AlbumPhotoFormProps } from './utils/types'
import dayjs from "dayjs"
import { apiClient, useApi } from '@/api'

const message = useMessage()

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    pagination.page = page
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  },
})

const albumList = ref<any[]>([])
const albumPhotoList = ref<any[]>([])

const onSearchAlbum=async()=>{
  const res=await useApi(apiClient.api.v1AlbumListList)
  if (res.code === 200) {
    albumList.value = res.data || []
  }
}

const onSearchAlbumPhoto=async()=>{
  const res=await albumPhotoApi.getAlbumPhotoPage({
    page:pagination.page,
    pageSize:pagination.pageSize,
  })
  if (res.code === 200) {
    albumPhotoList.value = res.data.records || []
  }
}

const openAlbumDialog = (title='新增',row?:any)=>{
  const formRef = ref()
  addDialog<AlbumFormProps>({
    title:`${title}相册`,
    props:{
      formInline:{
        name:row?.name ??"",
        description:row?.description ?? "",
        sort:row?.sort ?? 0,
      },
    },
    contentRenderer:({options})=>h(albumForm,{ref:formRef,formInline:options.props!.formInline}),
    beforeSure:async(done)=>{
      try{
        const data=await formRef.value.getData()
        const chores = ()=>{
          message.success('操作成功')
          done()
        }
        if(title=='新增'){
          albumApi.createAlbum(data).then(()=>{
            chores()
            onSearchAlbum()
          })
        }else{
          albumApi.updateAlbum(row?.id,data).then(()=>{
            chores()
            onSearchAlbum()
          })
        }
      }catch{
        done()
      }
    }
  })
}

const openAlbumPhotoDialog=(title='新增',row?:any)=>{
  const formRef = ref()
  addDialog<AlbumPhotoFormProps>({
    title:`${title}新增相片`,
    props:{
      formInline:{
        name:row?.name ??"",
        image_url:row?.image_url ??"",
        description:row?.description ?? "",
        album_id:row?.album_id ?? 0,
      }
    },
    contentRenderer:({options})=>h(albumPhotoForm,{ref:formRef,formInline:options.props!.formInline}),
    beforeSure:async(done)=>{
      try{
        const data=await formRef.value.getData()
        const chores = ()=>{
          message.success('操作成功')
          done()
        }
        if(title=='新增'){
          albumPhotoApi.createAlbumPhoto(data).then(()=>{
            chores()
            onSearchAlbumPhoto()
          })
        }else{
          albumPhotoApi.updateAlbumPhoto(row?.id,data).then(()=>{
            chores()
            onSearchAlbumPhoto()
          })
        }
      }catch{
        done()
      }
    }
  })
}

onMounted(()=>{
  onSearchAlbum()
  onSearchAlbumPhoto()
})
</script>
<template>
  <div class="p-6">
    <div class="album-card">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-6">相册管理</h2>

      <n-grid x-gap="6" cols="3">
        <n-gi span="1">
          <n-card title="相册列表">
            <template #header-extra>
              <n-button type="primary" style="margin-right: 5px;" @click="openAlbumDialog('新增')"> <i class="fas fa-plus mr-2"></i>新增相册 </n-button>
              <n-button @click="onSearchAlbum()">刷新</n-button>
            </template>
            <ul class="space-y-2">
              <!-- 示例相册 -->
              <li
                v-for="(item,index) in albumList"
                :key="index"
                class="flex justify-between items-center p-2 bg-gray-50 dark:bg-gray-700 rounded-md cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600 album-item"
                :data-album-id="item.id"
              >
                <span class="text-gray-800 dark:text-gray-200">{{ item.name }}</span>
                <span class="text-sm text-gray-500 dark:text-gray-400">0 张</span>
              </li>

            </ul>
          </n-card>
        </n-gi>
        <n-gi span="2">
          <n-card title="相片列表">
            <template #header-extra>
              <n-button type="primary" @click="openAlbumPhotoDialog('新增')" style="margin-right: 5px;"> <i class="fas fa-upload mr-2"></i>新增 </n-button>
              <n-button @click="onSearchAlbumPhoto()">刷新</n-button>
            </template>
            <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
              <!-- 示例照片卡片 -->
              <div v-for="(item,index) in albumPhotoList" :key="index" class="bg-gray-50 dark:bg-gray-700 rounded-lg shadow-sm overflow-hidden">
                <n-image
                  :src="item.image_url"
                  :alt="item.name"
                  width="100%"
                  class="h-32 w-full"
                />
                <div class="p-3">
                  <p class="text-sm font-medium text-gray-800 dark:text-gray-200 truncate">
                    {{ item.name }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">{{ dayjs(item.updated_at).format('YYYY-MM-DD HH:mm:ss') }}</p>
                </div>
              </div>

            </div>
            <div class="mt-4 flex justify-end">
              <n-pagination
                :page="pagination.page"
                :page-size="pagination.pageSize"
                :show-size-picker="pagination.showSizePicker"
                :page-sizes="pagination.pageSizes"
                @on-change="pagination.onChange"
                @on-update-page-size="pagination.onUpdatePageSize"
              />
            </div>
          </n-card>
        </n-gi>
      </n-grid>
    </div>


  </div>
</template>
<style scoped>
.album-card{
  max-width: 1600px;
  margin: 0 auto;
}
</style>
