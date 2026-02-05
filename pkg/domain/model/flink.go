package model

type FlinkCreateReq struct {
	// 名称
	Name string `json:"name"`
	// 链接
	URL string `json:"url"`
	// logo
	AvatarURL string `json:"avatar_url"`
	// 简介
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	// 快照
	SnapshotURL string `json:"snapshot_url"`
	// 邮箱
	Email string `json:"email"`
	// 是否开启朋友圈
	EnableFriendCircle bool `json:"enable_friend_circle"`
	// 朋友圈解析规则
	FriendCircleRuleID *int `json:"friend_circle_rule_id"`
	GroupID            int  `json:"group_id" validate:"required"`
}

type FlinkUpdateReq struct {
	// ID of the ent.
	ID int `json:"id"`
	// 名称
	Name string `json:"name"`
	// 链接
	URL string `json:"url"`
	// logo
	AvatarURL string `json:"avatar_url"`
	CoverURL  string `json:"cover_url"`
	// 简介
	Description string `json:"description"`
	// 快照
	SnapshotURL string `json:"snapshot_url"`
	// 邮箱
	Email string `json:"email"`
	// 是否开启朋友圈
	EnableFriendCircle bool `json:"enable_friend_circle"`
	// 朋友圈解析规则
	FriendCircleRuleID int `json:"friend_circle_rule_id"`
	GroupID            int `json:"group_id" validate:"required"`
}

type FlinkResp struct {
	// ID of the ent.
	ID int `json:"id"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt LocalTime `json:"created_at"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt LocalTime `json:"updated_at"`
	// 名称
	Name string `json:"name"`
	// 链接
	URL string `json:"url"`
	// logo
	AvatarURL string `json:"avatar_url"`
	// 简介
	Description string `json:"description"`
	// 状态
	Status   int    `json:"status"`
	CoverURL string `json:"cover_url"`
	// 快照
	SnapshotUrl string `json:"snapshot_url"`
	// 邮箱
	Email string `json:"email"`
	// 是否开启朋友圈
	EnableFriendCircle bool `json:"enable_friend_circle"`
	// 朋友圈解析规则
	FriendCircleRuleID *int            `json:"friend_circle_rule_id"`
	Group              *FlinkGroupResp `json:"group,omitempty"`
}

type FlinkListReq struct {
	GroupId   *int    `json:"group_id" query:"group_id"`
	GroupName *string `json:"group_name" query:"group_name"`
}

type FlinkPageReq struct {
	GroupId *int `json:"group_id" query:"group_id" form:"group_id"`
	Page    int  `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size    int  `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
}

type FlinkRandomReq struct {
	Limit int `json:"limit" query:"limit" form:"limit" validate:"required,min=1,max=100"`
}
