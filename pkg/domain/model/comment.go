package model

// CommentCreateReq represents the request body for creating a comment.

type CommentCreateReq struct {
	PostID  int    `json:"post_id" validate:"required_without=PageID"`
	PageID  int    `json:"page_id" validate:"required_without=PostID"`
	URL     string `json:"url"`
	Content string `json:"content" validate:"required,max=1024"`
	UserID  int    `json:"user_id" validate:"required"`
}

// CommentUpdateReq represents the request body for updating a comment.

type CommentUpdateReq struct {
	Content *string `json:"content,omitempty" validate:"omitempty,max=1024"`
	Status  *int    `json:"status,omitempty"`
	Pinned  *bool   `json:"pinned,omitempty"`
}

// CommentResp represents the response body for a comment.

type CommentResp struct {
	ID         int       `json:"id"`
	CreatedAt  LocalTime `json:"created_at"`
	PostID     int       `json:"post_id"`
	PageID     int       `json:"page_id"`
	Content    string    `json:"content"`
	UserID     int       `json:"user_id"`
	Status     int       `json:"status"`
	UserAgent  *string   `json:"user_agent,omitempty"`
	IPAddress  string    `json:"ip_address"`
	IPLocation *string   `json:"ip_location,omitempty"`
	Pinned     bool      `json:"pinned"`
}
