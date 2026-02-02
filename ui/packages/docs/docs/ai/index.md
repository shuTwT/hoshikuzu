# AI 模块

AI 模块是 Hoshikuzu 的核心功能之一，提供强大的 AI 能力，包括 AI 摘要、AI 播客、AI 写作、AI 知识库等功能。

## 功能特性

- **AI 摘要**：自动生成文章摘要，快速了解文章内容
- **AI 播客**：将文本内容转换为音频播客
- **AI 写作**：智能辅助写作，提供写作建议和内容生成
- **AI 知识库**：构建和管理知识库，智能检索和问答

## 快速开始

### 1. 配置 AI 参数

在 `config.yaml` 中配置 AI 参数：

```yaml
ai:
  api_key: your_ai_api_key
  model: gpt-4
  max_tokens: 2000
```

### 2. 使用 AI 功能

```javascript
import { aiApi } from '@/api/ai'

const generateSummary = async () => {
  const data = {
    content: '文章内容...'
  }
  
  const result = await aiApi.summary(data)
  console.log(result)
}
```

## AI 功能

### AI 摘要

自动生成文章摘要，帮助用户快速了解文章内容。

### AI 播客

将文本内容转换为音频播客，提供更好的阅读体验。

### AI 写作

智能辅助写作，提供写作建议和内容生成。

### AI 知识库

构建和管理知识库，智能检索和问答。

## API 接口

### AI 摘要

```http
POST /api/v1/ai/summary
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文章内容 |
| max_length | int | 否 | 摘要最大长度 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "summary": "文章摘要..."
  }
}
```

### AI 播客

```http
POST /api/v1/ai/podcast
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文章内容 |
| voice | string | 否 | 语音类型 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "audio_url": "https://your-domain.com/audio/xxx.mp3"
  }
}
```

### AI 写作

```http
POST /api/v1/ai/writing
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 是 | 写作提示 |
| type | string | 否 | 写作类型 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content": "生成的内容..."
  }
}
```

### AI 知识库

```http
POST /api/v1/ai/knowledge
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
    "answer": "回答..."
  }
}
```

## 使用示例

### 生成文章摘要

```javascript
import { aiApi } from '@/api/ai'

const generateSummary = async () => {
  const data = {
    content: 'Hoshikuzu 是一款高性能、高稳定性、易扩展的内容管理系统...',
    max_length: 100
  }
  
  const result = await aiApi.summary(data)
  console.log(result.summary)
}
```

### 生成播客

```javascript
import { aiApi } from '@/api/ai'

const generatePodcast = async () => {
  const data = {
    content: 'Hoshikuzu 是一款高性能、高稳定性、易扩展的内容管理系统...',
    voice: 'male'
  }
  
  const result = await aiApi.podcast(data)
  console.log(result.audio_url)
}
```

### AI 写作

```javascript
import { aiApi } from '@/api/ai'

const aiWriting = async () => {
  const data = {
    prompt: '写一篇关于 Hoshikuzu 的介绍文章',
    type: 'article'
  }
  
  const result = await aiApi.writing(data)
  console.log(result.content)
}
```

### 知识库问答

```javascript
import { aiApi } from '@/api/ai'

const knowledgeQuery = async () => {
  const data = {
    question: 'Hoshikuzu 支持哪些数据库？'
  }
  
  const result = await aiApi.knowledge(data)
  console.log(result.answer)
}
```

## 最佳实践

### AI 摘要

1. **内容长度**：建议文章内容在 500 字以上
2. **摘要长度**：根据需求设置合适的摘要长度
3. **内容质量**：高质量的内容能生成更好的摘要

### AI 播客

1. **语音选择**：根据内容类型选择合适的语音
2. **音频质量**：使用高质量的音频格式
3. **文件大小**：注意音频文件大小，避免过大

### AI 写作

1. **提示词**：使用清晰明确的提示词
2. **内容类型**：指定内容类型以获得更好的结果
3. **人工审核**：AI 生成的内容需要人工审核

### AI 知识库

1. **知识库构建**：定期更新知识库内容
2. **问题质量**：使用清晰明确的问题
3. **上下文信息**：提供足够的上下文信息

## 相关文档

- [AI 摘要](summary.md) - AI 摘要详细文档
- [AI 播客](podcast.md) - AI 播客详细文档
- [AI 写作](writing.md) - AI 写作详细文档
- [AI 知识库](knowledge.md) - AI 知识库详细文档
- [内容管理模块](../content/article.md) - 内容管理模块文档
