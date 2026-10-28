package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// 提示词画廊-提示词
type Prompt struct {
	ent.Schema
}

func (Prompt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Prompt.
func (Prompt) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty().MaxLen(255).Comment("标题"),
		field.Text("content").NotEmpty().Comment("提示词正文"),
		field.Text("negative_prompt").Optional().Comment("负向提示词"),
		field.String("description").Optional().MaxLen(500).Comment("描述"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the Prompt.
func (Prompt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("images", PromptImage.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("tags", PromptTag.Type).
			StorageKey(edge.Table("prompt_tag_relations"), edge.Columns("prompt_id", "prompt_tag_id")),
	}
}
