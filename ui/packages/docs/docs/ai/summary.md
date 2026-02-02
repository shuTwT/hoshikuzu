# AI 摘要

AI 摘要功能可以自动生成文章摘要，帮助用户快速了解文章内容。

## 功能特性

- 自动生成摘要：根据文章内容自动生成摘要
- 可控长度：支持设置摘要的最大长度
- 多语言支持：支持多种语言的摘要生成
- 高质量摘要：使用先进的 AI 模型生成高质量摘要

## API 接口

### 生成摘要

```http
POST /api/v1/ai/summary
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文章内容 |
| max_length | int | 否 | 摘要最大长度，默认 200 |
| language | string | 否 | 语言，默认 zh |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "summary": "Hoshikuzu 是一款高性能、高稳定性、易扩展的内容管理系统，采用现代化的技术栈构建。",
    "original_length": 1000,
    "summary_length": 50
  }
}
```

## 使用示例

### 基本使用

```javascript
import { aiApi } from '@/api/ai'

const generateSummary = async () => {
  const data = {
    content: 'Hoshikuzu 是一款高性能、高稳定性、易扩展的内容管理系统，采用现代化的技术栈构建，为开发者提供强大的内容管理和业务扩展能力。Hoshikuzu 基于 GoFiber 框架构建，提供出色的并发处理能力和响应速度。'
  }
  
  const result = await aiApi.summary(data)
  console.log(result.summary)
}
```

### 设置摘要长度

```javascript
import { aiApi } from '@/api/ai'

const generateSummary = async () => {
  const data = {
    content: '文章内容...',
    max_length: 100
  }
  
  const result = await aiApi.summary(data)
  console.log(result.summary)
}
```

### 指定语言

```javascript
import { aiApi } from '@/api/ai'

const generateSummary = async () => {
  const data = {
    content: 'Article content...',
    language: 'en'
  }
  
  const result = await aiApi.summary(data)
  console.log(result.summary)
}
```

## 最佳实践

### 内容长度

- 建议文章内容在 500 字以上
- 过短的内容可能无法生成有意义的摘要
- 过长的内容可以分段处理

### 摘要长度

- 根据使用场景设置合适的摘要长度
- 一般建议摘要长度在 50-200 字之间
- 摘要长度应包含文章的核心信息

### 内容质量

- 高质量的内容能生成更好的摘要
- 避免使用过多的重复内容
- 保持内容的逻辑性和连贯性

## 应用场景

### 文章列表

在文章列表中显示摘要，帮助用户快速了解文章内容。

### 搜索结果

在搜索结果中显示摘要，提高搜索体验。

### 内容推荐

在内容推荐中显示摘要，吸引用户点击。

### 社交分享

在社交分享时显示摘要，提高分享效果。

## 相关文档

- [AI 模块](index.md) - AI 模块总览
- [内容管理模块](../content/article.md) - 内容管理模块文档
