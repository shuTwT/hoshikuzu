package model

import "time"

type FriendCircleRuleResp struct {
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// 规则名称
	Name string `json:"name,omitempty"`
	// 标题选择器
	TitleSelector *string `json:"title_selector"`
	// 链接选择器
	LinkSelector *string `json:"link_selector"`
	// 创建时间选择器
	CreatedSelector *string `json:"created_selector"`
	// 更新时间选择器
	UpdatedSelector *string `json:"updated_selector"`
}

type FriendCircleRecordResp struct {
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt *LocalTime `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt *LocalTime `json:"updated_at,omitempty"`
	// 作者
	Author string `json:"author,omitempty"`
	// 标题
	Title string `json:"title,omitempty"`
	// 链接
	LinkURL string `json:"link_url,omitempty"`
	// 头像地址
	AvatarURL   string `json:"avatar_url,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
}

type FriendCircleRuleSaveReq struct {
	// ID of the ent.
	ID *int `json:"id"`
	// 规则名称
	Name string `json:"name"`
	// 标题选择器
	TitleSelector string `json:"title_selector"`
	// 链接选择器
	LinkSelector string `json:"link_selector"`
	// 创建时间选择器
	CreatedSelector string `json:"created_selector"`
	// 更新时间选择器
	UpdatedSelector string `json:"updated_selector"`
}

type FriendCircleRecordSaveReq struct {
	// ID of the ent.
	ID *int `json:"id"`
	// 作者
	Author string `json:"author"`
	// 标题
	Title string `json:"title"`
	// 链接
	LinkURL string `json:"link_url"`
	// 头像地址
	AvatarURL string `json:"avatar_url"`
}
