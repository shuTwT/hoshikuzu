# MailTemplate Handler

邮件模板模块的 HTTP 处理器，提供邮件模板的管理接口。

## 功能

- 分页查询邮件模板（支持按名称、编码模糊搜索）
- 查询邮件模板详情
- 创建邮件模板
- 更新邮件模板
- 删除邮件模板

## 接口

所有接口位于 `/api/v1/mail-template` 下，需要管理员登录态。

### GET /api/v1/mail-template/page

分页查询邮件模板，查询参数：`page`、`page_size`、`name`（模糊）、`code`（模糊）。
返回 `model.PageResult[ent.MailTemplate]` 结构（`total` + `records`）。

### GET /api/v1/mail-template/query/{id}

查询指定 ID 的邮件模板详情，不存在时返回 404。

### POST /api/v1/mail-template/create

创建邮件模板，请求体：

```json
{
  "name": "注册欢迎邮件",
  "code": "register_welcome",
  "subject": "欢迎加入星屑，{nickname}",
  "content": "<h1>你好 {nickname}</h1><p>感谢你的注册。</p>",
  "content_type": "html",
  "enabled": true,
  "remark": "注册场景使用"
}
```

`name`、`code`、`subject`、`content` 必填；`code` 全局唯一，重复时返回 400 及 `邮件模板编码已存在: {code}`。

### PUT /api/v1/mail-template/update/{id}

更新指定 ID 的邮件模板，请求体字段均可选（仅更新传入字段），
`code` 与其他模板冲突时返回 400，模板不存在时返回 404。

### DELETE /api/v1/mail-template/delete/{id}

删除指定 ID 的邮件模板，不存在时返回 404。
