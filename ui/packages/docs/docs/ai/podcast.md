# AI 播客

AI 播客功能可以将文本内容转换为音频播客，提供更好的阅读体验。

## 功能特性

- 文本转语音：将文本内容转换为音频
- 多种语音：支持多种语音类型
- 高质量音频：生成高质量的音频文件
- 多语言支持：支持多种语言的语音合成

## API 接口

### 生成播客

```http
POST /api/v1/ai/podcast
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文章内容 |
| voice | string | 否 | 语音类型 (male/female)，默认 male |
| speed | float | 否 | 语速，默认 1.0 |
| language | string | 否 | 语言，默认 zh |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "audio_url": "https://your-domain.com/audio/xxx.mp3",
    "duration": 120,
    "file_size": 1024000
  }
}
```

## 使用示例

### 基本使用

```javascript
import { aiApi } from '@/api/ai'

const generatePodcast = async () => {
  const data = {
    content: 'Hoshikuzu 是一款高性能、高稳定性、易扩展的内容管理系统。'
  }
  
  const result = await aiApi.podcast(data)
  console.log(result.audio_url)
}
```

### 选择语音类型

```javascript
import { aiApi } from '@/api/ai'

const generatePodcast = async () => {
  const data = {
    content: '文章内容...',
    voice: 'female'
  }
  
  const result = await aiApi.podcast(data)
  console.log(result.audio_url)
}
```

### 调整语速

```javascript
import { aiApi } from '@/api/ai'

const generatePodcast = async () => {
  const data = {
    content: '文章内容...',
    speed: 1.5
  }
  
  const result = await aiApi.podcast(data)
  console.log(result.audio_url)
}
```

### 指定语言

```javascript
import { aiApi } from '@/api/ai'

const generatePodcast = async () => {
  const data = {
    content: 'Article content...',
    language: 'en'
  }
  
  const result = await aiApi.podcast(data)
  console.log(result.audio_url)
}
```

## 最佳实践

### 语音选择

- 根据内容类型选择合适的语音
- 正式内容建议使用男声
- 轻松内容建议使用女声

### 语速调整

- 一般内容建议使用正常语速（1.0）
- 快速阅读可以调整语速为 1.2-1.5
- 慢速学习可以调整语速为 0.8-1.0

### 音频质量

- 使用高质量的音频格式（MP3）
- 注意音频文件大小，避免过大
- 定期清理旧的音频文件

### 内容准备

- 确保文本内容格式正确
- 避免使用过多的特殊符号
- 适当添加标点符号以改善语音效果

## 应用场景

### 文章朗读

为文章提供语音朗读功能，方便用户收听。

### 有声书

将文本内容转换为有声书，提供更好的阅读体验。

### 教程视频

为教程视频添加语音解说，提高学习效果。

### 无障碍访问

为视障用户提供语音服务，提高无障碍访问性。

## 相关文档

- [AI 模块](index.md) - AI 模块总览
- [内容管理模块](../content/article.md) - 内容管理模块文档
