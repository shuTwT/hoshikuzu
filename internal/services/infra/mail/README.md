# Mail Service

邮件模块的服务层，负责读取已保存的 SMTP 配置并发送邮件。

## 功能

- 发送测试邮件（验证 SMTP 配置是否可用）
- 读取当前保存的邮件配置（email 组）

## 配置来源

邮件配置存储于通用设置的 `email` 键（JSON 结构），由前端「通知设置」页保存，字段包括：

| 字段 | 说明 |
| --- | --- |
| smtpHost | SMTP 服务器地址 |
| smtpPort | SMTP 端口 |
| smtpUsername | 登录用户名 |
| smtpPassword | 登录密码 |
| smtpEncryption | 加密方式：tls / ssl / none |
| senderEmail | 发件人邮箱 |
| senderName | 发件人名称 |

## 方法

### SendTestEmail
使用当前保存的 SMTP 配置发送一封测试邮件到指定邮箱，主题与正文为空时使用默认值。

### GetMailConfig
读取 `email` 键并解析为 `EmailConfig` 结构。

## 复用发送能力

其他模块可通过 `EmailConfig.ToMailInfraConfig()` 转换为 `internal/infra/mail.Config`，
再调用 `Send` / `SendMail` 发送业务邮件。
