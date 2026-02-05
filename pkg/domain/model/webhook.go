package model

// WebHookCreateReq represents the request body for creating a webhook.
type WebHookCreateReq struct {
	Name  string `json:"name" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
	Event string `json:"event" validate:"required"`
}

// WebHookUpdateReq represents the request body for updating a webhook.
type WebHookUpdateReq struct {
	Name  *string `json:"name,omitempty"`
	URL   *string `json:"url,omitempty" validate:"url"`
	Event *string `json:"event,omitempty"`
}

// WebHookResp represents the response body for a webhook.
type WebHookResp struct {
	ID        int       `json:"id"`
	CreatedAt LocalTime `json:"created_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Event     string    `json:"event"`
}
