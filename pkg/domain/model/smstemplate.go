package model

type SmsTemplatePageQuery struct {
	Page int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
	Name string `json:"name" query:"name" form:"name"`
	Code string `json:"code" query:"code" form:"code"`
}

type SmsTemplateCreateReq struct {
	Name               string  `json:"name" validate:"required"`
	Code               string  `json:"code" validate:"required"`
	Content            string  `json:"content" validate:"required"`
	SignName           *string `json:"sign_name,omitempty"`
	ProviderTemplateID *string `json:"provider_template_id,omitempty"`
	Enabled            *bool   `json:"enabled,omitempty"`
	Remark             *string `json:"remark,omitempty"`
}

type SmsTemplateUpdateReq struct {
	Name               *string `json:"name,omitempty"`
	Code               *string `json:"code,omitempty"`
	Content            *string `json:"content,omitempty"`
	SignName           *string `json:"sign_name,omitempty"`
	ProviderTemplateID *string `json:"provider_template_id,omitempty"`
	Enabled            *bool   `json:"enabled,omitempty"`
	Remark             *string `json:"remark,omitempty"`
}
