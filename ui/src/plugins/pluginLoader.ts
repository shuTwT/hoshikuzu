import * as Vue from 'vue'
import * as VueRouter from 'vue-router'
import * as Pinia from 'pinia'
import * as naive from 'naive-ui'
import * as ApiClientPackage from '@hoshikuzu/api-client'
import type { RouteRecordRaw } from 'vue-router'
import { apiClient, useApi } from '@/api'

/**
 * 插件前端加载器：登录后拉取运行中的插件，加载其 bundle.js / style.css，
 * 读取 bundle 导出的 routes 并转换为宿主可注册、可显示在菜单里的路由。
 *
 * 约定（插件侧）：
 * - bundle 为 IIFE，全局变量名等于插件 key，默认导出 { routes }
 * - vue / vue-router / pinia / naive-ui / @hoshikuzu/api-client 设为 external，由宿主全局提供
 */

interface PluginRouteDeclaration {
  /** 保留字段：宿主统一为每个插件生成容器路由，不使用该值 */
  parentName?: string
  route: RouteRecordRaw
}

interface PluginBundle {
  routes?: PluginRouteDeclaration[]
}

interface RunningPlugin {
  id: number
  key: string
  name: string
}

const BUNDLE_FILE = 'bundle.js'
const STYLE_FILE = 'style.css'
const PAGE_SIZE = 100

/** 已加载过 bundle 的插件 key，避免重复执行脚本 */
const loadedKeys = new Set<string>()

/** 暴露给插件 bundle 的宿主全局依赖 */
function exposeHostGlobals() {
  const w = window as any
  w.Vue = Vue
  w.VueRouter = VueRouter
  w.Pinia = Pinia
  w.naive = naive
  w.HoshikuzuApiClient = ApiClientPackage
  w.__HOSHIKUZU__ = {
    /** 已携带登录态、刷新 token 拦截器的宿主 API 客户端 */
    apiClient,
    useApi,
    /** 按插件 key 调用其 capability 方法，内部走 POST /api/v1/plugin/:id/call */
    callPlugin: callPluginByKey,
  }
}

const pluginIds = new Map<string, number>()

async function callPluginByKey<T = unknown>(
  key: string,
  method: string,
  params: Record<string, any> = {},
): Promise<T> {
  const id = pluginIds.get(key)
  if (id === undefined) {
    throw new Error(`插件 '${key}' 未运行`)
  }
  const res = await useApi(apiClient.api.v1PluginCallCreate, id, { method, params })
  return res.data as T
}

async function fetchResource(id: number, file: string): Promise<string | null> {
  try {
    const res = await apiClient.api.v1PluginResourceDetail(id, file, { format: 'text' })
    return typeof res.data === 'string' ? res.data : null
  } catch {
    return null
  }
}

function injectStyle(key: string, css: string) {
  const el = document.createElement('style')
  el.dataset.plugin = key
  el.textContent = css
  document.head.appendChild(el)
}

function runScript(key: string, code: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(new Blob([code], { type: 'text/javascript' }))
    const el = document.createElement('script')
    el.dataset.plugin = key
    el.src = url
    el.onload = () => {
      URL.revokeObjectURL(url)
      resolve()
    }
    el.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new Error(`插件 '${key}' 的 ${BUNDLE_FILE} 执行失败`))
    }
    document.head.appendChild(el)
  })
}

async function loadBundle(plugin: RunningPlugin): Promise<PluginBundle | null> {
  if (!loadedKeys.has(plugin.key)) {
    const code = await fetchResource(plugin.id, BUNDLE_FILE)
    if (!code) return null
    const css = await fetchResource(plugin.id, STYLE_FILE)
    if (css) injectStyle(plugin.key, css)
    await runScript(plugin.key, code)
    loadedKeys.add(plugin.key)
  }
  const exported = (window as any)[plugin.key]
  return (exported?.default ?? exported) as PluginBundle | null
}

/** 把插件声明的路由挂到 /plugin/{key} 容器下，避免与宿主路由冲突 */
function buildPluginRoute(plugin: RunningPlugin, bundle: PluginBundle): RouteRecordRaw | null {
  const declared = bundle.routes ?? []
  if (declared.length === 0) return null

  const prefix = `/plugin/${plugin.key}`
  const children = declared.map(({ route }) => {
    const path = route.path.startsWith('/') ? route.path : `/${route.path}`
    return {
      ...route,
      path: `${prefix}${path}`,
      name: `${plugin.key}:${String(route.name ?? path)}`,
      meta: {
        showLink: true,
        roles: ['admin', 'common'],
        ...route.meta,
        plugin: plugin.key,
      },
    } as RouteRecordRaw
  })

  return {
    path: prefix,
    name: `Plugin:${plugin.key}`,
    redirect: children[0].path,
    meta: {
      title: plugin.name,
      showLink: true,
      roles: ['admin', 'common'],
      rank: 90,
      plugin: plugin.key,
    },
    children,
  }
}

async function listRunningPlugins(): Promise<RunningPlugin[]> {
  const res = await useApi(apiClient.api.v1PluginPageList, {
    page: 1,
    page_size: PAGE_SIZE,
    status: 'running',
  })
  return (res.data?.records ?? []).map((p: any) => ({ id: p.id, key: p.key, name: p.name }))
}

/**
 * 加载所有运行中插件的前端，返回要并入宿主菜单的路由。
 * 单个插件失败只跳过该插件，不影响宿主和其他插件。
 */
export async function loadPluginRoutes(): Promise<RouteRecordRaw[]> {
  exposeHostGlobals()

  let plugins: RunningPlugin[]
  try {
    plugins = await listRunningPlugins()
  } catch (error) {
    console.error('获取运行中的插件失败', error)
    return []
  }

  plugins.forEach((p) => pluginIds.set(p.key, p.id))

  const routes: RouteRecordRaw[] = []
  for (const plugin of plugins) {
    try {
      const bundle = await loadBundle(plugin)
      const route = bundle && buildPluginRoute(plugin, bundle)
      if (route) routes.push(route)
    } catch (error) {
      console.error(`加载插件 '${plugin.key}' 前端失败`, error)
    }
  }
  return routes
}
