# 提示词画廊 Handler

内容管理（content）下的二级模块"提示词画廊"，负责提示词的 HTTP 接口。

一套提示词（`Prompt`）对应多张模型生成图片（`PromptImage`，每张记录图片地址与生成模型名），并可关联多个标签（`PromptTag`，多对多）。

| 方法 | 路由 | 说明 |
| --- | --- | --- |
| GET | `/api/v1/prompt/page` | 分页查询，支持 `keyword`（标题/正文）与 `tag_id` 过滤 |
| POST | `/api/v1/prompt/create` | 创建提示词，同时保存图片与标签 |
| PUT | `/api/v1/prompt/update/:id` | 更新提示词，`images`、`tags` 全量覆盖 |
| GET | `/api/v1/prompt/query/:id` | 查询详情（含图片与标签） |
| DELETE | `/api/v1/prompt/delete/:id` | 删除提示词及其图片 |
| GET | `/api/v1/prompt-tag/list` | 标签列表，供前端下拉选择 |

图片文件本身通过 `/api/v1/file/upload` 上传至基础设施中配置的存储策略，本模块只保存返回的图片地址。
