package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Menu struct {
	ent.Schema
}

func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(100).
			Comment("菜单名称"),
		field.String("title").
			MaxLen(200).
			Default("").
			Comment("菜单标题"),
		field.String("path").
			MaxLen(500).
			Default("").
			Comment("菜单路径"),
		field.String("icon").
			Optional().
			MaxLen(100).
			Comment("菜单图标"),
		field.Int("parent_id").
			Default(0).
			Comment("父菜单ID"),
		field.Int("sort_order").
			Default(0).
			Comment("排序"),
		field.Bool("visible").
			Default(true).
			Comment("是否可见"),
		field.String("target").
			Default("_self").
			MaxLen(20).
			Comment("打开方式"),
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{}
}
