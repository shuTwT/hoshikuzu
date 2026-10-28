package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type SmsTemplate struct {
	ent.Schema
}

func (SmsTemplate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (SmsTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(100).
			Comment("模板名称"),
		field.String("code").
			NotEmpty().
			MaxLen(100).
			Unique().
			Comment("模板编码"),
		field.String("sign_name").
			Optional().
			MaxLen(50).
			Comment("短信签名"),
		field.String("provider_template_id").
			Optional().
			MaxLen(100).
			Comment("服务商模板ID"),
		field.String("content").
			NotEmpty().
			MaxLen(500).
			Comment("模板内容，支持 {变量} 占位符"),
		field.Bool("enabled").
			Default(true).
			Comment("是否启用"),
		field.String("remark").
			Optional().
			MaxLen(255).
			Comment("备注"),
	}
}

func (SmsTemplate) Edges() []ent.Edge {
	return []ent.Edge{}
}
