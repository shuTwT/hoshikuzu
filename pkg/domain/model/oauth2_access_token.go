package model

import "time"

// Oauth2AccessTokenCreateReq represents the request body for creating an OAuth2 access token.

type Oauth2AccessTokenCreateReq struct {
	UserID       int       `json:"user_id" validate:"required"`
	AccessToken  string    `json:"access_token" validate:"required"`
	RefreshToken string    `json:"refresh_token" validate:"required"`
	ClientID     string    `json:"client_id" validate:"required"`
	Scope        string    `json:"scope" validate:"required"`
	ExpireAt     time.Time `json:"expire_at" validate:"required"`
}

// Oauth2AccessTokenUpdateReq represents the request body for updating an OAuth2 access token.

type Oauth2AccessTokenUpdateReq struct {
	AccessToken  *string    `json:"access_token,omitempty"`
	RefreshToken *string    `json:"refresh_token,omitempty"`
	Scope        *string    `json:"scope,omitempty"`
	ExpireAt     *time.Time `json:"expire_at,omitempty"`
}

// Oauth2AccessTokenResp represents the response body for an OAuth2 access token.

type Oauth2AccessTokenResp struct {
	ID           int       `json:"id"`
	CreatedAt    LocalTime `json:"created_at"`
	UserID       int       `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ClientID     string    `json:"client_id"`
	Scope        string    `json:"scope"`
	ExpireAt     LocalTime `json:"expire_at"`
}
