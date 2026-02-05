package model

type VisitLogReq struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Path      string `json:"path"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
	Device    string `json:"device"`
}

type VisitLogPageQuery struct {
	Page int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
	IP   string `json:"ip" query:"ip" form:"ip"`
	Path string `json:"path" query:"path" form:"path"`
}

type VisitLogBatchDeleteReq struct {
	IDs []int `json:"ids" validate:"required,min=1"`
}

type VisitLogResp struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt LocalTime `json:"created_at,omitempty"`
	IP        string    `json:"ip,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
	Path      string    `json:"path,omitempty"`
	Os        *string   `json:"os,omitempty"`
	Browser   *string   `json:"browser,omitempty"`
	Device    *string   `json:"device,omitempty"`
}
