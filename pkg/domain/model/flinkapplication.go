package model

type FlinkApplicationCreateReq struct {
	WebsiteURL         string `json:"website_url" validate:"required"`
	ApplicationType    string `json:"application_type" validate:"required,oneof=create update"`
	WebsiteName        string `json:"website_name" validate:"required"`
	WebsiteLogo        string `json:"website_logo" validate:"required"`
	WebsiteDescription string `json:"website_description" validate:"required"`
	ContactEmail       string `json:"contact_email" validate:"required,email"`
	SnapshotURL        string `json:"snapshot_url"`
	OriginalWebsiteURL string `json:"original_website_url"`
	ModificationReason string `json:"modification_reason"`
}

type FlinkApplicationUpdateReq struct {
	ID           int    `json:"id"`
	Status       int    `json:"status" validate:"required,oneof=1 2"`
	RejectReason string `json:"reject_reason"`
}

type FlinkApplicationResp struct {
	ID                 int       `json:"id"`
	CreatedAt          LocalTime `json:"created_at"`
	UpdatedAt          LocalTime `json:"updated_at"`
	WebsiteURL         string    `json:"website_url"`
	ApplicationType    string    `json:"application_type"`
	WebsiteName        string    `json:"website_name"`
	WebsiteLogo        string    `json:"website_logo"`
	WebsiteDescription string    `json:"website_description"`
	ContactEmail       string    `json:"contact_email"`
	SnapshotURL        string    `json:"snapshot_url"`
	OriginalWebsiteURL string    `json:"original_website_url"`
	ModificationReason string    `json:"modification_reason"`
	Status             int       `json:"status"`
	RejectReason       string    `json:"reject_reason"`
}

type FlinkApplicationPageReq struct {
	Status          *int    `json:"status" query:"status" form:"status"`
	ApplicationType *string `json:"application_type" query:"application_type" form:"application_type"`
	Page            int     `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size            int     `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
}
