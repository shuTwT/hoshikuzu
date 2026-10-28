package model

import "errors"

type HttpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewError(code int, msg string) HttpError {
	return HttpError{
		Code: code,
		Msg:  msg,
	}
}

// 邮件相关业务错误
var (
	// ErrMailConfigIncomplete SMTP 配置不完整
	ErrMailConfigIncomplete = errors.New("邮件配置不完整")
	// ErrMailEmptyRecipient 收件人为空
	ErrMailEmptyRecipient = errors.New("收件人为空")
	// ErrMailSendFailed 邮件发送失败
	ErrMailSendFailed = errors.New("邮件发送失败")
)

// 短信模板相关业务错误
var (
	// ErrSmsTemplateNotFound 短信模板不存在
	ErrSmsTemplateNotFound = errors.New("短信模板不存在")
	// ErrSmsTemplateCodeExists 短信模板编码已存在
	ErrSmsTemplateCodeExists = errors.New("短信模板编码已存在")
)

// 邮件模板相关业务错误
var (
	// ErrMailTemplateNotFound 邮件模板不存在
	ErrMailTemplateNotFound = errors.New("邮件模板不存在")
	// ErrMailTemplateCodeExists 邮件模板编码已存在
	ErrMailTemplateCodeExists = errors.New("邮件模板编码已存在")
)

// 提示词画廊相关业务错误
var (
	// ErrPromptNotFound 提示词不存在
	ErrPromptNotFound = errors.New("提示词不存在")
	// ErrPromptImageURLEmpty 提示词图片地址为空
	ErrPromptImageURLEmpty = errors.New("提示词图片地址不能为空")
)

// 模拟支付相关业务错误
var (
	// ErrMockPayDisabled 模拟支付未开启
	ErrMockPayDisabled = errors.New("模拟支付未开启")
	// ErrMockPayOrderInvalid 订单不是待支付的模拟支付订单
	ErrMockPayOrderInvalid = errors.New("订单不是待支付的模拟支付订单")
)

// 系统初始化相关业务错误
var (
	// ErrSystemInitialized 系统已初始化
	ErrSystemInitialized = errors.New("系统已初始化，请勿重复操作")
)
