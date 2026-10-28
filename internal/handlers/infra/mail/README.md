# Mail Handler

邮件模块的 HTTP 处理器，提供 SMTP 配置连通性验证接口。

## 功能

- 发送测试邮件

## 接口

### POST /api/v1/mail/test

使用当前保存的 SMTP 配置向指定邮箱发送测试邮件，用于验证邮件配置是否正确。

请求体：

```json
{
  "to": "user@example.com",
  "subject": "测试主题（可选）",
  "content": "测试正文（可选）"
}
```

成功返回 `{ "code": 200, "msg": "success" }`；
SMTP 配置缺失、未保存或发送失败时返回对应错误信息。
