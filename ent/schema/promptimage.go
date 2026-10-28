package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// 提示词画廊-模型生成图片
type PromptImage struct {
	ent.Schema
}

func (PromptImage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the PromptImage.
func (PromptImage) Fields() []ent.Field {
	return []ent.Field{
		field.Int("prompt_id").Comment("提示词ID"),
		field.String("image_url").NotEmpty().Comment("图片地址"),
		field.String("model_name").Optional().MaxLen(100).Comment("生成所用模型名"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the PromptImage.
func (PromptImage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("prompt", Prompt.Type).
			Ref("images").
			Field("prompt_id").
			Unique().
			Required(),
	}
}
