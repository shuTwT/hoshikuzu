import axios, { AxiosError, type AxiosRequestConfig, type Method } from 'axios'
import { formatToken, getToken } from './auth'
import { useUserStoreHook } from '@/stores/modules/user'


export type RequestMethods = Extract<
  Method,
  'get' | 'post' | 'put' | 'delete' | 'patch' | 'option' | 'head'
>

class HttpService {

  constructor(){
    this.httpInterceptorsRequest()
    this.httpInterceptorsResponse()
  }
  /** `token`过期后，暂存待执行的请求 */
  private static requests: any[] = []

  /** 防止重复刷新`token` */
  private static isRefreshing = false

  /** 保存当前`Axios`实例对象 */
  private static axiosInstance = axios.create({
    timeout: 10000,
    headers: {
      Accept: 'application/json, text/plain, */*',
      'Content-Type': 'application/json',
      'X-Requested-With': 'XMLHttpRequest',
    },
  })

  private static retryOriginRequest(config: AxiosRequestConfig) {
    return new Promise((resolve) => {
      HttpService.requests.push((token: string) => {
        if (!config.headers) {
          config.headers = {}
        }
        config.headers['Authorization'] = formatToken(token)
        resolve(config)
      })
    })
  }

  /** 请求拦截 */
  private httpInterceptorsRequest(): void {
    HttpService.axiosInstance.interceptors.request.use(async (config):Promise<any> => {
      /** 白名单不添加token */
      const whiteList = ['/refresh-token', '/login']
      
      // 如果数据是 FormData，删除 Content-Type 让浏览器自动设置
      if (config.data instanceof FormData && config.headers) {
        delete config.headers['Content-Type']
      }
      
      return whiteList.some((url) => config.url?.endsWith(url))
        ? config
        : new Promise((resolve) => {
            const data = getToken()
            if (data) {
              const now = new Date().getTime()
              const expires = Number(data.expires) - now <= 0
              if (expires) {
                if (!HttpService.isRefreshing) {
                  HttpService.isRefreshing = true
                  // token过期刷新
                  useUserStoreHook()
                    .handleRefreshToken({ refreshToken: data.refreshToken })
                    .catch(() => {
                      useUserStoreHook().logOut()
                    })
                    .finally(() => {
                      HttpService.isRefreshing = false
                    })
                }
                resolve(HttpService.retryOriginRequest(config))
              } else {
                config.headers['Authorization'] = formatToken(data.accessToken)
                resolve(config)
              }
            }else{
              resolve(config)
            }
          })
    },
    error=>{
      return Promise.reject(error)
    }
  )
  }

  private  httpInterceptorsResponse():void{
    const instance = HttpService.axiosInstance;
    instance.interceptors.response.use(
      async (response) =>{
        const code = response.data.code || 200
        let msg = response.data.msg || response.statusText || "未知错误"

        if(code ===401){
          msg = "登录过期，请重新登录"
          console.error(msg)
          if(window.$dialog) window.$dialog.warning({
            title: "登录过期",
            content: msg,
            positiveText:"确定",
            onPositiveClick: () => {
              useUserStoreHook().logOut()
            }
          })
        }else if(code!==200){
          if(window.$message) window.$message.error(msg)
          return Promise.reject(new Error(msg))
        }
        return response.data
      },
      (error:AxiosError)=>{
        const $error = error
        console.log(error)
        if(window.$message) window.$message.error(error.response?.data+"")
        return Promise.reject($error)
      }
    )
  }

  public request<T=any>(method: RequestMethods, url: string, params?: AxiosRequestConfig): Promise<T> {
    const config = {
      method,
      url,
      ...params,
    }
    return new Promise((resolve, reject) => {
      HttpService.axiosInstance
        .request<T>(config)
        .then((response:any) => {
          resolve(response)
        })
        .catch((error) => {
          reject(error)
        })
    })
  }
  public post<T=any, P=any>(url: string, params?: AxiosRequestConfig<P>): Promise<T> {
    return this.request<T>('post', url, params)
  }
  public get<T=any, P=any>(url: string, params?: AxiosRequestConfig<P>): Promise<T> {
    return this.request<T>('get', url, params)
  }
}

export const http = new HttpService()
