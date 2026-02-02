# AI 写作

AI 写作功能提供智能辅助写作，帮助用户快速生成高质量的内容。

## 功能特性

- 内容生成：根据提示词生成内容
- 写作建议：提供写作建议和优化建议
- 多种类型：支持文章、博客、新闻等多种内容类型
- 智能续写：根据已有内容智能续写

## API 接口

### 生成内容

```http
POST /api/v1/ai/writing
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 是 | 写作提示 |
| type | string | 否 | 内容类型 (article/blog/news)，默认 article |
| length | int | 否 | 内容长度，默认 500 |
| language | string | 否 | 语言，默认 zh |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content": "生成的内容...",
    "word_count": 500
  }
}
```

### 写作建议

```http
POST /api/v1/ai/writing/suggestion
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文章内容 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "suggestions": [
      "建议增加更多具体案例",
      "可以添加数据支持",
      "建议优化段落结构"
    ]
  }
}
```

### 智能续写

```http
POST /api/v1/ai/writing/continue
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 已有内容 |
| length | int | 否 | 续写长度，默认 200 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "continued_content": "续写的内容...",
    "word_count": 200
  }
}
```

## 使用示例

### 生成文章

```javascript
import { aiApi } from '@/api/ai'

const generateArticle = async () => {
  const data = {
    prompt: '写一篇关于 Hoshikuzu 的介绍文章',
    type: 'article',
    length: 1000
  }
  
  const result = await aiApi.writing(data)
  console.log(result.content)
}
```

### 生成博客

```javascript
import { aiApi } from '@/api/ai'

const generateBlog = async () => {
  const data = {
    prompt: '写一篇关于技术趋势的博客',
    type: 'blog',
    length: 800
  }
  
  const result = await aiApi.writing(data)
  console.log(result.content)
}
```

### 获取写作建议

```javascript
import { aiApi } from '@/api/ai'

const getSuggestions = async () => {
  const data = {
    content: '文章内容...'
  }
  
  const result = await aiApi.writingSuggestion(data)
  console.log(result.suggestions)
}
```

### 智能续写

```javascript
import { aiApi } from '@/api/ai'

const continueWriting = async () => {
  const data = {
    content: '已有内容...',
    length: 300
  }
  
  const result = await aiApi.writingContinue(data)
  console.log(result.continued_content)
}
```

## 最佳实践

### 提示词编写

- 使用清晰明确的提示词
- 提供足够的上下文信息
- 指定内容类型和风格

### 内容类型

- 根据需求选择合适的内容类型
- 不同类型的内容有不同的风格和结构
- 可以组合使用多种类型

### 内容长度

- 根据使用场景设置合适的内容长度
- 一般文章建议 500-2000 字
- 博客建议 300-1000 字

### 人工审核

- AI 生成的内容需要人工审核
- 检查内容的准确性和合理性
- 根据需要进行修改和优化

## 应用场景

### 内容创作

帮助内容创作者快速生成内容，提高创作效率。

### 营销文案

生成营销文案，提高营销效果。

### 产品描述

生成产品描述，提高产品吸引力。

### 教程编写

编写教程文档，提高学习效果。

## 相关文档

- [AI 模块](index.md) - AI 模块总览
- [内容管理模块](../content/article.md) - 内容管理模块文档
