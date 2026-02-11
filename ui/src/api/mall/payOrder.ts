import { http } from "@/utils/http";
import { BASE_URL, type ApiResponse, type TableResponse } from "@/api/utils";

export const getPayOrderPage = (params?: any) => {
  return http.request<TableResponse<any>>('get', `${BASE_URL}/v1/pay-order/page`, { params })
}

export const getTodayStats = () => {
  return http.request<ApiResponse<any>>('get', `${BASE_URL}/v1/pay-order/today-stats`)
}

export const submitPayOrder = (data: any) => {
  return http.request<ApiResponse<any>>('post', `${BASE_URL}/v1/pay-order/submit`, { data })
}

export const updatePayOrder = (id: number, data: any) => {
  return http.request<ApiResponse<any>>('put', `${BASE_URL}/v1/pay-order/update/${id}`, { data })
}

export const queryPayOrder = (id: number) => {
  return http.request<ApiResponse<any>>('get', `${BASE_URL}/v1/pay-order/query/${id}`)
}

export const deletePayOrder = (id: number) => {
  return http.request<ApiResponse<any>>('delete', `${BASE_URL}/v1/pay-order/delete/${id}`)
}
