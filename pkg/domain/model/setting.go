package model

// SettingCreateReq represents the request body for creating a setting.
type SettingCreateReq struct {
	Key     string  `json:"key" validate:"required"`
	Value   string  `json:"value" validate:"required"`
	Comment *string `json:"comment,omitempty"`
}

// SettingUpdateReq represents the request body for updating a setting.
type SettingUpdateReq struct {
	Value   *string `json:"value,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

// SettingResp represents the response body for a setting.
type SettingResp struct {
	ID        int       `json:"id"`
	CreatedAt LocalTime `json:"created_at"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Comment   *string   `json:"comment,omitempty"`
}
