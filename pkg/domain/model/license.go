package model

import "time"

type LicensePageReq struct {
	Page         int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size         int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
	Domain       string `json:"domain" query:"domain" form:"domain"`
	CustomerName string `json:"customer_name" query:"customer_name" form:"customer_name"`
	Status       *int   `json:"status" query:"status" form:"status"`
}

type LicenseCreateReq struct {
	Domain       string    `json:"domain" validate:"required"`
	CustomerName string    `json:"customer_name"`
	ExpireDate   time.Time `json:"expire_date" validate:"required"`
}

type LicenseUpdateReq struct {
	Domain       string    `json:"domain,omitempty"`
	LicenseKey   string    `json:"license_key,omitempty"`
	CustomerName string    `json:"customer_name,omitempty"`
	ExpireDate   time.Time `json:"expire_date,omitempty"`
	Status       int       `json:"status,omitempty"`
}

type LicenseResp struct {
	ID           int       `json:"id"`
	CreatedAt    LocalTime `json:"created_at"`
	UpdatedAt    LocalTime `json:"updated_at"`
	Domain       string    `json:"domain"`
	LicenseKey   string    `json:"license_key"`
	CustomerName string    `json:"customer_name"`
	ExpireDate   LocalTime `json:"expire_date"`
	Status       int       `json:"status"`
}

type LicenseVerifyReq struct {
	Domain string `json:"domain" validate:"required"`
}

type LicenseVerifyResp struct {
	Valid        bool   `json:"valid"`
	CustomerName string `json:"customer_name"`
	ExpireDate   string `json:"expire_date"`
	Message      string `json:"message"`
}
