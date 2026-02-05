import {http} from '@/utils/http'
import { BASE_URL, type ApiResponse, type TableResponse } from '@/api/utils'

export interface Theme {
  id: number
  created_at: string
  updated_at: string
  type: string
  name: string
  display_name: string
  description?: string
  author_name?: string
  author_email?: string
  logo?: string
  homepage?: string
  repo?: string
  issue?: string
  setting_name?: string
  config_map_name?: string
  version: string
  require?: string
  license?: string
  path?: string
  external_url?: string
  enabled: boolean
}

export interface ThemePageParams {
  page: number
  page_size: number
}

export interface CreateThemeParams {
  type: 'internal' | 'external'
  file_path?: string
  name?: string
  display_name?: string
  description?: string
  external_url?: string
  version?: string
}

export const getThemePage = (params: ThemePageParams) => {
  return http.request<TableResponse<Theme>>('get', `${BASE_URL}/v1/theme/page`, { params })
}

export const queryTheme = (id: number) => {
  return http.request<ApiResponse<Theme>>('get', `${BASE_URL}/v1/theme/query/${id}`)
}

export const uploadThemeFile = (file: FormData) => {
  return http.request<ApiResponse<{ file_path: string }>>('post', `${BASE_URL}/v1/theme/upload`, { 
    data: file,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export const createTheme = (params: CreateThemeParams) => {
  return http.request<ApiResponse<Theme>>('post', `${BASE_URL}/v1/theme/create`, { data: params })
}

export const deleteTheme = (id: number) => {
  return http.request<ApiResponse<null>>('delete', `${BASE_URL}/v1/theme/delete/${id}`)
}

export const enableTheme = (id: number) => {
  return http.request<ApiResponse<Theme>>('post', `${BASE_URL}/v1/theme/${id}/enable`)
}

export const disableTheme = (id: number) => {
  return http.request<ApiResponse<Theme>>('post', `${BASE_URL}/v1/theme/${id}/disable`)
}
