# AI 知识库

AI 知识库功能可以构建和管理知识库，提供智能检索和问答能力。

## 功能特性

- 知识库构建：上传和管理知识库内容
- 智能检索：基于语义的智能检索
- 智能问答：基于知识库的智能问答
- 多格式支持：支持多种文档格式

## API 接口

### 上传知识

```http
POST /api/v1/ai/knowledge/upload
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 知识标题 |
| content | string | 是 | 知识内容 |
| category | string | 否 | 知识分类 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "知识标题",
    "category": "分类名称"
  }
}
```

### 检索知识

```http
POST /api/v1/ai/knowledge/search
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| query | string | 是 | 检索关键词 |
| limit | int | 否 | 返回数量，默认 5 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "results": [
      {
        "id": 1,
        "title": "知识标题",
        "content": "知识内容...",
        "score": 0.95
      }
    ]
  }
}
```

### 智能问答

```http
POST /api/v1/ai/knowledge/qa
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question | string | 是 | 问题 |
| context | string | 否 | 上下文 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "answer": "回答...",
    "sources": [
      {
        "id": 1,
        "title": "知识标题"
      }
    ]
  }
}
```

### 删除知识

```http
DELETE /api/v1/ai/knowledge/delete/:id
```

## 使用示例

### 上传知识

```javascript
import { aiApi } from '@/api/ai'

const uploadKnowledge = async () => {
  const data = {
    title: 'Hoshikuzu 介绍',
    content: 'Hoshikuzu 是一款高性能、高稳定性、易扩展的内容管理系统...',
    category: '产品介绍'
  }
  
  const result = await aiApi.knowledgeUpload(data)
  console.log(result)
}
```

### 检索知识

```javascript
import { aiApi } from '@/api/ai'

const searchKnowledge = async () => {
  const data = {
    query: 'Hoshikuzu 的特性',
    limit: 5
  }
  
  const result = await aiApi.knowledgeSearch(data)
  console.log(result.results)
}
```

### 智能问答

```javascript
import { aiApi } from '@/api/ai'

const askQuestion = async () => {
  const data = {
    question: 'Hoshikuzu 支持哪些数据库？'
  }
  
  const result = await aiApi.knowledgeQa(data)
  console.log(result.answer)
}
```

### 删除知识

```javascript
import { aiApi } from '@/api/ai'

const deleteKnowledge = async () => {
  const result = await aiApi.knowledgeDelete(1)
  console.log(result)
}
```

## 最佳实践

### 知识库构建

- 定期更新知识库内容
- 使用清晰的标题和分类
- 保持内容的准确性和时效性

### 知识检索

- 使用清晰明确的检索关键词
- 根据需要调整返回数量
- 结合上下文信息提高检索准确性

### 智能问答

- 使用清晰明确的问题
- 提供足够的上下文信息
- 根据回答质量调整问题表述

### 知识管理

- 定期清理过时的知识
- 建立知识审核机制
- 保持知识库的结构化

## 应用场景

### 客服系统

为客服系统提供智能问答能力，提高客服效率。

### 文档搜索

为文档系统提供智能检索能力，提高文档查找效率。

### 培训系统

为培训系统提供智能问答能力，提高学习效果。

### 内部知识

构建内部知识库，提高团队协作效率。

## 相关文档

- [AI 模块](index.md) - AI 模块总览
- [内容管理模块](../content/article.md) - 内容管理模块文档
