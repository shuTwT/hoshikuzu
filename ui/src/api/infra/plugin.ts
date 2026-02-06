import {http} from '@/utils/http'
import { BASE_URL, type ApiResponse, type TableResponse } from '@/api/utils'

export interface Plugin {
  id: number
  created_at: string
  updated_at: string
  key: string
  name: string
  version: string
  description: string
  bin_path: string
  protocol_version: string
  magic_cookie_key: string
  magic_cookie_value: string
  dependencies: string[]
  config: string
  enabled: boolean
  auto_start: boolean
  status: 'stopped' | 'running' | 'error' | 'loading'
  last_error: string
  last_started_at?: string
  last_stopped_at?: string
}

export interface PluginPageParams {
  page: number
  page_size: number
  name?: string
  key?: string
  status?: string
  enabled?: boolean
  auto_start?: boolean
}

export const getPluginPage = (params: PluginPageParams) => {
  return http.request<TableResponse<Plugin>>('get', `${BASE_URL}/v1/plugin/page`, { params })
}

export const queryPlugin = (id: number) => {
  return http.request<ApiResponse<Plugin>>('get', `${BASE_URL}/v1/plugin/query/${id}`)
}

export const createPlugin = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return http.request<ApiResponse<Plugin>>('post', `${BASE_URL}/v1/plugin/create`, { 
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export const deletePlugin = (id: number) => {
  return http.request<ApiResponse<null>>('delete', `${BASE_URL}/v1/plugin/delete/${id}`)
}

export const startPlugin = (id: number) => {
  return http.request<ApiResponse<Plugin>>('post', `${BASE_URL}/v1/plugin/${id}/start`)
}

export const stopPlugin = (id: number) => {
  return http.request<ApiResponse<Plugin>>('post', `${BASE_URL}/v1/plugin/${id}/stop`)
}

export const restartPlugin = (id: number) => {
  return http.request<ApiResponse<Plugin>>('post', `${BASE_URL}/v1/plugin/${id}/restart`)
}
