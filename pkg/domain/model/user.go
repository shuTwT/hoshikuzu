/**
 * 用户
 */
package model

import (
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
)

// UserCreateReq represents the request body for creating a user.

type UserCreateReq struct {
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number,omitempty"`
	RoleID      int    `json:"role_id,omitempty"`
}

// UserUpdateReq represents the request body for updating a user.

type UserUpdateReq struct {
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty" validate:"min=8"`
	PhoneNumber string `json:"phone_number,omitempty"`
	RoleID      *int   `json:"role_id,omitempty"`
}

// UserResp represents the response body for a user.

type UserResp struct {
	ID                  int        `json:"id"`
	CreatedAt           *time.Time `json:"created_at"`
	Email               string     `json:"email"`
	EmailVerified       bool       `json:"email_verified"`
	Name                string     `json:"name"`
	PhoneNumber         *string    `json:"phone_number"`
	PhoneNumberVerified bool       `json:"phone_number_verified"`
	RoleID              *int       `json:"role_id,omitempty"`
	Role                *ent.Role  `json:"role,omitempty"`
}

type UserProfileResp struct {
	UserID              int       `json:"user_id"`
	Email               string    `json:"email"`
	EmailVerified       bool      `json:"email_verified"`
	Name                string    `json:"name"`
	PhoneNumber         *string   `json:"phone_number"`
	PhoneNumberVerified bool      `json:"phone_number_verified"`
	RoleID              *int      `json:"role_id,omitempty"`
	Role                *ent.Role `json:"role,omitempty"`
}

type UserSearchReq struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"required"`
	Page    int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size    int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
}

type UserSearchResp struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
}
