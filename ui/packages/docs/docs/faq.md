# 常见问题

本文档汇总了使用 Hoshikuzu 过程中常见的问题和解决方案。

## 安装和配置

### Q: 如何安装 Hoshikuzu？

A: 请参考 [快速开始](quickstart.md) 文档，按照步骤进行安装。

### Q: 数据库连接失败怎么办？

A: 请检查以下内容：
1. 数据库服务是否启动
2. 数据库配置是否正确（host、port、username、password）
3. 数据库是否已创建
4. 用户是否有访问数据库的权限

### Q: 端口被占用怎么办？

A: 可以修改配置文件中的端口号：
- 后端端口：在 `config.yaml` 中修改 `server.port`
- 前端端口：在 `ui/vite.config.ts` 中修改 `server.port`

### Q: 如何配置支付功能？

A: 请参考 [支付模块](payment/index.md) 文档，按照步骤配置支付宝、微信或易支付。

## API 使用

### Q: 如何获取 API Token？

A: 使用登录接口获取 Token：

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "your_username",
  "password": "your_password"
}
```

### Q: Token 过期了怎么办？

A: Token 过期后需要重新登录获取新的 Token。建议在请求拦截器中处理 Token 过期的情况。

### Q: 如何处理 API 错误？

A: API 错误会返回标准的错误响应格式：

```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

请根据错误码和错误信息进行相应的处理。

### Q: 如何上传文件？

A: 使用文件上传接口：

```http
POST /api/v1/upload
Content-Type: multipart/form-data

file: [文件]
```

## 功能模块

### Q: 如何创建文章？

A: 使用文章创建接口：

```http
POST /api/v1/content/article/create
Content-Type: application/json

{
  "title": "文章标题",
  "content": "文章内容",
  "category_id": 1,
  "status": "published"
}
```

### Q: 如何创建支付订单？

A: 使用支付创建接口：

```http
POST /api/v1/payment/create
Content-Type: application/json

{
  "order_id": "ORDER_20240101",
  "amount": 100,
  "subject": "商品名称",
  "type": "alipay"
}
```

### Q: 如何使用 AI 功能？

A: 请参考 [AI 模块](ai/index.md) 文档，了解如何使用 AI 摘要、AI 播客、AI 写作、AI 知识库等功能。

### Q: 如何创建商品？

A: 使用商品创建接口：

```http
POST /api/v1/mall/product/create
Content-Type: application/json

{
  "name": "商品名称",
  "price": 100,
  "stock": 100,
  "description": "商品描述"
}
```

### Q: 如何发布说说？

A: 使用说说创建接口：

```http
POST /api/v1/social/moment/create
Content-Type: application/json

{
  "content": "今天天气真好！",
  "images": ["image1.jpg", "image2.jpg"]
}
```

## 性能优化

### Q: 如何提高 API 响应速度？

A: 可以从以下几个方面优化：
1. 使用缓存（Redis）
2. 优化数据库查询
3. 使用 CDN 加速静态资源
4. 启用 Gzip 压缩

### Q: 如何处理大量数据？

A: 可以使用分页接口获取数据，避免一次性获取大量数据：

```http
GET /api/v1/content/article/page?page=1&page_size=10
```

### Q: 如何优化图片加载？

A: 可以使用以下方法：
1. 使用图片压缩
2. 使用懒加载
3. 使用 CDN 加速
4. 使用 WebP 格式

## 安全问题

### Q: 如何保护 API 安全？

A: 建议采取以下措施：
1. 使用 HTTPS
2. 使用 JWT 认证
3. 验证请求参数
4. 设置请求频率限制
5. 定期更新依赖

### Q: 如何防止 SQL 注入？

A: Hoshikuzu 使用 Ent ORM，自动处理 SQL 注入问题。请确保使用 ORM 提供的方法进行数据库操作。

### Q: 如何防止 XSS 攻击？

A: 建议采取以下措施：
1. 对用户输入进行过滤和转义
2. 使用 CSP（Content Security Policy）
3. 设置 HttpOnly Cookie
4. 验证和过滤用户上传的文件

## 开发问题

### Q: 如何调试 API？

A: 可以使用以下工具：
1. Postman
2. cURL
3. 浏览器开发者工具

### Q: 如何查看日志？

A: 日志文件位于项目根目录的 `logs` 目录下。可以根据配置文件中的日志配置查看日志。

### Q: 如何生成 Ent 代码？

A: 运行以下命令生成 Ent 代码：

```bash
go generate ./ent
```

### Q: 如何运行测试？

A: 运行以下命令运行测试：

```bash
go test ./...
```

## 其他问题

### Q: 如何获取技术支持？

A: 您可以通过以下方式获取技术支持：
1. 查看文档
2. 在 GitHub 上提交 Issue
3. 联系技术支持团队

### Q: 如何参与项目开发？

A: 欢迎参与项目开发！请参考以下步骤：
1. Fork 项目
2. 创建分支
3. 提交代码
4. 发起 Pull Request

### Q: 如何报告 Bug？

A: 请在 GitHub 上提交 Issue，详细描述 Bug 的复现步骤和预期行为。

## 相关文档

- [快速开始](quickstart.md) - 快速开始指南
- [配置](config.md) - 系统配置说明
- [开发指南](dev/integration.md) - 前台项目接入指南
