package model

import "time"

// Oauth2RefreshTokenCreateReq represents the request body for creating an OAuth2 refresh token.

type Oauth2RefreshTokenCreateReq struct {
	UserID       int       `json:"user_id" validate:"required"`
	RefreshToken string    `json:"refresh_token" validate:"required"`
	ClientID     string    `json:"client_id" validate:"required"`
	Scope        string    `json:"scope" validate:"required"`
	ExpireAt     time.Time `json:"expire_at" validate:"required"`
}

// Oauth2RefreshTokenUpdateReq represents the request body for updating an OAuth2 refresh token.

type Oauth2RefreshTokenUpdateReq struct {
	ExpireAt *time.Time `json:"expire_at,omitempty"`
	Scope    *string    `json:"scope,omitempty"`
}

// Oauth2RefreshTokenResp represents the response body for an OAuth2 refresh token.

type Oauth2RefreshTokenResp struct {
	ID           int       `json:"id"`
	CreatedAt    LocalTime `json:"created_at"`
	UserID       int       `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ClientID     string    `json:"client_id"`
	Scope        string    `json:"scope"`
	ExpireAt     LocalTime `json:"expire_at"`
}
