import {http} from '@/utils/http'
import { BASE_URL, type ApiResponse, type TableResponse } from '@/api/utils'

export interface License {
  id: number
  created_at: string
  updated_at: string
  domain: string
  license_key: string
  customer_name: string
  expire_date: string
  status: number
}

export interface LicensePageParams {
  page: number
  page_size: number
}

export interface LicenseCreateReq {
  domain: string
  customer_name: string
  expire_date: string
}

export interface LicenseUpdateReq {
  domain?: string
  license_key?: string
  customer_name?: string
  expire_date?: string
  status?: number
}

export interface LicenseVerifyReq {
  domain: string
}

export interface LicenseVerifyResp {
  valid: boolean
  customer_name: string
  expire_date: string
  message: string
}

export const getLicensePage = (params: LicensePageParams) => {
  return http.request<TableResponse<License>>('get', `${BASE_URL}/v1/license/page`, { params })
}

export const queryLicense = (id: number) => {
  return http.request<ApiResponse<License>>('get', `${BASE_URL}/v1/license/query/${id}`)
}

export const createLicense = (data: LicenseCreateReq) => {
  return http.request<ApiResponse<License>>('post', `${BASE_URL}/v1/license/create`, { data })
}

export const updateLicense = (id: number, data: LicenseUpdateReq) => {
  return http.request<ApiResponse<License>>('put', `${BASE_URL}/v1/license/update/${id}`, { data })
}

export const deleteLicense = (id: number) => {
  return http.request<ApiResponse<null>>('delete', `${BASE_URL}/v1/license/delete/${id}`)
}

export const verifyLicense = (data: LicenseVerifyReq) => {
  return http.request<ApiResponse<LicenseVerifyResp>>('post', `${BASE_URL}/v1/license/verify`, { data })
}
