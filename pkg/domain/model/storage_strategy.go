package model

type StorageStrategyPageReq struct {
	Page  int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size  int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
	Name  string `json:"name" query:"name" form:"name"`
	Type  string `json:"type" query:"type" form:"type" validate:"omitempty,oneof=local s3"`
	Master *bool  `json:"master" query:"master" form:"master"`
}

// StorageStrategyCreateReq represents the request body for creating a storage strategy.
type StorageStrategyCreateReq struct {
	Name      string `json:"name" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=local s3"`
	NodeID    string `json:"node_id"`
	Endpoint  string `json:"endpoint" validate:"required_if=Type s3"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket" validate:"required_if=Type s3"`
	AccessKey string `json:"access_key" validate:"required_if=Type s3"`
	SecretKey string `json:"secret_key" validate:"required_if=Type s3"`
	BasePath  string `json:"base_path"`
	Domain    string `json:"domain" validate:"required,url"`
	Master    bool   `json:"master"`
}

// StorageStrategyUpdateReq represents the request body for updating a storage strategy.
type StorageStrategyUpdateReq struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty" validate:"required,oneof=local s3"`
	NodeID    string `json:"node_id"`
	Endpoint  string `json:"endpoint,omitempty"`
	Region    string `json:"region,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	BasePath  string `json:"base_path,omitempty"`
	Domain    string `json:"domain,omitempty" validate:"required,url"`
	Master    bool   `json:"master,omitempty"`
}

// StorageStrategyResp represents the response body for a storage strategy.
type StorageStrategyResp struct {
	ID        int       `json:"id"`
	CreatedAt LocalTime `json:"created_at"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	NodeID    string    `json:"node_id"`
	Endpoint  string    `json:"endpoint"`
	Region    string    `json:"region"`
	Bucket    string    `json:"bucket"`
	BasePath  string    `json:"base_path"`
	Domain    string    `json:"domain"`
	Master    bool      `json:"master"`
}

type StorageStrategyListResp struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Master bool   `json:"master"`
}
