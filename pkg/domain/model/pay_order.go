package model

// PayOrderCreateReq represents the request body for creating a payment order.
type PayOrderCreateReq struct {
	ChannelType string `json:"channel_id" validate:"required"`
	OrderID     string `json:"order_id"`
	OutTradeNo  string `json:"out_trade_no" validate:"required"`
	TotalFee    string `json:"total_fee" validate:"required"`
	Subject     string `json:"subject" validate:"required"`
	Body        string `json:"body" validate:"required"`
	NotifyURL   string `json:"notify_url" validate:"required,url"`
	ReturnURL   string `json:"return_url" validate:"required,url"`
	Extra       string `json:"extra"`
	PayURL      string `json:"pay_url,omitempty"`
	State       string `json:"state"`
	ErrorMsg    string `json:"error_msg,omitempty"`
	Raw         string `json:"raw,omitempty"`
}

// PayOrderUpdateReq represents the request body for updating a payment order.
type PayOrderUpdateReq struct {
	ChannelType string `json:"channel_id"`
	OrderID     string `json:"order_id"`
	OutTradeNo  string `json:"out_trade_no"`
	TotalFee    string `json:"total_fee"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	NotifyURL   string `json:"notify_url"`
	ReturnURL   string `json:"return_url"`
	Extra       string `json:"extra"`
	PayURL      string `json:"pay_url,omitempty"`
	State       string `json:"state"`
	ErrorMsg    string `json:"error_msg,omitempty"`
	Raw         string `json:"raw,omitempty"`
}

// PayOrderResp represents the response body for a payment order.
type PayOrderResp struct {
	ID         int       `json:"id"`
	CreatedAt  LocalTime `json:"created_at"`
	ChannelID  string    `json:"channel_id"`
	OrderID    string    `json:"order_id"`
	OutTradeNo string    `json:"out_trade_no"`
	TotalFee   string    `json:"total_fee"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	NotifyURL  string    `json:"notify_url"`
	ReturnURL  string    `json:"return_url"`
	Extra      string    `json:"extra"`
	PayURL     *string   `json:"pay_url,omitempty"`
	State      string    `json:"state"`
	ErrorMsg   *string   `json:"error_msg,omitempty"`
}

const (
	PayOrderTypePost    = "1"
	PayOrderTypeProduct = "2"
)

type PayOrderSubmitReq struct {
	// 渠道类型 1 支付宝 2 微信 3 银联
	ChannelType string `json:"channel_type" validate:"required"`
	// 返回地址
	ReturnUrl string `json:"return_url" validate:"required,url"`
	// 订单类型 1 文章付费 2 商品购买
	OrderType string `json:"order_type" validate:"required"`
	// 商品名称
	Name string `json:"name"`
	// 金额
	Money int `json:"money"`
	// 文章 id，可选
	PostId int `json:"post_id"`
	// 商品 id，可选
	ProductId int `json:"product_id"`
}
