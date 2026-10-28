package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// 提示词画廊-标签
type PromptTag struct {
	ent.Schema
}

func (PromptTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the PromptTag.
func (PromptTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(50).Unique().Comment("标签名称"),
	}
}

// Edges of the PromptTag.
func (PromptTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("prompts", Prompt.Type).Ref("tags"),
	}
}
