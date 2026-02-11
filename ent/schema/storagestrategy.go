package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// 存储策略
type StorageStrategy struct {
	ent.Schema
}

func (StorageStrategy) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StorageStrategy.
func (StorageStrategy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty().Comment("策略名称"),
		field.Enum("type").Values("local", "s3").Default("local").Comment("策略类型"),
		field.String("node_id").Default("").Comment("节点 ID"),
		field.String("endpoint").Default("").Comment("端点"),
		field.String("region").Default("").Comment("region"),
		field.String("bucket").Default("").Comment("存储桶名称"),
		field.String("access_key").Default("").Comment("accessKey"),
		field.String("secret_key").Default("").Comment("secret_key"),
		field.String("base_path").Default("").Comment("local 基础路径"),
		field.String("domain").Default("").Comment("访问域名"),
		field.Bool("master").Default(false).Comment("是否为默认策略"),
	}
}

// Edges of the StorageStrategy.
func (StorageStrategy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("files", File.Type).
			Ref("storage_strategy"),
	}
}
