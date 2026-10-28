# SmsTemplate Service

短信模板模块的服务层，负责短信模板的增删改查与编码唯一性校验。

## 功能

- 分页查询短信模板（支持按模板名称、模板编码模糊搜索）
- 查询指定 ID 的短信模板详情
- 创建短信模板
- 更新短信模板
- 删除短信模板

## 模板字段

| 字段 | 说明 |
| --- | --- |
| name | 模板名称（必填） |
| code | 模板编码（必填，全局唯一，业务方按此编码查找模板） |
| sign_name | 短信签名（可选） |
| provider_template_id | 服务商模板 ID（可选，接入短信服务商时使用） |
| content | 模板内容（必填，支持 `{变量}` 占位符） |
| enabled | 是否启用（默认启用） |
| remark | 备注（可选） |

## 方法

### QuerySmsTemplate
查询指定 ID 的短信模板，不存在时返回 `model.ErrSmsTemplateNotFound`。

### QuerySmsTemplatePage
按 `SmsTemplatePageQuery` 分页查询，name / code 为模糊匹配条件，按 ID 倒序返回。

### CreateSmsTemplate
创建前按 code 查重，编码重复（含并发下数据库唯一约束冲突）返回 `model.ErrSmsTemplateCodeExists`。

### UpdateSmsTemplate
更新指定 ID 的模板；仅当传入新 code 时校验唯一性（排除自身 ID）。

### DeleteSmsTemplate
删除指定 ID 的模板，不存在时返回 `model.ErrSmsTemplateNotFound`。

## 后续扩展

短信发送能力接入后，可按 code 查找启用状态的模板，
将 `content` 中的 `{变量}` 占位符替换为实际参数后交由短信客户端发送。
