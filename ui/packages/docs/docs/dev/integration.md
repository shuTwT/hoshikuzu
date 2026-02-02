# 开发指南 - 前台项目接入

本指南将帮助您将 Hoshikuzu 的功能接入到前台项目中。

## 概述

Hoshikuzu 提供了完整的 RESTful API，可以轻松接入到各种前台项目中，包括 Web 应用、移动应用等。

## API 基础

### API 地址

- 开发环境：`http://localhost:13000/api`
- 生产环境：`https://your-domain.com/api`

### API 版本

当前 API 版本为 `v1`，所有 API 路径以 `/api/v1` 开头。

### 认证方式

Hoshikuzu 使用 JWT 进行身份认证。

#### 获取 Token

```http
POST /api/v1/auth/login
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "jwt_token",
    "user": {
      "id": 1,
      "username": "用户名"
    }
  }
}
```

#### 使用 Token

在请求头中添加 `Authorization` 字段：

```http
Authorization: Bearer jwt_token
```

## 接入步骤

### 1. 配置 API 地址

在前台项目中配置 API 地址：

```typescript
// config/api.ts
export const API_BASE_URL = 'http://localhost:13000/api'
```

### 2. 创建 API 客户端

创建统一的 API 客户端，处理请求和响应：

```typescript
// utils/request.ts
import axios from 'axios'
import { API_BASE_URL } from '@/config/api'

const request = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default request
```

### 3. 创建 API 接口

根据业务需求创建对应的 API 接口：

```typescript
// api/content/article.ts
import request from '@/utils/request'

export const articleApi = {
  page: (params: any) => {
    return request.get('/v1/content/article/page', { params })
  },
  create: (data: any) => {
    return request.post('/v1/content/article/create', data)
  },
  update: (id: number, data: any) => {
    return request.put(`/v1/content/article/update/${id}`, data)
  },
  delete: (id: number) => {
    return request.delete(`/v1/content/article/delete/${id}`)
  },
  query: (id: number) => {
    return request.get(`/v1/content/article/query/${id}`)
  }
}
```

### 4. 使用 API 接口

在组件中使用 API 接口：

```typescript
// views/article/index.vue
import { articleApi } from '@/api/content/article'
import { ref, onMounted } from 'vue'

const articles = ref([])
const loading = ref(false)

const fetchArticles = async () => {
  loading.value = true
  try {
    const result = await articleApi.page({ page: 1, page_size: 10 })
    articles.value = result.data.items
  } catch (error) {
    console.error('获取文章列表失败', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchArticles()
})
```

## 功能模块接入

### 内容管理模块

#### 获取文章列表

```typescript
import { articleApi } from '@/api/content/article'

const getArticles = async () => {
  const result = await articleApi.page({
    page: 1,
    page_size: 10,
    status: 'published'
  })
  return result.data.items
}
```

#### 创建文章

```typescript
import { articleApi } from '@/api/content/article'

const createArticle = async (data: any) => {
  const result = await articleApi.create(data)
  return result.data
}
```

### 支付模块

#### 创建支付订单

```typescript
import { paymentApi } from '@/api/payment'

const createPayment = async (data: any) => {
  const result = await paymentApi.create(data)
  return result.data
}
```

### AI 模块

#### 生成文章摘要

```typescript
import { aiApi } from '@/api/ai'

const generateSummary = async (content: string) => {
  const result = await aiApi.summary({ content })
  return result.data.summary
}
```

### 商城模块

#### 获取商品列表

```typescript
import { productApi } from '@/api/mall/product'

const getProducts = async () => {
  const result = await productApi.page({
    page: 1,
    page_size: 10,
    status: 'on_sale'
  })
  return result.data.items
}
```

### 社交模块

#### 发布说说

```typescript
import { momentApi } from '@/api/social/moment'

const publishMoment = async (data: any) => {
  const result = await momentApi.create(data)
  return result.data
}
```

## 错误处理

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

### 错误处理示例

```typescript
import { articleApi } from '@/api/content/article'

const fetchArticles = async () => {
  try {
    const result = await articleApi.page({ page: 1, page_size: 10 })
    return result.data.items
  } catch (error: any) {
    if (error.response) {
      switch (error.response.status) {
        case 400:
          console.error('请求参数错误')
          break
        case 401:
          console.error('未授权，请重新登录')
          break
        case 404:
          console.error('资源不存在')
          break
        case 500:
          console.error('服务器错误')
          break
        default:
          console.error('未知错误')
      }
    } else {
      console.error('网络错误')
    }
    throw error
  }
}
```

## 最佳实践

### 请求优化

1. **请求缓存**：对不经常变化的数据进行缓存
2. **请求节流**：对频繁的请求进行节流
3. **请求取消**：在组件卸载时取消未完成的请求

### 错误处理

1. **统一错误处理**：使用拦截器统一处理错误
2. **错误提示**：给用户友好的错误提示
3. **错误日志**：记录错误日志便于排查

### 安全建议

1. **HTTPS**：生产环境使用 HTTPS
2. **Token 存储**：安全存储 Token
3. **数据验证**：对请求数据进行验证

## 相关文档

- [快速开始](../quickstart.md) - 快速开始指南
- [配置](../config.md) - 系统配置说明
- [常见问题](../faq.md) - 常见问题解答
