# 提示词画廊 Service

提示词画廊的业务逻辑层。

- 主表 `Prompt`：标题、提示词正文、负向提示词、描述、排序。
- 子表 `PromptImage`：图片地址、生成所用模型名、排序；通过 `prompt_id` 关联主表。
- 标签 `PromptTag`：按名称唯一，通过 `prompt_tag_relations` 与提示词多对多关联。

创建/更新/删除均在 ent 事务内完成：更新时先清空旧标签和旧图片，再按请求全量重建；
标签按名称查找，不存在时自动创建。业务错误定义在 `pkg/domain/model/error.go`（`ErrPromptNotFound` 等）。
