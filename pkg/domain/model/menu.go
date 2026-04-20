package model

type MenuCreateReq struct {
	Name      string  `json:"name" validate:"required"`
	Title     *string `json:"title,omitempty"`
	Path      *string `json:"path,omitempty"`
	Icon      *string `json:"icon,omitempty"`
	ParentID  *int    `json:"parent_id,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
	Visible   *bool   `json:"visible,omitempty"`
	Target    *string `json:"target,omitempty"`
}

type MenuUpdateReq struct {
	Name      *string `json:"name,omitempty"`
	Title     *string `json:"title,omitempty"`
	Path      *string `json:"path,omitempty"`
	Icon      *string `json:"icon,omitempty"`
	ParentID  *int    `json:"parent_id,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
	Visible   *bool   `json:"visible,omitempty"`
	Target    *string `json:"target,omitempty"`
}

type MenuResp struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	ParentID  int    `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	Visible   bool   `json:"visible"`
	Target    string `json:"target"`
}
