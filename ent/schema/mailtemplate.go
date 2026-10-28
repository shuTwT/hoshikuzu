package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type MailTemplate struct {
	ent.Schema
}

func (MailTemplate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (MailTemplate) Fields() []ent.Field {
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
		field.String("subject").
			NotEmpty().
			MaxLen(200).
			Comment("邮件主题，支持 {变量} 占位符"),
		field.String("content").
			NotEmpty().
			MaxLen(5000).
			Comment("邮件正文，支持 {变量} 占位符与 HTML"),
		field.String("content_type").
			Default("html").
			MaxLen(10).
			Comment("内容类型：html / text"),
		field.Bool("enabled").
			Default(true).
			Comment("是否启用"),
		field.String("remark").
			Optional().
			MaxLen(255).
			Comment("备注"),
	}
}

func (MailTemplate) Edges() []ent.Edge {
	return []ent.Edge{}
}
