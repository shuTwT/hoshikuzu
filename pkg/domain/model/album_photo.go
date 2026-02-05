package model

// AlbumPhotoCreateReq represents the request body for creating an album photo.

type AlbumPhotoCreateReq struct {
	ImageURL    string `json:"image_url" validate:"required,url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AlbumID     int    `json:"album_id" validate:"required"`
}

// AlbumPhotoUpdateReq represents the request body for updating an album photo.

type AlbumPhotoUpdateReq struct {
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AlbumID     int    `json:"album_id,omitempty"`
}

// AlbumPhotoResp represents the response body for an album photo.

type AlbumPhotoResp struct {
	ID        int       `json:"id"`
	CreatedAt LocalTime `json:"created_at"`
	ImageURL  string    `json:"image_url"`
	ViewCount int       `json:"view_count"`
	AlbumID   int       `json:"album_id"`
}
