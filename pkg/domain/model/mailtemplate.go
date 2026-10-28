package model

type MailTemplatePageQuery struct {
	Page int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
	Name string `json:"name" query:"name" form:"name"`
	Code string `json:"code" query:"code" form:"code"`
}

type MailTemplateCreateReq struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Subject     string  `json:"subject" validate:"required"`
	Content     string  `json:"content" validate:"required"`
	ContentType *string `json:"content_type,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Remark      *string `json:"remark,omitempty"`
}

type MailTemplateUpdateReq struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Subject     *string `json:"subject,omitempty"`
	Content     *string `json:"content,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Remark      *string `json:"remark,omitempty"`
}
