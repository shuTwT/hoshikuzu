package albumphoto

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/albumphoto"
)

type AlbumPhotoService interface {
	ListAlbumPhoto(ctx context.Context) ([]*ent.AlbumPhoto, error)
	ListAlbumPhotoPage(ctx context.Context, page, size int) (int, []*ent.AlbumPhoto, error)
	CreateAlbumPhoto(ctx context.Context, albumID int, name, imageURL, description string) (*ent.AlbumPhoto, error)
	UpdateAlbumPhoto(ctx context.Context, id int, name, imageURL, description string) (*ent.AlbumPhoto, error)
	QueryAlbumPhoto(ctx context.Context, id int) (*ent.AlbumPhoto, error)
	DeleteAlbumPhoto(ctx context.Context, id int) error
}

type AlbumPhotoServiceImpl struct {
	client *ent.Client
}

func NewAlbumPhotoServiceImpl(client *ent.Client) *AlbumPhotoServiceImpl {
	return &AlbumPhotoServiceImpl{client: client}
}

func (s *AlbumPhotoServiceImpl) ListAlbumPhoto(ctx context.Context) ([]*ent.AlbumPhoto, error) {
	albumPhotos, err := s.client.AlbumPhoto.Query().
		Order(ent.Desc(albumphoto.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return albumPhotos, nil
}

func (s *AlbumPhotoServiceImpl) ListAlbumPhotoPage(ctx context.Context, page, size int) (int, []*ent.AlbumPhoto, error) {
	count, err := s.client.AlbumPhoto.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	albumPhotos, err := s.client.AlbumPhoto.Query().
		Order(ent.Desc(albumphoto.FieldID)).
		Limit(size).
		Offset((page - 1) * size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, albumPhotos, nil
}

func (s *AlbumPhotoServiceImpl) CreateAlbumPhoto(ctx context.Context, albumID int, name, imageURL, description string) (*ent.AlbumPhoto, error) {
	newAlbumPhoto, err := s.client.AlbumPhoto.Create().
		SetAlbumID(albumID).
		SetName(name).
		SetImageURL(imageURL).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newAlbumPhoto, nil
}

func (s *AlbumPhotoServiceImpl) UpdateAlbumPhoto(ctx context.Context, id int, name, imageURL, description string) (*ent.AlbumPhoto, error) {
	updatedAlbumPhoto, err := s.client.AlbumPhoto.UpdateOneID(id).
		SetName(name).
		SetImageURL(imageURL).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedAlbumPhoto, nil
}

func (s *AlbumPhotoServiceImpl) QueryAlbumPhoto(ctx context.Context, id int) (*ent.AlbumPhoto, error) {
	albumPhoto, err := s.client.AlbumPhoto.Query().
		Where(albumphoto.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return albumPhoto, nil
}

func (s *AlbumPhotoServiceImpl) DeleteAlbumPhoto(ctx context.Context, id int) error {
	err := s.client.AlbumPhoto.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
