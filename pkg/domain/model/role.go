package model

// RoleCreateReq represents the request body for creating a role.

type RoleCreateReq struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"is_default"`
}

// RoleUpdateReq represents the request body for updating a role.

type RoleUpdateReq struct {
	Name        string `json:"name,omitempty"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"is_default,omitempty"`
}

// RoleResp represents the response body for a role.

type RoleResp struct {
	ID          int       `json:"id"`
	CreatedAt   LocalTime `json:"created_at"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	IsDefault   bool      `json:"is_default"`
}
