# 社交模块

社交模块提供完整的社交功能，包括说说、朋友圈、社交登录等功能。

## 功能特性

- **说说**：用户动态发布和互动
- **朋友圈**：好友动态分享
- **社交登录**：支持第三方社交平台登录

## 快速开始

### 1. 发布说说

```javascript
import { momentApi } from '@/api/social/moment'

const publishMoment = async () => {
  const data = {
    content: '今天天气真好！',
    images: ['image1.jpg', 'image2.jpg']
  }
  
  const result = await momentApi.create(data)
  console.log(result)
}
```

### 2. 发布朋友圈

```javascript
import { friendApi } from '@/api/social/friend'

const publishFriend = async () => {
  const data = {
    content: '分享我的生活',
    images: ['image1.jpg']
  }
  
  const result = await friendApi.create(data)
  console.log(result)
}
```

### 3. 社交登录

```javascript
import { loginApi } from '@/api/social/login'

const socialLogin = async () => {
  const data = {
    provider: 'wechat',
    code: 'auth_code'
  }
  
  const result = await loginApi.social(data)
  console.log(result)
}
```

## 模块概览

### 说说

- 发布说说
- 查看说说列表
- 点赞和评论
- 删除说说

### 朋友圈

- 发布朋友圈
- 查看朋友圈动态
- 点赞和评论
- 设置可见范围

### 社交登录

- 微信登录
- QQ 登录
- 微博登录

## API 接口

### 说说接口

- `POST /api/v1/social/moment/create` - 发布说说
- `GET /api/v1/social/moment/page` - 获取说说列表
- `DELETE /api/v1/social/moment/delete/:id` - 删除说说
- `POST /api/v1/social/moment/like/:id` - 点赞说说
- `POST /api/v1/social/moment/comment/:id` - 评论说说

### 朋友圈接口

- `POST /api/v1/social/friend/create` - 发布朋友圈
- `GET /api/v1/social/friend/page` - 获取朋友圈列表
- `DELETE /api/v1/social/friend/delete/:id` - 删除朋友圈
- `POST /api/v1/social/friend/like/:id` - 点赞朋友圈
- `POST /api/v1/social/friend/comment/:id` - 评论朋友圈

### 社交登录接口

- `POST /api/v1/social/login` - 社交登录
- `GET /api/v1/social/login/url` - 获取登录授权链接

## 相关文档

- [说说](moment.md) - 说说详细文档
- [朋友圈](friend.md) - 朋友圈详细文档
- [社交登录](login.md) - 社交登录详细文档
