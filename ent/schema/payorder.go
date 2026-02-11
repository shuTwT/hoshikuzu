package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// 支付订单
type PayOrder struct {
	ent.Schema
}

func (PayOrder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the PayOrder.
func (PayOrder) Fields() []ent.Field {
	return []ent.Field{
		field.String("channel_type").
			Optional().
			Nillable().
			Comment("支付渠道类型"),
		field.String("order_id").
			Optional().
			Nillable().
			Comment("支付订单ID"),
		field.String("merchant_order_id").
			Optional().
			Nillable().
			Comment("系统内部商户订单号，按规则生成，唯一"),
		field.String("out_trade_no").
			Optional().
			Nillable().
			Comment("外部订单号,对应支付渠道那里的"),
		field.Int("order_price").
			Optional().
			Comment("订单金额,单位分"),
		field.Int("price").
			Optional().
			Comment("支付金额,单位分"),
		field.Int("channel_fee_price").
			Optional().
			Comment("渠道手续费金额,单位分"),
		field.String("subject").
			Optional().
			Comment("订单标题"),
		field.String("body").
			Optional().
			Comment("订单描述"),
		field.String("notify_url").
			Optional().
			Comment("异步通知地址"),
		field.String("return_url").
			Optional().
			Nillable().
			Comment("同步通知地址"),
		field.String("extra").
			Optional().
			Nillable().
			Comment("额外参数"),
		field.String("pay_url").
			Optional().
			Nillable().
			Comment("支付链接"),
		field.String("state").
			Optional().
			Default("1").
			Comment("支付状态"),
		field.String("error_msg").
			Optional().
			Comment("错误信息"),
		field.String("raw").
			Optional().
			Comment("原始响应"),
	}
}

// Edges of the PayOrder.
func (PayOrder) Edges() []ent.Edge {
	return nil
}
