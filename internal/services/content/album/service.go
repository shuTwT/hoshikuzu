package album

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/album"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type AlbumService interface {
	ListAlbum(ctx context.Context) ([]*ent.Album, error)
	ListAlbumPage(ctx context.Context, page, size int) (int, []*ent.Album, error)
	CreateAlbum(ctx context.Context, req *model.AlbumCreateReq) (*ent.Album, error)
	UpdateAlbum(ctx context.Context, id int, req *model.AlbumUpdateReq) (*ent.Album, error)
	QueryAlbum(ctx context.Context, id int) (*ent.Album, error)
	DeleteAlbum(ctx context.Context, id int) error
}

type AlbumServiceImpl struct {
	client *ent.Client
}

func NewAlbumServiceImpl(client *ent.Client) *AlbumServiceImpl {
	return &AlbumServiceImpl{client: client}
}

func (s *AlbumServiceImpl) ListAlbum(ctx context.Context) ([]*ent.Album, error) {
	albums, err := s.client.Album.Query().
		Order(ent.Desc(album.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *AlbumServiceImpl) ListAlbumPage(ctx context.Context, page, size int) (int, []*ent.Album, error) {
	count, err := s.client.Album.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	albums, err := s.client.Album.Query().
		Order(ent.Desc(album.FieldID)).
		Limit(size).
		Offset((page - 1) * size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, albums, nil
}

func (s *AlbumServiceImpl) CreateAlbum(ctx context.Context, req *model.AlbumCreateReq) (*ent.Album, error) {
	newAlbum, err := s.client.Album.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetSort(req.Sort).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newAlbum, nil
}

func (s *AlbumServiceImpl) UpdateAlbum(ctx context.Context, id int, req *model.AlbumUpdateReq) (*ent.Album, error) {
	updatedAlbum, err := s.client.Album.UpdateOneID(id).
		SetName(req.Name).
		SetDescription(req.Description).
		SetSort(req.Sort).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedAlbum, nil
}

func (s *AlbumServiceImpl) QueryAlbum(ctx context.Context, id int) (*ent.Album, error) {
	album, err := s.client.Album.Query().
		Where(album.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *AlbumServiceImpl) DeleteAlbum(ctx context.Context, id int) error {
	err := s.client.Album.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
