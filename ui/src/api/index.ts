import { Api, HttpClient, type RequestParams } from '@hoshikuzu/api-client'
import type { AxiosResponse } from 'axios'
import { getToken } from "@/utils/auth"

const httpClient = new HttpClient({
  baseURL: '',
})

httpClient.instance.interceptors.request.use((config) => {
  config.headers['Authorization'] = `Bearer ${getToken()?.accessToken}`
  return config
})


const apiClient = new Api(httpClient)

type RequestFn<P extends any[], T=any> = (...args:[...P,RequestParams|undefined]) => Promise<AxiosResponse<T>>

/**
 * 调用 API 函数
 * @param requestFn API 函数
 * @param req 请求参数
 * @param params 请求参数
 * @returns API 响应数据
 */
async function useApi<P extends any[], T = any>(
  requestFn: RequestFn<P, T>,
  ...args:P
): Promise<T> {
  const res = await requestFn(...args,undefined)
  return res.data
}

export { apiClient, useApi }
