package model

// PromptImageReq 提示词模型图片请求
type PromptImageReq struct {
	ImageURL  string `json:"image_url" validate:"required"`
	ModelName string `json:"model_name,omitempty"`
	Sort      int    `json:"sort"`
}

// PromptCreateReq 创建提示词请求
type PromptCreateReq struct {
	Title          string           `json:"title" validate:"required"`
	Content        string           `json:"content" validate:"required"`
	NegativePrompt string           `json:"negative_prompt,omitempty"`
	Description    string           `json:"description,omitempty"`
	Sort           int              `json:"sort"`
	Images         []PromptImageReq `json:"images"`
	Tags           []string         `json:"tags"`
}

// PromptUpdateReq 更新提示词请求，images 与 tags 为全量覆盖
type PromptUpdateReq struct {
	Title          string           `json:"title" validate:"required"`
	Content        string           `json:"content" validate:"required"`
	NegativePrompt string           `json:"negative_prompt,omitempty"`
	Description    string           `json:"description,omitempty"`
	Sort           int              `json:"sort"`
	Images         []PromptImageReq `json:"images"`
	Tags           []string         `json:"tags"`
}

// PromptPageQuery 提示词分页查询
type PromptPageQuery struct {
	Page    int    `json:"page" query:"page"`
	Size    int    `json:"page_size" query:"page_size"`
	Keyword string `json:"keyword" query:"keyword"`
	TagID   int    `json:"tag_id" query:"tag_id"`
}
