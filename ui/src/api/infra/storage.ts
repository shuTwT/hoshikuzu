import {http} from '@/utils/http'
import { BASE_URL, type ApiResponse } from '@/api/utils'

export interface StorageStrategy {
  id: number
  name: string
  type: 'local' | 's3'
  domain: string
  access_key?: string
  secret_key?: string
  bucket?: string
  region?: string
  master: boolean
  created_at: string
  updated_at: string
}

export interface StorageStrategyPageParams {
  page: number
  page_size: number
  name?: string
  type?: string
  master?: boolean
}

export interface CreateStorageStrategyParams {
  name: string
  type: StorageStrategy['type']
  domain: string
  access_key?: string
  secret_key?: string
  bucket?: string
  region?: string
  master: boolean
}

export interface UpdateStorageStrategyParams extends Partial<CreateStorageStrategyParams> {
  id: number
}

// 获取存储策略列表
export const getStorageStrategyList = (params?: StorageStrategyPageParams) => {
  return http.request<ApiResponse<StorageStrategy[]>>('get',`${BASE_URL}/v1/storage-strategy/page`, { params })
}

// 获取所有存储策略
export const getStorageStrategyListAll = () => {
  return http.request<ApiResponse<StorageStrategy[]>>('get',`${BASE_URL}/v1/storage-strategy/list`)
}

// 创建存储策略
export const createStorageStrategy = (data: CreateStorageStrategyParams) => {
  return http.request('post',`${BASE_URL}/v1/storage-strategy/create`, {data})
}

// 更新存储策略
export const updateStorageStrategy = (data: UpdateStorageStrategyParams) => {
  return http.request('put',`${BASE_URL}/v1/storage-strategy/update/${data.id}`, {data})
}

// 删除存储策略
export const deleteStorageStrategy = (id: number) => {
  return http.request('delete',`${BASE_URL}/v1/storage-strategy/delete/${id}`)
}

// 设置默认存储策略
export const setDefaultStorageStrategy = (id: number) => {
  return http.request('put',`${BASE_URL}/v1/storage-strategy/default/${id}`)
}
