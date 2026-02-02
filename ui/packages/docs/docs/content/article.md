# 内容管理模块 - 文章

内容管理模块提供完整的文章管理能力，包括文章的创建、编辑、发布、分类、标签等功能。

## 功能特性

- 文章创建与编辑：支持富文本编辑器，支持 Markdown 格式
- 文章分类管理：多级分类支持
- 文章标签管理：灵活的标签系统
- 文章发布与草稿：支持文章草稿保存和定时发布
- 内容审核：完善的内容审核流程
- SEO 优化：内置 SEO 优化功能
- 文章统计：阅读量、点赞数等统计数据
- 评论管理：文章评论功能

## API 接口

### 获取文章列表

```http
GET /api/v1/content/article/page
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 10 |
| category_id | int | 否 | 分类 ID |
| tag_id | int | 否 | 标签 ID |
| status | string | 否 | 状态 (draft/published) |
| keyword | string | 否 | 搜索关键词 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 1,
        "title": "文章标题",
        "content": "文章内容",
        "summary": "文章摘要",
        "category_id": 1,
        "category_name": "分类名称",
        "tags": [
          {
            "id": 1,
            "name": "标签名称"
          }
        ],
        "status": "published",
        "view_count": 100,
        "like_count": 10,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 创建文章

```http
POST /api/v1/content/article/create
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 文章标题 |
| content | string | 是 | 文章内容 |
| summary | string | 否 | 文章摘要 |
| category_id | int | 否 | 分类 ID |
| tag_ids | array | 否 | 标签 ID 列表 |
| status | string | 否 | 状态 (draft/published)，默认 draft |
| seo_title | string | 否 | SEO 标题 |
| seo_keywords | string | 否 | SEO 关键词 |
| seo_description | string | 否 | SEO 描述 |

**请求示例：**

```json
{
  "title": "文章标题",
  "content": "文章内容",
  "summary": "文章摘要",
  "category_id": 1,
  "tag_ids": [1, 2],
  "status": "published"
}
```

### 更新文章

```http
PUT /api/v1/content/article/update/:id
```

**请求参数：** 同创建文章

### 删除文章

```http
DELETE /api/v1/content/article/delete/:id
```

### 查询文章详情

```http
GET /api/v1/content/article/query/:id
```

## 文章分类管理

### 获取分类列表

```http
GET /api/v1/content/category/list
```

### 创建分类

```http
POST /api/v1/content/category/create
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 分类名称 |
| parent_id | int | 否 | 父分类 ID |
| description | string | 否 | 分类描述 |
| sort | int | 否 | 排序 |

## 文章标签管理

### 获取标签列表

```http
GET /api/v1/content/tag/list
```

### 创建标签

```http
POST /api/v1/content/tag/create
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 标签名称 |
| color | string | 否 | 标签颜色 |

## 使用示例

### 创建文章

```javascript
import { articleApi } from '@/api/content/article'

const createArticle = async () => {
  const data = {
    title: '如何使用 Hoshikuzu',
    content: '# Hoshikuzu 使用指南\n\nHoshikuzu 是一款高性能的内容管理系统...',
    summary: 'Hoshikuzu 使用指南',
    category_id: 1,
    tag_ids: [1, 2],
    status: 'published'
  }
  
  const result = await articleApi.create(data)
  console.log(result)
}
```

### 获取文章列表

```javascript
import { articleApi } from '@/api/content/article'

const getArticles = async () => {
  const params = {
    page: 1,
    page_size: 10,
    status: 'published'
  }
  
  const result = await articleApi.page(params)
  console.log(result)
}
```

## 最佳实践

### 文章写作建议

1. **标题优化**：使用简洁明了的标题，包含关键词
2. **内容结构**：使用清晰的段落结构和标题层级
3. **摘要编写**：摘要应该简明扼要地概括文章内容
4. **图片优化**：使用适当的图片格式和尺寸
5. **SEO 优化**：合理设置 SEO 标题、关键词和描述

### 分类管理建议

1. **分类层级**：建议不超过 3 级分类
2. **分类命名**：使用简洁明了的分类名称
3. **分类数量**：避免创建过多的分类

### 标签管理建议

1. **标签命名**：使用简短的标签名称
2. **标签数量**：每篇文章建议使用 3-5 个标签
3. **标签颜色**：使用不同的颜色区分标签类型

## 相关文档

- [AI 摘要](../ai/summary.md) - 自动生成文章摘要
- [AI 写作](../ai/writing.md) - 智能辅助写作
- [开发指南](../dev/integration.md) - 前台项目接入指南
