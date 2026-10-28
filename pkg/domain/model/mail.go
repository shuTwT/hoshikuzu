package model

// SendTestMailReq 发送测试邮件请求
type SendTestMailReq struct {
	To      string `json:"to" validate:"required"` // 收件人邮箱
	Subject string `json:"subject"`                // 主题（可选）
	Content string `json:"content"`                // 正文（可选）
}
