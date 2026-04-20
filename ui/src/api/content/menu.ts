import { http } from "@/utils/http";
import { BASE_URL, type ApiResponse, type TableResponse } from "@/api/utils";

export const getMenuList = () => {
  return http.request<ApiResponse<any>>('get', `${BASE_URL}/v1/menu/list`)
}

export const getMenuPage = (params?: any) => {
  return http.request<TableResponse<any>>('get', `${BASE_URL}/v1/menu/page`, { params })
}

export const queryMenu = (id: number) => {
  return http.request<ApiResponse<any>>('get', `${BASE_URL}/v1/menu/query/${id}`)
}

export const createMenu = (data: any) => {
  return http.request<ApiResponse<any>>('post', `${BASE_URL}/v1/menu/create`, { data })
}

export const updateMenu = (id: number, data: any) => {
  return http.request<ApiResponse<any>>('put', `${BASE_URL}/v1/menu/update/${id}`, { data })
}

export const deleteMenu = (id: number) => {
  return http.request<ApiResponse<any>>('delete', `${BASE_URL}/v1/menu/delete/${id}`)
}
