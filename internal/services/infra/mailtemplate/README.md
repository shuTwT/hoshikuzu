# MailTemplate Service

邮件模板模块的服务层，负责邮件模板的增删改查与编码唯一性校验。

## 功能

- 分页查询邮件模板（支持按模板名称、模板编码模糊搜索）
- 查询指定 ID 的邮件模板详情
- 创建邮件模板
- 更新邮件模板
- 删除邮件模板

## 模板字段

| 字段 | 说明 |
| --- | --- |
| name | 模板名称（必填） |
| code | 模板编码（必填，全局唯一，业务方按此编码查找模板） |
| subject | 邮件主题（必填，支持 `{变量}` 占位符） |
| content | 邮件正文（必填，支持 `{变量}` 占位符与 HTML） |
| content_type | 内容类型：html / text（默认 html） |
| enabled | 是否启用（默认启用） |
| remark | 备注（可选） |

## 方法

### QueryMailTemplate
查询指定 ID 的邮件模板，不存在时返回 `model.ErrMailTemplateNotFound`。

### QueryMailTemplatePage
按 `MailTemplatePageQuery` 分页查询，name / code 为模糊匹配条件，按 ID 倒序返回。

### CreateMailTemplate
创建前按 code 查重，编码重复（含并发下数据库唯一约束冲突）返回 `model.ErrMailTemplateCodeExists`；
`content_type` 为空时默认 `html`。

### UpdateMailTemplate
更新指定 ID 的模板；仅当传入新 code 时校验唯一性（排除自身 ID）。

### DeleteMailTemplate
删除指定 ID 的模板，不存在时返回 `model.ErrMailTemplateNotFound`。

## 后续扩展

邮件发送能力接入后，可按 code 查找启用状态的模板，
将 `subject` / `content` 中的 `{变量}` 占位符替换为实际参数，
按 `content_type` 决定以 HTML 或纯文本方式调用 `internal/infra/mail` 发送。
