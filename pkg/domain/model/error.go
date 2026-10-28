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
