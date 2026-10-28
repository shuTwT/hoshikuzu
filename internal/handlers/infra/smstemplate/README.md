# SmsTemplate Handler

短信模板模块的 HTTP 处理器，提供短信模板的管理接口。

## 功能

- 分页查询短信模板（支持按名称、编码模糊搜索）
- 查询短信模板详情
- 创建短信模板
- 更新短信模板
- 删除短信模板

## 接口

所有接口位于 `/api/v1/sms-template` 下，需要管理员登录态。

### GET /api/v1/sms-template/page

分页查询短信模板，查询参数：`page`、`page_size`、`name`（模糊）、`code`（模糊）。
返回 `model.PageResult[ent.SmsTemplate]` 结构（`total` + `records`）。

### GET /api/v1/sms-template/query/{id}

查询指定 ID 的短信模板详情，不存在时返回 404。

### POST /api/v1/sms-template/create

创建短信模板，请求体：

```json
{
  "name": "注册验证码",
  "code": "register_verify_code",
  "sign_name": "星屑",
  "provider_template_id": "SMS_123456",
  "content": "您的验证码为{code}，{minutes}分钟内有效",
  "enabled": true,
  "remark": "注册场景使用"
}
```

`name`、`code`、`content` 必填；`code` 全局唯一，重复时返回 400 及 `短信模板编码已存在: {code}`。

### PUT /api/v1/sms-template/update/{id}

更新指定 ID 的短信模板，请求体字段均可选（仅更新传入字段），
`code` 与其他模板冲突时返回 400，模板不存在时返回 404。

### DELETE /api/v1/sms-template/delete/{id}

删除指定 ID 的短信模板，不存在时返回 404。
